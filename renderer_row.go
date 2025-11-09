package tbl

// PreparedRow contains all rendering instructions for one cell row.
type PreparedRow struct {
	row        int
	borderOps  []RenderOp   // nil if no border above this row
	contentOps [][]RenderOp // one slice per content line (height)
}

// prepareRow analyzes cell row and builds rendering instructions.
func prepareRow(grid [][]*Cell, colWidths []int, row int) *PreparedRow {
	seen := make(map[ID]bool)
	var cells []*Cell

	for col := 0; col < len(grid[row]); col++ {
		c := grid[row][col]
		if c != nil && !seen[c.id] {
			seen[c.id] = true
			cells = append(cells, c)
		}
	}

	pr := &PreparedRow{row: row}

	// Build border instructions (logic determines if needed)
	pr.borderOps = buildBorderOps(grid, colWidths, row, cells)

	// Build content instructions (single line for now)
	pr.contentOps = [][]RenderOp{
		buildContentOps(grid, colWidths, row, cells),
	}

	return pr
}

// buildBorderOps constructs border instruction sequence.
// Returns nil if no border should be rendered.
func buildBorderOps(grid [][]*Cell, colWidths []int, row int, cells []*Cell) []RenderOp {
	// Top border always rendered
	isTop := row == 0

	// Check if structure differs from previous row
	needsBorder := isTop
	if !isTop {
		prev := row - 1
		for col := 0; col < len(grid[row]); col++ {
			prevCell := grid[prev][col]
			currCell := grid[row][col]

			if prevCell != currCell {
				needsBorder = true
				break
			}

			if prevCell != nil && col == prevCell.c+prevCell.cSpan-1 {
				if col+1 < len(grid[row]) {
					if grid[prev][col+1] != grid[row][col+1] {
						needsBorder = true
						break
					}
				}
			}
		}
	}

	if !needsBorder {
		return nil
	}

	// TODO: build actual instruction sequence
	return []RenderOp{CornerTL{}, HLine{10}, CornerTR{}}
}

// buildContentOps constructs content instruction sequence for one line.
func buildContentOps(grid [][]*Cell, colWidths []int, row int, cells []*Cell) []RenderOp {
	// TODO: implement
	return []RenderOp{VLine{}, Content{"TODO", 10, HAlignLeft}, VLine{}}
}
