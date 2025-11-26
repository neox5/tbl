package main

import (
	"github.com/neox5/tbl"
)

func main() {
	tbl.Simple(
		tbl.Row("Product", "Price", "Stock"),
		tbl.Row("Widget", "$10.00", "100"),
		tbl.Row("Gadget", "$25.50", "50"),
		tbl.Row("Doohickey", "$5.99", "200"),
	).SetDefaultStyle(tbl.BAll()).Print()
}
