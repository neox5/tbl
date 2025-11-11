package tbl

import (
	"fmt"
	"strings"
)

// PreparedRow contains all rendering instructions for one cell row.
type PreparedRow struct {
	row          int
	borderNeeded bool         // borders needed when cell alignment differs between rows
	borderOps    []RenderOp   // nil if no border above this row
	contentOps   [][]RenderOp // one slice per content line (height)
}

// prepareRow analyzes cell row and builds rendering instructions.
func prepareRow(
	grid [][]*Cell,
	colWidths []int,
	cellLayouts map[ID][]string,
	rowHeights []int,
	row int,
	resolveStyle func(*Cell) CellStyle,
) *PreparedRow {
	pr := &PreparedRow{row: row}

	// Build border instructions
	pr.borderNeeded, pr.borderOps = buildBorderOps(grid, colWidths, row)

	// Build content instructions
	pr.contentOps = buildContentOps(grid, colWidths, cellLayouts, rowHeights, row, resolveStyle)

	return pr
}

// buildBorderOps constructs border instruction sequence between two rows.
// Returns borderNeeded flag and instruction sequence.
//
// Parameters:
//   - grid: dense cell grid
//   - colWidths: column widths (no padding)
//   - row: target row index (can be len(grid) for bottom border)
//
// Border types:
//   - Top border: row == 0 (compares against empty row above)
//   - Mid border: 0 < row < len(grid) (compares row-1 vs row)
//   - Bottom border: row == len(grid) (compares last row vs empty row below)
func buildBorderOps(grid [][]*Cell, colWidths []int, row int) (bool, []RenderOp) {
	if len(grid) == 0 || len(colWidths) == 0 {
		return false, nil
	}

	cols := len(colWidths)

	var borderNeeded bool
	if row == 0 || row == len(grid) {
		borderNeeded = true
	}

	ops := make([]RenderOp, 0, cols*2+1)

	// Column segments: left junction + top edge
	for col := range cols {
		// Check if boundary changed from row above
		if !borderNeeded && row > 0 && colBoundaryAt(grid, row, col) != colBoundaryAt(grid, row-1, col) {
			borderNeeded = true
		}

		// Left junction
		leftJunction := selectJunction(grid, row, col)
		ops = append(ops, leftJunction)

		// Top edge (HLine or Space)
		// Width includes content padding (2 chars from writeAlignedContent)
		var topEdge RenderOp
		if rowBoundaryAt(grid, row, col) {
			topEdge = HLine{Width: colWidths[col]}
		} else {
			topEdge = Space{Width: colWidths[col]}
		}
		ops = append(ops, topEdge)
	}

	// Right junction
	rightJunction := selectJunction(grid, row, cols)
	ops = append(ops, rightJunction)

	return borderNeeded, ops
}

// rowBoundaryAt reports whether there is a horizontal border at (row,col).
// A row boundary exists when cells differ between row-1 and row.
func rowBoundaryAt(grid [][]*Cell, row, col int) bool {
	// left and right outside table
	if col < 0 || col >= len(grid[0]) {
		return false
	}
	// first or row below table
	if row == 0 || row == len(grid) {
		return true
	}
	return grid[row-1][col] != grid[row][col]
}

// colBoundaryAt reports whether there is a vertical border at (row,col).
// A col boundary exists when cells differ between col-1 and col.
func colBoundaryAt(grid [][]*Cell, row, col int) bool {
	// rows outside table or column at -1
	if row < 0 || row >= len(grid) || col < 0 {
		return false
	}
	// left and right edge of the table
	if col == 0 || col == len(grid[0]) {
		return true
	}
	return grid[row][col-1] != grid[row][col]
}

// selectJunction determines the junction character at a column position.
// A junction sits at the intersection of 4 edges where horizontal and vertical
// borders meet. The function examines cell boundaries in all four directions
// (top, right, bottom, left) to determine the appropriate junction type.
//
// Parameters:
//   - grid: dense cell grid where each position points to its owning cell
//   - row: border line position (0 = top border, len(grid) = bottom border)
//   - col: column position (0 = left edge, len(grid[0]) = right edge)
//
// The junction logic checks four boundaries meeting at position (row, col):
//
//	+---+---+
//	|   |   |
//	+---o---+  <- junction 'o' at (row, col)
//	|   | X |
//	+---+---+
//
// Boundaries checked:
//   - top: vertical boundary above junction (between row-1 columns)
//   - right: horizontal boundary right of junction (between row cells)
//   - bottom: vertical boundary below junction (between row columns)
//   - left: horizontal boundary left of junction (between row cells)
//
// Returns:
//   - Corner characters (┌ ┐ └ ┘) for grid edges
//   - Junction characters (┬ ┴ ├ ┤ ┼) for T and cross intersections
//   - HLine for horizontal continuation
//   - VLine for vertical continuation
//   - Space when no border exists at junction
func selectJunction(grid [][]*Cell, row, col int) RenderOp {
	rows := len(grid)
	cols := len(grid[0])

	// Corners
	if row == 0 && col == 0 {
		return CornerTL{} // ┌
	}
	if row == 0 && col == cols {
		return CornerTR{} // ┐
	}
	if row == rows && col == 0 {
		return CornerBL{} // └
	}
	if row == rows && col == cols {
		return CornerBR{} // ┘
	}

	// Calculate all 4 boundaries meeting at junction
	top := colBoundaryAt(grid, row-1, col)
	right := rowBoundaryAt(grid, row, col)
	bottom := colBoundaryAt(grid, row, col)
	left := rowBoundaryAt(grid, row, col-1)

	// HLine (horizontal continuation)
	if !top && right && !bottom && left {
		return HLine{Width: 1}
	}

	// VLine (vertical continuation)
	if top && !right && bottom && !left {
		return VLine{}
	}

	// T-junctions
	if !top && right && bottom && left {
		return CornerT{} // ┬
	}
	if top && right && !bottom && left {
		return CornerB{} // ┴
	}
	if top && right && bottom && !left {
		return CornerL{} // ├
	}
	if top && !right && bottom && left {
		return CornerR{} // ┤
	}

	// Cross junction
	if top && right && bottom && left {
		return CornerX{} // ┼
	}

	// No border at this junction
	return Space{Width: 1}
}

// buildContentOps constructs content ops for one grid row.
// Applies padding to content lines and distributes across rowHeight lines.
func buildContentOps(
	grid [][]*Cell,
	colWidths []int,
	cellLayouts map[ID][]string,
	rowHeights []int,
	row int,
	resolveStyle func(*Cell) CellStyle,
) [][]RenderOp {
	height := rowHeights[row]
	result := make([][]RenderOp, height)

	for lineIdx := range height {
		var ops []RenderOp

		for col := 0; col < len(colWidths); {
			cell := grid[row][col]
			if cell == nil {
				panic(fmt.Sprintf("tbl: nil cell at (%d,%d)", row, col))
			}

			// Only at cell start column
			if cell.c == col {
				ops = append(ops, VLine{})

				style := resolveStyle(cell)
				contentLines := cellLayouts[cell.id]

				// Calculate which line of cell to render
				// For multi-row cells: need offset based on previous rows
				var cellLineIdx int
				if cell.r == row {
					// Cell starts this row
					cellLineIdx = lineIdx
				} else {
					// Cell started in previous row(s)
					// Sum up lines from previous rows
					prevLines := 0
					for r := cell.r; r < row; r++ {
						prevLines += rowHeights[r]
					}
					cellLineIdx = prevLines + lineIdx
				}

				// Determine text for this line
				var text string
				totalLines := style.Padding.Top + len(contentLines) + style.Padding.Bottom

				if cellLineIdx < style.Padding.Top {
					// Top padding line
					totalWidth := cellWidth(colWidths, cell)
					text = strings.Repeat(" ", totalWidth)
				} else if cellLineIdx < style.Padding.Top+len(contentLines) {
					// Content line with horizontal padding
					contentIdx := cellLineIdx - style.Padding.Top
					contentLine := contentLines[contentIdx]

					leftPad := strings.Repeat(" ", style.Padding.Left)
					rightPad := strings.Repeat(" ", style.Padding.Right)
					text = leftPad + contentLine + rightPad
				} else if cellLineIdx < totalLines {
					// Bottom padding line
					totalWidth := cellWidth(colWidths, cell)
					text = strings.Repeat(" ", totalWidth)
				} else {
					// Beyond cell height (row taller than this cell)
					totalWidth := cellWidth(colWidths, cell)
					text = strings.Repeat(" ", totalWidth)
				}

				ops = append(ops, Content{Text: text})

				col += cell.cSpan
			} else {
				col++
			}
		}

		ops = append(ops, VLine{})
		result[lineIdx] = ops
	}

	return result
}
