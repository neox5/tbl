// Package tbl provides CLI table rendering.
//
// Features:
//   - Cell-level row/col spanning
//   - Flex and Static cell types for dynamic column sizing
//   - Multi-line content with word wrapping
//   - Configurable borders, padding, and alignment
//   - Dynamic column resolution
//
// Usage Patterns:
//
// 1. Simple tables (helper function):
//
//	tbl.Simple(
//	    tbl.Row("Name", "Age"),
//	    tbl.Row("Alice", "30"),
//	).SetDefaultStyle(tbl.BAll()).Print()
//
// 2. Advanced tables (builder with IDs):
//
//	t := tbl.New()
//	headerRow := t.AddRow()
//	t.AddCell(tbl.Static, 1, 1, "Name")
//	ageCell := t.AddCell(tbl.Static, 1, 1, "Age")
//
//	t.AddRow()
//	t.AddCell(tbl.Static, 1, 1, "Alice")
//	t.AddCell(tbl.Static, 1, 1, "30")
//
//	t.SetRowStyle(headerRow, tbl.BBottom())
//	t.SetCellStyle(ageCell, tbl.Right())
//	t.Print()
//
// 3. Bulk operations (struct data):
//
//	employees := []Employee{...}
//	t := tbl.New()
//	t.AddRowsFromStructs(employees, "Name", "Age", "Salary")
//	t.SetRowStyle(0, tbl.BBottom()) // header
//	t.Print()
package tbl
