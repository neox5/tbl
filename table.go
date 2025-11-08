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
		cells:      make(map[ID]*Cell),
		nextCellID: 1,
	}

	if cols > 0 {
		t.colsFixed = true
	}

	return t
}

// AddRow advances to next row with validation and cursor positioning.
func (t *Table) AddRow() *Table {
	row := t.c.Row()

	// Validate previous row if not first row
	if row >= 0 {
		if !t.isRowComplete(row) {
			panic(fmt.Sprintf("tbl: incomplete row %d before AddRow", row))
		}

		// Check if we can fix columns
		if !t.colsFixed && t.isRowStatic(row) {
			t.colsFixed = true
		}
	}

	// Ensure next row exists
	t.ensureRows(row + 1)

	// Advance cursor and get new row
	row = t.c.NextRow()

	// Position cursor at first free column
	freeCol := t.findFirstFreeCol(row)
	t.c.Advance(freeCol)

	return t
}

// AddCell adds a cell at cursor position with specified type and span.
// Expands columns if needed (when not fixed). Validates span fits in grid.
// Panics if: no row active, span invalid, insufficient columns (when fixed),
// or space occupied.
func (t *Table) AddCell(ct CellType, rowSpan, colSpan int) *Table {
	if rowSpan <= 0 || colSpan <= 0 {
		panic(fmt.Sprintf("tbl: invalid span rowSpan=%d colSpan=%d at cursor (%d,%d)", rowSpan, colSpan, t.c.Row(), t.c.Col()))
	}

	if t.c.Row() == -1 {
		panic(fmt.Sprintf("tbl: no row to add cell at cursor (%d,%d)", t.c.Row(), t.c.Col()))
	}

	row, col := t.c.Pos()

	// Ensure sufficient rows for cell span
	t.ensureRows(row + rowSpan)

	// Simple check on first row
	if row == 0 && col+colSpan > t.g.Cols() {
		if err := t.ensureCols(col + colSpan); err != nil {
			panic(err.Error())
		}
	}

	// when cell does not fit in the current position, we try to expand
	if colSpan > t.g.CountZerosFromInRow(row, col) {
		// Space occupied - wall blocking
		if t.colsFixed {
			panic(fmt.Sprintf("tbl: space occupied at cursor (%d,%d), cannot expand", row, col))
		}

		// Attempt expansion
		ok, flexCells := t.traverseFlex(row, col)
		if !ok {
			panic(fmt.Sprintf("tbl: no flex cells available for expansion at cursor (%d,%d)", row, col))
		}

		// Calculate needed columns
		needed := t.calculateNeeded(row, col, colSpan)

		// Add columns to grid
		t.g.GrowCols(needed)

		// Process rows top to bottom
		for r := 0; r <= row; r++ {
			if rowFlexCells, exists := flexCells[r]; exists && len(rowFlexCells) > 0 {
				t.distributeAndExpand(r, rowFlexCells, needed)
			}
		}
	}

	// Create cell
	id := t.nextCellID
	t.nextCellID++
	c := NewCell(id, ct, row, col, rowSpan, colSpan)

	// Store cell
	t.cells[id] = c

	// Set in grid
	t.g.SetRect(row, col, rowSpan, colSpan)

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
