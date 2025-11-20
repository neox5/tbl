package main

import (
	"fmt"
	"runtime"

	"github.com/neox5/tbl"
)

func main() {
	fmt.Println("=== Simple Table Example ===")
	fmt.Println()

	t := tbl.New()

	defer func() {
		if r := recover(); r != nil {
			fmt.Println()
			fmt.Println("=== PANIC OCCURRED ===")
			fmt.Printf("Error: %v\n", r)
			fmt.Println()
			fmt.Println("Stack trace:")
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			fmt.Println(string(buf[:n]))
			fmt.Println()
			fmt.Println("Table state at panic:")
			fmt.Println(t.PrintDebug())
		}
	}()

	// t.AddRow()
	// t.AddCell(tbl.Flex, 1, 1, "TEST")
	// t.AddCell(tbl.Static, 1, 4, "tbl")

	// Row 0: Flex cells (will expand)
	t.AddRow().
		AddCell(tbl.Static, 2, 1, "tbl").
		AddCell(tbl.Flex, 1, 1, "Name\nSurname").
		AddCell(tbl.Flex, 1, 1, "Role")

	// Row 1: Static cells with different structure
	t.AddRow().
		AddCell(tbl.Static, 1, 1, "Alice").
		AddCell(tbl.Static, 1, 1, "Farmer").
		AddCell(tbl.Static, 1, 1, "Backend\nEngineer").
		AddCell(tbl.Static, 1, 1, "Remote")

	// Row 2: Full span summary cells
	t.AddRow().
		AddCell(tbl.Flex, 1, 1, "Summary").
		AddCell(tbl.Static, 1, 1, "END")

	// t.BorderAll()
	t.SetDefaultStyle(tbl.CellStyle{
		Padding: tbl.Padding{Left: 1, Right: 1},
		HAlign:  tbl.HAlignCenter,
		VAlign:  tbl.VAlignMiddle,
		Border:  tbl.Border{Sides: tbl.BorderAll},
	})

	t.Print()
}
