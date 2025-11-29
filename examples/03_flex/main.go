package main

import (
	"fmt"
	"runtime"

	"github.com/neox5/tbl"
)

func main() {
	fmt.Println("=== Flex Table Example ===")
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
	t.AddRow(
		tbl.Cx(2, 1, "tbl"),
		tbl.F("Name\nSurname"),
		tbl.F("Role"),
	)

	// Row 1: Static cells with different structure
	t.AddRow(
		tbl.C("Alice"),
		tbl.C("Farmer"),
		tbl.C("Backend\nEngineer"),
		tbl.C("Remote"),
	)

	// Row 2: Full span summary cells
	t.AddRow(
		tbl.F("Summary"),
		tbl.C("END"),
	)

	t.SetDefaultStyle(tbl.Pad(0, 1), tbl.Center(), tbl.Middle(), tbl.BAll())

	t.Print()
}
