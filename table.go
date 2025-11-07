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
			if rowFlexCells, exists := flexCells[r]; exists && len(rowFlexCells) > 0 {
				t.distributeAndExpand(r, rowFlexCells, needed)
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

// getCellAt finds the cell containing position (row, col).
// Returns nil if position is empty.
func (t *Table) getCellAt(row, col int) *Cell {
	for _, cell := range t.cells {
		if cell.Contains(row, col) {
			return cell
		}
	}
	return nil
}

// getCellsInRow returns all cells that touch the specified row.
func (t *Table) getCellsInRow(row int) []*Cell {
	var result []*Cell
	for _, cell := range t.cells {
		if cell.TouchesRow(row) {
			result = append(result, cell)
		}
	}
	return result
}

// isFlex reports whether the cell at (row, col) is a Flex type.
func (t *Table) isFlex(row, col int) bool {
	cell := t.getCellAt(row, col)
	return cell != nil && cell.typ == Flex
}

// isWall reports whether the cell at (row, col) acts as a wall.
// A cell is a wall if it spans multiple rows and originates above row.
func (t *Table) isWall(row, col int) bool {
	cell := t.getCellAt(row, col)
	return cell != nil && cell.rSpan > 1 && cell.r < row
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
		cell := t.getCellAt(row, col)
		if cell == nil || cell.typ != Static {
			return false
		}
		col += cell.cSpan
	}
	return true
}

// ensureRows grows grid up to targetRow (inclusive).
func (t *Table) ensureRows(targetRow int) {
	if targetRow < t.g.Rows() {
		return
	}

	delta := targetRow - t.g.Rows() + 1
	t.g.GrowRows(delta)
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

// flexCell represents a flex cell found during traversal.
type flexCell struct {
	cell      *Cell
	addedSpan int
}

// traverseFlex scans row for flex cells in both directions from col position.
// Recursively processes flex cells found in row above (row-1).
// Returns success status and flex cells organized by row.
//
// Return structure: map[int][]flexCell where key = row
func (t *Table) traverseFlex(row, col int) (bool, map[int][]flexCell) {
	if row < 0 {
		return true, nil
	}

	result := make(map[int][]flexCell)
	seen := make(map[ID]bool)

	// Scan left
	okLeft := t.traverseFlexDir(row, col, DirLeft, result, seen)
	if !okLeft {
		return false, nil
	}

	// Scan right (includes origin col)
	okRight := t.traverseFlexDir(row, col, DirRight, result, seen)
	if !okRight {
		return false, nil
	}

	return true, result
}

// traverseFlexDir scans single direction from col position in row.
// Stops at walls or grid boundaries.
// For each flex cell found, recursively calls traverseFlex(row-1, flexCol).
func (t *Table) traverseFlexDir(row, col int, dir Direction, result map[int][]flexCell, seen map[ID]bool) bool {
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
			cell := t.getCellAt(row, c)

			// Skip if already seen
			if seen[cell.id] {
				continue
			}
			seen[cell.id] = true

			// Add to current row result
			result[row] = append(result[row], flexCell{
				cell:      cell,
				addedSpan: cell.AddedSpan(),
			})

			// Recurse to row above
			ok, aboveResult := t.traverseFlex(row-1, c)
			if !ok {
				return false
			}

			// Merge results
			for r, flexCells := range aboveResult {
				for _, fc := range flexCells {
					if !seen[fc.cell.id] {
						seen[fc.cell.id] = true
						result[r] = append(result[r], fc)
					}
				}
			}
		}
	}

	return true
}

// distributeAndExpand handles movement and expansion for one row.
// Distributes needed cols fairly: base amount to all, remainder to cells with least expansion.
// Tie-breaking: leftmost cells get priority.
func (t *Table) distributeAndExpand(row int, flexCells []flexCell, needed int) {
	if len(flexCells) == 0 || needed <= 0 {
		return
	}

	// Sort by: 1) addedSpan ascending, 2) col ascending (left to right)
	sort.Slice(flexCells, func(i, j int) bool {
		if flexCells[i].addedSpan != flexCells[j].addedSpan {
			return flexCells[i].addedSpan < flexCells[j].addedSpan
		}
		return flexCells[i].cell.c < flexCells[j].cell.c
	})

	n := len(flexCells)

	// Calculate base distribution
	base := needed / n
	remainder := needed % n

	// Plan all movements and expansions
	type action struct {
		cell   *Cell
		expand int
		moveBy int
	}
	var actions []action

	// Calculate expansions
	for i, fc := range flexCells {
		expandAmount := base
		if i < remainder {
			expandAmount++
		}
		if expandAmount > 0 {
			actions = append(actions, action{
				cell:   fc.cell,
				expand: expandAmount,
			})
		}
	}

	// Execute movements and expansions
	t.executeExpansions(row, actions)
}

// executeExpansions applies expansions with necessary shifts.
func (t *Table) executeExpansions(row int, actions []action) {
	// Sort actions by column position for processing
	sort.Slice(actions, func(i, j int) bool {
		return actions[i].cell.c < actions[j].cell.c
	})

	// Process each expansion
	for _, a := range actions {
		if a.expand == 0 {
			continue
		}

		// Check if there's an adjacent cell to move
		adjacentCol := a.cell.c + a.cell.cSpan

		// Find all cells that need to shift right
		shiftCells := t.findCellsToShift(row, adjacentCol, a.expand)

		// Execute shifts from right to left
		t.shiftCellsRight(shiftCells, a.expand)

		// Expand the flex cell
		t.expandCell(a.cell, a.expand)
	}
}

// findCellsToShift identifies cells that need to move for expansion.
func (t *Table) findCellsToShift(row, fromCol, delta int) []*Cell {
	var cells []*Cell
	seen := make(map[ID]bool)

	// Find all cells in row starting at fromCol
	for col := fromCol; col < t.g.Cols(); col++ {
		cell := t.getCellAt(row, col)
		if cell != nil && !seen[cell.id] {
			seen[cell.id] = true
			cells = append(cells, cell)
			col = cell.c + cell.cSpan - 1 // Skip to end of cell
		}
	}

	return cells
}

// shiftCellsRight moves cells right by delta columns.
func (t *Table) shiftCellsRight(cells []*Cell, delta int) {
	if delta <= 0 || len(cells) == 0 {
		return
	}

	// Sort by column descending for right-to-left processing
	sort.Slice(cells, func(i, j int) bool {
		return cells[i].c > cells[j].c
	})

	// Clear all cells from grid
	for _, cell := range cells {
		t.g.ClearRect(cell.r, cell.c, cell.rSpan, cell.cSpan)
	}

	// Move cells and re-set in grid
	for _, cell := range cells {
		cell.c += delta
		t.g.SetRect(cell.r, cell.c, cell.rSpan, cell.cSpan)
	}
}

// expandCell expands a flex cell and updates grid state.
func (t *Table) expandCell(cell *Cell, cols int) {
	if cols <= 0 {
		return
	}
	if cell.typ != Flex {
		panic(fmt.Sprintf("cannot expand Static cell id=%d", cell.id))
	}

	// Clear old position
	t.g.ClearRect(cell.r, cell.c, cell.rSpan, cell.cSpan)

	// Update cell span
	cell.cSpan += cols

	// Set new position
	t.g.SetRect(cell.r, cell.c, cell.rSpan, cell.cSpan)
}
