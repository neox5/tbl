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
	t.AddRow()
	t.AddCell(tbl.Static, 1, 1, "Item")
	t.AddCell(tbl.Static, 1, 1, "Qty")
	t.AddCell(tbl.Static, 1, 1, "Price")

	// Content rows
	t.AddRow()
	t.AddCell(tbl.Static, 1, 1, "Apples")
	t.AddCell(tbl.Static, 1, 1, "10")
	t.AddCell(tbl.Static, 1, 1, "$2.50")

	t.AddRow()
	t.AddCell(tbl.Static, 1, 1, "Bananas")
	t.AddCell(tbl.Static, 1, 1, "5")
	t.AddCell(tbl.Static, 1, 1, "$3.10")

	t.AddRow()
	t.AddCell(tbl.Static, 1, 1, "Oranges")
	t.AddCell(tbl.Static, 1, 1, "8")
	t.AddCell(tbl.Static, 1, 1, "$4.00")

	// Footer
	t.AddRow()
	t.AddCell(tbl.Static, 1, 1, "Total")
	t.AddCell(tbl.Static, 1, 1, "23")
	t.AddCell(tbl.Static, 1, 1, "$9.60")

	// Column styles for alignment: labels left, numbers right
	t.SetColStyle(0, tbl.Left())
	t.SetColStyle(1, tbl.Right())
	t.SetColStyle(2, tbl.Right())

	// Default: no borders; SetStyleFunc defines all borders.
	t.SetDefaultStyle(tbl.Pad(0, 1), tbl.BNone())

	lastCol := cols - 1

	// SetStyleFunc controls border shape:
	// - header and footer: full border
	// - content rows: only table borders left and right, no column separation
	t.SetStyleFunc(func(row, col int) tbl.CellStyle {
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
