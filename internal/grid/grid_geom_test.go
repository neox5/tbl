package grid_test

import (
	"testing"

	"github.com/neox5/tbl/internal/grid"
)

func TestGeometry_CellAndArea(t *testing.T) {
	g := grid.New([]int{2, 3}, []int{1, 4})

	// cell origin
	x, y := g.CellOrigin(grid.Cell{Col: 1, Row: 1})
	if x != 2 || y != 1 {
		t.Fatalf("cell origin got (%d,%d)", x, y)
	}

	// area rect
	a := grid.NewArea(1, 0, 1, 2) // uses col[1]=3, rows[0..1]=1+4
	ax, ay, aw, ah := g.AreaRect(a)
	if ax != 2 || ay != 0 || aw != 3 || ah != 5 {
		t.Fatalf("area rect got (%d,%d,%d,%d)", ax, ay, aw, ah)
	}
}
