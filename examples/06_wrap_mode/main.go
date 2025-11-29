package main

import (
	"github.com/neox5/tbl"
)

func main() {
	longText := "This is a very long sentence that demonstrates different wrapping behaviors in table cells when content exceeds the available width"

	println("=== WrapMode Comparison ===")
	println()

	t := tbl.New()
	t.AddCol(0, 15, 20) // Mode column
	t.AddCol(0, 30, 40) // Content column

	t.AddRow()
	t.AddCell(tbl.Static, 1, 1, "Mode")
	t.AddCell(tbl.Static, 1, 1, "Behavior")

	t.AddRow()
	t.AddCell(tbl.Static, 1, 1, "WrapWord")
	t.AddCell(tbl.Static, 1, 1, longText)

	t.AddRow()
	t.AddCell(tbl.Static, 1, 1, "WrapChar")
	t.AddCell(tbl.Static, 1, 1, longText)

	t.AddRow()
	t.AddCell(tbl.Static, 1, 1, "WrapTruncate")
	t.AddCell(tbl.Static, 1, 1, longText)

	t.SetDefaultStyle(tbl.BAll())
	t.SetRowStyle(0, tbl.BBottom(), tbl.Center())
	t.SetRowStyle(1, tbl.WrapWord)
	t.SetRowStyle(2, tbl.WrapChar)
	t.SetRowStyle(3, tbl.WrapTruncate)
	t.Print()
}
