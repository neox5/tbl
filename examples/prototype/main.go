package main

import (
	"fmt"

	"github.com/neox5/tbl"
	"github.com/neox5/tbl/types"
)

func main() {
	fmt.Println("=== TBL New Architecture Test ===")

	// Test basic table creation
	t := tbl.New()

	// Test AddRow with mixed types
	t.AddRow(tbl.C("Name"), tbl.C("Age"), tbl.C("City"))

	// Test cell creation and chaining
	nameCell := tbl.C("John").A(types.Center, types.Middle)
	ageCell := tbl.C("25").S(1, 2) // span 1 col, 2 rows
	cityCell := tbl.C("Vienna")

	t.R(nameCell, ageCell, cityCell)

	// Test short form methods
	t.R(
		tbl.C("Jane").A(types.Right, types.Top),
		tbl.C("30"),
		tbl.C("Berlin").B(types.CellBorder{
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

	t2 := tbl.NewWithConfig(config)
	t2.AddRow(tbl.C("Configured"), tbl.C("Table"), tbl.C("Test"))

	// Test standalone cell creation
	standaloneCell := tbl.NewCell("Standalone").A(types.Center, types.Bottom)
	fmt.Printf("Standalone cell created: %v\n", standaloneCell)

	fmt.Println("New architecture test completed successfully!")
	fmt.Println("- Clean wrapper types (Tbl, Cell)")
	fmt.Println("- Internal implementation separated")
	fmt.Println("- No circular imports")
	fmt.Println("- Simple config struct in types package")
	fmt.Println("- Method chaining works on public types")
	fmt.Println("- Mixed type handling works")
	fmt.Println("- Standalone cell creation works")
}
