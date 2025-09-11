package tbl

import "fmt"

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

func (t *Table) initWithFirstRow(row Row) error {
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

	w := cell.Width() / cell.ColSpan
	r := cell.Width() % cell.ColSpan
	for j := range cell.ColSpan {
		t.colLevels[t.col+j] = cell.RowSpan
		t.colWidths[t.col+j] = w

		if r > 0 {
			t.colWidths[t.col+j]++
			r--
		}
	}
	t.col += cell.ColSpan
}
