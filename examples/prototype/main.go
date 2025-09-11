package main

import (
	"fmt"

	"github.com/neox5/tbl"
)

func main() {
	fmt.Println("=== TBL Playground ===")

	t := tbl.New()
	t.AddRow(
		t.C("Name").S(1, 2),
		t.C("Age").S(1, 2),
		t.C("City").S(1, 2),
	)
	t.AddRow("chris", 35, "vienna")
	t.PrintState()
}
