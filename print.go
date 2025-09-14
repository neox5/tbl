package tbl

import "fmt"

func (t *Table) PrintState() {
	fmt.Println("=== Table State ===")
	fmt.Printf("Cells: %d\n", len(t.cells))
	fmt.Printf("Rows: %d\n", len(t.rowStarts))
	fmt.Printf("Current Index: %d\n", t.currIndex)
	fmt.Printf("Width: %d\n", t.width)

	if len(t.rowStarts) > 0 {
		fmt.Printf("Row Starts: %v\n", t.rowStarts)
	}

	if len(t.colLevels) > 0 {
		fmt.Printf("Col Levels: %v\n", t.colLevels)
	}

	if len(t.colWidths) > 0 {
		fmt.Printf("Col Widths: %v\n", t.colWidths)
	}

	if len(t.openFlexCells) > 0 {
		fmt.Printf("Open Flex Cells: %v\n", t.openFlexCells)
	}

	if len(t.hLines) > 0 {
		fmt.Printf("H Lines: %v\n", t.hLines)
	}

	fmt.Println("==================")
}
