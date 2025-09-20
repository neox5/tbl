// Package grid defines a 2D grid by column widths and row heights.
// Columns and rows are indexed from 0. Widths/heights must be >= 0.
package grid

import "fmt"

type Grid struct {
	cols []int // column widths
	rows []int // row heights
}

// New creates a grid. Panics if any width/height < 0.
func New(cols, rows []int) *Grid {
	c := make([]int, len(cols))
	r := make([]int, len(rows))
	copy(c, cols)
	copy(r, rows)
	for i, w := range c {
		if w < 0 {
			panic(fmt.Sprintf("grid: negative column width at index %d", i))
		}
	}
	for j, h := range r {
		if h < 0 {
			panic(fmt.Sprintf("grid: negative row height at index %d", j))
		}
	}
	return &Grid{cols: c, rows: r}
}

// Cols returns the number of columns.
func (g *Grid) Cols() int { return len(g.cols) }

// Rows returns the number of rows.
func (g *Grid) Rows() int { return len(g.rows) }

// ColW returns the width of column i. Panics if out of range.
func (g *Grid) ColW(i int) int { return g.cols[i] }

// RowH returns the height of row j. Panics if out of range.
func (g *Grid) RowH(j int) int { return g.rows[j] }

// TotalWidth is the sum of all column widths.
func (g *Grid) TotalWidth() int {
	s := 0
	for _, w := range g.cols {
		s += w
	}
	return s
}

// TotalHeight is the sum of all row heights.
func (g *Grid) TotalHeight() int {
	s := 0
	for _, h := range g.rows {
		s += h
	}
	return s
}
