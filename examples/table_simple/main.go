package main

import (
	"fmt"

	"github.com/neox5/tbl"
)

func main() {
	fmt.Println("=== Simple Table Example ===")
	fmt.Println()

	// Create table
	t := tbl.New()

	// Defer panic recovery to show debug output
	defer func() {
		if r := recover(); r != nil {
			fmt.Println()
			fmt.Println("=== PANIC OCCURRED ===")
			fmt.Printf("Error: %v\n", r)
			fmt.Println()
			fmt.Println("Table state at panic:")
			fmt.Println(t.PrintDebug())
		}
	}()

	// Row 0: Static + Flex + Static
	t.AddRow().
		AddCell(tbl.Static, 1, 1).
		AddCell(tbl.Flex, 1, 1).
		AddCell(tbl.Static, 1, 1)

	// Row 1: Static(span 2) + Flex
	t.AddRow().
		AddCell(tbl.Static, 1, 2).
		AddCell(tbl.Flex, 1, 1)

	// Row 2: Flex + Static + Static
	t.AddRow().
		AddCell(tbl.Flex, 1, 1).
		AddCell(tbl.Static, 1, 1).
		AddCell(tbl.Static, 1, 2)

	fmt.Println("Success!")
	fmt.Println(t.PrintDebug())
}
