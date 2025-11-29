package main

import (
	"github.com/neox5/tbl"
)

type Employee struct {
	Name   string
	Age    int
	City   string
	Salary float64
	Active bool
}

func main() {
	employees := []Employee{
		{"Alice Johnson", 30, "New York", 85000.50, true},
		{"Bob Smith", 25, "San Francisco", 92000.00, true},
		{"Carol White", 35, "Boston", 78000.75, false},
		{"David Brown", 28, "Seattle", 88500.25, true},
	}

	println("=== AddRowsFromStructs Example ===")
	println()

	tbl.New().
		SetDefaultStyle(tbl.BAll()).
		SetRowStyle(0, tbl.BBottom(), tbl.Center()).
		SetColStyle(1, tbl.Right()). // Age
		SetColStyle(3, tbl.Right()). // Salary
		AddRowsFromStructs(employees, "Name", "Age", "City", "Salary", "Active").
		Print()
}
