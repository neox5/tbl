package main

import (
	"github.com/neox5/tbl"
)

func main() {
	// Example: Column configuration with width, minWidth, maxWidth
	t := tbl.New()
	t.AddCol(0, 10, 0) // Name: minWidth 10, auto width
	t.AddCol(7, 0, 0)  // Age: fixed width 7
	t.AddCol(0, 0, 20) // Bio: maxWidth 20, auto width

	t.AddRow()
	t.AddCell(tbl.Static, 1, 1, "Name")
	t.AddCell(tbl.Static, 1, 1, "Age")
	t.AddCell(tbl.Static, 1, 1, "Bio")

	t.AddRow()
	t.AddCell(tbl.Static, 1, 1, "Alice")
	t.AddCell(tbl.Static, 1, 1, "30")
	t.AddCell(tbl.Static, 1, 1, "Software engineer with 10 years of experience")

	t.AddRow()
	t.AddCell(tbl.Static, 1, 1, "Bob")
	t.AddCell(tbl.Static, 1, 1, "25")
	t.AddCell(tbl.Static, 1, 1, "Designer")

	t.SetRowStyle(0, tbl.BBottom())

	t.Print()
}
