package main

import (
	"fmt"

	"github.com/neox5/tbl"
)

func main() {
	fmt.Println("=== TBL Library Test ===")

	// Test basic config creation
	cfg := tbl.Config{
		Border: tbl.TableBorder{
			All:   true,
			Style: tbl.Single,
		},
		DefaultCell: tbl.Cell{
			Content: "default",
			ColSpan: 1,
			RowSpan: 1,
		},
	}

	fmt.Printf("Config created: %+v\n", cfg)

	// Test table creation
	table := tbl.New(cfg)
	fmt.Printf("Table created: %p\n", table)

	// Test border chars access
	chars := tbl.DefaultBorderChars[tbl.Single]
	fmt.Printf("Single border chars - H:'%c' V:'%c'\n", chars.Horizontal, chars.Vertical)

	// Test FLEX constant
	fmt.Printf("FLEX value: %d\n", tbl.FLEX)

	fmt.Println("✓ Basic functionality working")
}
