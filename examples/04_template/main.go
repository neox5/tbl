package main

import "github.com/neox5/tbl"

func main() {
	// Demo: Different templates

	// Thin (default)
	tbl.Simple(
		tbl.Row("Template", "Style"),
		tbl.Row("Thin", "Default"),
	).SetDefaultStyle(tbl.BAll()).Print()

	println()

	// ASCII
	tbl.Simple(
		tbl.Row("Template", "Style"),
		tbl.Row("ASCII", "Compatible"),
	).SetDefaultStyle(tbl.ASCII(), tbl.BAll()).Print()

	println()

	// Thick
	tbl.Simple(
		tbl.Row("Template", "Style"),
		tbl.Row("Thick", "Bold"),
	).SetDefaultStyle(tbl.Thick(), tbl.BAll()).Print()

	println()

	// Double
	tbl.Simple(
		tbl.Row("Template", "Style"),
		tbl.Row("Double", "Formal"),
	).SetDefaultStyle(tbl.Double(), tbl.BAll()).Print()

	println()

	// Validation: This should panic
	// t := tbl.New()
	// t.SetRowStyle(0, tbl.ASCII()) // panic: CharTemplate only supported via SetDefaultStyle
}
