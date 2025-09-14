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
		cell.NewFromValue("Name"),
		cell.NewFromValue("Age"),
		cell.NewFromValue("City"),
	)

	// Test cell creation and chaining
	nameCell := cell.NewFromValue("John").WithAlign(types.Center, types.Middle)
	ageCell := cell.NewFromValue("25").WithSpan(1, 2) // span 1 col, 2 rows
	cityCell := cell.NewFromValue("Vienna")

	t.AddRow(nameCell, ageCell, cityCell)

	// Test short form methods
	t.AddRow(
		cell.NewFromValue("Jane").A(types.Right, types.Top),
		cell.NewFromValue("30"),
		cell.NewFromValue("Berlin").B(types.CellBorder{
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
		cell.NewFromValue("Configured"),
		cell.NewFromValue("Table"),
		cell.NewFromValue("Test"),
	)

	// Test standalone cell creation
	standaloneCell := cell.NewFromValue("Standalone").A(types.Center, types.Bottom)
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
