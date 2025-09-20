package grid

import "fmt"

// Get returns the reference at cell c. Panics if out of range.
func (g *Grid) Get(c Cell) Ref {
	return g.cells[c.Row][c.Col]
}

// Set stores ref at cell c. Panics if out of range.
func (g *Grid) Set(c Cell, ref Ref) {
	g.cells[c.Row][c.Col] = ref
}

// Clear sets cell c to Nil.
func (g *Grid) Clear(c Cell) {
	g.cells[c.Row][c.Col] = Nil
}

// PlaceArea writes the area's ref into each covered cell. Fails on overlap or OOB.
func (g *Grid) PlaceArea(ref Ref) error {
	a := g.areaByRef(ref)
	// bounds
	if a.Col() < 0 || a.Row() < 0 || a.ColEnd() > g.Cols() || a.RowEnd() > g.Rows() {
		return fmt.Errorf("grid: area out of bounds")
	}
	// overlap check
	var overlap error
	a.ForEachCell(func(c, r int) {
		if overlap != nil {
			return
		}
		if g.cells[r][c] != Nil {
			overlap = fmt.Errorf("grid: overlap at row=%d col=%d", r, c)
		}
	})
	if overlap != nil {
		return overlap
	}
	// write
	a.ForEachCell(func(c, r int) {
		g.cells[r][c] = ref
	})
	return nil
}

// ClearArea sets the area's cells to Nil.
func (g *Grid) ClearArea(ref Ref) {
	a := g.areaByRef(ref)
	a.ForEachCell(func(c, r int) {
		g.cells[r][c] = Nil
	})
}
