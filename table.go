package tbl

import (
	"fmt"
	"io"

	"github.com/neox5/btmp"
)

// ID identifies a cell in the table.
type ID int64

// Table manages incremental table construction with flex/static cells.
type Table struct {
	g     *btmp.Grid
	cells map[ID]*Cell

	// Cursor state
	row int // current row position, -1 before first AddRow
	col int // current column position within row

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
		g:          btmp.NewGridWithSize(1, max(cols, 1)),
		cells:      make(map[ID]*Cell),
		row:        -1,
		col:        0,
		nextCellID: 1,
	}

	if cols > 0 {
		t.colsFixed = true
	}

	return t
}

// AddRow advances to next row with validation and cursor positioning.
func (t *Table) AddRow() *Table {
	// Validate previous row if not first row
	if t.row >= 0 {
		if !t.isRowComplete(t.row) {
			panic(fmt.Sprintf("tbl: incomplete row %d before AddRow", t.row))
		}

		// Check if we can fix columns
		if !t.colsFixed && t.isRowStatic(t.row) {
			t.colsFixed = true
		}
	}

	// Ensure next row exists
	t.ensureRows(t.row + 2) // t.row = current row index; +1 current row count; +2 next row count

	// Advance to next row
	t.nextRow()

	return t
}

// AddCell adds a cell at cursor position with specified type and span.
// Expands columns if needed (when not fixed). Validates span fits in grid.
// Panics if: no row active, span invalid, insufficient columns (when fixed),
// or space occupied.
func (t *Table) AddCell(ct CellType, rowSpan, colSpan int, content string) *Table {
	if rowSpan <= 0 || colSpan <= 0 {
		panic(fmt.Sprintf("tbl: invalid span rowSpan=%d colSpan=%d at cursor (%d,%d)", rowSpan, colSpan, t.row, t.col))
	}

	if t.row == -1 {
		panic(fmt.Sprintf("tbl: no row to add cell at cursor (%d,%d)", t.row, t.col))
	}

	// Ensure sufficient rows for cell span
	t.ensureRows(t.row + rowSpan)

	// Simple check on first row
	if t.row == 0 && t.col+colSpan > t.g.Cols() {
		if err := t.ensureCols(t.col + colSpan); err != nil {
			panic(err.Error())
		}
	}

	// when cell does not fit in the current position, we try to expand
	if colSpan > t.g.CountZerosFromInRow(t.row, t.col) {
		// Space occupied - wall blocking
		if t.colsFixed {
			panic(fmt.Sprintf("tbl: space occupied at cursor (%d,%d), cannot expand", t.row, t.col))
		}

		// Attempt expansion
		ok, flexCells := t.traverseFlex(t.row, t.col)
		if !ok {
			panic(fmt.Sprintf("tbl: no flex cells available for expansion at cursor (%d,%d)", t.row, t.col))
		}

		// Calculate needed columns
		needed := t.calculateNeeded(t.row, t.col, colSpan)

		// Add columns to grid
		t.g.GrowCols(needed)

		// Process rows top to bottom
		for r := 0; r <= t.row; r++ {
			if rowFlexCells, exists := flexCells[r]; exists && len(rowFlexCells) > 0 {
				t.distributeAndExpand(r, rowFlexCells, needed)
			}
		}
	}

	if t.g.B.Test(t.g.Index(t.row, t.col)) {
		t.col = t.findFirstFreeCol(t.row)
	}

	// Create cell
	id := t.nextCellID
	t.nextCellID++

	// Store cell
	c := NewCell(id, ct, t.row, t.col, rowSpan, colSpan, content)
	t.cells[id] = c

	// Set in grid
	t.g.SetRect(t.row, t.col, rowSpan, colSpan)
	t.advance()
	return t
}

// Render returns the psql-style ASCII table as a string.
func (t *Table) Render() string {
	return newRenderer(t).render()
}

// RenderTo writes the table to w.
func (t *Table) RenderTo(w io.Writer) error {
	s := newRenderer(t).render()
	_, err := io.WriteString(w, s)
	return err
}

// PrintDebug renders table structure in TBL Grid Notation format.
// Shows grid layout with cell types and current cursor position.
// Returns empty string if table has no rows.
// For debug/development purposes.
func (t *Table) PrintDebug() string {
	return t.printDebug()
}
