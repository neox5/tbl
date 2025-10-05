package tbl

// addRow ensures cursor row exists in grid.
func (t *Table) addRow() {
	t.cur.NextRow()

	t.g.AddRow() // add row to grid
	t.rows = append(t.rows, nil)

	// reduce colLevels
	for i := range len(t.cols) {
		if t.colLevels[i] > 0 {
			t.colLevels[i]--
		}
	}
}

// addCols grows grid columns to accommodate cell at cursor.
func (t *Table) addCols(colSpan int) {
	needed := t.cur.Col() + colSpan
	for needed > t.g.Cols() {
		t.g.AddCol()
		t.cols = append(t.cols, nil)
		t.colWidths = append(t.colWidths, 0)
		t.colLevels = append(t.colLevels, 0)
	}
}
