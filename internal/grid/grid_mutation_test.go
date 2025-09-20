package grid_test

import (
	"testing"

	grid "github.com/neox5/tbl/internal/grid"
)

// --- helpers ---

func mustAdd(t *testing.T, g *grid.Grid, a *grid.Area) grid.ID {
	t.Helper()
	id, err := g.AddArea(a)
	if err != nil {
		t.Fatalf("AddArea failed: %v", err)
	}
	return id
}

func cellsEqual(t *testing.T, got, want [][]grid.ID) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("rows mismatch: got %d want %d", len(got), len(want))
	}
	for r := range want {
		if len(got[r]) != len(want[r]) {
			t.Fatalf("cols mismatch at row %d: got %d want %d", r, len(got[r]), len(want[r]))
		}
		for c := range want[r] {
			if got[r][c] != want[r][c] {
				t.Fatalf("cell[%d][%d]: got %d want %d", r, c, got[r][c], want[r][c])
			}
		}
	}
}

// --- column insertion workflow tests ---

func TestColumnInsertion_SingleRow_AtStart(t *testing.T) {
	// [A|B|C] -> AddCol -> [A|B|C|Nil] -> Shift(0,0) -> [Nil|A|B|C]
	g := grid.New(3, 1)
	idA := mustAdd(t, g, grid.NewArea(0, 0, 1, 1))
	idB := mustAdd(t, g, grid.NewArea(1, 0, 1, 1))
	idC := mustAdd(t, g, grid.NewArea(2, 0, 1, 1))

	g.AddCol()
	if err := g.ShiftRightRow(0, 0); err != nil {
		t.Fatalf("ShiftRightRow failed: %v", err)
	}

	want := [][]grid.ID{
		{grid.Nil, idA, idB, idC},
	}
	cellsEqual(t, g.Cells(), want)
}

func TestColumnInsertion_SingleRow_AtMiddle(t *testing.T) {
	// [A|B|C] -> AddCol -> [A|B|C|Nil] -> Shift(0,1) -> [A|Nil|B|C]
	g := grid.New(3, 1)
	idA := mustAdd(t, g, grid.NewArea(0, 0, 1, 1))
	idB := mustAdd(t, g, grid.NewArea(1, 0, 1, 1))
	idC := mustAdd(t, g, grid.NewArea(2, 0, 1, 1))

	g.AddCol()
	if err := g.ShiftRightRow(0, 1); err != nil {
		t.Fatalf("ShiftRightRow failed: %v", err)
	}

	want := [][]grid.ID{
		{idA, grid.Nil, idB, idC},
	}
	cellsEqual(t, g.Cells(), want)
}

func TestColumnInsertion_SingleRow_AtEnd(t *testing.T) {
	// [A|B|C] -> AddCol -> [A|B|C|Nil] -> Shift(0,2) -> [A|B|Nil|C]
	g := grid.New(3, 1)
	idA := mustAdd(t, g, grid.NewArea(0, 0, 1, 1))
	idB := mustAdd(t, g, grid.NewArea(1, 0, 1, 1))
	idC := mustAdd(t, g, grid.NewArea(2, 0, 1, 1))

	g.AddCol()
	if err := g.ShiftRightRow(0, 2); err != nil {
		t.Fatalf("ShiftRightRow failed: %v", err)
	}

	want := [][]grid.ID{
		{idA, idB, grid.Nil, idC},
	}
	cellsEqual(t, g.Cells(), want)
}

// --- multi-row propagation tests ---

func TestColumnInsertion_MultiRowSpan_Propagation(t *testing.T) {
	// [A|B|C] where A spans 2 rows
	// [A|D|E]
	// Insert at start -> shift A triggers shift of D,E on row 1
	g := grid.New(3, 2)
	idA := mustAdd(t, g, grid.NewArea(0, 0, 1, 2)) // spans rows 0-1
	idB := mustAdd(t, g, grid.NewArea(1, 0, 1, 1))
	idC := mustAdd(t, g, grid.NewArea(2, 0, 1, 1))
	idD := mustAdd(t, g, grid.NewArea(1, 1, 1, 1))
	idE := mustAdd(t, g, grid.NewArea(2, 1, 1, 1))

	g.AddCol()
	if err := g.ShiftRightRow(0, 0); err != nil {
		t.Fatalf("ShiftRightRow failed: %v", err)
	}

	want := [][]grid.ID{
		{grid.Nil, idA, idB, idC},
		{grid.Nil, idA, idD, idE},
	}
	cellsEqual(t, g.Cells(), want)
}

func TestColumnInsertion_MultiSpanArea_RightEdgePropagation(t *testing.T) {
	// [A|B B|C] where B spans 2 columns
	// [D|E E|F]
	// Insert at position 1 -> only affects B's right edge at position 3
	g := grid.New(4, 2)
	idA := mustAdd(t, g, grid.NewArea(0, 0, 1, 1))
	idB := mustAdd(t, g, grid.NewArea(1, 0, 2, 1)) // spans cols 1-2
	idC := mustAdd(t, g, grid.NewArea(3, 0, 1, 1))
	idD := mustAdd(t, g, grid.NewArea(0, 1, 1, 1))
	idE := mustAdd(t, g, grid.NewArea(1, 1, 2, 1)) // spans cols 1-2
	idF := mustAdd(t, g, grid.NewArea(3, 1, 1, 1))

	g.AddCol()
	if err := g.ShiftRightRow(0, 1); err != nil {
		t.Fatalf("ShiftRightRow failed: %v", err)
	}

	want := [][]grid.ID{
		{idA, grid.Nil, idB, idB, idC},
		{idD, grid.Nil, idE, idE, idF},
	}
	cellsEqual(t, g.Cells(), want)
}

// --- complex cascade tests ---

func TestColumnInsertion_ComplexCascade(t *testing.T) {
	// Layout:
	// [A|B|C] where A spans rows 0-1, C spans rows 0-2
	// [A|D|C]
	// [E|F|C]
	// Insert at start -> A moves, triggers cascades
	g := grid.New(3, 3)
	idA := mustAdd(t, g, grid.NewArea(0, 0, 1, 2)) // spans rows 0-1
	idB := mustAdd(t, g, grid.NewArea(1, 0, 1, 1))
	idC := mustAdd(t, g, grid.NewArea(2, 0, 1, 3)) // spans rows 0-2
	idD := mustAdd(t, g, grid.NewArea(1, 1, 1, 1))
	idE := mustAdd(t, g, grid.NewArea(0, 2, 1, 1))
	idF := mustAdd(t, g, grid.NewArea(1, 2, 1, 1))

	g.AddCol()
	if err := g.ShiftRightRow(0, 0); err != nil {
		t.Fatalf("ShiftRightRow failed: %v", err)
	}

	want := [][]grid.ID{
		{grid.Nil, idA, idB, idC},
		{grid.Nil, idA, idD, idC},
		{grid.Nil, idE, idF, idC},
	}
	cellsEqual(t, g.Cells(), want)
}

// --- edge cases ---

func TestColumnInsertion_EmptyRow(t *testing.T) {
	// Row 0 empty, row 1 has content
	g := grid.New(3, 2)
	mustAdd(t, g, grid.NewArea(0, 1, 3, 1)) // only in row 1

	g.AddCol()
	if err := g.ShiftRightRow(0, 0); err != nil {
		t.Fatalf("ShiftRightRow on empty row failed: %v", err)
	}

	// Row 0 should remain empty, row 1 unchanged
	want := [][]grid.ID{
		{grid.Nil, grid.Nil, grid.Nil, grid.Nil},
		{g.Cells()[1][0], g.Cells()[1][1], g.Cells()[1][2], grid.Nil},
	}
	cellsEqual(t, g.Cells(), want)
}

func TestColumnInsertion_NoShiftNeeded(t *testing.T) {
	// Insert beyond last occupied column
	g := grid.New(3, 1)
	idA := mustAdd(t, g, grid.NewArea(0, 0, 1, 1))
	// columns 1,2 empty

	g.AddCol()
	if err := g.ShiftRightRow(0, 3); err != nil {
		t.Fatalf("ShiftRightRow beyond content failed: %v", err)
	}

	// Should be no-op since at > lastSrc (which is 2)
	want := [][]grid.ID{
		{idA, grid.Nil, grid.Nil, grid.Nil},
	}
	cellsEqual(t, g.Cells(), want)
}

// --- error cases ---

func TestColumnInsertion_InvalidParameters(t *testing.T) {
	g := grid.New(3, 2)
	mustAdd(t, g, grid.NewArea(0, 0, 1, 1))
	g.AddCol() // now cols = 4

	tests := []struct {
		name string
		row  int
		at   int
	}{
		{"negative row", -1, 0},
		{"row out of bounds", 2, 0},
		{"negative at", 0, -1},
		{"at too large", 0, 3}, // at must be < cols-1 = 3
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := g.ShiftRightRow(tt.row, tt.at)
			if err == nil {
				t.Fatal("expected error for invalid parameters")
			}
		})
	}
}

func TestColumnInsertion_BlockedShift(t *testing.T) {
	// Create scenario where shift is blocked
	g := grid.New(2, 1)
	idA := mustAdd(t, g, grid.NewArea(0, 0, 1, 1))
	idB := mustAdd(t, g, grid.NewArea(1, 0, 1, 1))

	// Don't add column - no room to shift
	err := g.ShiftRightRow(0, 0)
	if err == nil {
		t.Fatal("expected error when no room to shift")
	}

	// Verify grid unchanged
	want := [][]grid.ID{
		{idA, idB},
	}
	cellsEqual(t, g.Cells(), want)
}

// --- workflow validation ---

func TestColumnInsertion_WithoutAddCol_Fails(t *testing.T) {
	// Demonstrate that ShiftRightRow needs AddCol first
	g := grid.New(3, 1)
	mustAdd(t, g, grid.NewArea(0, 0, 3, 1)) // fill entire row

	// Try to shift without adding column first
	err := g.ShiftRightRow(0, 0)
	if err == nil {
		t.Fatal("expected error when shifting without AddCol")
	}
}

func TestColumnInsertion_MultipleOperations(t *testing.T) {
	// Test multiple insert operations
	g := grid.New(2, 1)
	idA := mustAdd(t, g, grid.NewArea(0, 0, 1, 1))
	idB := mustAdd(t, g, grid.NewArea(1, 0, 1, 1))

	// First insertion at start
	g.AddCol()
	if err := g.ShiftRightRow(0, 0); err != nil {
		t.Fatalf("first insertion failed: %v", err)
	}

	// Verify intermediate state
	want1 := [][]grid.ID{
		{grid.Nil, idA, idB},
	}
	cellsEqual(t, g.Cells(), want1)

	// Second insertion at position 2
	g.AddCol()
	if err := g.ShiftRightRow(0, 2); err != nil {
		t.Fatalf("second insertion failed: %v", err)
	}

	// Final state
	want2 := [][]grid.ID{
		{grid.Nil, idA, grid.Nil, idB},
	}
	cellsEqual(t, g.Cells(), want2)
}
