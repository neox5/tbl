package tbl

import "fmt"

func (t *Table) PrintState() {
	fmt.Println("=== Table State ===")
	fmt.Printf("Rows: %d\n", len(t.rows))
	fmt.Printf("Current Row: %d\n", t.row)
	fmt.Printf("Flexible Cols: %t\n", t.flexibleCols)
	fmt.Printf("Width: %d\n", t.width)

	if len(t.virtualRows) > 0 {
		fmt.Printf("Virtual Rows: %v\n", t.virtualRows)
	}

	if len(t.colLevels) > 0 {
		fmt.Printf("Col Levels: %v\n", t.colLevels)
	}

	if len(t.colWidths) > 0 {
		fmt.Printf("Col Widths: %v\n", t.colWidths)
	}

	if len(t.hLines) > 0 {
		fmt.Printf("H Lines: %v\n", t.hLines)
	}

	fmt.Println("==================")
}
