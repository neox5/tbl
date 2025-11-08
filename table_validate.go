package tbl

import "fmt"

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
