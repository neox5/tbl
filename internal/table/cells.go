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
	t.resolveRow(t.row)
	
	// Update levels for current row after resolution
	flexCells := t.collectFlexCells(t.row)
	for _, c := range flexCells {
		t.updateLevels(c)
	}
	
	t.validateColLevels()

	// Early return if no flex rows exist
	if len(t.flexRows) == 0 {
		t.colsFixed = true
		return
	}

	// Skip resolution if current row has flex cells
	if t.isFlexRow(t.row) {
		return
	}

	// Set colsFixed = true BEFORE resolving flex rows
	t.colsFixed = true

	// Resolve all previous flex rows
	for row := range t.flexRows {
		t.resolveFlexRow(row)
	}
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
		t.addFlexRow(t.row)
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
	if c.IsColFlex() && !t.colsFixed {
		w = 0
	}

	// Distribute width across spanned columns
	dw := w / c.ColSpan()
	r := w % c.ColSpan()
	for i := range c.ColSpan() {
		colW := dw
		if i < r {
			colW++ // Distribute remainder to first columns
		}

		if t.colWidths[c.ColStart()+i] < colW {
			t.colWidths[c.ColStart()+i] = colW
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

// resolveRow resolves flexible column spans for the specified row
func (t *Table) resolveRow(row int) {
	flexCells := t.collectFlexCells(row)
	if len(flexCells) == 0 {
		return
	}

	colsAvailable := t.calcAvailableCols(row)
	if colsAvailable <= 0 {
		return
	}

	distributeCols(flexCells, colsAvailable)
	t.recalculatePositions(row)
}

// resolveFlexRow resolves flexible column spans and updates widths for the specified row
func (t *Table) resolveFlexRow(row int) {
	t.resolveRow(row)
	t.removeFlexRow(row)

	// Update widths for all cells in this row
	for _, idx := range t.CellsInRow(row) {
		t.updateWidths(t.cells[idx])
	}
}

// validateColLevels ensures all columns are filled in the current row
func (t *Table) validateColLevels() {
	for i, level := range t.colLevels {
		if level == 0 {
			panic(fmt.Sprintf("incomplete row %d: column %d not filled", t.row, i))
		}
	}
}

// collectFlexCells returns flex cells in specified row
func (t *Table) collectFlexCells(row int) []*cell.Cell {
	flexCells := make([]*cell.Cell, 0)
	for _, idx := range t.CellsInRow(row) {
		c := t.cells[idx]
		if c.IsColFlex() {
			flexCells = append(flexCells, c)
		}
	}
	return flexCells
}

// calcAvailableCols returns columns available for distribution in specified row
func (t *Table) calcAvailableCols(row int) int {
	colCount := 0
	for _, idx := range t.CellsInRow(row) {
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

// recalculatePositions updates ColStart positions for all cells in a row after span changes
func (t *Table) recalculatePositions(row int) {
	cellIndices := t.CellsInRow(row)
	if len(cellIndices) == 0 {
		return
	}

	// Sort cell indices by their current ColStart to maintain order
	for i := range len(cellIndices)-1 {
		for j := i + 1; j < len(cellIndices); j++ {
			if t.cells[cellIndices[i]].ColStart() > t.cells[cellIndices[j]].ColStart() {
				cellIndices[i], cellIndices[j] = cellIndices[j], cellIndices[i]
			}
		}
	}

	// Recalculate positions sequentially
	pos := 0
	for _, idx := range cellIndices {
		c := t.cells[idx]
		c.SetColStart(pos)
		pos += c.ColSpan()
	}
}
