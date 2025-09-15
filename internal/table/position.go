package table

// advanceCol advances the column index to the next free column (colLevel = 0)
func (t *Table) advanceCol() {
	// Find next free column
	for t.col < t.ColCount() && t.colLevels[t.col] > 0 {
		t.col++
	}

	// If we've reached the end and have flex cells, expand
	if t.col >= t.ColCount() && !t.colsFixed {
		t.colWidths = append(t.colWidths, 0)
		t.colLevels = append(t.colLevels, 0)
		// t.col is already at the new column position
		return
	}

	// If no flex cells available, reset to start of next row
	if t.col >= t.ColCount() {
		t.col = 0
	}
}

// advanceRow advances to the next row, starts new row, and decrements colLevels
func (t *Table) advanceRow() {
	t.row++
	t.startNewRow()

	// If no columns exist yet (first row), no need to check levels
	if len(t.colLevels) == 0 {
		return
	}

	hasFreeCol := false
	for i := range t.colLevels {
		t.colLevels[i]--
		if t.colLevels[i] == 0 {
			hasFreeCol = true
		}
	}

	if !hasFreeCol {
		t.advanceRow() // Recursive call if no columns are free
	}

	t.col = 0 // reset to the beginning
}

// startNewRow starts a new row by recording the current cell position
func (t *Table) startNewRow() {
	t.rowStarts = append(t.rowStarts, len(t.cells))
}

// spanFit returns true if the span can fit at the current row position
func (t *Table) spanFit(span int) bool {
	r := 0
	for i := t.col; i < t.ColCount(); i++ {
		if t.colLevels[i] == 0 {
			r++
			continue
		}
		break
	}
	return r >= span
}
