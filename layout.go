package tbl

import "fmt"

const (
	NO_COL = -1
)

func (t *Table) processRow(rowIdx int) error {
	if rowIdx == 0 {
		return t.initWithFirstRow()
	}

	start := t.rowStarts[rowIdx]
	end := len(t.cells)
	if rowIdx+1 < len(t.rowStarts) {
		end = t.rowStarts[rowIdx+1]
	}

	for i := start; i < end; i++ {
		c := &t.cells[i]
		if c.Col.IsFlex() {
			if t.availableSpan() < 1 {
				return fmt.Errorf("cell[%d]: flex cell does not fit", i)
			}
		} else {
			if c.Col.Span > t.availableSpan() {
				return fmt.Errorf("cell[%d]: cell does not fit - needs: %d got: %d", i, c.Col.Span, t.availableSpan())
			}
		}
		t.processCell(c, i)
	}

	t.advanceRow()
	return nil
}

func (t *Table) initWithFirstRow() error {
	if len(t.rowStarts) == 0 || len(t.cells) == 0 {
		return fmt.Errorf("no cells to process")
	}

	colCount := 0
	start := t.rowStarts[0]
	end := len(t.cells)
	if len(t.rowStarts) > 1 {
		end = t.rowStarts[1]
	}

	for i := start; i < end; i++ {
		c := &t.cells[i]
		if c.Col.IsFlex() {
			colCount++
			continue
		}
		colCount += c.Col.Span
	}

	t.colLevels = make([]int, colCount)
	t.colWidths = make([]int, colCount)

	for i := start; i < end; i++ {
		c := &t.cells[i]
		t.processCell(c, i)
	}

	t.advanceRow()
	return nil
}

func (t *Table) processCell(cell *Cell, cellIdx int) {
	if cell.Col.IsFlex() {
		t.colLevels[t.currIndex] = cell.Row.Span
		t.colWidths[t.currIndex] = 1
		t.openFlexCells = append(t.openFlexCells, cellIdx)
		t.currIndex++
		return
	}

	w := cell.Width() / cell.Col.Span
	r := cell.Width() % cell.Col.Span
	for j := range cell.Col.Span {
		t.colLevels[t.currIndex+j] = cell.Row.Span
		t.colWidths[t.currIndex+j] = w

		if r > 0 {
			t.colWidths[t.currIndex+j]++
			r--
		}
	}
	t.currIndex += cell.Col.Span
}
