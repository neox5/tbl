package tbl

import (
	"strings"
)

// renderer holds everything needed to produce the final ASCII table.
type renderer struct {
	t           *Table
	colWidths   []int
	rowHeights  []int
	grid        [][]*Cell
	cellLayouts map[ID][]string // pre-computed content lines per cell
}

// newRenderer constructs a renderer for the given table.
// Builds grid structure, calculates dimensions, and pre-computes cell layouts.
func newRenderer(t *Table) *renderer {
	if !t.g.B.All() {
		panic("tbl: incomplete table")
	}

	rows := t.g.Rows()
	cols := t.g.Cols()

	r := &renderer{
		t:           t,
		colWidths:   make([]int, cols),
		rowHeights:  make([]int, rows),
		grid:        make([][]*Cell, rows),
		cellLayouts: make(map[ID][]string),
	}

	// Initialize grid rows
	for i := range r.grid {
		r.grid[i] = make([]*Cell, cols)
	}

	// Single pass: build grid + calculate dimensions
	for _, cell := range t.cells {
		// Populate grid
		for rr := cell.r; rr < cell.r+cell.rSpan; rr++ {
			for cc := cell.c; cc < cell.c+cell.cSpan; cc++ {
				r.grid[rr][cc] = cell
			}
		}

		// Update column widths
		r.updateColWidths(cell)

		// Update row heights
		r.updateRowHeights(cell)
	}

	// Generate layouts with finalized dimensions
	r.buildCellLayouts()

	return r
}

// updateColWidths updates column widths for cell.
// Handles both single-column (span=1) and multi-column (span>1) cells.
// Multi-column cells distribute width shortfall evenly across spanned columns.
func (r *renderer) updateColWidths(cell *Cell) {
	style := r.t.resolveStyle(cell)
	contentWidth := cell.Width()
	required := contentWidth + style.Padding.Left + style.Padding.Right

	if cell.cSpan == 1 {
		if required > r.colWidths[cell.c] {
			r.colWidths[cell.c] = required
		}
		return
	}

	// Multi-column: check current total and distribute shortfall
	current := 0
	for i := cell.c; i < cell.c+cell.cSpan; i++ {
		current += r.colWidths[i]
	}

	if current >= required {
		return
	}

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

// updateRowHeights updates row heights for cell.
// Handles both single-row (span=1) and multi-row (span>1) cells.
// Multi-row cells distribute height shortfall evenly across spanned rows.
func (r *renderer) updateRowHeights(cell *Cell) {
	style := r.t.resolveStyle(cell)
	contentHeight := cell.Height()
	required := style.Padding.Top + contentHeight + style.Padding.Bottom

	if cell.rSpan == 1 {
		if required > r.rowHeights[cell.r] {
			r.rowHeights[cell.r] = required
		}
		return
	}

	// Multi-row: check current total and distribute shortfall
	current := 0
	for i := cell.r; i < cell.r+cell.rSpan; i++ {
		current += r.rowHeights[i]
	}

	if current >= required {
		return
	}

	short := required - current
	perRow := short / cell.rSpan
	extra := short % cell.rSpan

	for i := 0; i < cell.rSpan; i++ {
		r.rowHeights[cell.r+i] += perRow
		if i < extra {
			r.rowHeights[cell.r+i]++
		}
	}
}

// cellWidth calculates total layout width for cell spanning multiple columns.
// colWidths already include padding, so only add space for removed VLines.
// Formula: sum(colWidths) + (colSpan - 1)
func cellWidth(colWidths []int, cell *Cell) int {
	width := 0
	for i := 0; i < cell.cSpan; i++ {
		width += colWidths[cell.c+i]
	}
	// Add space from removed VLines: 1 char per span junction
	width += cell.cSpan - 1
	return width
}

// buildCellLayouts generates content lines for all cells using finalized dimensions.
// Content includes horizontal and vertical alignment but excludes padding.
// Padding applied during rendering via buildContentOps.
func (r *renderer) buildCellLayouts() {
	for _, cell := range r.t.cells {
		style := r.t.resolveStyle(cell)

		// Calculate content width (total width minus padding)
		totalWidth := cellWidth(r.colWidths, cell)
		contentWidth := totalWidth - style.Padding.Left - style.Padding.Right

		// Calculate content height (total height minus padding)
		totalHeight := 0
		for i := cell.r; i < cell.r+cell.rSpan; i++ {
			totalHeight += r.rowHeights[i]
		}
		contentHeight := totalHeight - style.Padding.Top - style.Padding.Bottom

		// Generate content lines with resolved alignment
		contentLines := cell.Layout(contentWidth, contentHeight, style.HAlign, style.VAlign)
		r.cellLayouts[cell.id] = contentLines
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
		pr := prepareRow(
			r.grid,
			r.colWidths,
			r.cellLayouts,
			r.rowHeights,
			row,
			r.t.resolveStyle,
		)

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
