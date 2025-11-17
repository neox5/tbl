package tbl

import (
	"fmt"
	"strings"
)

// renderer holds everything needed to produce the final ASCII table.
type renderer struct {
	t             *Table
	colMaxPadding []int // max(paddingLeft + paddingRight) per column
	colWidths     []int
	rowHeights    []int
	grid          [][]*Cell
	cellLayouts   map[ID][]string // pre-computed content lines per cell
}

// newRenderer constructs a renderer for the given table.
// Builds grid structure, calculates dimensions, and pre-computes cell layouts.
//
// Calculation order:
//  1. Build grid and track max padding per column
//  2. Calculate natural column widths with constraints
//  3. Enforce global table width constraint
//  4. Calculate row heights with final column widths
//  5. Generate cell layouts
func newRenderer(t *Table) *renderer {
	if !t.g.B.All() {
		panic("tbl: incomplete table")
	}

	rows := t.g.Rows()
	cols := t.g.Cols()

	r := &renderer{
		t:             t,
		colMaxPadding: make([]int, cols),
		colWidths:     make([]int, cols),
		rowHeights:    make([]int, rows),
		grid:          make([][]*Cell, rows),
		cellLayouts:   make(map[ID][]string),
	}

	// Initialize grid rows
	for i := range r.grid {
		r.grid[i] = make([]*Cell, cols)
	}

	// Pass 1: Grid population, padding tracking, and width calculation
	for _, cell := range t.cells {
		// Resolve style once per cell
		style := t.resolveStyle(cell)

		// Populate grid
		for rr := cell.r; rr < cell.r+cell.rSpan; rr++ {
			for cc := cell.c; cc < cell.c+cell.cSpan; cc++ {
				r.grid[rr][cc] = cell
			}
		}

		// Track max padding for origin column only
		r.updateColMaxPadding(cell, style)

		// Update column widths
		r.updateColWidths(cell, style)
	}

	// Enforce global table width constraint
	r.enforceTableMaxWidth()

	// Pass 2: Height calculation with final column widths
	for _, cell := range t.cells {
		style := t.resolveStyle(cell)
		r.updateRowHeights(cell, style)
	}

	// Generate layouts with finalized dimensions
	r.buildCellLayouts()

	return r
}

// updateColMaxPadding tracks maximum padding for origin column only.
func (r *renderer) updateColMaxPadding(cell *Cell, style CellStyle) {
	padding := style.Padding.Left + style.Padding.Right

	if padding > r.colMaxPadding[cell.c] {
		r.colMaxPadding[cell.c] = padding
	}
}

// updateColWidths updates column widths for cell.
// Applies ColConfig constraints after natural calculation.
// Handles both single-column (span=1) and multi-column (span>1) cells.
func (r *renderer) updateColWidths(cell *Cell, style CellStyle) {
	contentWidth := cell.Width()
	required := contentWidth + style.Padding.Left + style.Padding.Right

	// Apply column constraints to origin column only
	if cfg, ok := r.t.colConfigs[cell.c]; ok {
		required = r.applyColConstraints(required, cfg)
	}

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

// applyColConstraints enforces ColConfig limits on calculated width.
// Priority: Width (fixed) > MaxWidth > MinWidth
func (r *renderer) applyColConstraints(width int, cfg ColConfig) int {
	// Fixed width overrides everything
	if cfg.Width > 0 {
		return cfg.Width
	}

	// Apply max constraint
	if cfg.MaxWidth > 0 && width > cfg.MaxWidth {
		width = cfg.MaxWidth
	}

	// Apply min constraint
	if cfg.MinWidth > 0 && width < cfg.MinWidth {
		width = cfg.MinWidth
	}

	return width
}

// updateRowHeights calculates and updates row heights for cell.
// Uses final column widths to determine content wrapping.
// Handles both single-row (span=1) and multi-row (span>1) cells.
func (r *renderer) updateRowHeights(cell *Cell, style CellStyle) {
	// Calculate final width for this cell
	totalWidth := cellWidth(r.colWidths, cell)
	contentWidth := totalWidth - style.Padding.Left - style.Padding.Right

	// Rebuild content lines with final width constraint
	contentLines := buildRawLines(cell.content, contentWidth)
	contentHeight := len(contentLines)

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

// colMinWidth calculates minimum width for column.
// Returns minWidth and whether column is fixed (cannot be reduced).
func (r *renderer) colMinWidth(col int) (int, bool) {
	cfg, hasCfg := r.t.colConfigs[col]

	// Fixed width
	if hasCfg && cfg.Width > 0 {
		return cfg.Width, true
	}

	// Start with padding + 1 char
	minWidth := r.colMaxPadding[col] + 1

	// Apply configured MinWidth if larger
	if hasCfg && cfg.MinWidth > 0 {
		minWidth = max(minWidth, cfg.MinWidth)
	}

	return minWidth, false
}

// enforceTableMaxWidth reduces column widths if total exceeds TableConfig.MaxWidth.
// Reduction strategy: left-to-right until hitting minimums, remainder distributed evenly.
//
// Algorithm:
//  1. Calculate minimum possible table width and identify reducible columns
//  2. Panic if minimum exceeds MaxWidth
//  3. Calculate current total width
//  4. If total <= MaxWidth: no action
//  5. Reduce columns left-to-right until hitting minimums
//  6. Distribute remainder evenly across reducible columns
func (r *renderer) enforceTableMaxWidth() {
	if r.t.tableConfig.MaxWidth <= 0 {
		return
	}

	// Calculate minimum table width and identify reducible columns
	minPossible := 0
	reducible := make([]int, 0, len(r.colWidths))

	for col, width := range r.colWidths {
		minWidth, fixed := r.colMinWidth(col)
		minPossible += minWidth

		if !fixed && width > minWidth {
			reducible = append(reducible, col)
		}
	}

	if minPossible > r.t.tableConfig.MaxWidth {
		panic(fmt.Sprintf("tbl: impossible MaxWidth constraint: minimum possible width %d exceeds MaxWidth %d", minPossible, r.t.tableConfig.MaxWidth))
	}

	// Calculate current total
	total := 0
	for _, w := range r.colWidths {
		total += w
	}

	if total <= r.t.tableConfig.MaxWidth {
		return
	}

	excess := total - r.t.tableConfig.MaxWidth

	// Reduce left-to-right until hitting minimums
	for _, col := range reducible {
		if excess == 0 {
			break
		}

		minWidth, _ := r.colMinWidth(col)
		available := r.colWidths[col] - minWidth

		if available <= 0 {
			continue
		}

		reduction := min(available, excess)
		r.colWidths[col] -= reduction
		excess -= reduction
	}

	// Distribute remainder evenly
	if excess > 0 && len(reducible) > 0 {
		perCol := excess / len(reducible)
		remainder := excess % len(reducible)

		for i, col := range reducible {
			minWidth, _ := r.colMinWidth(col)

			reduction := perCol
			if i < remainder {
				reduction++
			}

			// Ensure we don't go below minimum
			newWidth := r.colWidths[col] - reduction
			if newWidth < minWidth {
				reduction = r.colWidths[col] - minWidth
			}

			r.colWidths[col] -= reduction
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
