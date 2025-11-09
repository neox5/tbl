package tbl

// nextRow advances cursor to next row and positions at first free column.
// Returns new row index.
func (t *Table) nextRow() {
	t.row++
	t.col = t.findFirstFreeCol(t.row)
}

// advance moves cursor forward by colSpan columns.
func (t *Table) advance(colSpan int) {
	if colSpan <= 0 {
		return
	}
	t.col = min(t.col+colSpan, t.g.Cols()-1)
}
