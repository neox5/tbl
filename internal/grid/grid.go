package grid

import (
	"fmt"

	"github.com/neox5/btmp"
)

// Grid with explicit sizing and simple storage.
type Grid struct {
	occ      *btmp.Grid
	areas    map[ID]*Area
	rowIndex map[int][]ID
	toMove   map[ID]struct{}
	nextID   ID
}

// New builds a grid. Zero counts allowed. Negatives panic with values.
func New(numCols, numRows int) *Grid {
	if numCols < 0 || numRows < 0 {
		panic(fmt.Errorf("grid: invalid size cols=%d rows=%d", numCols, numRows))
	}
	return &Grid{
		occ:      btmp.NewGridWithSize(numCols, numRows),
		areas:    make(map[ID]*Area),
		rowIndex: make(map[int][]ID),
		toMove:   make(map[ID]struct{}),
		nextID:   1,
	}
}

func (g *Grid) Cols() int { return g.occ.Cols() }
func (g *Grid) Rows() int { return g.occ.Rows() }

func (g *Grid) AddCol() {
	g.occ.GrowCols(1)
}

func (g *Grid) AddRow() {
	g.occ.GrowRows(1)
}

// IsFree reports whether the area is unoccupied in the grid.
func (g *Grid) IsFree(a *Area) bool {
	return g.occ.IsFree(a.Col(), a.Row(), a.ColSpan(), a.RowSpan())
}

// ValidateArea validates that area fits within grid bounds.
// Returns error if area exceeds grid dimensions.
func (g *Grid) ValidateArea(a *Area) error {
	return g.occ.ValidateRect(a.Col(), a.Row(), a.ColSpan(), a.RowSpan())
}

func (g *Grid) Print() string {
	return g.occ.Print()
}

// AddArea registers and places the area.
// Side effects order: validate -> registry -> placeArea -> indexAdd.
func (g *Grid) AddArea(a *Area) (ID, error) {
	if a == nil {
		return Nil, fmt.Errorf("grid: AddArea nil area")
	}

	if err := g.ValidateArea(a); err != nil {
		return Nil, err
	}
	if !g.IsFree(a) {
		return Nil, fmt.Errorf("grid: area overlaps existing area")
	}

	id := g.nextID
	g.nextID++
	g.areas[id] = a
	g.occ.SetRect(a.Col(), a.Row(), a.ColSpan(), a.RowSpan())
	g.indexAdd(id)

	return id, nil
}
