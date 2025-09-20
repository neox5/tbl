// Package grid defines a 2D grid by column widths and row heights.
// A Cell represents a rectangular region of one or more spanned cells.
package grid

import "fmt"

// Cell represents one or more spanned cells in a Grid.
// It refers back to its parent grid to compute geometry.
type Cell struct {
	g       *Grid
	col     int
	row     int
	colSpan int
	rowSpan int
}

// Cell returns a Cell for (col,row) with span 1x1. Panics if out of range.
func (g *Grid) Cell(col, row int) Cell {
	return g.CellSpan(col, row, 1, 1)
}

// CellSpan returns a Cell for (col,row) spanning colSpan × rowSpan cells.
// Panics if indices or spans are out of range.
func (g *Grid) CellSpan(col, row, colSpan, rowSpan int) Cell {
	if col < 0 || col >= len(g.cols) {
		panic(fmt.Sprintf("grid: column %d out of range [0,%d)", col, len(g.cols)))
	}
	if row < 0 || row >= len(g.rows) {
		panic(fmt.Sprintf("grid: row %d out of range [0,%d)", row, len(g.rows)))
	}
	if colSpan <= 0 || rowSpan <= 0 {
		panic(fmt.Sprintf("grid: invalid span colSpan=%d rowSpan=%d (must be >0)", colSpan, rowSpan))
	}
	if col+colSpan > len(g.cols) {
		panic(fmt.Sprintf("grid: colSpan=%d at col=%d exceeds column count %d", colSpan, col, len(g.cols)))
	}
	if row+rowSpan > len(g.rows) {
		panic(fmt.Sprintf("grid: rowSpan=%d at row=%d exceeds row count %d", rowSpan, row, len(g.rows)))
	}

	return Cell{
		g:       g,
		col:     col,
		row:     row,
		colSpan: colSpan,
		rowSpan: rowSpan,
	}
}

// Col returns the starting column index.
func (c Cell) Col() int { return c.col }

// Row returns the starting row index.
func (c Cell) Row() int { return c.row }

// ColSpan returns the number of columns spanned.
func (c Cell) ColSpan() int { return c.colSpan }

// RowSpan returns the number of rows spanned.
func (c Cell) RowSpan() int { return c.rowSpan }

// W returns the total width of the cell.
func (c Cell) W() int {
	w := 0
	for _, cw := range c.g.cols[c.col : c.col+c.colSpan] {
		w += cw
	}
	return w
}

// H returns the total height of the cell.
func (c Cell) H() int {
	h := 0
	for _, rh := range c.g.rows[c.row : c.row+c.rowSpan] {
		h += rh
	}
	return h
}
