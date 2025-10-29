package tbl

import (
	"fmt"

	"github.com/neox5/btmp"
	"github.com/neox5/tbl/internal/cursor"
	"github.com/neox5/tbl/internal/grid"
)

// ID identifies a cell in the table.
type ID grid.ID

// Table manages incremental table construction with flex/static cells.
type Table struct {
	g     *btmp.Grid
	c     *cursor.Cursor
	rows  [][]ID
	cols  [][]ID
	cells map[ID]*Cell

	colsFixed  bool
	nextCellID ID
}

// New creates a new Table with zero columns.
func New() *Table {
	return NewWithCols(0)
}

// NewWithCols creates a new Table with initial column capacity.
func NewWithCols(cols int) *Table {
	if cols < 0 {
		panic("tbl: invalid cols value")
	}

	t := &Table{
		g:          btmp.NewGridWithSize(0, cols),
		c:          cursor.New(),
		rows:       make([][]ID, 0),
		cols:       make([][]ID, cols),
		cells:      make(map[ID]*Cell),
		nextCellID: 1,
	}

	if cols > 0 {
		t.colsFixed = true
	}

	return t
}

// AddRow advances to the next row for cell placement.
func (t *Table) AddRow() *Table {
	t.c.NextRow()

	t.rows = append(t.rows, []ID{})
	return t
}

// addCols grows grid columns to accommodate cell at cursor.
func (t *Table) addCols(colSpan int) {
	needed := t.c.Col() + colSpan
	for needed > t.g.Cols() {
		// TODO
	}
}

// AddCell adds a cell at the.c.ent.c.or position.
// rowSpan and colSpan define cell dimensions (must be > 0).
// Returns *Table for chaining. Panics on invalid input.
//
// Column expansion behavior:
//   - If no columns exist, expands by colSpan
//   - If colsFixed is false, expands when no free columns available
//   - If colsFixed is true and no space, panics
//
// Placement validation:
//   - Requires sufficient contiguous free columns for colSpan
//   - Panics if cell cannot fit in available space
func (t *Table) AddCell(ct CellType, rowSpan, colSpan int) *Table {
	if rowSpan <= 0 || colSpan <= 0 {
		panic(fmt.Sprintf("tbl: invalid span rowSpan=%d colSpan=%d at.c.or (%d,%d)", rowSpan, colSpan, t.c.Row(), t.c.Col()))
	}

	if t.c.Row() == -1 {
		panic(fmt.Sprintf("tbl: no row to add cell at.c.or (%d,%d)", t.c.Row(), t.c.Col()))
	}

	// Step 1: Ensure columns exist
	if len(t.cols) == 0 {
		t.g.EnsureCols(colSpan)
		t.g.EnsureRows(rowSpan)
		t.addCell(ct, rowSpan, colSpan)
		return t
	}

	return t
}

// addCell places cell in grid, creates metadata, and advances cursor.
func (t *Table) addCell(typ CellType, rowSpan, colSpan int) {
	row, col := t.c.Pos()
	id := t.nextCellID
	c := NewCell(typ, row, col, rowSpan, colSpan)
	t.g.SetRect(row, col, rowSpan, colSpan)

	for i := row; i < row+rowSpan; i++ {
		for j := col; j < col+colSpan; j++ {
		}
	}

	t.cells[id] = c

	// Advance cursor
	t.c.Advance(colSpan)
}

// PrintDebug renders table structure in TBL Grid Notation format.
// Shows grid layout with cell types and.c.ent.c.or position.
// Returns empty string if table has no rows.
// For debug/development purposes.
func (t *Table) PrintDebug() string {
	return t.printDebug()
}
