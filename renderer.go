package tbl

import (
	"fmt"
	"strings"
)

// renderer holds everything needed to produce the final ASCII table.
type renderer struct {
	t             *Table
	colMaxPadding []int
	colWidths     []int
	rowHeights    []int
	colBorders    []bool // Physical column border presence
	rowBorders    []bool // Physical row border presence
	grid          [][]*Cell
	cellLayouts   map[ID][]string  // pre-computed content lines per cell
	styles        map[ID]CellStyle // cached resolved styles
}

// newRenderer constructs a renderer for the given table.
// Builds grid structure, calculates dimensions, and pre-computes cell layouts.
//
// Pipeline:
//  1. Single pass: cache styles, build grid, track dimensions and borders
//  2. Enforce global table width constraint
//  3. Calculate row heights with finalized column widths
//  4. Generate cell layouts
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
		colBorders:    make([]bool, cols+1),
		rowBorders:    make([]bool, rows+1),
		grid:          make([][]*Cell, rows),
		cellLayouts:   make(map[ID][]string),
		styles:        make(map[ID]CellStyle),
	}

	// Initialize grid rows
	for i := range r.grid {
		r.grid[i] = make([]*Cell, cols)
	}

	// Single pass: Cache styles and track dimensions
	for _, cell := range t.cells {
		r.cacheStyle(cell)
		r.populateGrid(cell)
		r.trackDimensions(cell)
	}

	// Enforce global table width constraint
	r.enforceTableMaxWidth()

	// Calculate heights with finalized column widths
	r.calculateHeights()

	// Generate layouts with finalized dimensions
	r.buildCellLayouts()

	return r
}

// cacheStyle resolves and stores cell style.
func (r *renderer) cacheStyle(cell *Cell) {
	r.styles[cell.id] = r.t.resolveStyle(cell)
}

// populateGrid fills grid positions for cell.
func (r *renderer) populateGrid(cell *Cell) {
	for rr := cell.r; rr < cell.r+cell.rSpan; rr++ {
		for cc := cell.c; cc < cell.c+cell.cSpan; cc++ {
			r.grid[rr][cc] = cell
		}
	}
}

// trackDimensions calculates widths, padding, and borders in one pass.
func (r *renderer) trackDimensions(cell *Cell) {
	r.updateColMaxPadding(cell)
	r.updateColWidths(cell)
	r.updateBorders(cell)
}

// updateColMaxPadding tracks maximum padding for origin column only.
func (r *renderer) updateColMaxPadding(cell *Cell) {
	style := r.styles[cell.id]
	padding := style.Padding.Left + style.Padding.Right

	if padding > r.colMaxPadding[cell.c] {
		r.colMaxPadding[cell.c] = padding
	}
}

// updateColWidths updates column widths for cell.
// Applies ColConfig constraints after natural calculation.
// Handles both single-column (span=1) and multi-column (span>1) cells.
func (r *renderer) updateColWidths(cell *Cell) {
	style := r.styles[cell.id]
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

// updateBorders marks border presence based on cell style.
// Border exists if Sides requests visual OR Physical requests space.
func (r *renderer) updateBorders(cell *Cell) {
	style := r.styles[cell.id]

	if style.Border.Has(BorderTop) {
		r.rowBorders[cell.r] = true
	}

	if style.Border.Has(BorderBottom) {
		r.rowBorders[cell.r+cell.rSpan] = true
	}

	if style.Border.Has(BorderLeft) {
		r.colBorders[cell.c] = true
	}

	if style.Border.Has(BorderRight) {
		r.colBorders[cell.c+cell.cSpan] = true
	}
}

// calculateHeights processes all cells with finalized column widths.
func (r *renderer) calculateHeights() {
	for _, cell := range r.t.cells {
		r.updateRowHeights(cell)
	}
}

// updateRowHeights calculates and updates row heights for cell.
// Uses final column widths to determine content wrapping.
// Handles both single-row (span=1) and multi-row (span>1) cells.
func (r *renderer) updateRowHeights(cell *Cell) {
	style := r.styles[cell.id]

	// Calculate final width for this cell
	totalWidth := r.cellWidth(cell)
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
// Accounts for physical borders in total width calculation.
// Reduction strategy: left-to-right until hitting minimums, remainder distributed evenly.
//
// Algorithm:
//  1. Calculate physical border count
//  2. Calculate available width for content (MaxWidth - borders)
//  3. Calculate minimum possible table width and identify reducible columns
//  4. Panic if minimum exceeds available width
//  5. Calculate current total width
//  6. If total <= available width: no action
//  7. Reduce columns left-to-right until hitting minimums
//  8. Distribute remainder evenly across reducible columns
func (r *renderer) enforceTableMaxWidth() {
	if r.t.tableConfig.MaxWidth <= 0 {
		return
	}

	// Calculate physical border count
	borderCount := 0
	for _, exists := range r.colBorders {
		if exists {
			borderCount++
		}
	}

	// Available width for content
	availableWidth := r.t.tableConfig.MaxWidth - borderCount

	if availableWidth <= 0 {
		panic(fmt.Sprintf("tbl: MaxWidth %d insufficient for %d borders",
			r.t.tableConfig.MaxWidth, borderCount))
	}

	// Calculate minimum possible width and identify reducible columns
	minPossible := 0
	reducible := make([]int, 0, len(r.colWidths))

	for col, width := range r.colWidths {
		minWidth, fixed := r.colMinWidth(col)
		minPossible += minWidth

		if !fixed && width > minWidth {
			reducible = append(reducible, col)
		}
	}

	if minPossible > availableWidth {
		panic(fmt.Sprintf("tbl: impossible MaxWidth constraint: minimum possible width %d exceeds available width %d (MaxWidth %d - %d borders)",
			minPossible, availableWidth, r.t.tableConfig.MaxWidth, borderCount))
	}

	// Calculate current total
	total := 0
	for _, w := range r.colWidths {
		total += w
	}

	if total <= availableWidth {
		return
	}

	excess := total - availableWidth

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
// Accounts for physical borders that are removed within cell span.
func (r *renderer) cellWidth(cell *Cell) int {
	width := 0

	// Sum column widths
	for i := range cell.cSpan {
		width += r.colWidths[cell.c+i]
	}

	// Add space from removed internal vertical borders
	// Check borders between columns within cell span
	for i := 1; i < cell.cSpan; i++ {
		if r.colBorders[cell.c+i] {
			width++ // Border column removed, add its width
		}
	}

	return width
}

// cellHeight calculates total layout height for cell spanning multiple rows.
// Accounts for physical borders that are removed within cell span.
func (r *renderer) cellHeight(cell *Cell) int {
	height := 0

	// Sum row heights
	for i := range cell.rSpan {
		height += r.rowHeights[cell.r+i]
	}

	// Add space from removed internal horizontal borders
	// Check borders between rows within cell span
	for i := 1; i < cell.rSpan; i++ {
		if r.rowBorders[cell.r+i] {
			height++ // Border line removed, add its height
		}
	}

	return height
}

// buildCellLayouts generates content lines for all cells using finalized dimensions.
// Content includes padding, horizontal and vertical alignment - ready to render.
func (r *renderer) buildCellLayouts() {
	for _, cell := range r.t.cells {
		style := r.styles[cell.id]

		// Calculate total dimensions (including padding)
		totalWidth := r.cellWidth(cell)
		totalHeight := r.cellHeight(cell)

		// Generate complete lines with padding
		r.cellLayouts[cell.id] = cell.Layout(totalWidth, totalHeight, style)
	}
}

// render returns the complete ASCII table as a string.
func (r *renderer) render() string {
	rows := r.t.g.Rows()
	if rows == 0 || r.t.g.Cols() == 0 {
		return ""
	}

	lines := r.buildRenderLines()

	var b strings.Builder
	for _, lineOps := range lines {
		r.writeLine(&b, lineOps)
	}

	return b.String()
}
