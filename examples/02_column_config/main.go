package main

import (
	"github.com/neox5/tbl"
)

func main() {
	// Example: Column configuration with width, minWidth, maxWidth
	t := tbl.New().
		AddCol(0, 10, 0). // Name: minWidth 10, auto width
		AddCol(7, 0, 0).  // Age: fixed width 7
		AddCol(0, 0, 20)  // Bio: maxWidth 20, auto width

	t.AddRow().
		AddCell(tbl.Static, 1, 1, "Name").
		AddCell(tbl.Static, 1, 1, "Age").
		AddCell(tbl.Static, 1, 1, "Bio")

	t.AddRow().
		AddCell(tbl.Static, 1, 1, "Alice").
		AddCell(tbl.Static, 1, 1, "30").
		AddCell(tbl.Static, 1, 1, "Software engineer with 10 years of experience")

	t.AddRow().
		AddCell(tbl.Static, 1, 1, "Bob").
		AddCell(tbl.Static, 1, 1, "25").
		AddCell(tbl.Static, 1, 1, "Designer")

	t.SetRowStyle(0, tbl.BBottom())

	t.Print()
}
