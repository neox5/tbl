package tbl

// getCellAt finds the cell containing position (row, col).
// Returns nil if position is empty.
func (t *Table) getCellAt(row, col int) *Cell {
	for _, cell := range t.cells {
		if cell.Contains(row, col) {
			return cell
		}
	}
	return nil
}

// isFlex reports whether the cell at (row, col) is a Flex type.
func (t *Table) isFlex(row, col int) bool {
	cell := t.getCellAt(row, col)
	return cell != nil && cell.typ == Flex
}

// isWall reports whether the cell at (row, col) acts as a wall.
// A cell is a wall if it spans multiple rows and originates above row.
func (t *Table) isWall(row, col int) bool {
	cell := t.getCellAt(row, col)
	return cell != nil && cell.rSpan > 1 && cell.r < row
}

// findFirstFreeCol locates first unoccupied column in row.
// Returns column index of first free position, or Cols() if row full.
// Accounts for cells with rowSpan > 1 from previous rows.
func (t *Table) findFirstFreeCol(row int) int {
	if row < 0 || row >= t.g.Rows() {
		return 0
	}
	return t.g.NextZeroInRow(row, 0)
}

// calculateNeeded determines how many columns required to fit cell.
// Returns columns needed beyond current grid width or first blocking position.
func (t *Table) calculateNeeded(row, col, colSpan int) int {
	// Find first blocking position in row
	firstBlocked := t.g.NextOneInRow(row, col)
	if firstBlocked == -1 {
		firstBlocked = t.g.Cols()
	}

	// Calculate shortage
	required := col + colSpan
	if required <= firstBlocked {
		return 0
	}

	return required - firstBlocked
}
