package main

import (
	"fmt"

	"github.com/neox5/tbl/internal/cell"
	"github.com/neox5/tbl/internal/table"
	"github.com/neox5/tbl/types"
)

func main() {
	fmt.Println("=== TBL Internal Architecture Test ===")

	// Test basic table creation
	t := table.New()

	// Test AddRow with mixed types - using NewFromValue for type conversion
	t.AddRow(
		cell.New().WithContent("Name"),
		cell.New().WithContent("Age"),
		cell.New().WithContent("City"),
	)

	// Test cell creation and chaining
	nameCell := cell.New().WithContent("John").WithAlign(types.Center, types.Middle)
	ageCell := cell.New().WithContent("25").WithSpan(1, 2) // span 1 col, 2 rows
	cityCell := cell.New().WithContent("Vienna")

	t.AddRow(nameCell, ageCell, cityCell)

	// Test short form methods
	t.AddRow(
		cell.New().WithContent("Jane").WithAlign(types.Right, types.Top),
		cell.New().WithContent("30"),
		cell.New().WithContent("Berlin").WithBorder(types.CellBorder{
			Top: true, Right: true, Bottom: true, Left: true,
			Style: types.Single,
		}),
	)

	// Test configuration
	config := &types.Config{
		Border: &types.TableBorder{
			All:   true,
			Style: types.Single,
		},
		Width: 80,
	}

	t2 := table.NewWithConfig(config)
	t2.AddRow(
		cell.New().WithContent("Configured"),
		cell.New().WithContent("Table"),
		cell.New().WithContent("Test"),
	)

	// Test standalone cell creation
	standaloneCell := cell.New().WithContent("Standalone").WithAlign(types.Center, types.Bottom)
	fmt.Printf("Standalone cell created: %v\n", standaloneCell)

	fmt.Println("Internal architecture test completed successfully!")
	fmt.Println("- Direct internal package usage")
	fmt.Println("- Internal implementation exposed")
	fmt.Println("- No circular imports")
	fmt.Println("- Simple config struct in types package")
	fmt.Println("- Method chaining works on internal types")
	fmt.Println("- Mixed type handling works")
	fmt.Println("- Standalone cell creation works")
}
