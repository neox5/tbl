// example/funcstyler/main.go
package main

import (
	"github.com/neox5/tbl"
)

func main() {
	const (
		cols      = 3
		headerRow = 0
		footerRow = 4
	)

	t := tbl.NewWithCols(cols)

	// Header
	t.AddRow().
		AddCell(tbl.Static, 1, 1, "Item").
		AddCell(tbl.Static, 1, 1, "Qty").
		AddCell(tbl.Static, 1, 1, "Price")

	// Content rows
	t.AddRow().
		AddCell(tbl.Static, 1, 1, "Apples").
		AddCell(tbl.Static, 1, 1, "10").
		AddCell(tbl.Static, 1, 1, "$2.50")

	t.AddRow().
		AddCell(tbl.Static, 1, 1, "Bananas").
		AddCell(tbl.Static, 1, 1, "5").
		AddCell(tbl.Static, 1, 1, "$3.10")

	t.AddRow().
		AddCell(tbl.Static, 1, 1, "Oranges").
		AddCell(tbl.Static, 1, 1, "8").
		AddCell(tbl.Static, 1, 1, "$4.00")

	// Footer
	t.AddRow().
		AddCell(tbl.Static, 1, 1, "Total").
		AddCell(tbl.Static, 1, 1, "23").
		AddCell(tbl.Static, 1, 1, "$9.60")

	// Column styles for alignment: labels left, numbers right
	t.SetColStyle(0, tbl.Left()).
		SetColStyle(1, tbl.Right()).
		SetColStyle(2, tbl.Right())

	// Default: no borders; Funcstyler defines all borders.
	t.SetDefaultStyle(tbl.Pad(0, 1), tbl.BNone())

	lastCol := cols - 1

	// Funcstyler controls border shape:
	// - header and footer: full border
	// - content rows: only table borders left and right, no column separation
	t.SetFuncStyle(func(row, col int) tbl.CellStyle {
		style := tbl.CellStyle{}

		switch {
		case row == headerRow || row == footerRow:
			// Header and footer: full border around all cells
			style.Border = tbl.BAll()
		case row > headerRow && row < footerRow:
			// Content rows: only outer table borders on left and right
			if col == 0 {
				style.Border = tbl.BLeft()
			} else if col == lastCol {
				style.Border = tbl.BRight()
			} else {
				style.Border = tbl.BNone()
			}
		}

		return style
	})

	t.Print()
}
