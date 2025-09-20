package grid

import "fmt"

// Grid defines a grid of rows and columns with integer widths/heights.
type Grid struct {
	cols    []int // column widths
	rows    []int // row heights
	colPref []int // prefix sums, len = len(cols)+1
	rowPref []int // prefix sums, len = len(rows)+1
}

// New creates a grid. Panics if any width/height < 0.
func New(cols, rows []int) *Grid {
	if len(cols) == 0 || len(rows) == 0 {
		panic("grid: need at least one column and one row")
	}
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
	g := &Grid{cols: c, rows: r}
	g.buildPrefixes()
	return g
}

func (g *Grid) buildPrefixes() {
	g.colPref = make([]int, len(g.cols)+1)
	for i := range len(g.cols) {
		g.colPref[i+1] = g.colPref[i] + g.cols[i]
	}
	g.rowPref = make([]int, len(g.rows)+1)
	for j := range len(g.rows) {
		g.rowPref[j+1] = g.rowPref[j] + g.rows[j]
	}
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
func (g *Grid) TotalWidth() int { return g.colPref[len(g.colPref)-1] }

// TotalHeight is the sum of all row heights.
func (g *Grid) TotalHeight() int { return g.rowPref[len(g.rowPref)-1] }

// ColOffset returns the x-offset of the left edge of column i.
func (g *Grid) ColOffset(i int) int { return g.colPref[i] }

// RowOffset returns the y-offset of the top edge of row j.
func (g *Grid) RowOffset(j int) int { return g.rowPref[j] }

// CellOrigin returns the top-left coordinate of a cell.
func (g *Grid) CellOrigin(c Cell) (x, y int) {
	return g.colPref[c.Col], g.rowPref[c.Row]
}

// AreaOrigin returns the top-left coordinate of an area.
func (g *Grid) AreaOrigin(a Area) (x, y int) {
	return g.colPref[a.Col()], g.rowPref[a.Row()]
}

// AreaSize returns the width and height of an area.
func (g *Grid) AreaSize(a Area) (w, h int) {
	w = g.colPref[a.Col()+a.ColSpan()] - g.colPref[a.Col()]
	h = g.rowPref[a.Row()+a.RowSpan()] - g.rowPref[a.Row()]
	return
}

// AreaRect returns x, y, w, h for an area in absolute units.
func (g *Grid) AreaRect(a Area) (x, y, w, h int) {
	x, y = g.AreaOrigin(a)
	w, h = g.AreaSize(a)
	return
}
