package tbl

import (
	"io"
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
	var b strings.Builder
	r.renderTo(&b)
	return b.String()
}

// renderTo writes the table to w.
func (r *renderer) renderTo(w io.Writer) error {
	var b strings.Builder
	r.renderToBuilder(&b)
	_, err := io.WriteString(w, b.String())
	return err
}

// renderToBuilder does the actual rendering into a strings.Builder.
func (r *renderer) renderToBuilder(b *strings.Builder) {
	if r.t.g.Cols() == 0 || r.t.g.Rows() == 0 {
		return
	}
	r.writeHLine(b)
	for row := 0; row < r.t.g.Rows(); row++ {
		r.writeRow(b, row)
	}
	r.writeHLine(b)
}

// writeHLine emits a full horizontal border (+---+-...-+).
func (r *renderer) writeHLine(b *strings.Builder) {
	b.WriteByte('+')
	for col, w := range r.colWidths {
		if col > 0 {
			b.WriteByte('+')
		}
		for i := 0; i < w+2; i++ { // +2 for the two padding spaces
			b.WriteByte('-')
		}
	}
	b.WriteByte('+')
	b.WriteByte('\n')
}

// writeRow emits one visual content row.
func (r *renderer) writeRow(b *strings.Builder, row int) {
	b.WriteByte('|')
	col := 0
	for col < r.t.g.Cols() {
		cell := r.grid[row][col]
		span := cell.cSpan
		total := 0
		for i := range span {
			total += r.colWidths[col+i]
		}
		// add padding spaces but remove internal separators
		total += 2          // left + right padding
		total -= (span - 1) // remove (span-1) vertical bars

		content := cell.Content()
		padRight := total - len(content) - 2 // -2 because we already added both paddings
		b.WriteByte(' ')
		b.WriteString(content)
		for range padRight {
			b.WriteByte(' ')
		}
		b.WriteByte(' ')

		col += span
		if col < r.t.g.Cols() {
			b.WriteByte('|')
		}
	}
	b.WriteByte('|')
	b.WriteByte('\n')
}
