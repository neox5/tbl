package table

import (
	"fmt"

	"github.com/neox5/tbl/internal/cell"
)

// AddRow adds a new row with the specified cells
func (t *Table) AddRow(cells ...*cell.Cell) {
	t.advanceRow()
	if t.row == 0 {
		t.initWithFirstRow(cells)
	}
	t.addCells(cells)
	t.resolveRow()
	t.validateRow()
}

// addCells appends multiple cells to the current row
func (t *Table) addCells(cells []*cell.Cell) {
	for _, c := range cells {
		t.addCell(c)
	}
}

// addCell adds a single cell to the current position
func (t *Table) addCell(c *cell.Cell) {
	t.advanceCol()
	
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

// initWithFirstRow initializes column arrays from first row
func (t *Table) initWithFirstRow(cells []*cell.Cell) {
	totalCols := 0
	for _, c := range cells {
		totalCols += c.ColSpan()
	}
	t.colWidths = make([]int, totalCols)
	t.colLevels = make([]int, totalCols)
}

// resolveRow resolves flexible column spans for the current row
func (t *Table) resolveRow() {
	if len(t.openFlexCells) == 0 {
		return
	}
	
	flexCells := t.collectFlexCells()
	if len(flexCells) == 0 {
		return
	}
	
	colsAvailable := t.calcAvailableCols()
	if colsAvailable <= 0 {
		return
	}
	
	distributeCols(flexCells, colsAvailable)

	// Update levels for expanded flex cells
	for _, c := range flexCells {
		t.updateLevels(c)
	}
	
	// Clear open flex cells after resolution
	t.openFlexCells = t.openFlexCells[:0]
}

// validateRow ensures all columns are filled
func (t *Table) validateRow() {
	for i, level := range t.colLevels {
		if level == 0 {
			panic(fmt.Sprintf("incomplete row %d: column %d not filled", t.row, i))
		}
	}
}

// collectFlexCells returns flex cells in current row
func (t *Table) collectFlexCells() []*cell.Cell {
	flexCells := make([]*cell.Cell, 0)
	for _, idx := range t.CellsInRow(t.row) {
		c := t.cells[idx]
		if c.IsColFlex() {
			flexCells = append(flexCells, c)
		}
	}
	return flexCells
}

// calcAvailableCols returns columns available for distribution
func (t *Table) calcAvailableCols() int {
	colCount := 0
	for _, idx := range t.CellsInRow(t.row) {
		c := t.cells[idx]
		colCount += c.ColSpan()
	}
	return t.ColCount() - colCount
}

// distributeCols distributes available columns to flex cells respecting maxSpan
func distributeCols(flexCells []*cell.Cell, remaining int) {
	activeCells := make([]*cell.Cell, 0, len(flexCells))
	for _, c := range flexCells {
		if c.ColCanGrow() {
			activeCells = append(activeCells, c)
		}
	}
	
	for remaining > 0 && len(activeCells) > 0 {
		// Filter out cells that can no longer grow
		newActive := activeCells[:0]
		for _, c := range activeCells {
			if remaining <= 0 {
				break
			}
			c.AddColSpan(1)
			remaining--
			
			if c.ColCanGrow() {
				newActive = append(newActive, c)
			}
		}
		activeCells = newActive
	}
}
