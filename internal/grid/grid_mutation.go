package grid

import "fmt"

// ShiftRightRow shifts cells on a row one step to the right starting at column 'at'.
// Uses recursive propagation so multi-row areas remain rectangular.
// Side effects: moves occupancy, then for each touched area moves its descriptor by (+1,0).
// Index order stays valid because relative order is preserved.
func (g *Grid) ShiftRightRow(row, at int) error {
	// bounds
	if row < 0 || row >= g.Rows() {
		return fmt.Errorf("grid: row out of range row=%d rows=%d", row, g.Rows())
	}
	if at < 0 || at >= g.Cols()-1 {
		return fmt.Errorf("grid: at out of range at=%d cols=%d", at, g.Cols())
	}

	// Clear toMove for this operation
	clear(g.toMove)

	if err := g.shiftRow(row, at); err != nil {
		return err
	}

	// Apply one logical descriptor move per touched area.
	for id := range g.toMove {
		a := g.areas[id]
		if a == nil {
			return fmt.Errorf("grid: moved unknown id=%d", id)
		}
		a.MoveBy(1, 0)
	}
	return nil
}

// shiftRow performs the in-row shifts and propagates to other rows spanned by areas.
func (g *Grid) shiftRow(row, at int) error {
	// Get IDs in this row, iterate right to left from at
	rowIDs := g.rowIndex[row]

	for i := len(rowIDs) - 1; i >= 0; i-- {
		id := rowIDs[i]
		a := g.areas[id]
		if a == nil {
			return fmt.Errorf("grid: shiftRow unknown id=%d", id)
		}

		// Skip if already processed or doesn't overlap shift region
		if _, seen := g.toMove[id]; seen {
			continue
		}
		if a.ColEnd() <= at {
			continue
		}

		// Mark as seen
		g.toMove[id] = struct{}{}

		// 1) Recursively shift other rows spanned by this area
		for _, r := range a.Rows() {
			if r == row {
				continue
			}
			if err := g.shiftRow(r, a.ColEnd()); err != nil {
				return err
			}
		}

		// 2) Validate: this area must be able to move now
		if !g.occ.CanShiftRight(a.Col(), a.Row(), a.ColSpan(), a.RowSpan()) {
			return fmt.Errorf("grid: cannot shift right")
		}

		// 3) Shift occupancy bitmap
		g.occ.ShiftRectRight(a.Col(), a.Row(), a.ColSpan(), a.RowSpan())
	}

	return nil
}
