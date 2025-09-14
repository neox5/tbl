package tbl

func (t *Table) AddRow(cells ...any) {
	t.R(cells...)
}

func (t *Table) R(cells ...any) {
	t.newRow()
	t.addCells(cells)
}

func (t *Table) C(value any) Cell {
	return t.newCell(value)
}
