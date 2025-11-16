package tbl

import "fmt"

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

	// Ensure next row exists
	t.ensureRows(t.row + 2) // t.row = current row index; +1 current row count; +2 next row count

	// Advance to next row
	t.nextRow()
}
