package grid_test

import (
	"testing"

	"github.com/neox5/tbl/internal/grid"
)

func mkGrid() *grid.Grid { return grid.New([]int{2, 2, 2}, []int{1, 1, 1}) }

func TestGetSetClear(t *testing.T) {
	g := mkGrid()
	c := grid.Cell{Col: 1, Row: 2}
	if g.Get(c) != grid.Nil {
		t.Fatal("expected Nil at start")
	}
	g.Set(c, grid.Ref(7))
	if g.Get(c) != grid.Ref(7) {
		t.Fatal("set/get mismatch")
	}
	g.Clear(c)
	if g.Get(c) != grid.Nil {
		t.Fatal("clear failed")
	}
}

func TestPlaceArea_OkOverlapAndOOB(t *testing.T) {
	g := mkGrid()

	// ok place
	a := grid.NewArea(0, 0, 2, 2) // covers (0..1,0..1)
	r := g.RegisterArea(&a)
	if err := g.PlaceArea(r); err != nil {
		t.Fatalf("place ok: %v", err)
	}
	// verify
	for rr := range 2 {
		for cc := range 2 {
			if g.Get(grid.Cell{Col: cc, Row: rr}) != r {
				t.Fatalf("cell not set at (%d,%d)", cc, rr)
			}
		}
	}

	// overlap
	b := grid.NewArea(1, 1, 2, 2) // overlaps a
	rb := g.RegisterArea(&b)
	if err := g.PlaceArea(rb); err == nil {
		t.Fatal("expected overlap error")
	}

	// out of bounds
	c := grid.NewArea(2, 2, 2, 1) // col end=4 > 3
	rc := g.RegisterArea(&c)
	if err := g.PlaceArea(rc); err == nil {
		t.Fatal("expected OOB error")
	}

	// clear area
	g.ClearArea(r)
	for rr := range 2 {
		for cc := range 2 {
			if g.Get(grid.Cell{Col: cc, Row: rr}) != grid.Nil {
				t.Fatalf("clear area failed at (%d,%d)", cc, rr)
			}
		}
	}
}
