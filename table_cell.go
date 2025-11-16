package tbl

import "fmt"

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
