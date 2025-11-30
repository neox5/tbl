// Package tbl provides CLI table rendering.
//
// Features:
//   - Cell-level row/col spanning
//   - Flex and Static cell types for dynamic column sizing
//   - Multi-line content with word wrapping
//   - Configurable borders, padding, and alignment
//   - Dynamic column resolution
//   - Content-aware styling with StyleContext
//   - Composable predicate logic with And, Or, Not
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
//
// 4. Content-based styling:
//
//	t := tbl.New()
//	// ... add rows ...
//
//	negative := func(ctx tbl.StyleContext) bool {
//	    val, err := strconv.ParseFloat(ctx.Content, 64)
//	    return err == nil && val < 0
//	}
//	t.SetStyleFunc(
//	    tbl.FirstRow(tbl.BBottom()),
//	    tbl.Predicate(negative, tbl.Red()),
//	)
//	t.Print()
//
// 5. Composable predicates with logical operators:
//
//	t := tbl.New()
//	// ... add rows ...
//
//	// Style even rows excluding header
//	evenNotFirst := tbl.And(
//	    func(ctx tbl.StyleContext) bool { return ctx.Row%2 == 0 },
//	    tbl.Not(func(ctx tbl.StyleContext) bool { return ctx.Row == 0 }),
//	)
//	t.SetStyleFunc(
//	    tbl.FirstRow(tbl.BBottom(), tbl.Bold()),
//	    tbl.Predicate(evenNotFirst, tbl.Pad(1)),
//	)
//	t.Print()
package tbl
