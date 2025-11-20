package tbl

// buildRenderLines constructs all rendering lines for the table.
// Returns complete line sequence as [][]RenderOp.
func (r *renderer) buildRenderLines() [][]RenderOp {
	var lines [][]RenderOp

	for row := range r.rowCount() {
		rowLines := r.buildLinesForRow(row)
		lines = append(lines, rowLines...)
	}

	// Add final bottom boundary if present
	if r.hBoundaries[r.rowCount()] {
		bottomLine := r.buildBoundaryLine(r.rowCount())
		lines = append(lines, bottomLine)
	}

	return lines
}

// buildLinesForRow constructs all lines for a single row.
// Combines boundary line (if present) with content lines.
func (r *renderer) buildLinesForRow(row int) [][]RenderOp {
	hasHBoundary := r.hBoundaries[row]

	var lines [][]RenderOp

	// Add boundary line if present
	if hasHBoundary {
		boundaryLine := r.buildBoundaryLine(row)
		lines = append(lines, boundaryLine)
	}

	// Add content lines
	contentLines := r.buildContentLines(row)
	lines = append(lines, contentLines...)

	return lines
}

// cellOffset calculates layout line index for cell at given row position.
// Returns index to the first line at row position in cell's layout.
// Includes all previous rows' content lines and boundary lines between them.
//
// Note: Does not include boundary line at row position itself.
// Caller must adjust offset if boundary exists at current row.
func (r *renderer) cellOffset(cell *Cell, row int) int {
	offset := 0
	for i := cell.r; i < row; i++ {
		offset += r.rowHeights[i]
		if i > cell.r && r.hBoundaries[i] {
			offset++
		}
	}
	return offset
}

// buildBoundarySegment constructs segment for single column in boundary line.
// Handles both normal cells (starting at row) and spanning cells (cutting through).
//
// Returns:
//   - RenderOp for this segment
//   - colsToSkip: number of columns to skip (1 for normal, cell.cSpan for spanning)
func (r *renderer) buildBoundarySegment(row, col int) (RenderOp, int) {
	// Bottom boundary: no spanning cells possible
	if row == r.rowCount() {
		op := RenderOp(Space{Width: r.colWidths[col]})
		if r.hBorderAt(row, col) {
			op = HLine{Width: r.colWidths[col]}
		}
		return op, 1
	}

	cell := r.grid[row][col]

	// Spanning cell: return content
	if cell.r < row {
		content := r.cellLayouts[cell.id]
		offset := r.cellOffset(cell, row)
		op := Content{Text: content[offset]}
		return op, cell.cSpan
	}

	// Normal cell: return horizontal segment
	op := RenderOp(Space{Width: r.colWidths[col]})
	if r.hBorderAt(row, col) {
		op = HLine{Width: r.colWidths[col]}
	}
	return op, 1
}

// buildBoundaryLine constructs horizontal boundary line at row position.
//
// Boundary semantics: row N = boundary before row N (top of row N).
// Special case: row = rowCount() constructs table's bottom boundary.
//
// Only called when r.hBoundaries[row] is true.
//
// Handles mixed content: cells with rowSpan > 1 show content on boundary line,
// while cells starting at this row show horizontal border segments.
//
// Returns complete boundary line as []RenderOp sequence.
func (r *renderer) buildBoundaryLine(row int) []RenderOp {
	var op RenderOp
	var ops []RenderOp

	for col := 0; col < r.colCount(); {
		// Add junction if vertical boundary exists
		if r.vBoundaries[col] {
			op = r.selectJunction(row, col)
			ops = append(ops, op)
		}

		// Build segment and get columns to skip
		segment, skip := r.buildBoundarySegment(row, col)
		ops = append(ops, segment)
		col += skip
	}

	// Add rightmost junction if boundary exists
	if r.vBoundaries[r.colCount()] {
		op = r.selectJunction(row, r.colCount())
		ops = append(ops, op)
	}

	return ops
}

// buildContentLines constructs all content lines for a row.
// Content lines are lines containing cell content (no horizontal boundary).
//
// Processes all cells once, building vertical borders and content
// for each line in parallel.
//
// Parameters:
//
//	row: grid row position
//
// Returns all content lines as [][]RenderOp.
func (r *renderer) buildContentLines(row int) [][]RenderOp {
	lineCount := r.rowHeights[row]
	lines := make([][]RenderOp, lineCount)
	for i := range lines {
		lines[i] = make([]RenderOp, 0)
	}

	var op RenderOp // reused throughout

	// Process each cell in the row
	for col := 0; col < r.colCount(); {
		cell := r.grid[row][col]
		content := r.cellLayouts[cell.id]

		// Calculate content offset for spanning cells
		offset := r.cellOffset(cell, row)

		// Adjust offset only for spanning cells with boundary at current row
		if cell.r < row && r.hBoundaries[row] {
			offset++
		}

		// Add left vertical boundary to all lines
		if r.vBoundaries[col] {
			op = Space{Width: 1}
			if r.vBorderAt(row, col) {
				op = VLine{}
			}
			for i := range lineCount {
				lines[i] = append(lines[i], op)
			}
		}

		// Add cell content to all lines
		for i := range lineCount {
			op = Content{Text: content[offset+i]}
			lines[i] = append(lines[i], op)
		}

		// Add right boundary on last cell
		if r.lastCellInRow(cell) && r.vBoundaries[r.colCount()] {
			op = Space{Width: 1}
			if r.vBorderAt(row, r.colCount()) {
				op = VLine{}
			}
			for i := range lineCount {
				lines[i] = append(lines[i], op)
			}
		}

		col += cell.cSpan
	}

	return lines
}

// hBorderAt reports whether horizontal border renders as character at position (row, col).
// Checks cells that meet at this position.
func (r *renderer) hBorderAt(row, col int) bool {
	// Out of bounds
	if row < 0 || row > r.rowCount() || col < 0 || col >= r.colCount() {
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
	if row == r.rowCount() {
		cell := r.grid[r.rowCount()-1][col]
		if cell != nil {
			style := r.styles[cell.id]
			return style.Border.IsVisual(BorderBottom)
		}
		return false
	}

	// Between rows: check both cells
	cellAbove := r.grid[row-1][col]
	cellBelow := r.grid[row][col]

	if cellAbove == cellBelow { // cell with rowSpan > 1
		return false
	}

	if r.styles[cellAbove.id].Border.IsVisual(BorderBottom) || r.styles[cellBelow.id].Border.IsVisual(BorderTop) {
		return true
	}

	return false
}

// vBorderAt reports whether vertical border renders as character at position (row, col).
// Checks cells that meet at this position.
func (r *renderer) vBorderAt(row, col int) bool {
	// Out of bounds
	if row < 0 || row >= r.rowCount() || col < 0 || col > r.colCount() {
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
	if col == r.colCount() {
		cell := r.grid[row][r.colCount()-1]
		if cell != nil {
			style := r.styles[cell.id]
			return style.Border.IsVisual(BorderRight)
		}
		return false
	}

	// Between columns: check both cells
	cellLeft := r.grid[row][col-1]
	cellRight := r.grid[row][col]

	if cellLeft == cellRight { // cell with colSpan > 1
		return false
	}

	if r.styles[cellLeft.id].Border.IsVisual(BorderRight) || r.styles[cellRight.id].Border.IsVisual(BorderLeft) {
		return true
	}

	return false
}

// selectJunction determines the junction character at a column position.
// Examines borders in all four directions (top, right, bottom, left).
//
// Parameters:
//   - row: border line position (0 = top border, rowCount() = bottom border)
//   - col: column position (0 = left edge, colCount() = right edge)
//
// Returns appropriate junction RenderOp based on border directions.
func (r *renderer) selectJunction(row, col int) RenderOp {
	// Table corners
	if row == 0 && col == 0 {
		return CornerTL{} // top-left: '┌'
	}
	if row == 0 && col == r.colCount() {
		return CornerTR{} // top-right: '┐'
	}
	if row == r.rowCount() && col == 0 {
		return CornerBL{} // bottom-left: '└'
	}
	if row == r.rowCount() && col == r.colCount() {
		return CornerBR{} // bottom-right: '┘'
	}

	// Calculate all 4 boundaries meeting at junction
	top := r.vBorderAt(row-1, col)
	right := r.hBorderAt(row, col)
	bottom := r.vBorderAt(row, col)
	left := r.hBorderAt(row, col-1)

	// Continuations
	if !top && right && !bottom && left {
		return HLine{Width: 1} // horizontal: '─'
	}
	if top && !right && bottom && !left {
		return VLine{} // vertical: '│'
	}

	// T-junctions
	if !top && right && bottom && left {
		return CornerT{} // top T: '┬'
	}
	if top && right && !bottom && left {
		return CornerB{} // bottom T: '┴'
	}
	if top && right && bottom && !left {
		return CornerL{} // left T: '├'
	}
	if top && !right && bottom && left {
		return CornerR{} // right T: '┤'
	}

	// Cross junction
	if top && right && bottom && left {
		return CornerX{} // cross: '┼'
	}

	// Interior corners (for cells with partial borders)
	if !top && right && bottom && !left {
		return CornerTL{} // interior top-left: '┌'
	}
	if !top && !right && bottom && left {
		return CornerTR{} // interior top-right: '┐'
	}
	if top && right && !bottom && !left {
		return CornerBL{} // interior bottom-left: '└'
	}
	if top && !right && !bottom && left {
		return CornerBR{} // interior bottom-right: '┘'
	}

	// No border at this junction
	return Space{Width: 1}
}
