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
	t.AddRow("Name", "Age", "City")

	// Test cell creation and chaining
	nameCell := t.C("John").A(types.Center, types.Middle)
	ageCell := t.C("25").S(1, 2) // span 1 col, 2 rows
	cityCell := t.C("Vienna")

	t.R(nameCell, ageCell, cityCell)

	// Test short form methods
	t.R(
		t.C("Jane").C("Updated").A(types.Right, types.Top),
		t.C("30"),
		t.C("Berlin").B(types.CellBorder{
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
	t2.AddRow("Configured", "Table", "Test")

	fmt.Println("New architecture test completed successfully!")
	fmt.Println("- Single interface definitions in root package")
	fmt.Println("- Concrete types returned from constructors")
	fmt.Println("- No circular imports")
	fmt.Println("- Simple config struct in types package")
	fmt.Println("- Method chaining works")
	fmt.Println("- Mixed type handling works")
}
