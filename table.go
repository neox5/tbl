package tbl

import "fmt"

const (
	NO_COL = -1
)

type Table struct {
	config       Config
	rows         []Row
	virtualRows  []int
	colLevels    []int
	colWidths    []int
	hLines       []bool
	flexibleCols bool
	width        int
	// current possition
	col, row int
}

func New() *Table {
	return NewWithConfig(DefaultConfig)
}

func NewWithConfig(cfg Config) *Table {
	mergedCfg := DefaultConfig.Merge(cfg)
	return &Table{
		config:      mergedCfg,
		rows:        []Row{},
		virtualRows: []int{},
		colLevels:   []int{},
		colWidths:   []int{},
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

func (t *Table) nextCol() int {
	for i, l := range t.colLevels {
		if l == 0 {
			return i
		}
	}
	return NO_COL
}

func (t *Table) reduceColLevels() {
	for i, l := range t.colLevels {
		if l == FLEX {
			continue
		} // skip flex colomns
		t.colLevels[i] -= 1
	}
}

func (t *Table) advanceRow() {
	t.row++
	t.col = 0
	t.reduceColLevels()

	for t.nextCol() == NO_COL {
		t.virtualRows = append(t.virtualRows, t.row)
		t.row++
		t.reduceColLevels()
	}
}

func (t *Table) availableSpan() int {
	max, span := 0, 0
	for _, l := range t.colLevels {
		if l == 0 {
			span++
		} else {
			if span > max {
				max = span
			}
			span = 0
		}
	}
	if span > max {
		max = span
	}
	return max
}

func (t *Table) initWithFirstRow(row Row) error {
	// Calculate total column count from cell spans
	colCount := 0
	for _, c := range row.Cells {
		if c.ColSpan == FLEX {
			colCount++
			continue
		}
		colCount += c.ColSpan
	}
	t.colLevels = make([]int, colCount)
	t.colWidths = make([]int, colCount)

	for _, c := range row.Cells {
		t.processCell(c)

		if row.Height < c.Height() {
			row.Height = c.Height()
		}
	}

	t.rows = append(t.rows, row)
	t.advanceRow()
	return nil
}

func (t *Table) processCell(cell Cell) {
	if cell.ColSpan == FLEX {
		t.colLevels[t.col] = cell.RowSpan
		t.colWidths[t.col] = 1
		t.col++
		return
	}

	w := cell.Width() / cell.ColSpan // width contribution
	r := cell.Width() % cell.ColSpan // width remainder
	for j := range cell.ColSpan {
		t.colLevels[t.col+j] = cell.RowSpan // rowSpan or FLEX
		t.colWidths[t.col+j] = w

		if r > 0 { // adding remainder from left to right
			t.colWidths[t.col+j]++
			r-- // reduce remainder for next col until it runs out
		}
	}
	t.col += cell.ColSpan
}

func (t *Table) processRow(row Row) error {
	if len(t.rows) == 0 {
		return t.initWithFirstRow(row)
	}

	for _, c := range row.Cells {
		if c.ColSpan > t.availableSpan() {
			return fmt.Errorf("cell[%d:%d]: cell does not fit - needs: %d got: %d", t.row, t.col, c.ColSpan, t.availableSpan())
		}
		t.processCell(c)
		if row.Height < c.Height() {
			row.Height = c.Height()
		}
	}

	t.rows = append(t.rows, row)
	t.advanceRow()
	return nil
}
