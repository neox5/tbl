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
	t.AddRow(tbl.C("Item"), tbl.C("Qty"), tbl.C("Price"))

	// Content rows
	t.AddRow(tbl.C("Apples"), tbl.C("10"), tbl.C("$2.50"))
	t.AddRow(tbl.C("Bananas"), tbl.C("5"), tbl.C("$3.10"))
	t.AddRow(tbl.C("Oranges"), tbl.C("8"), tbl.C("$4.00"))

	// Footer
	t.AddRow(tbl.C("Total"), tbl.C("23"), tbl.C("$9.60"))

	// Column styles for alignment: labels left, numbers right
	t.SetColStyle(0, tbl.Left())
	t.SetColStyle(1, tbl.Right())
	t.SetColStyle(2, tbl.Right())

	// Default: no borders; SetStyleFunc defines all borders.
	t.SetDefaultStyle(tbl.Pad(0, 1), tbl.BNone())

	const lastCol = cols - 1

	// SetStyleFunc controls border shape:
	// - header and footer: full border
	// - content rows: only table borders left and right, no column separation
	t.SetStyleFunc(func(row, col int) tbl.CellStyle {
		style := tbl.CellStyle{}

		switch row {
		case headerRow, footerRow:
			// Header and footer: full border around all cells
			style.Border = tbl.BAll()
		default:
			// Content rows: only outer table borders on left and right
			switch col {
			case 0:
				style.Border = tbl.BLeft()
			case lastCol:
				style.Border = tbl.BRight()
			default:
				style.Border = tbl.BNone()
			}
		}

		return style
	})

	t.Print()
}
