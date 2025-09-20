package grid_test

import (
	"slices"
	"testing"

	"github.com/neox5/tbl/internal/grid"
)

func TestArea_ConstructAndMutate(t *testing.T) {
	// invalid ctor
	expectPanic := func(fn func(), msg string) {
		t.Helper()
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("expected panic: %s", msg)
			}
		}()
		fn()
	}
	expectPanic(func() { grid.NewArea(-1, 0, 1, 1) }, "neg col")
	expectPanic(func() { grid.NewArea(0, 0, 0, 1) }, "nonpositive span")

	// ok ctor and accessors
	a := grid.NewArea(1, 2, 3, 4)
	if a.Col() != 1 || a.Row() != 2 || a.ColSpan() != 3 || a.RowSpan() != 4 {
		t.Fatal("accessors wrong")
	}
	if a.ColEnd() != 4 || a.RowEnd() != 6 {
		t.Fatal("ends wrong")
	}

	// MoveTo and MoveBy guards
	expectPanic(func() { aa := a; aa.MoveTo(-1, 0) }, "MoveTo neg")
	expectPanic(func() { aa := a; aa.MoveBy(-2, 0) }, "MoveBy neg")
	aa := a
	aa.MoveTo(0, 0)
	aa.MoveBy(2, 1)
	if aa.Col() != 2 || aa.Row() != 1 {
		t.Fatal("move failed")
	}

	// Resize guards and ok
	expectPanic(func() { rr := a; rr.Resize(0, 1) }, "Resize nonpositive")
	r2 := a
	r2.Resize(5, 6)
	if r2.ColSpan() != 5 || r2.RowSpan() != 6 {
		t.Fatal("resize failed")
	}
}

func TestArea_Iterators(t *testing.T) {
	a := grid.NewArea(1, 2, 2, 3) // cols 1..2, rows 2..4
	rows := []int{}
	a.ForRows(func(r int) { rows = append(rows, r) })
	if !slices.Equal(rows, []int{2, 3, 4}) {
		t.Fatalf("ForRows got %v", rows)
	}
	cells := [][2]int{}
	a.ForEachCell(func(c, r int) { cells = append(cells, [2]int{c, r}) })
	want := [][2]int{
		{1, 2},
		{2, 2},
		{1, 3},
		{2, 3},
		{1, 4},
		{2, 4},
	}
	if len(cells) != len(want) {
		t.Fatalf("ForEachCell count %d want %d", len(cells), len(want))
	}
	for i := range want {
		if cells[i] != want[i] {
			t.Fatalf("idx %d got %v want %v", i, cells[i], want[i])
		}
	}
}
