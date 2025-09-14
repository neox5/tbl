package tbl

import "fmt"

func (t *Table) AddRow(cells ...any) {
	t.R(cells...)
}

func (t *Table) R(cells ...any) {
	t.rowStarts = append(t.rowStarts, len(t.cells))

	for _, c := range cells {
		t.cells = append(t.cells, t.C(c))
	}
}

func (t *Table) C(value any) Cell {
	switch v := value.(type) {
	case string:
		return t.cellDefault.WithContent(v)
	case Cell:
		return v
	default:
		return t.cellDefault.WithContent(fmt.Sprintf("%v", v))
	}
}
