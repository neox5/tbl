package tbl

import (
	"fmt"

	"github.com/neox5/tbl/internal/grid"
)

// cell holds metadata for a table cell.
// Position and span data is stored in the grid.
type cell struct {
	id  ID
	typ CellType
}

// addCell places cell in grid, creates metadata, and advances cursor.
func (t *Table) addCell(typ CellType, rowSpan, colSpan int) {
	row, col := t.cur.position()

	// Create and place area in grid
	a := grid.NewArea(row, col, rowSpan, colSpan)
	gid, err := t.g.AddArea(a)
	if err != nil {
		panic(fmt.Errorf("tbl: failed to add cell: %w", err))
	}

	// Create cell metadata
	c := &cell{id: ID(gid), typ: typ}
	t.cells[ID(gid)] = c

	// Advance cursor
	t.cur.advance(colSpan)
}
