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

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic recovered: %v\n", r)
			t.PrintDebug()
		}
	}()

	// Test header row only
	t.AddRow(
		cell.New().WithContent("Name").WithAlign(types.Center, types.Top),
		cell.New().WithContent("Age").WithAlign(types.Center, types.Top),
		cell.New().WithContent("City").WithAlign(types.Center, types.Top),
	)


	fmt.Println("Internal architecture test completed successfully!")
}
