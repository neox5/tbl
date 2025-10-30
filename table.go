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

// isRowComplete validates row has no holes and all columns filled.
// Returns true if entire row range [0, Cols) is occupied.
func (t *Table) isRowComplete(row int) bool {
	if row < 0 || row >= t.g.Rows() {
		return false
	}
	if t.g.Cols() == 0 {
		return true
	}
	return t.g.AllRow(row)
}

// hasOnlyStaticCells checks if all cells in row are Static type.
// Returns true only if row exists and contains exclusively Static cells.
func (t *Table) hasOnlyStaticCells(row int) bool {
	if row < 0 || row >= len(t.rows) {
		return false
	}
	if len(t.rows[row]) == 0 {
		return false
	}
	for _, id := range t.rows[row] {
		cell := t.cells[id]
		if cell == nil || cell.typ != Static {
			return false
		}
	}
	return true
}

// findFirstFreeCol locates first unoccupied column in row.
// Returns column index of first free position, or Cols() if row full.
// Accounts for cells with rowSpan > 1 from previous rows.
func (t *Table) findFirstFreeCol(row int) int {
	if row < 0 || row >= t.g.Rows() {
		return 0
	}
	return t.g.NextFreeCol(row, 0)
}

// AddRow advances to next row with validation and cursor positioning.
// Validates previous row completeness if not first row.
// Sets colsFixed=true if row 0 complete and contains only Static cells.
// Positions cursor at first free column in new row.
// Panics if previous row incomplete or has holes.
func (t *Table) AddRow() *Table {
	prevRow := t.c.Row()

	// Validate previous row if not first row
	if prevRow >= 0 {
		if !t.isRowComplete(prevRow) {
			panic(fmt.Sprintf("tbl: incomplete row %d before AddRow", prevRow))
		}

		// Check if we can fix columns
		if !t.colsFixed && t.hasOnlyStaticCells(prevRow) {
			t.colsFixed = true
		}
	}

	t.c.NextRow()
	t.rows = append(t.rows, []ID{})

	// Position cursor at first free column in new row
	row := t.c.Row()
	if row < t.g.Rows() {
		freeCol := t.findFirstFreeCol(row)
		// Update cursor column directly
		t.c.Advance(freeCol)
	}

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

	// Check dimensional feasibility
	fitRow, fitCol := t.g.CanFit(row, col, rowSpan, colSpan)

	// Handle column constraint
	if !fitCol {
		if t.colsFixed {
			panic(fmt.Sprintf("tbl: insufficient columns for cell colSpan=%d at cursor (%d,%d) cols=%d", colSpan, row, col, t.g.Cols()))
		}
		// Expand columns if not fixed
		t.g.EnsureCols(col + colSpan)
		for i := len(t.cols); i < col+colSpan; i++ {
			t.cols = append(t.cols, []ID{})
		}
	}

	// Expand rows if needed
	if !fitRow {
		t.g.EnsureRows(row + rowSpan)
		for i := len(t.rows); i < row+rowSpan; i++ {
			t.rows = append(t.rows, []ID{})
		}
	}

	// Validate space available - now guaranteed to be in bounds
	if !t.g.IsFree(row, col, rowSpan, colSpan) {
		panic(fmt.Sprintf("tbl: space not free for cell at cursor (%d,%d) span=(%d,%d)", row, col, rowSpan, colSpan))
	}

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
