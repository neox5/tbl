package tbl

import "strings"

// buildRenderLines constructs all rendering lines for the table.
// Returns complete line sequence as [][]RenderOp.
func (r *renderer) buildRenderLines() [][]RenderOp {
	var lines [][]RenderOp

	for row := range len(r.rowHeights) {
		rowLines := r.buildLinesForRow(row)
		lines = append(lines, rowLines...)
	}

	return lines
}

// buildLinesForRow constructs all lines for a single row.
// Processes each cell once and populates all lines simultaneously.
func (r *renderer) buildLinesForRow(row int) [][]RenderOp {
	lineCount := r.rowHeights[row]
	hasTopBorder := r.rowBorders[row]

	lines := make([][]RenderOp, lineCount)
	for i := range lines {
		lines[i] = make([]RenderOp, 0)
	}

	// Process each cell in the row
	for col := 0; col < len(r.colWidths); {
		cell := r.grid[row][col]

		if cell != nil && cell.c == col {
			cellOps := r.buildCellSegments(cell, row, lineCount, hasTopBorder)

			// Append cell ops to each line
			for lineIdx := range lineCount {
				lines[lineIdx] = append(lines[lineIdx], cellOps[lineIdx]...)
			}

			col += cell.cSpan
		} else {
			col++
		}
	}

	// Add right edge to all lines
	for lineIdx := range lineCount {
		if r.colBorders[len(r.colWidths)] {
			if lineIdx == 0 && hasTopBorder {
				junction := r.selectJunction(row, len(r.colWidths))
				lines[lineIdx] = append(lines[lineIdx], junction)
			} else {
				if r.leftBorderAt(row, len(r.colWidths)) {
					lines[lineIdx] = append(lines[lineIdx], VLine{})
				} else {
					lines[lineIdx] = append(lines[lineIdx], Space{Width: 1})
				}
			}
		}
	}

	return lines
}

// buildCellSegments creates RenderOps for cell across all lines in row.
// Returns slice where each element is RenderOps for one line.
func (r *renderer) buildCellSegments(cell *Cell, row, lineCount int, hasTopBorder bool) [][]RenderOp {
	style := r.styles[cell.id]
	col := cell.c

	hasLeftBorder := r.colBorders[col]
	isCellOriginRow := (cell.r == row)

	// Calculate line offset for multi-row cells
	lineOffset := 0
	if cell.r < row {
		for rr := cell.r; rr < row; rr++ {
			lineOffset += r.rowHeights[rr]
		}
	}

	// Get complete cell layout (with padding already applied)
	cellLines := r.cellLayouts[cell.id]

	segments := make([][]RenderOp, lineCount)

	for lineInRow := range lineCount {
		var ops []RenderOp

		// Left border/junction
		if hasLeftBorder {
			if lineInRow == 0 && hasTopBorder {
				junction := r.selectJunction(row, col)
				ops = append(ops, junction)
			} else {
				if r.leftBorderAt(row, col) {
					ops = append(ops, VLine{})
				} else {
					ops = append(ops, Space{Width: 1})
				}
			}
		}

		// Cell content or top border
		if lineInRow == 0 && hasTopBorder && isCellOriginRow {
			// Top border
			totalWidth := r.cellWidth(cell)
			if style.Border.IsVisual(BorderTop) {
				ops = append(ops, HLine{Width: totalWidth})
			} else {
				ops = append(ops, Space{Width: totalWidth})
			}
		} else {
			// Cell content - directly from cellLayouts
			cellLineIdx := lineOffset + lineInRow

			// Adjust if top border consumed line 0
			if hasTopBorder && isCellOriginRow && lineInRow > 0 {
				cellLineIdx = lineOffset + lineInRow - 1
			}

			// Get text from pre-computed layout
			var text string
			if cellLineIdx < len(cellLines) {
				text = cellLines[cellLineIdx]
			} else {
				// Beyond cell height
				totalWidth := r.cellWidth(cell)
				text = strings.Repeat(" ", totalWidth)
			}

			ops = append(ops, Content{Text: text})
		}

		segments[lineInRow] = ops
	}

	return segments
}

// topBorderAt reports whether horizontal border renders as character at position (row, col).
// Checks cells that meet at this position.
func (r *renderer) topBorderAt(row, col int) bool {
	rows := len(r.grid)
	cols := len(r.grid[0])

	// Out of bounds
	if row < 0 || row > rows || col < 0 || col >= cols {
		return false
	}

	// Top edge of grid
	if row == 0 {
		cell := r.grid[0][col]
		if cell != nil {
			style := r.styles[cell.id]
			return style.Border.IsVisual(BorderTop)
		}
		return false
	}

	// Bottom edge of grid
	if row == rows {
		cell := r.grid[rows-1][col]
		if cell != nil {
			style := r.styles[cell.id]
			return style.Border.IsVisual(BorderBottom)
		}
		return false
	}

	// Between rows: check both cells
	cellAbove := r.grid[row-1][col]
	cellBelow := r.grid[row][col]

	if r.styles[cellAbove.id].Border.IsVisual(BorderBottom) || r.styles[cellBelow.id].Border.IsVisual(BorderTop) {
		return true
	}

	return false
}

// leftBorderAt reports whether vertical border renders as character at position (row, col).
// Checks cells that meet at this position.
func (r *renderer) leftBorderAt(row, col int) bool {
	rows := len(r.grid)
	cols := len(r.grid[0])

	// Out of bounds
	if row < 0 || row >= rows || col < 0 || col > cols {
		return false
	}

	// Left edge of grid
	if col == 0 {
		cell := r.grid[row][0]
		if cell != nil {
			style := r.styles[cell.id]
			return style.Border.IsVisual(BorderLeft)
		}
		return false
	}

	// Right edge of grid
	if col == cols {
		cell := r.grid[row][cols-1]
		if cell != nil {
			style := r.styles[cell.id]
			return style.Border.IsVisual(BorderRight)
		}
		return false
	}

	// Between columns: check both cells
	cellLeft := r.grid[row][col-1]
	cellRight := r.grid[row][col]

	if r.styles[cellLeft.id].Border.IsVisual(BorderRight) || r.styles[cellRight.id].Border.IsVisual(BorderLeft) {
		return true
	}

	return false
}

// selectJunction determines the junction character at a column position.
// Examines borders in all four directions (top, right, bottom, left).
//
// Parameters:
//   - grid: dense cell grid
//   - row: border line position (0 = top border, len(grid) = bottom border)
//   - col: column position (0 = left edge, len(grid[0]) = right edge)
//   - resolveStyle: function to resolve cell style
//
// Returns appropriate junction RenderOp based on border directions.
func (r *renderer) selectJunction(row, col int) RenderOp {
	rows := len(r.grid)
	cols := len(r.grid[0])

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
	top := r.leftBorderAt(row-1, col)
	right := r.topBorderAt(row, col)
	bottom := r.leftBorderAt(row, col)
	left := r.topBorderAt(row, col-1)

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
