package grid_test

import (
	"fmt"
	"testing"

	"github.com/neox5/tbl/internal/grid"
)

func TestAddColAndAddRow(t *testing.T) {
	g := grid.New([]int{2, 3}, []int{1})
	if g.Cols() != 2 || g.TotalWidth() != 5 {
		t.Fatalf("pre got cols=%d w=%d", g.Cols(), g.TotalWidth())
	}
	g.AddCol(4)
	if g.Cols() != 3 || g.TotalWidth() != 9 {
		t.Fatalf("addcol cols=%d w=%d", g.Cols(), g.TotalWidth())
	}
	// each row extended with Nil
	for rr := range g.Rows() {
		if got := g.Get(grid.Cell{Col: g.Cols() - 1, Row: rr}); got != grid.Nil {
			t.Fatalf("new cell not Nil at row %d", rr)
		}
	}
	// add row
	g.AddRow(7)
	if g.Rows() != 2 || g.TotalHeight() != 8 {
		t.Fatalf("addrow rows=%d h=%d", g.Rows(), g.TotalHeight())
	}
	// new row filled with Nil
	for cc := range g.Cols() {
		if got := g.Get(grid.Cell{Col: cc, Row: g.Rows() - 1}); got != grid.Nil {
			t.Fatalf("new row cell not Nil at col %d", cc)
		}
	}
}

func TestAddCol_PanicsOnNegative(t *testing.T) {
	g := grid.New([]int{1}, []int{1})
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	g.AddCol(-1)
}

func TestAddRow_PanicsOnNegative(t *testing.T) {
	g := grid.New([]int{1}, []int{1})
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	g.AddRow(-1)
}

func TestShiftRightRow_Errors(t *testing.T) {
	g := grid.New([]int{1, 1, 1}, []int{1}) // 3 cols

	// occupy the last column (index 2)
	a := grid.NewArea(2, 0, 1, 1)
	r := g.RegisterArea(&a)
	if err := g.PlaceArea(r); err != nil {
		t.Fatalf("place: %v", err)
	}

	// last column occupied -> error
	if err := g.ShiftRightRow(0, 0); err == nil {
		t.Fatal("expected last column occupied error")
	}

	// row OOB
	if err := g.ShiftRightRow(5, 0); err == nil {
		t.Fatal("expected row OOB error")
	}
	// at OOB or no space to the right
	if err := g.ShiftRightRow(0, 2); err == nil {
		t.Fatal("expected at OOB error")
	}
}

// --- Arrangement test ---

// buildArrangement returns a grid with the requested layout:
//
// [ A | B | C | D D ]
// [ E | B | F | G | H ]
// [ I | B | F | J | H ]
//
// All cols width=1, rows height=1.
func buildArrangement(t *testing.T) (*grid.Grid, map[string]grid.Ref, map[string]*grid.Area) {
	t.Helper()
	g := grid.New([]int{1, 1, 1, 1, 1}, []int{1, 1, 1})

	reg := func(a *grid.Area) grid.Ref {
		return g.RegisterArea(a)
	}

	areas := map[string]*grid.Area{
		"A": ptrArea(grid.NewArea(0, 0, 1, 1)),
		"B": ptrArea(grid.NewArea(1, 0, 1, 3)), // spans all rows at col=1
		"C": ptrArea(grid.NewArea(2, 0, 1, 1)),
		"D": ptrArea(grid.NewArea(3, 0, 2, 1)), // width 2 on row 0
		"E": ptrArea(grid.NewArea(0, 1, 1, 1)),
		"F": ptrArea(grid.NewArea(2, 1, 1, 2)), // rows 1..2 at col=2
		"G": ptrArea(grid.NewArea(3, 1, 1, 1)),
		"H": ptrArea(grid.NewArea(4, 1, 1, 2)), // rows 1..2 at col=4
		"I": ptrArea(grid.NewArea(0, 2, 1, 1)),
		"J": ptrArea(grid.NewArea(3, 2, 1, 1)),
	}
	refs := map[string]grid.Ref{}
	for k, a := range areas {
		refs[k] = reg(a)
		if err := g.PlaceArea(refs[k]); err != nil {
			t.Fatalf("place %s: %v", k, err)
		}
	}
	return g, refs, areas
}

func ptrArea(a grid.Area) *grid.Area { return &a }

// snapshotCols records Area.Col() by ref.
func snapshotCols(g *grid.Grid, refs map[string]grid.Ref) map[grid.Ref]int {
	m := make(map[grid.Ref]int, len(refs))
	for _, r := range refs {
		m[r] = gCol(g, r)
	}
	return m
}

// gCol gets Area.Col() for a ref.
func gCol(g *grid.Grid, r grid.Ref) int {
	// uses unexported area via behavior; safe through methods we call
	// we infer by scanning a leftmost occurrence in any row
	// but we can rely on the fact areas are stored pointer and moved by the grid;
	// use AreaRect to find origin x from col pref is not helpful for col.
	// Workaround: find minimal column where ref appears.
	minc := -1
	for rr := range g.Rows() {
		for cc := range g.Cols() {
			if g.Get(grid.Cell{Col: cc, Row: rr}) == r {
				if minc == -1 || cc < minc {
					minc = cc
				}
			}
		}
	}
	return minc
}

// countCellsByRef counts occupied cells per ref.
func countCellsByRef(g *grid.Grid) map[grid.Ref]int {
	m := map[grid.Ref]int{}
	for rr := range g.Rows() {
		for cc := range g.Cols() {
			r := g.Get(grid.Cell{Col: cc, Row: rr})
			if r != grid.Nil {
				m[r]++
			}
		}
	}
	return m
}

func TestShiftRightRow_Arrangement_ShiftAtEachCol(t *testing.T) {
	for at := 0; at <= 4; at++ {
		t.Run(fmt.Sprintf("at_%d", at), func(t *testing.T) {
			g, refs, _ := buildArrangement(t)

			// add a free last column
			g.AddCol(1)
			if g.Cols() != 6 {
				t.Fatalf("want 6 cols got %d", g.Cols())
			}
			// one Nil per row at last column before shift
			for rr := range 3 {
				if g.Get(grid.Cell{Col: 5, Row: rr}) != grid.Nil {
					t.Fatalf("expected Nil at last column before shift, row %d", rr)
				}
			}

			beforeCols := snapshotCols(g, refs)
			beforeCounts := countCellsByRef(g)

			// compute refs touched on row 0 from at..lastSrc
			moved := map[grid.Ref]struct{}{}
			for c := 4; c >= at; c-- { // lastSrc = 4
				if r := g.Get(grid.Cell{Col: c, Row: 0}); r != grid.Nil {
					moved[r] = struct{}{}
				}
			}

			// perform shift
			if err := g.ShiftRightRow(0, at); err != nil {
				t.Fatalf("shift at %d: %v", at, err)
			}

			// invariants:
			afterCounts := countCellsByRef(g)
			if len(beforeCounts) != len(afterCounts) {
				t.Fatal("ref set changed")
			}
			for r, n := range beforeCounts {
				if afterCounts[r] != n {
					t.Fatalf("ref %v cell count changed: %d -> %d", r, n, afterCounts[r])
				}
			}
			for _, r := range refs {
				bc := beforeCols[r]
				ac := gCol(g, r)
				_, isMoved := moved[r]
				if isMoved {
					if ac != bc+1 {
						t.Fatalf("ref %v col not advanced: before %d after %d", r, bc, ac)
					}
				} else {
					if ac != bc {
						t.Fatalf("ref %v col changed unexpectedly: before %d after %d", r, bc, ac)
					}
				}
			}
			if g.Get(grid.Cell{Col: at, Row: 0}) != grid.Nil {
				t.Fatalf("row0 col %d should be Nil after shift", at)
			}
			if g.Get(grid.Cell{Col: 5, Row: 0}) == grid.Nil {
				t.Fatalf("row0 last column should be filled after shift")
			}
			for rr := range 3 {
				nilCount := 0
				for cc := range g.Cols() {
					if g.Get(grid.Cell{Col: cc, Row: rr}) == grid.Nil {
						nilCount++
					}
				}
				if nilCount != 1 {
					t.Fatalf("row %d Nil count = %d want 1", rr, nilCount)
				}
			}
		})
	}
}
