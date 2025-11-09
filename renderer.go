package tbl

import (
	"strings"
)

// renderer holds everything needed to produce the final ASCII table.
type renderer struct {
	t         *Table
	colWidths []int     // raw width (no padding) of each column
	grid      [][]*Cell // dense visual grid: grid[row][col] = owning cell
}

// newRenderer constructs a renderer for the given table.
func newRenderer(t *Table) *renderer {
	if !t.g.B.All() {
		panic("tbl: incomplete table")
	}

	r := &renderer{t: t}
	r.buildGrid()
	r.calculateColumnWidths()
	return r
}

// buildGrid creates a dense [][]*Cell where every visual position points to its owning cell.
func (r *renderer) buildGrid() {
	rows := r.t.g.Rows()
	cols := r.t.g.Cols()
	r.grid = make([][]*Cell, rows)
	for i := range r.grid {
		r.grid[i] = make([]*Cell, cols)
	}
	for _, cell := range r.t.cells {
		for rr := cell.r; rr < cell.r+cell.rSpan; rr++ {
			for cc := cell.c; cc < cell.c+cell.cSpan; cc++ {
				r.grid[rr][cc] = cell
			}
		}
	}
}

// calculateColumnWidths computes raw column widths.
// 1. base width = maximum cell.Width() that starts in this column
// 2. for every cell with cSpan>1 enforce minimum by equal distribution
func (r *renderer) calculateColumnWidths() {
	cols := r.t.g.Cols()
	if cols == 0 {
		return
	}
	r.colWidths = make([]int, cols)

	// 1. measure per-column maximum
	for _, cell := range r.t.cells {
		if cell.cSpan == 1 {
			if w := cell.Width(); w > r.colWidths[cell.c] {
				r.colWidths[cell.c] = w
			}
		}
	}

	// 2. enforce colSpan > 1 minimums
	for _, cell := range r.t.cells {
		if cell.cSpan <= 1 {
			continue
		}
		required := cell.Width()
		current := 0
		for i := cell.c; i < cell.c+cell.cSpan; i++ {
			current += r.colWidths[i]
		}
		if current < required {
			short := required - current
			perCol := short / cell.cSpan
			extra := short % cell.cSpan
			for i := 0; i < cell.cSpan; i++ {
				r.colWidths[cell.c+i] += perCol
				if i < extra {
					r.colWidths[cell.c+i]++
				}
			}
		}
	}
}

// render returns the complete ASCII table as a string.
func (r *renderer) render() string {
	rows := r.t.g.Rows()
	if rows == 0 || r.t.g.Cols() == 0 {
		return ""
	}

	var b strings.Builder

	for row := range rows {
		pr := prepareRow(r.grid, r.colWidths, row)

		if pr.borderNeeded {
			r.writeBorder(&b, pr.borderOps)
		}

		for _, lineOps := range pr.contentOps {
			r.writeContent(&b, lineOps)
		}
	}

	// Bottom border
	_, bottomOps := buildBorderOps(r.grid, r.colWidths, r.t.g.Rows())
	r.writeBorder(&b, bottomOps)

	return b.String()
}
