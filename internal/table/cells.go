package table

import (
	"fmt"

	"github.com/neox5/tbl/internal/cell"
)

// AddRow adds a new row with the specified cells
func (t *Table) AddRow(cells ...*cell.Cell) {
	t.startNewRow()
	if t.row == 0 {
		t.initWithFirstRow(cells)
	}
	t.addCells(cells)
	t.resolveRow()
	t.validateRow()
	t.advanceRow()
}

// addCells appends multiple cells to the current row
func (t *Table) addCells(cells []*cell.Cell) {
	for _, c := range cells {
		t.addCell(c)
	}
}

// addCell adds a single cell to the current position
func (t *Table) addCell(c *cell.Cell) {
	idx := t.nextIndex
	span := c.ColSpan()
	
	if !t.spanFit(span) {
		panic(fmt.Sprintf("cell(%d) span %d does not fit at row %d col %d", idx, span, t.row, t.col))
	}

	// Position cell
	c.SetColStart(t.col)
	c.SetRowStart(t.row)
	
	// Track flex cells
	if c.IsColFlex() {
		t.openFlexCells = append(t.openFlexCells, idx)
	}

	// Add to table
	t.cells = append(t.cells, c)
	t.nextIndex++

	// Update state
	t.updateWidths(c)
	t.updateLevels(c)
	t.addIndices(idx, c)
	
	t.advanceCol()
}

// updateWidths updates column widths based on cell content
func (t *Table) updateWidths(c *cell.Cell) {
	w := c.Width()
	for i := range c.ColSpan() {
		if t.colWidths[c.ColStart()+i] < w {
			t.colWidths[c.ColStart()+i] = w
		}
	}
}

// updateLevels updates column levels for spanning cells
func (t *Table) updateLevels(c *cell.Cell) {
	for i := range c.ColSpan() {
		t.colLevels[c.ColStart()+i] = c.RowSpan()
	}
}

// addIndices adds cell to column and row indices
func (t *Table) addIndices(idx int, c *cell.Cell) {
	t.addColIndex(idx, c.ColStart(), c.ColSpan())
	t.addRowIndex(idx, c.RowStart(), c.RowSpan())
}

// resolveRow resolves flexible column spans for the current row
func (t *Table) resolveRow() {
	if len(t.openFlexCells) == 0 {
		return
	}
	
	// Calculate total column span used in current row
	colCount := 0
	flexCells := make([]*cell.Cell, 0)
	flexWeights := 0
	for _, idx := range t.CellsInRow(t.row) {
		c := t.cells[idx]
		colCount += c.ColSpan()

		if c.IsColFlex() {
			flexCells = append(flexCells, c)
			flexWeights += c.ColWeight()
		}
	}
	
	// Only resolve flex cells if row is under-filled
	// Over-filled rows will be caught by validation
	colsAvailable := t.ColCount() - colCount
	if colsAvailable <= 0 {
		return
	}

	// TODO: Distribute available columns among flex cells by weight
	// TODO: Respect maxSpan constraints
	// TODO: Update cell spans and table state
}

// initWithFirstRow initializes column arrays from first row
func (t *Table) initWithFirstRow(cells []*cell.Cell) {
	totalCols := 0
	for _, c := range cells {
		totalCols += c.ColSpan()
	}
	t.colWidths = make([]int, totalCols)
	t.colLevels = make([]int, totalCols)
}

// validateRow ensures all columns are filled
func (t *Table) validateRow() {
	for i, level := range t.colLevels {
		if level == 0 {
			panic(fmt.Sprintf("incomplete row %d: column %d not filled", t.row, i))
		}
	}
}
