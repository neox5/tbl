package tbl

import (
	"fmt"

	"github.com/neox5/btmp"
	"github.com/neox5/tbl/internal/cursor"
)

// ID identifies a cell in the table.
type ID int64

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

// AddCell adds a cell at the current cursor position.
// rowSpan and colSpan define cell dimensions (must be > 0).
// Returns *Table for chaining. Panics on invalid input.
//
// Column expansion behavior:
//   - If colsFixed is false, expands grid when needed
//   - If colsFixed is true and no space, panics
//
// Placement validation:
//   - Requires sufficient contiguous free columns for colSpan
//   - Panics if cell cannot fit in available space
func (t *Table) AddCell(ct CellType, rowSpan, colSpan int) *Table {
	if rowSpan <= 0 || colSpan <= 0 {
		panic(fmt.Sprintf("tbl: invalid span rowSpan=%d colSpan=%d at cursor (%d,%d)", rowSpan, colSpan, t.c.Row(), t.c.Col()))
	}

	if t.c.Row() == -1 {
		panic(fmt.Sprintf("tbl: no row to add cell at cursor (%d,%d)", t.c.Row(), t.c.Col()))
	}

	row, col := t.c.Pos()

	// Expand columns if needed and not fixed
	if !t.colsFixed {
		needed := col + colSpan
		if needed > t.g.Cols() {
			delta := needed - t.g.Cols()
			t.g.GrowCols(delta)
			for range delta {
				t.cols = append(t.cols, []ID{})
			}
		}
	}

	// Expand rows if needed
	needed := row + rowSpan
	if needed > t.g.Rows() {
		delta := needed - t.g.Rows()
		t.g.GrowRows(delta)
		for range delta {
			t.rows = append(t.rows, []ID{})
		}
	}

	// Validate space available
	if col+colSpan > t.g.Cols() {
		panic(fmt.Sprintf("tbl: insufficient columns for cell colSpan=%d at cursor (%d,%d) cols=%d", colSpan, row, col, t.g.Cols()))
	}

	// TODO: Check if space is free using t.g.IsFree

	// Create cell
	id := t.nextCellID
	t.nextCellID++
	c := NewCell(ct, row, col, rowSpan, colSpan)

	// Set in grid
	t.g.SetRect(row, col, rowSpan, colSpan)

	// Index cell by rows
	for i := row; i < row+rowSpan; i++ {
		t.rows[i] = append(t.rows[i], id)
	}

	// Index cell by cols
	for j := col; j < col+colSpan; j++ {
		t.cols[j] = append(t.cols[j], id)
	}

	// Store cell
	t.cells[id] = c

	// Advance cursor
	t.c.Advance(colSpan)

	return t
}

// PrintDebug renders table structure in TBL Grid Notation format.
// Shows grid layout with cell types and current cursor position.
// Returns empty string if table has no rows.
// For debug/development purposes.
func (t *Table) PrintDebug() string {
	return t.printDebug()
}
