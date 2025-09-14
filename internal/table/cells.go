package table

import (
	"fmt"

	"github.com/neox5/tbl/internal/cell"
)

// AddRow adds a new row with the specified cells
func (t *Table) AddRow(cells ...*cell.Cell) {
	t.startNewRow()
	if t.row == 0 {
		t.initializeWithFirstRow(cells)
	}
	t.addCells(cells)
	t.validateRow()
	t.advanceRow()
}

// addCells appends multiple cells to the current row
func (t *Table) addCells(cells []*cell.Cell) {
	for _, c := range cells {
		cSpan, rSpan := c.Span()
		if !t.spanFit(cSpan) {
			panic(fmt.Sprintf("cell span %d does not fit at row %d col %d", cSpan, t.row, t.col))
		}

		w := c.ColWidth()
		for i := range cSpan {
			// update colWidths when cell width is wider then the current value
			if t.colWidths[t.col+i] < w {
				t.colWidths[t.col+i] = w
			}
			t.colLevels[t.col+i] = rSpan // colLevels will always be set; >0 or -1 (FLEX)
		}

		t.advanceCol()

		t.cells = append(t.cells, c)
		t.nextIndex++
	}
}

// initializeWithFirstRow initializes colWidths and colLevels so that addCells
// works also with the first row
func (t *Table) initializeWithFirstRow(cells []*cell.Cell) {
	totalCols := 0
	for _, c := range cells {
		cSpan, _ := c.Span()
		if cSpan == cell.FLEX {
			cSpan = 1
		}
		totalCols += cSpan
	}

	t.colWidths = make([]int, totalCols)
	t.colLevels = make([]int, totalCols)
}

// validateRow validates if all 0 colLevels were increased.
// If not we need to panic as the row is incomplete.
func (t *Table) validateRow() {
	for i, level := range t.colLevels {
		if level == 0 {
			panic(fmt.Sprintf("incomplete row %d: column %d was not filled", t.row, i))
		}
	}
}
