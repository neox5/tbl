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
	t.AddRow().
		AddCell(tbl.Static, 2, 1, "tbl").
		AddCell(tbl.Flex, 1, 1, "Name").
		AddCell(tbl.Flex, 1, 1, "Role")

	// Row 1: Static cells with different structure
	t.AddRow().
		AddCell(tbl.Static, 1, 1, "Alice").
		AddCell(tbl.Static, 1, 1, "Farmer").
		AddCell(tbl.Static, 1, 1, "Engineer").
		AddCell(tbl.Static, 1, 1, "Remote")

	fmt.Println("Rendered table:")
	fmt.Println(t.Render())
	fmt.Println()
	fmt.Println("Debug view:")
	fmt.Println(t.PrintDebug())
}
