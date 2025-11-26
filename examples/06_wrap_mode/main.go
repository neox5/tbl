package main

import (
	"github.com/neox5/tbl"
)

func main() {
	longText := "This is a very long sentence that demonstrates different wrapping behaviors in table cells when content exceeds the available width"

	println("=== WrapMode Comparison ===")
	println()

	t := tbl.New().
		AddCol(0, 15, 20). // Mode column
		AddCol(0, 30, 40)  // Content column

	t.AddRow().
		AddCell(tbl.Static, 1, 1, "Mode").
		AddCell(tbl.Static, 1, 1, "Behavior")

	t.AddRow().
		AddCell(tbl.Static, 1, 1, "WrapWord").
		AddCell(tbl.Static, 1, 1, longText)

	t.AddRow().
		AddCell(tbl.Static, 1, 1, "WrapChar").
		AddCell(tbl.Static, 1, 1, longText)

	t.AddRow().
		AddCell(tbl.Static, 1, 1, "WrapTruncate").
		AddCell(tbl.Static, 1, 1, longText)

	t.SetDefaultStyle(tbl.BAll()).
		SetRowStyle(0, tbl.BBottom(), tbl.Center()).
		SetRowStyle(1, tbl.WrapWord).
		SetRowStyle(2, tbl.WrapChar).
		SetRowStyle(3, tbl.WrapTruncate).
		Print()
}
