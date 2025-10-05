package grid

import (
	"fmt"

	"github.com/neox5/btmp"
)

type ID int64

const Nil ID = 0

// Grid with explicit sizing and simple storage. Row-major wrapper.
type Grid struct {
	occ *btmp.Grid

	areas    map[ID]*Area
	rowIndex map[int][]ID
	toMove   map[ID]struct{}
	nextID   ID
}

// New builds a grid. Parameters are (rows, cols).
func New(numRows, numCols int) *Grid {
	if numRows < 0 || numCols < 0 {
		panic(fmt.Errorf("grid: invalid size rows=%d cols=%d", numRows, numCols))
	}
	return &Grid{
		occ:      btmp.NewGridWithSize(numRows, numCols),
		areas:    make(map[ID]*Area),
		rowIndex: make(map[int][]ID),
		toMove:   make(map[ID]struct{}),
		nextID:   1,
	}
}

func (g *Grid) Rows() int { return g.occ.Rows() }
func (g *Grid) Cols() int { return g.occ.Cols() }

func (g *Grid) AddRow() { g.occ.GrowRows(1) }
func (g *Grid) AddCol() { g.occ.GrowCols(1) }

// IsFree reports whether the area is unoccupied in the grid.
func (g *Grid) IsFree(a *Area) bool {
	return g.occ.IsFree(a.Row(), a.Col(), a.RowSpan(), a.ColSpan())
}

// ValidateArea checks boundaries.
func (g *Grid) ValidateArea(a *Area) error {
	return g.occ.ValidateRect(a.Row(), a.Col(), a.RowSpan(), a.ColSpan())
}

// AddArea reserves cells and returns an ID.
func (g *Grid) AddArea(a *Area) (ID, error) {
	if err := g.ValidateArea(a); err != nil {
		return Nil, fmt.Errorf("grid: invalid area: %w", err)
	}
	if !g.IsFree(a) {
		return Nil, fmt.Errorf("grid: area overlaps existing area")
	}
	id := g.allocID()
	g.areas[id] = a
	g.occ.SetRect(a.Row(), a.Col(), a.RowSpan(), a.ColSpan())
	g.indexAdd(id, a)
	return id, nil
}

func (g *Grid) Get(id ID) (*Area, bool) { a, ok := g.areas[id]; return a, ok }

func (g *Grid) allocID() ID {
	id := g.nextID
	g.nextID++
	return id
}

func (g *Grid) indexAdd(id ID, a *Area) {
	for _, r := range a.Rows() {
		g.rowIndex[r] = append(g.rowIndex[r], id)
	}
}
