package main

import (
	"fmt"

	"github.com/neox5/tbl"
)

func main() {
	fmt.Println("=== Simple Flex Example ===")
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
		AddCell(tbl.Flex, 1, 1).
		AddCell(tbl.Flex, 1, 1)

	t.AddRow().
		AddCell(tbl.Static, 1, 1).
		AddCell(tbl.Static, 1, 1).
		AddCell(tbl.Static, 1, 1).
		AddCell(tbl.Static, 1, 1)

	fmt.Println("Success!")
	fmt.Println(t.PrintDebug())
}
