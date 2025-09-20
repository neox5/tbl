package grid

import "fmt"

// Grid with explicit sizing and simple storage.
type Grid struct {
	cells    [][]ID
	areas    map[ID]*Area
	rowIndex map[int][]ID
	nextID   ID

	cols []int // column widths (default 0)
	rows []int // row heights (default 0)
}

// New builds a grid. Zero counts allowed. Negatives panic with values.
func New(numCols, numRows int) *Grid {
	if numCols < 0 || numRows < 0 {
		panic(fmt.Errorf("grid: invalid size cols=%d rows=%d", numCols, numRows))
	}
	g := &Grid{
		cells:    make([][]ID, numRows),
		areas:    make(map[ID]*Area),
		rowIndex: make(map[int][]ID),
		nextID:   1,
		cols:     make([]int, numCols),
		rows:     make([]int, numRows),
	}
	for r := range numRows {
		g.cells[r] = make([]ID, numCols)
	}
	return g
}

// Dimensions from sizing arrays.
func (g *Grid) Cols() int { return len(g.cols) }
func (g *Grid) Rows() int { return len(g.rows) }

// Extending grid.

func (g *Grid) AddCol() {
	g.cols = append(g.cols, 0) // width default 0
	for r := range g.cells {
		g.cells[r] = append(g.cells[r], Nil)
	}
}

func (g *Grid) AddRow() {
	g.rows = append(g.rows, 0) // height default 0
	g.cells = append(g.cells, make([]ID, g.Cols()))
}

// AddArea registers and places the area. On failure it removes the registry entry.
// Side effects order: registry -> validate -> placeArea -> indexAdd.
func (g *Grid) AddArea(a *Area) (ID, error) {
	if a == nil {
		return Nil, fmt.Errorf("grid: AddArea nil area")
	}
	id := g.nextID
	g.nextID++
	g.areas[id] = a

	if err := g.validateBounds(a); err != nil {
		delete(g.areas, id)
		return Nil, err
	}
	if err := g.validateOverlap(a); err != nil {
		delete(g.areas, id)
		return Nil, err
	}

	a.ForEachCell(func(c, r int) { g.cells[r][c] = id }) // occupy grid cells
	g.indexAdd(id)
	return id, nil
}
