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

	t.AddRow(tbl.C("Name"), tbl.C("Age"), tbl.C("Bio"))
	t.AddRow(tbl.C("Alice"), tbl.C("30"), tbl.C("Software engineer with 10 years of experience"))
	t.AddRow(tbl.C("Bob"), tbl.C("25"), tbl.C("Designer"))

	t.SetRowStyle(0, tbl.BBottom())

	t.Print()
}
