package tbl

import "fmt"

type Table struct {
	config Config
	state  state
	rows   []Row
}

type state struct {
	colLevels         []int
	colWidths         []int
	rowHeights        []int
	hLines            []bool
	currentRow        int
	stillFlexibleCols bool
	width             int
}

func New() *Table {
	return NewWithConfig(DefaultConfig)
}

func NewWithConfig(cfg Config) *Table {
	mergedCfg := DefaultConfig.Merge(cfg)
	return &Table{
		config: mergedCfg,
		state:  state{colLevels: []int{}, colWidths: []int{}},
		rows:   []Row{},
	}
}

func (t *Table) R(cells ...any) Row {
	row := *t.config.RowDefault
	row.Cells = make([]Cell, 0, len(cells))

	for _, cell := range cells {
		row.Cells = append(row.Cells, t.C(cell))
	}

	return row
}

func (t *Table) C(v any) Cell {
	switch val := v.(type) {
	case string:
		return t.config.CellDefault.WithContent(val)
	case Cell:
		return val
	default:
		return t.config.CellDefault.WithContent(fmt.Sprintf("%v", val))
	}
}

func (t *Table) AddRow(cells ...any) {
	row := t.R(cells...)
	t.rows = append(t.rows, row)
	t.state.currentRow++
}

func (t *Table) Add(row Row) {
	t.rows = append(t.rows, row)
	t.state.currentRow++
}
