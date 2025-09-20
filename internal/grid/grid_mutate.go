package grid

import "fmt"

// AddCol appends a column of width w. Fills new cells with Nil.
func (g *Grid) AddCol(w int) {
	if w < 0 {
		panic("grid: AddCol negative width")
	}
	g.cols = append(g.cols, w)
	for row := range g.cells {
		g.cells[row] = append(g.cells[row], Nil)
	}
	g.buildPrefixes()
}

// AddRow appends a row of height h. Fills new row with Nil.
func (g *Grid) AddRow(h int) {
	if h < 0 {
		panic("grid: AddRow negative height")
	}
	g.rows = append(g.rows, h)
	g.cells = append(g.cells, make([]Ref, g.Cols()))
	g.buildPrefixes()
}

// ShiftRightRow shifts one step to the right on 'row' starting at column 'at'.
// Moves refs only. Preserves rectangles via recursion.
// Preconditions per call: last column on the initiating row must be Nil.
func (g *Grid) ShiftRightRow(row, at int) error {
	if row < 0 || row >= g.Rows() {
		return fmt.Errorf("grid: row out of range")
	}
	if at < 0 || at >= g.Cols()-1 {
		return fmt.Errorf("grid: 'at' out of range or no space to the right")
	}
	// Top-level preflight: require spare space on the initiating row only.
	if g.cells[row][g.Cols()-1] != Nil {
		return fmt.Errorf("grid: last column occupied on row %d", row)
	}

	moved := make(map[Ref]struct{})              // areas scheduled to mutate once at end
	movedOnRow := make(map[int]map[Ref]struct{}) // per-row refs already shifted

	if err := g.shiftRow(row, at, moved, movedOnRow); err != nil {
		return err
	}
	for ref := range moved {
		g.areaByRef(ref).MoveBy(1, 0)
	}
	return nil
}

func (g *Grid) shiftRow(row, at int, moved map[Ref]struct{}, movedOnRow map[int]map[Ref]struct{}) error {
	r := g.cells[row]

	lastSrc := g.Cols() - 2
	if at > lastSrc {
		return nil
	}
	if movedOnRow[row] == nil {
		movedOnRow[row] = make(map[Ref]struct{})
	}

	for col := lastSrc; col >= at; col-- {
		ref := r[col]
		if ref == Nil {
			continue
		}

		// First touch of this area across the whole operation.
		if _, seen := moved[ref]; !seen {
			a := g.areaByRef(ref)
			// Mark before recursing to break cycles across rows.
			moved[ref] = struct{}{}
			// Recurse for other rows spanned by this area at its left edge.
			for rr := a.Row(); rr < a.RowEnd(); rr++ {
				if rr == row {
					continue
				}
				if err := g.shiftRow(rr, a.Col(), moved, movedOnRow); err != nil {
					return err
				}
			}
		}

		// Avoid shifting the same ref twice on this row.
		if _, done := movedOnRow[row][ref]; done {
			continue
		}

		r[col+1] = ref
		r[col] = Nil
		movedOnRow[row][ref] = struct{}{}
	}
	return nil
}
