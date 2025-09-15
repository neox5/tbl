package main

import (
	"fmt"

	"github.com/neox5/tbl/internal/cell"
	"github.com/neox5/tbl/internal/table"
)

func main() {
	fmt.Println("=== TBL Comparison Table Test ===")

	// Test comparison table with flex rows followed by fixed rows
	t := table.New()

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\nPANIC DETECTED\n")
			fmt.Printf("Error: %v\n\n", r)
			fmt.Println("Current Grid State:")
			t.PrintGrid()
			fmt.Println()
			t.PrintDebug()
			return
		}
	}()

	// Row 0: Flex title spanning entire table
	t.AddRow(
		cell.NewColFlex().WithContent("Employee Comparison Table - Q4 Performance"),
	)

	// Row 1: Two flex sections - "Left Side" and "Right Side"
	t.AddRow(
		cell.NewColFlex().WithContent("Team Alpha"),
		cell.NewColFlex().WithContent("Team Beta"),
	)

	// Row 2: Headers for each section (fixed cells - triggers colsFixed = true)
	t.AddRow(
		cell.New().WithContent("Name"),
		cell.New().WithContent("Score"),
		cell.New().WithContent("Name"),
		cell.New().WithContent("Score"),
	)

	// Row 3: Person 1 data vs Person 2 data
	t.AddRow(
		cell.New().WithContent("Alice Johnson"),
		cell.New().WithContent("95"),
		cell.New().WithContent("Bob Wilson"),
		cell.New().WithContent("87"),
	)

	fmt.Println("\nGrid Visualization:")
	t.PrintGrid()

	fmt.Println("\nTable State:")
	t.PrintDebug()

	fmt.Println("\nComparison table test completed successfully!")
}
