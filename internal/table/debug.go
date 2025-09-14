package table

import "fmt"

// PrintDebug prints the internal state of the table for debugging
func (t *Table) PrintDebug() {
	fmt.Printf("=== Table Debug Info ===\n")
	fmt.Printf("Rows: %d, Cols: %d\n", t.row, t.ColCount())
	fmt.Printf("Next Index: %d, Next Col: %d\n", t.nextIndex, t.col)
	fmt.Printf("Cells: %d\n", len(t.cells))
	fmt.Printf("ColWidths: %v\n", t.colWidths)
	fmt.Printf("ColLevels: %v\n", t.colLevels)
	fmt.Printf("RowStarts: %v\n", t.rowStarts)
	fmt.Printf("OpenFlexCells: %v\n", t.openFlexCells)
	fmt.Printf("========================\n")
}
