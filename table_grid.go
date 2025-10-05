package tbl

// addRow ensures cursor row exists in grid.
func (t *Table) addRow() {
	for t.cur.row >= t.g.Rows() {
		t.g.AddRow()
	}
}

// addCols grows grid columns to accommodate cell at cursor.
func (t *Table) addCols(colSpan int) {
	needed := t.cur.col + colSpan
	for needed > t.g.Cols() {
		t.g.AddCol()
	}
}
