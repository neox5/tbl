package tbl

import (
	"fmt"
	"sort"

	"github.com/neox5/btmp"
	"github.com/neox5/tbl/internal/cursor"
)

// ID identifies a cell in the table.
type ID int64

// Direction indicates traversal direction for flex cell scanning.
type Direction int

const (
	DirLeft Direction = iota
	DirRight
)

// Table manages incremental table construction with flex/static cells.
type Table struct {
	g     *btmp.Grid
	c     *cursor.Cursor
	cells map[ID]*Cell
	pos   map[int]map[int]ID // pos[row][col] = cell ID

	colsFixed  bool
	nextCellID ID
}

// New creates a new Table with zero columns.
func New() *Table {
	return NewWithCols(0)
}

// NewWithCols creates a new Table with initial column capacity.
func NewWithCols(cols int) *Table {
	if cols < 0 {
		panic("tbl: invalid cols value")
	}

	t := &Table{
		g:          btmp.NewGridWithSize(0, cols),
		c:          cursor.New(),
		cells:      make(map[ID]*Cell),
		pos:        make(map[int]map[int]ID),
		nextCellID: 1,
	}

	if cols > 0 {
		t.colsFixed = true
	}

	return t
}

// AddRow advances to next row with validation and cursor positioning.
func (t *Table) AddRow() *Table {
	row := t.c.Row()

	// Validate previous row if not first row
	if row >= 0 {
		if !t.isRowComplete(row) {
			panic(fmt.Sprintf("tbl: incomplete row %d before AddRow", row))
		}

		// Check if we can fix columns
		if !t.colsFixed && t.isRowStatic(row) {
			t.colsFixed = true
		}
	}

	// Ensure next row exists
	t.ensureRows(row + 1)

	// Advance cursor and get new row
	row = t.c.NextRow()

	// Position cursor at first free column
	freeCol := t.findFirstFreeCol(row)
	t.c.Advance(freeCol)

	return t
}

// AddCell adds a cell at cursor position with specified type and span.
// Expands columns if needed (when not fixed). Validates span fits in grid.
// Panics if: no row active, span invalid, insufficient columns (when fixed),
// or space occupied.
func (t *Table) AddCell(ct CellType, rowSpan, colSpan int) *Table {
	if rowSpan <= 0 || colSpan <= 0 {
		panic(fmt.Sprintf("tbl: invalid span rowSpan=%d colSpan=%d at cursor (%d,%d)", rowSpan, colSpan, t.c.Row(), t.c.Col()))
	}

	if t.c.Row() == -1 {
		panic(fmt.Sprintf("tbl: no row to add cell at cursor (%d,%d)", t.c.Row(), t.c.Col()))
	}

	row, col := t.c.Pos()

	// Ensure sufficient rows for cell span
	t.ensureRows(row + rowSpan - 1)

	// Step 1: Check if enough columns exist
	needed := col + colSpan
	if needed > t.g.Cols() {
		if err := t.ensureCols(col, colSpan); err != nil {
			panic(err.Error())
		}
	}

	// Step 2: Check if space is free
	if !t.g.CanFitWidth(row, col, colSpan) {
		// Space occupied - wall blocking
		if t.colsFixed {
			panic(fmt.Sprintf("tbl: space occupied at cursor (%d,%d), cannot expand", row, col))
		}

		// Attempt expansion
		ok, flexCells := t.traverseFlex(row, col)
		if !ok {
			panic(fmt.Sprintf("tbl: no flex cells available for expansion at cursor (%d,%d)", row, col))
		}

		// Calculate needed columns
		needed := t.calculateNeeded(row, col, colSpan)

		// Add columns to grid
		t.g.GrowCols(needed)

		// Process rows top to bottom (0 → row)
		for r := 0; r <= row; r++ {
			if rowFlexMap, exists := flexCells[r]; exists && len(rowFlexMap) > 0 {
				t.distributeAndExpand(r, rowFlexMap, needed)
			}
		}
	}

	// Create cell
	id := t.nextCellID
	t.nextCellID++
	c := NewCell(id, ct, row, col, rowSpan, colSpan)

	// Store cell
	t.cells[id] = c

	// Set in grid
	t.g.SetRect(row, col, rowSpan, colSpan)

	// Populate position map for all occupied cells
	for r := row; r < row+rowSpan; r++ {
		for c := col; c < col+colSpan; c++ {
			t.pos[r][c] = id
		}
	}

	// Advance cursor
	t.c.Advance(colSpan)

	return t
}

// PrintDebug renders table structure in TBL Grid Notation format.
// Shows grid layout with cell types and current cursor position.
// Returns empty string if table has no rows.
// For debug/development purposes.
func (t *Table) PrintDebug() string {
	return t.printDebug()
}

// moveCellBy moves cell at (row, col) right by delta columns.
// Checks each row, cascades next cell if needed, then moves.
// Internal method - no validation.
func (t *Table) moveCellBy(row, col, delta int) {
	cell := t.getCell(row, col)

	// Check each row spanned by cell
	for r := cell.r; r < cell.r+cell.rSpan; r++ {
		if !t.canMoveBy(cell, delta) {
			// Cascade to next cell
			t.moveCellBy(r, cell.c+cell.cSpan, delta)
		}
	}

	// Move cell
	t.clearCell(cell)
	cell.MoveBy(delta)
	t.addCell(cell)
}

// canMoveBy checks if cell can move by delta columns.
// Calculates non-overlapping rectangle: destination minus origin.
func (t *Table) canMoveBy(cell *Cell, delta int) bool {
	if delta <= 0 {
		return true // only handle rightward
	}

	// Non-overlapping rectangle when moving right
	startCol := cell.c + cell.cSpan
	endCol := cell.c + delta + cell.cSpan

	// Check if non-overlapping region is free
	return t.g.IsFree(cell.r, startCol, cell.rSpan, endCol-startCol)
}

// addCell adds cell to pos map and grid.
func (t *Table) addCell(cell *Cell) {
	for r := cell.r; r < cell.r+cell.rSpan; r++ {
		for c := cell.c; c < cell.c+cell.cSpan; c++ {
			t.pos[r][c] = cell.id
		}
	}
	t.g.SetRect(cell.r, cell.c, cell.rSpan, cell.cSpan)
}

// clearCell removes cell from pos map and grid.
func (t *Table) clearCell(cell *Cell) {
	for r := cell.r; r < cell.r+cell.rSpan; r++ {
		for c := cell.c; c < cell.c+cell.cSpan; c++ {
			delete(t.pos[r], c)
		}
	}
	t.g.ClearRect(cell.r, cell.c, cell.rSpan, cell.cSpan)
}

// expandCell expands a flex cell and updates all table state.
// Uses clearCell/addCell to maintain consistency with existing patterns.
func (t *Table) expandCell(cell *Cell, cols int) {
	if cols <= 0 {
		return
	}
	if cell.typ != Flex {
		panic(fmt.Sprintf("cannot expand Static cell id=%d", cell.id))
	}

	// Remove cell from table state
	t.clearCell(cell)

	// Update cell span
	cell.cSpan += cols

	// Re-add cell with new span
	t.addCell(cell)
}

// isRowComplete validates row has no holes and all columns filled.
// Returns true if entire row range [0, Cols) is occupied.
func (t *Table) isRowComplete(row int) bool {
	if row < 0 || row >= t.g.Rows() {
		return false
	}
	if t.g.Cols() == 0 {
		return true
	}
	return t.g.AllRow(row)
}

// isRowStatic checks if all cells in row are Static type.
// Returns false if row has no columns.
func (t *Table) isRowStatic(row int) bool {
	if t.g.Cols() == 0 {
		return false
	}

	for col := 0; col < t.g.Cols(); {
		cell := t.getCell(row, col)
		if cell.typ != Static {
			return false
		}
		col += cell.cSpan
	}
	return true
}

// ensureRows grows grid and initializes position maps up to targetRow (inclusive).
// Centralizes all position map initialization logic.
func (t *Table) ensureRows(targetRow int) {
	if targetRow < t.g.Rows() {
		return
	}

	currentRows := t.g.Rows()
	delta := targetRow - currentRows + 1
	t.g.GrowRows(delta)

	// Initialize position maps for all new rows
	for r := currentRows; r < t.g.Rows(); r++ {
		if t.pos[r] == nil {
			t.pos[r] = make(map[int]ID)
		}
	}
}

// ensureCols ensures sufficient columns for span at position.
// Returns error if colsFixed and insufficient columns.
// Expands columns if not fixed.
func (t *Table) ensureCols(col, colSpan int) error {
	needed := col + colSpan
	if needed <= t.g.Cols() {
		return nil
	}

	if t.colsFixed {
		return fmt.Errorf("tbl: insufficient columns for cell colSpan=%d at col=%d, cols=%d", colSpan, col, t.g.Cols())
	}

	// Expand columns to fit
	delta := needed - t.g.Cols()
	t.g.GrowCols(delta)
	return nil
}

// getCell returns cell at position using position map lookup.
// Panics if position is empty - indicates invalid grid state.
func (t *Table) getCell(row, col int) *Cell {
	rowMap := t.pos[row]
	if rowMap == nil {
		panic(fmt.Sprintf("tbl: no row map at row %d", row))
	}

	id := rowMap[col]
	if id == 0 {
		panic(fmt.Sprintf("tbl: no cell at position (%d,%d)", row, col))
	}

	cell := t.cells[id]
	if cell == nil {
		panic(fmt.Sprintf("tbl: cell id=%d not found in cells map", id))
	}

	return cell
}

// isFlex reports whether the cell at (row, col) is a Flex type.
func (t *Table) isFlex(row, col int) bool {
	cell := t.getCell(row, col)
	return cell.typ == Flex
}

// isWall reports whether the cell at (row, col) acts as a wall.
// A cell is a wall if it spans multiple rows and originates above row.
func (t *Table) isWall(row, col int) bool {
	cell := t.getCell(row, col)
	return cell.rSpan > 1 && cell.r < row
}

// findFirstFreeCol locates first unoccupied column in row.
// Returns column index of first free position, or Cols() if row full.
// Accounts for cells with rowSpan > 1 from previous rows.
func (t *Table) findFirstFreeCol(row int) int {
	if row < 0 || row >= t.g.Rows() {
		return 0
	}
	return t.g.NextFreeCol(row, 0)
}

// calculateNeeded determines how many columns required to fit cell.
// Returns columns needed beyond current grid width or first blocking position.
func (t *Table) calculateNeeded(row, col, colSpan int) int {
	// Find first blocking position in row
	firstBlocked := t.g.NextOccupiedCol(row, col)
	if firstBlocked == -1 {
		firstBlocked = t.g.Cols()
	}

	// Calculate shortage
	required := col + colSpan
	if required <= firstBlocked {
		return 0
	}

	return required - firstBlocked
}

// flexEntry represents a flex cell found during traversal.
type flexEntry struct {
	col       int
	id        ID
	addedSpan int
}

// traverseFlex scans row for flex cells in both directions from col position.
// Recursively processes flex cells found in row above (row-1).
// Returns success status and flex cells organized by row and column.
//
// Return structure: map[int]map[int]flexEntry where outer key = row, inner key = col
func (t *Table) traverseFlex(row, col int) (bool, map[int]map[int]flexEntry) {
	if row < 0 {
		return true, nil
	}

	result := make(map[int]map[int]flexEntry)

	// Scan left
	okLeft, leftResult := t.traverseFlexDir(row, col, DirLeft)
	if !okLeft {
		return false, nil
	}
	mergeFlexResults(result, leftResult)

	// Scan right (includes origin col)
	okRight, rightResult := t.traverseFlexDir(row, col, DirRight)
	if !okRight {
		return false, nil
	}
	mergeFlexResults(result, rightResult)

	return true, result
}

// traverseFlexDir scans single direction from col position in row.
// Stops at walls or grid boundaries.
// For each flex cell found, recursively calls traverseFlex(row-1, flexCol).
func (t *Table) traverseFlexDir(row, col int, dir Direction) (bool, map[int]map[int]flexEntry) {
	result := make(map[int]map[int]flexEntry)

	// Determine iteration bounds
	start, end, step := col, t.g.Cols(), 1
	if dir == DirLeft {
		start, end, step = col-1, -1, -1
	}

	// Iterate in direction
	for c := start; c != end; c += step {
		// Check for wall
		if t.isWall(row, c) {
			break
		}

		// Check if flex cell
		if t.isFlex(row, c) {
			cell := t.getCell(row, c)

			// Add to current row result (map auto-deduplicates)
			if result[row] == nil {
				result[row] = make(map[int]flexEntry)
			}
			result[row][c] = flexEntry{
				col:       c,
				id:        cell.id,
				addedSpan: cell.AddedSpan(),
			}

			// Recurse to row above
			ok, aboveResult := t.traverseFlex(row-1, c)
			if !ok {
				return false, nil
			}
			mergeFlexResults(result, aboveResult)
		}
	}

	return true, result
}

// mergeFlexResults combines source into dest (both map[int]map[int]flexEntry).
// Map structure naturally handles duplicates by overwriting.
func mergeFlexResults(dest, src map[int]map[int]flexEntry) {
	for row, colMap := range src {
		if dest[row] == nil {
			dest[row] = make(map[int]flexEntry)
		}
		for col, entry := range colMap {
			dest[row][col] = entry
		}
	}
}

// distributeAndExpand handles movement and expansion for one row.
// Distributes needed cols fairly: base amount to all, remainder to cells with least expansion.
// Tie-breaking: leftmost cells get priority.
func (t *Table) distributeAndExpand(row int, flexMap map[int]flexEntry, needed int) {
	if len(flexMap) == 0 || needed <= 0 {
		return
	}

	// Build sorted slice from map
	entries := make([]flexEntry, 0, len(flexMap))
	for _, entry := range flexMap {
		entries = append(entries, entry)
	}

	// Sort by: 1) addedSpan ascending, 2) col ascending (left to right)
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].addedSpan != entries[j].addedSpan {
			return entries[i].addedSpan < entries[j].addedSpan
		}
		return entries[i].col < entries[j].col
	})

	n := len(entries)

	// Calculate base distribution
	base := needed / n
	remainder := needed % n

	// Process each flex cell in sorted order
	for i, entry := range entries {
		flexCell := t.cells[entry.id]

		// Everyone gets base amount
		expandAmount := base

		// First 'remainder' cells get +1 (prioritizes least expanded)
		if i < remainder {
			expandAmount++
		}

		if expandAmount == 0 {
			continue
		}

		// Check if there's an adjacent cell to move
		adjacentCol := entry.col + flexCell.cSpan
		adjacentOccupied := !t.g.IsFree(row, adjacentCol, 1, expandAmount)

		if adjacentOccupied {
			// Move adjacent cell right to make space
			t.moveCellBy(row, adjacentCol, expandAmount)
		}

		// Expand flex cell (updates grid + position map)
		t.expandCell(flexCell, expandAmount)
	}
}
