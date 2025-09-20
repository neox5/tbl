// Package grid defines a 2D grid by column widths and row heights.
// An Area represents a rectangular region of one or more spanned cells.
package grid

import "fmt"

// Area represents one or more spanned cells in a Grid.
// It refers back to its parent grid to compute geometry.
type Area struct {
	g       *Grid
	col     int
	row     int
	colSpan int
	rowSpan int
}

// Area returns an Area for (col,row) with span 1x1. Panics if out of range.
func (g *Grid) Area(col, row int) Area {
	return g.AreaSpan(col, row, 1, 1)
}

// AreaSpan returns an Area for (col,row) spanning colSpan × rowSpan cells.
// Panics if indices or spans are out of range.
func (g *Grid) AreaSpan(col, row, colSpan, rowSpan int) Area {
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

	return Area{
		g:       g,
		col:     col,
		row:     row,
		colSpan: colSpan,
		rowSpan: rowSpan,
	}
}

// Col returns the starting column index.
func (a Area) Col() int { return a.col }

// Row returns the starting row index.
func (a Area) Row() int { return a.row }

// ColSpan returns the number of columns spanned.
func (a Area) ColSpan() int { return a.colSpan }

// RowSpan returns the number of rows spanned.
func (a Area) RowSpan() int { return a.rowSpan }

// W returns the total width of the area.
func (a Area) W() int {
	w := 0
	for _, cw := range a.g.cols[a.col : a.col+a.colSpan] {
		w += cw
	}
	return w
}

// H returns the total height of the area.
func (a Area) H() int {
	h := 0
	for _, rh := range a.g.rows[a.row : a.row+a.rowSpan] {
		h += rh
	}
	return h
}
