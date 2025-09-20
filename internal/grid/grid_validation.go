package grid

import "fmt"

// validateBounds checks area fits into current grid. Does not check a.Col/a.Row.
func (g *Grid) validateBounds(a *Area) error {
	cols, rows := g.Cols(), g.Rows()
	if a.ColEnd() > cols || a.RowEnd() > rows {
		return fmt.Errorf(
			"grid: out of bounds colEnd=%d rowEnd=%d cols=%d rows=%d",
			a.ColEnd(), a.RowEnd(), cols, rows,
		)
	}
	return nil
}

// validateOverlap checks for any existing IDs in target cells.
func (g *Grid) validateOverlap(a *Area) error {
	for r := a.Row(); r < a.RowEnd(); r++ {
		for c := a.Col(); c < a.ColEnd(); c++ {
			if ex := g.cells[r][c]; ex != Nil {
				return fmt.Errorf("grid: overlap at row=%d col=%d existing=%d", r, c, ex)
			}
		}
	}
	return nil
}

// canShiftRight checks bounds + immediate occupancy for shifting 'a' right by 1.
func (g *Grid) canShiftRight(a *Area) error {
	if a == nil {
		return fmt.Errorf("grid: canShiftRight nil area")
	}

	destCol := a.ColEnd() // write index after shift
	if destCol >= g.Cols() {
		return fmt.Errorf("grid: cannot shift; destination column=%d exceeds last index=%d",
			destCol, g.Cols()-1)
	}

	return a.ForRowsWithError(func(r int) error {
		if occ := g.cells[r][destCol]; occ != Nil {
			return fmt.Errorf("grid: cannot shift; blocked at row=%d col=%d by id=%d",
				r, destCol, occ)
		}
		return nil
	})
}
