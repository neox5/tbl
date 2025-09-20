package grid

import "fmt"

// Grid holds discrete sizing, prefix sums, and a per-cell reference matrix.
type Grid struct {
	cols    []int
	rows    []int
	colPref []int
	rowPref []int
	cells   [][]Ref  // [row][col] -> Ref
	areas   []*Area  // 1-based store: Ref = index+1
}

// New creates a grid. Requires at least one column and one row. Sizes >= 0.
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
	g.initCells()
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

func (g *Grid) initCells() {
	g.cells = make([][]Ref, len(g.rows))
	for row := range len(g.rows) {
		g.cells[row] = make([]Ref, len(g.cols)) // zero = Nil
	}
}

// Accessors
func (g *Grid) Cols() int           { return len(g.cols) }
func (g *Grid) Rows() int           { return len(g.rows) }
func (g *Grid) ColW(i int) int      { return g.cols[i] }
func (g *Grid) RowH(j int) int      { return g.rows[j] }
func (g *Grid) TotalWidth() int     { return g.colPref[len(g.colPref)-1] }
func (g *Grid) TotalHeight() int    { return g.rowPref[len(g.rowPref)-1] }
func (g *Grid) ColOffset(i int) int { return g.colPref[i] }
func (g *Grid) RowOffset(j int) int { return g.rowPref[j] }
