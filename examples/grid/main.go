package main

import (
	"fmt"

	"github.com/neox5/tbl/internal/grid"
)

func main() {
	fmt.Println("=== Grid ShiftRightRow Example ===")
	fmt.Println()

	// Create 5x3 grid
	g := grid.New(5, 3)

	// Build grid:
	// [ A | B | C | D   D ]
	// [ E | B | F | G | H ]
	// [ I | B | F | J | H ]

	// Row 0
	g.AddArea(grid.NewArea(0, 0, 1, 1)) // A at (0,0)
	g.AddArea(grid.NewArea(1, 0, 1, 3)) // B at (1,0) spanning 3 rows
	g.AddArea(grid.NewArea(2, 0, 1, 1)) // C at (2,0)
	g.AddArea(grid.NewArea(3, 0, 2, 1)) // D at (3,0) spanning 2 cols

	// Row 1
	g.AddArea(grid.NewArea(0, 1, 1, 1)) // E at (0,1)
	g.AddArea(grid.NewArea(2, 1, 1, 2)) // F at (2,1) spanning 2 rows
	g.AddArea(grid.NewArea(3, 1, 1, 1)) // G at (3,1)
	g.AddArea(grid.NewArea(4, 1, 1, 2)) // H at (4,1) spanning 2 rows

	// Row 2
	g.AddArea(grid.NewArea(0, 2, 1, 1)) // I at (0,2)
	g.AddArea(grid.NewArea(3, 2, 1, 1)) // J at (3,2)

	fmt.Println("Initial Grid:")
	fmt.Println(g.Print())
	fmt.Println()

	// Add column for shift space
	g.AddCol()

	fmt.Println("After AddCol:")
	fmt.Println(g.Print())
	fmt.Println()

	// Shift row 0 right starting at column 1
	fmt.Println("Executing: ShiftRightRow(0, 1)")
	if err := g.ShiftRightRow(0, 1); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("After First Shift:")
	fmt.Println(g.Print())
	fmt.Println()

	// Add another column for second shift
	g.AddCol()

	fmt.Println("After Second AddCol:")
	fmt.Println(g.Print())
	fmt.Println()

	// Shift row 0 right again starting at column 1
	fmt.Println("Executing: ShiftRightRow(0, 1)")
	if err := g.ShiftRightRow(0, 1); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println("After Second Shift:")
	fmt.Println(g.Print())
}
