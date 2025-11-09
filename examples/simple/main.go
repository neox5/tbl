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

	t.AddRow().
		AddCell(tbl.Flex, 1, 1, "Title").
		AddCell(tbl.Static, 2, 1, "Author")

	t.AddRow().
		AddCell(tbl.Static, 1, 2, "The Go Programming Language")

	fmt.Println(t.Render())
	fmt.Println()
	fmt.Println("Debug view:")
	fmt.Println(t.PrintDebug())
}
