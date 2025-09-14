package table

// updateColIndex updates the column index for the given cell
func (t *Table) updateColIndex(cellIdx int, colStart, colSpan int) {
	for i := range colSpan {
		col := colStart + i
		if t.colIndex[col] == nil {
			t.colIndex[col] = make([]int, 0, 4)
		}
		t.colIndex[col] = append(t.colIndex[col], cellIdx)
	}
}

// updateRowIndex updates the row index for the given cell
func (t *Table) updateRowIndex(cellIdx int, rowStart, rowSpan int) {
	for i := range rowSpan {
		row := rowStart + i
		if t.rowIndex[row] == nil {
			t.rowIndex[row] = make([]int, 0)
		}
		t.rowIndex[row] = append(t.rowIndex[row], cellIdx)
	}
}

// CellsInCol returns the cell indices that overlap the specified column
func (t *Table) CellsInCol(col int) []int {
	return t.colIndex[col]
}

// CellsInRow returns the cell indices that overlap the specified row
func (t *Table) CellsInRow(row int) []int {
	return t.rowIndex[row]
}
