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

func (t *Table) initWithFirstRow(row Row) error {
	// Calculate total column count from cell spans
	colCount := 0
	for _, c := range row.Cells {
		colCount += c.ColSpan
	}
	t.state.colLevels = make([]int, colCount)
	t.state.colWidths = make([]int, colCount)

	rHeight := 0
	col := 0
	for _, c := range row.Cells {
		// columns
		w := c.Width() / c.ColSpan // width contribution
		r := c.Width() % c.ColSpan // width remainder
		for j := range c.ColSpan {
			t.state.colLevels[col+j] = c.RowSpan
			t.state.colWidths[col+j] = w

			if r > 0 { // adding remainder from left to right
				t.state.colWidths[col+j]++
				r-- // reduce remainder for next col until it runs out
			}
		}
		col += c.ColSpan

		// row
		if rHeight < c.Height() {
			rHeight = c.Height()
		}
	}

	t.rows = append(t.rows, row)
	t.state.rowHeights = append(t.state.rowHeights, rHeight)
	t.state.currentRow++
	return nil
}

func (t *Table) processRow(row Row) error {
	if len(t.rows) == 0 {
		return t.initWithFirstRow(row)
	}

	return nil
}
