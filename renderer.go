package tbl

import (
	"strings"
)

// renderer holds everything needed to produce the final ASCII table.
type renderer struct {
	t             *Table
	styles        map[ID]CellStyle // cached resolved styles
	grid          [][]*Cell
	colMaxPadding []int
	colWidths     []int
	rowHeights    []int
	vBoundaries   []bool          // vertical boundary presence
	hBoundaries   []bool          // horizontal boundary presence
	cellLayouts   map[ID][]string // pre-computed content lines per cell
	tpl           CharTemplate    // table-level template for rendering
}

// rowCount returns total number of rows in grid.
func (r *renderer) rowCount() int {
	return len(r.grid)
}

// colCount returns total number of columns in grid.
func (r *renderer) colCount() int {
	if len(r.grid) == 0 {
		return 0
	}
	return len(r.grid[0])
}

// lastCellInRow returns true if the cell is positioned at the end of the row.
func (r *renderer) lastCellInRow(cell *Cell) bool {
	return cell.c+cell.cSpan == r.colCount()
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
	r := &renderer{
		t:             t,
		styles:        make(map[ID]CellStyle),
		grid:          make([][]*Cell, t.g.Rows()),
		colMaxPadding: make([]int, t.g.Cols()),
		colWidths:     make([]int, t.g.Cols()),
		rowHeights:    make([]int, t.g.Rows()),
		vBoundaries:   make([]bool, t.g.Cols()+1),
		hBoundaries:   make([]bool, t.g.Rows()+1),
		cellLayouts:   make(map[ID][]string),
		tpl:           t.defaultStyle.Template,
	}

	// Initialize grid rows
	for i := range r.grid {
		r.grid[i] = make([]*Cell, t.g.Cols())
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

// render returns the complete ASCII table as a string.
func (r *renderer) render() string {
	if r.rowCount() == 0 || r.colCount() == 0 {
		return ""
	}

	lines := r.buildRenderLines()

	var b strings.Builder
	for _, lineOps := range lines {
		r.writeLine(&b, lineOps)
	}

	return b.String()
}
