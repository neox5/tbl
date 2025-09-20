package grid_test

import (
	"testing"

	"github.com/neox5/tbl/internal/grid"
)

func TestNew_BasicsAndPrefixes(t *testing.T) {
	g := grid.New([]int{2, 3}, []int{1, 4})
	if g.Cols() != 2 || g.Rows() != 2 {
		t.Fatalf("dims got %d×%d", g.Cols(), g.Rows())
	}
	if g.TotalWidth() != 5 || g.TotalHeight() != 5 {
		t.Fatalf("totals got %d×%d", g.TotalWidth(), g.TotalHeight())
	}
	// prefix sums
	if g.ColOffset(0) != 0 || g.ColOffset(1) != 2 || g.ColOffset(2) != 5 {
		t.Fatalf("col offsets wrong: %d,%d,%d", g.ColOffset(0), g.ColOffset(1), g.ColOffset(2))
	}
	if g.RowOffset(0) != 0 || g.RowOffset(1) != 1 || g.RowOffset(2) != 5 {
		t.Fatalf("row offsets wrong: %d,%d,%d", g.RowOffset(0), g.RowOffset(1), g.RowOffset(2))
	}
}

func TestNew_PanicsOnEmptyOrNegative(t *testing.T) {
	expectPanic := func(fn func(), name string) {
		t.Helper()
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("%s: expected panic", name)
			}
		}()
		fn()
	}
	expectPanic(func() { grid.New(nil, []int{1}) }, "no cols")
	expectPanic(func() { grid.New([]int{1}, nil) }, "no rows")
	expectPanic(func() { grid.New([]int{-1}, []int{1}) }, "neg col")
	expectPanic(func() { grid.New([]int{1}, []int{-1}) }, "neg row")
}
