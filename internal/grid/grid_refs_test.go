package grid_test

import (
	"testing"

	"github.com/neox5/tbl/internal/grid"
)

func TestRegisterArea_RefMapping(t *testing.T) {
	g := grid.New([]int{1}, []int{1})
	a1 := grid.NewArea(0, 0, 1, 1)
	a2 := grid.NewArea(0, 0, 1, 1)
	r1 := g.RegisterArea(&a1)
	r2 := g.RegisterArea(&a2)
	if r1 == grid.Nil || r2 == grid.Nil || r1 == r2 {
		t.Fatalf("bad refs r1=%v r2=%v", r1, r2)
	}
}

func TestPlaceArea_PanicsOnInvalidRef(t *testing.T) {
	g := grid.New([]int{1}, []int{1})
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on invalid Ref")
		}
	}()
	_ = g.PlaceArea(grid.Ref(123)) // triggers areaByRef panic
}
