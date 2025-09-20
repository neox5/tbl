package grid

import "fmt"

// ShiftRightRow shifts cells on a row one step to the right starting at column 'at'.
// Uses recursive propagation so multi-row areas remain rectangular.
// Side effects: moves cells, then for each touched area moves its descriptor by (+1,0).
// Index order stays valid because relative order is preserved.
func (g *Grid) ShiftRightRow(row, at int) error {
	// bounds
	if row < 0 || row >= g.Rows() {
		return fmt.Errorf("grid: row out of range row=%d rows=%d", row, g.Rows())
	}
	if at < 0 || at >= g.Cols()-1 {
		return fmt.Errorf("grid: at out of range at=%d cols=%d", at, g.Cols())
	}

	seen := make(map[ID]struct{})              // cycle guard, and final set to MoveBy(+1,0)
	seenOnRow := make(map[int]map[ID]struct{}) // per-row guard to shift each ID once per row

	if err := g.shiftRow(row, at, seen, seenOnRow); err != nil {
		return err
	}

	// Apply one logical descriptor move per touched area.
	for id := range seen {
		a := g.areas[id]
		if a == nil {
			return fmt.Errorf("grid: moved unknown id=%d", id)
		}
		a.MoveBy(1, 0)
	}
	return nil
}

// shiftRow performs the in-row cell shifts and propagates to other rows spanned by areas.
// It shifts cells from right to left to avoid overwrite.
func (g *Grid) shiftRow(row, at int, seen map[ID]struct{}, seenOnRow map[int]map[ID]struct{}) error {
	if seenOnRow[row] == nil {
		seenOnRow[row] = make(map[ID]struct{})
	}
	lastSrc := g.Cols() - 2
	if at > lastSrc {
		return nil
	}

	for col := lastSrc; col >= at; {
		id := g.cells[row][col]
		if id == Nil || hasID(seenOnRow[row], id) {
			col--
			continue
		}

		a := g.areas[id]
		if a == nil {
			return fmt.Errorf("grid: shiftRow unknown id=%d", id)
		}

		// Mark globally before recursing to avoid cycles.
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
		}

		// 1) shift other rows spanned by this area starting at its right edge
		shiftOtherRows := func(r int) error {
			if r == row {
				return nil
			}
			return g.shiftRow(r, a.ColEnd(), seen, seenOnRow)
		}
		if err := a.ForRowsWithError(shiftOtherRows); err != nil {
			return err
		}

		// 2) validate: this area must be able to move now
		if err := g.canShiftRight(a); err != nil {
			return err
		}

		// 3) apply shift on all its rows
		applyShift := func(r int) error {
			g.cells[r][a.ColEnd()] = id // destination
			g.cells[r][a.Col()] = Nil   // source vacated
			if seenOnRow[r] == nil {
				seenOnRow[r] = make(map[ID]struct{})
			}
			seenOnRow[r][id] = struct{}{}
			return nil
		}
		if err := a.ForRowsWithError(applyShift); err != nil {
			return err
		}

		col--
	}
	return nil
}

func hasID(s map[ID]struct{}, id ID) bool {
	_, ok := s[id]
	return ok
}
