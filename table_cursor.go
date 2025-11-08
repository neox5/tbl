package tbl

// nextRow advances cursor to next row and positions at first free column.
// Returns new row index.
func (t *Table) nextRow() int {
	t.row++
	t.col = 0
	if t.row > 0 {
		t.col = t.findFirstFreeCol(t.row)
	}
	return t.row
}

// advance moves cursor forward by colSpan columns.
// Caps cursor at grid column boundary if advancement exceeds limits.
func (t *Table) advance(colSpan int) {
	if colSpan <= 0 {
		return
	}

	newCol := t.col + colSpan
	if newCol > t.g.Cols() {
		t.col = t.g.Cols()
		return
	}

	t.col = newCol
}
