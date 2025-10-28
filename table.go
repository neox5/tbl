package tbl

import (
	"fmt"

	"github.com/neox5/tbl/internal/cursor"
	"github.com/neox5/tbl/internal/grid"
)

// ID identifies a cell in the table.
type ID grid.ID

// Table manages incremental table construction with flex/static cells.
type Table struct {
	g         *grid.Grid
	cur       *cursor.Cursor
	rows      [][]ID
	cols      [][]ID
	colWidths []int
	cells     map[ID]*cell

	colsFixed bool
	colLevels []int
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
		g:         grid.New(0, cols),
		cur:       cursor.New(),
		rows:      make([][]ID, 0),
		cols:      make([][]ID, cols),
		colWidths: make([]int, cols),
		cells:     make(map[ID]*cell),
		colLevels: make([]int, cols),
	}

	return t
}

// AddRow advances to the next row for cell placement.
func (t *Table) AddRow() *Table {
	t.addRow()
	return t
}

// AddCell adds a cell at the current cursor position.
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
		panic(fmt.Sprintf("tbl: invalid span rowSpan=%d colSpan=%d at cursor (%d,%d)", rowSpan, colSpan, t.cur.Row(), t.cur.Col()))
	}

	if t.cur.Row() == -1 {
		panic(fmt.Sprintf("tbl: no row to add cell at cursor (%d,%d)", t.cur.Row(), t.cur.Col()))
	}

	// Step 1: Ensure columns exist
	if len(t.cols) == 0 {
		t.addCols(colSpan)
	} else {
		// Step 2: Ensure column space available
		t.ensureColumnSpace(colSpan)

		// Step 3: Validate span fits
		t.validateSpanFit(colSpan)
	}

	t.addCell(ct, rowSpan, colSpan)

	return t
}

// ensureColumnSpace verifies column availability and expands if needed.
// Panics if expansion not possible and no free columns exist.
func (t *Table) ensureColumnSpace(colSpan int) {
	_, _, ok := t.nextZeroRun(t.cur.Col())
	if !ok {
		// No free columns found
		if t.colsFixed {
			panic("tbl: no free columns and column expansion disabled")
		}
		t.addCols(colSpan)
	}
}

// validateSpanFit checks if colSpan fits in next available free column run.
// Panics if insufficient contiguous space exists.
func (t *Table) validateSpanFit(colSpan int) {
	pos, span, ok := t.nextZeroRun(t.cur.Col())
	if !ok {
		panic("tbl: no free column run found")
	}

	if colSpan > span {
		panic(fmt.Sprintf("tbl: cell span %d exceeds available space %d at column %d",
			colSpan, span, pos))
	}
}

// PrintDebug renders table structure in TBL Grid Notation format.
// Shows grid layout with cell types and current cursor position.
// Returns empty string if table has no rows.
// For debug/development purposes.
func (t *Table) PrintDebug() string {
	return t.printDebug()
}

// nextZeroRun returns (pos, count, ok) for the next run of 0s in colLevels.
func (t *Table) nextZeroRun(from int) (pos, count int, ok bool) {
	if from < 0 {
		from = 0
	}
	n := len(t.colLevels)
	if n == 0 || from >= n {
		return -1, 0, false
	}

	i := from
	for i < n && t.colLevels[i] != 0 {
		i++
	}
	if i >= n {
		return -1, 0, false
	}

	j := i
	for j < n && t.colLevels[j] == 0 {
		j++
	}
	return i, j - i, true
}
