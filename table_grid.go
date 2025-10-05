package tbl

// addRow ensures cursor row exists in grid.
func (t *Table) addRow() {
	t.cur.nextRow()            // move cursor to next row
	t.rowHeight[t.cur.row] = 0 // initialize row tracking

	t.g.AddRow() // add row to grid
}

// addCols grows grid columns to accommodate cell at cursor.
func (t *Table) addCols(colSpan int) {
	needed := t.cur.col + colSpan
	for needed > t.g.Cols() {
		t.g.AddCol()
	}
}
