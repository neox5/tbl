package main

import (
	"fmt"

	"github.com/neox5/tbl/internal/cell"
	"github.com/neox5/tbl/internal/table"
)

func main() {
	fmt.Println("=== TBL Grid Visualization Test ===")

	// Test basic table creation
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

	// Flex title row
	t.AddRow(
		cell.NewColFlex().WithContent("Table Title"),
	)

	// Header row
	t.AddRow(
		cell.New().WithContent("Name"),
		cell.New().WithContent("Age"),
		cell.New().WithContent("City"),
	)

	// Data row
	t.AddRow(
		cell.New().WithContent("John"),
		cell.New().WithContent("25"),
		cell.New().WithContent("NYC"),
	)

	fmt.Println("\nGrid Visualization:")
	t.PrintGrid()

	fmt.Println("\nTable State:")
	t.PrintDebug()

	fmt.Println("\nInternal architecture test completed successfully!")
}
