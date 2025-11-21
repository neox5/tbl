package tbl

import "fmt"

// addCol adds column to grid and stores configuration.
// Internal implementation for AddCol().
func (t *Table) addCol(width, minWidth, maxWidth int) {
	// Get current column count
	col := t.g.Cols()

	// Expand grid by one column
	t.g.GrowCols(1)

	// Store configuration
	t.colConfigs[col] = ColConfig{
		Width:    width,
		MinWidth: minWidth,
		MaxWidth: maxWidth,
	}

	// Mark columns as fixed
	t.colsFixed = true
}

// addRow advances to next row with validation and cursor positioning.
// Internal implementation for AddRow().
func (t *Table) addRow() {
	// Validate previous row if not first row
	if t.row >= 0 {
		if !t.isRowComplete(t.row) {
			panic(fmt.Sprintf("tbl: incomplete row %d before AddRow", t.row))
		}

		// Check if we can fix columns
		if !t.colsFixed && t.isRowStatic(t.row) {
			t.colsFixed = true
		}
	}

	// Safeguard: ensure at least 1 column exists for first row
	if t.row == -1 && t.g.Cols() == 0 {
		t.g.GrowCols(1)
	}

	// Ensure next row exists
	t.ensureRows(t.row + 2) // t.row = current row index; +1 current row count; +2 next row count

	// Advance to next row
	t.nextRow()
}

// addCell adds a cell at cursor position with specified type and span.
// Internal implementation for AddCell().
// Expands columns if needed (when not fixed). Validates span fits in grid.
// Panics if: no row active, span invalid, insufficient columns (when fixed),
// or space occupied.
func (t *Table) addCell(ct CellType, rowSpan, colSpan int, content string) {
	// Ensure sufficient rows for cell span
	t.ensureRows(t.row + rowSpan)

	// Simple check on first row
	if t.row == 0 && t.col+colSpan > t.g.Cols() {
		if err := t.ensureCols(t.col + colSpan); err != nil {
			panic(err.Error())
		}
	}

	// when cell does not fit in the current position, we try to expand
	if colSpan > t.g.CountZerosFromInRow(t.row, t.col) {
		// Space occupied - wall blocking
		if t.colsFixed {
			panic(fmt.Sprintf("tbl: space occupied at cursor (%d,%d), cannot expand", t.row, t.col))
		}

		// Attempt expansion
		ok, flexCells := t.traverseFlex(t.row, t.col)
		if !ok {
			panic(fmt.Sprintf("tbl: no flex cells available for expansion at cursor (%d,%d)", t.row, t.col))
		}

		// Calculate needed columns
		needed := t.calculateNeeded(t.row, t.col, colSpan)

		// Add columns to grid
		t.g.GrowCols(needed)

		// Process rows top to bottom
		for r := 0; r <= t.row; r++ {
			if rowFlexCells, exists := flexCells[r]; exists && len(rowFlexCells) > 0 {
				t.distributeAndExpand(r, rowFlexCells, needed)
			}
		}
	}

	if t.g.B.Test(t.g.Index(t.row, t.col)) {
		t.col = t.findFirstFreeCol(t.row)
	}

	// Create cell
	id := t.nextCellID
	t.nextCellID++

	// Store cell
	c := NewCell(id, ct, t.row, t.col, rowSpan, colSpan, content)
	t.cells[id] = c

	// Set in grid
	t.g.SetRect(t.row, t.col, rowSpan, colSpan)
	t.advance()
}

// nextRow advances cursor to next row and positions at first free column.
func (t *Table) nextRow() {
	t.row++
	t.advance()
}

// advance moves cursor forward to first free column.
func (t *Table) advance() {
	col := t.findFirstFreeCol(t.row)
	if col == -1 {
		col = t.g.Cols() - 1
	}

	t.col = col
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

// findFirstFreeCol locates first unoccupied column in row.
// Returns column index of first free position, or Cols() if row full.
// Accounts for cells with rowSpan > 1 from previous rows.
func (t *Table) findFirstFreeCol(row int) int {
	if row < 0 || row >= t.g.Rows() {
		return 0
	}
	return t.g.NextZeroInRow(row, 0)
}

// calculateNeeded determines how many columns required to fit cell.
// Returns columns needed beyond current grid width or first blocking position.
func (t *Table) calculateNeeded(row, col, colSpan int) int {
	// Find first blocking position in row
	firstBlocked := t.g.NextOneInRow(row, col)
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

// ensureRows grows grid up to rowCount (inclusive).
func (t *Table) ensureRows(rowCount int) {
	if rowCount <= t.g.Rows() {
		return
	}

	delta := rowCount - t.g.Rows()
	t.g.GrowRows(delta)
}

// ensureCols ensures sufficient columns for span at position.
// Returns error if colsFixed and insufficient columns.
// Expands columns if not fixed.
func (t *Table) ensureCols(colCount int) error {
	if colCount <= t.g.Cols() {
		return nil
	}

	if t.colsFixed {
		return fmt.Errorf("tbl: cannot ensure cols %d on fixed table (%d)", colCount, t.g.Cols())
	}

	// Expand columns to fit
	delta := colCount - t.g.Cols()
	t.g.GrowCols(delta)
	return nil
}

// finalize validates and completes table before rendering.
// Checks last row completion and attempts flex expansion if incomplete.
// Panics if table cannot be completed.
//
// Validation sequence:
//  1. Check empty row (cursor at col 0)
//  2. Check last row completion
//  3. If incomplete: scan backward for flex cells and expand
//  4. Validate complete grid (all bits set)
func (t *Table) finalize() {
	if t.g.Rows() == 0 {
		panic("tbl: empty table")
	}

	// Check if cursor at start of row - empty row added
	if t.col == 0 {
		panic(fmt.Sprintf("tbl: empty last row %d\n%s", t.row, t.PrintDebug()))
	}

	// Check if last row complete
	if !t.isRowComplete(t.row) {
		t.finalizeLastRow()
	}

	// Final validation - entire grid must be filled
	if !t.g.B.All() {
		panic(fmt.Sprintf("tbl: incomplete table\n%s", t.PrintDebug()))
	}
}

// finalizeLastRow attempts to complete last row using flex expansion.
// Scans backward from end of row to find flex cells.
// Applies distributeAndExpand to fill incomplete columns.
// Panics if no flex cells available for expansion.
func (t *Table) finalizeLastRow() {
	// Calculate columns needed to complete row
	needed := t.g.Cols() - t.col

	// Scan row for flex cells
	flexCells, wallCol := t.scanRowForFlex(t.row)

	if len(flexCells) == 0 {
		if wallCol > 0 {
			panic(fmt.Sprintf("tbl: incomplete last row %d, wall at col %d blocks expansion\n%s", t.row, wallCol, t.PrintDebug()))
		}
		panic(fmt.Sprintf("tbl: incomplete last row %d, no flex cells for expansion\n%s", t.row, t.PrintDebug()))
	}

	// Expand flex cells to fill gap
	t.distributeAndExpand(t.row, flexCells, needed)

	// Verify expansion succeeded
	if !t.isRowComplete(t.row) {
		panic(fmt.Sprintf("tbl: flex expansion failed to complete last row %d\n%s", t.row, t.PrintDebug()))
	}
}

// scanRowForFlex scans backward from end of row for flex cells.
// Stops at walls or row start.
// Returns flex cells found in reverse order (rightmost first) and wall column (0 if no wall).
func (t *Table) scanRowForFlex(row int) ([]flexCell, int) {
	var flexCells []flexCell

	// Scan backward from end of row, jumping by cell span
	c := t.g.Cols() - 1
	for c >= 0 {
		// Stop at wall
		if t.isWall(row, c) {
			return flexCells, c
		}

		cell := t.getCellAt(row, c)
		if cell == nil {
			c--
			continue
		}

		// Collect flex cells
		if cell.typ == Flex {
			flexCells = append(flexCells, flexCell{
				cell:      cell,
				addedSpan: cell.AddedSpan(),
			})
		}

		// Jump backward by cell span
		c -= cell.cSpan
	}

	return flexCells, 0
}
