package tbl

import "fmt"

const (
	NO_COL = -1
)

func (t *Table) R(cells ...any) Row {
	r := *t.config.RowDefault
	r.Cells = make([]Cell, 0, len(cells))

	for _, c := range cells {
		r.Cells = append(r.Cells, t.C(c))
	}

	return r
}

func (t *Table) C(value any) Cell {
	switch v := value.(type) {
	case string:
		return t.config.CellDefault.WithContent(v)
	case Cell:
		return v
	default:
		return t.config.CellDefault.WithContent(fmt.Sprintf("%v", v))
	}
}

func (t *Table) Add(row Row) error {
	return t.processRow(row)
}

func (t *Table) AddRow(cells ...any) error {
	r := t.R(cells...)
	return t.Add(r)
}
