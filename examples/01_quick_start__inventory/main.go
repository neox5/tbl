package main

import (
	"os"
	"path/filepath"

	"github.com/neox5/tbl"
)

func main() {
	table := tbl.New().Simple(
		tbl.Row("Product", "SKU", "Stock", "Price"),
		tbl.Row("USB-C Cable", "ACC-001", "156", "$12.99"),
		tbl.Row("Wireless Mouse", "PER-042", "89", "$24.99"),
		tbl.Row("Mechanical Keyboard", "PER-103", "34", "$89.99"),
		tbl.Row("Monitor Stand", "ACC-028", "67", "$45.50"),
		tbl.Row("Laptop Sleeve", "ACC-015", "203", "$18.99"),
	).SetDefaultStyle(tbl.BAll(), tbl.Pad(0, 1))

	// Print to stdout
	table.Print()

	// Write to output.txt (run from project root)
	outputPath := filepath.Join("examples", "01_quick_start__inventory", "output.txt")
	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := table.RenderTo(f); err != nil {
		panic(err)
	}
}
