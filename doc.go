// Package tbl provides CLI table rendering.
//
// Features:
//   - Cell-level row/col spanning
//   - Flex and Static cell types for dynamic column sizing
//   - Multi-line content with word wrapping
//   - Configurable borders, padding, and alignment
//   - Dynamic column resolution
//
// Basic usage:
//
//	t := tbl.New()
//	t.AddRow().
//	    AddCell(tbl.Static, 1, 1, "Name").
//	    AddCell(tbl.Static, 1, 1, "Age")
//	t.AddRow().
//	    AddCell(tbl.Static, 1, 1, "Alice").
//	    AddCell(tbl.Static, 1, 1, "30")
//	t.Print()
package tbl
