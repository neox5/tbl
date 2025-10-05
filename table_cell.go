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
	row, col := t.cur.Pos()

	// Create and place area in grid
	a := grid.NewArea(row, col, rowSpan, colSpan)
	gid, err := t.g.AddArea(a)
	if err != nil {
		panic(fmt.Errorf("tbl: failed to add cell at cursor (%d,%d): %w", t.cur.Row(), t.cur.Col(), err))
	}

	t.cells[ID(gid)] = &cell{id: ID(gid), typ: typ}
	for i := col; i < col+colSpan; i++ {
		t.colLevels[i] = rowSpan
	}

	// Advance cursor
	t.cur.Advance(colSpan)
}
