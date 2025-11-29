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

	t.AddRow(tbl.C("Mode"), tbl.C("Behavior"))
	t.AddRow(tbl.C("WrapWord"), tbl.C(longText))
	t.AddRow(tbl.C("WrapChar"), tbl.C(longText))
	t.AddRow(tbl.C("WrapTruncate"), tbl.C(longText))

	t.SetDefaultStyle(tbl.BAll())
	t.SetRowStyle(0, tbl.BBottom(), tbl.Center())
	t.SetRowStyle(1, tbl.WrapWord)
	t.SetRowStyle(2, tbl.WrapChar)
	t.SetRowStyle(3, tbl.WrapTruncate)
	t.Print()
}
