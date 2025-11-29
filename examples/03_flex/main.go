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

	// Row 0: Flex cells (will expand)
	t.AddRow()
	t.AddCell(tbl.Static, 2, 1, "tbl")
	t.AddCell(tbl.Flex, 1, 1, "Name\nSurname")
	t.AddCell(tbl.Flex, 1, 1, "Role")

	// Row 1: Static cells with different structure
	t.AddRow()
	t.AddCell(tbl.Static, 1, 1, "Alice")
	t.AddCell(tbl.Static, 1, 1, "Farmer")
	t.AddCell(tbl.Static, 1, 1, "Backend\nEngineer")
	t.AddCell(tbl.Static, 1, 1, "Remote")

	// Row 2: Full span summary cells
	t.AddRow()
	t.AddCell(tbl.Flex, 1, 1, "Summary")
	t.AddCell(tbl.Static, 1, 1, "END")

	t.SetDefaultStyle(tbl.Pad(0, 1), tbl.Center(), tbl.Middle(), tbl.BAll())

	t.Print()
}
