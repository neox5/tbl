package tbl

// nextRow advances cursor to next row and positions at first free column.
// Returns new row index.
func (t *Table) nextRow() {
	t.row++
	t.advance()
}

// advance moves cursor forward by colSpan columns.
func (t *Table) advance() {
	col := t.findFirstFreeCol(t.row)
	if col == -1 {
		col = t.g.Cols() - 1
	}

	t.col = col
}
