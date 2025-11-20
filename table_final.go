package tbl

import "fmt"

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
