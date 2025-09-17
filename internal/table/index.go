package table

// addRowIndex adds cell to row index
func (t *Table) addRowIndex(cellIdx int, rowStart, rowSpan int) {
	for i := range rowSpan {
		row := rowStart + i
		if t.rowIndex[row] == nil {
			t.rowIndex[row] = make([]int, 0)
		}
		t.rowIndex[row] = append(t.rowIndex[row], cellIdx)
	}
}

// CellsInRow returns the cell indices that overlap the specified row
func (t *Table) CellsInRow(row int) []int {
	return t.rowIndex[row]
}
