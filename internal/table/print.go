package table

import (
	"fmt"
	"strings"
)

func (t *Table) PrintGrid() {
	// Build content matrix
	grid := make([][]string, len(t.rowStarts))
	for row := range len(t.rowStarts) {
		grid[row] = make([]string, t.ColCount())
	}

	// Fill using existing indices
	for rowIdx := range len(t.rowStarts) {
		cellIndices := t.CellsInRow(rowIdx)

		for _, cellIdx := range cellIndices {
			cell := t.cells[cellIdx]
			label := string(rune('A' + cellIdx))

			for r := range cell.RowSpan() {
				for c := range cell.ColSpan() {
					if cell.RowStart()+r < len(grid) {
						grid[cell.RowStart()+r][cell.ColStart()+c] = label
					}
				}
			}
		}
	}

	// Print with separators
	for _, row := range grid {
		fmt.Printf("[ %s ]\n", strings.Join(row, " | "))
	}
}

// PrintDebug prints the internal state of the table for debugging
func (t *Table) PrintDebug() {
	fmt.Printf("=== Table Debug Info ===\n")
	fmt.Printf("Rows: %d, Cols: %d\n", t.RowCount(), t.ColCount())
	fmt.Printf("Next Index: %d, Next Col: %d\n", t.nextIndex, t.col)
	fmt.Printf("Cells: %d\n", len(t.cells))
	fmt.Printf("ColsFixed: %t\n", t.colsFixed)
	fmt.Printf("ColWidths: %v\n", t.colWidths)
	fmt.Printf("ColLevels: %v\n", t.colLevels)
	fmt.Printf("RowStarts: %v\n", t.rowStarts)
	fmt.Printf("FlexRows: %v\n", t.flexRows)
	fmt.Printf("========================\n")
}
