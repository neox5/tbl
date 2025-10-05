package tbl

import (
	"fmt"

	"github.com/neox5/tbl/internal/grid"
)

// ID identifies a cell in the table.
type ID grid.ID

// CellType indicates whether a cell is static or flexible.
type CellType int

const (
	// Static cells have fixed column spans.
	Static CellType = iota

	// Flex cells can expand to fill available space.
	Flex
)

// Table manages incremental table construction with flex/static cells.
type Table struct {
	g         *grid.Grid
	cells     map[ID]*cell
	cur       cursor
	colWidth  map[int]int
	rowHeight map[int]int
}

// New creates a new Table with zero columns.
func New() *Table {
	return NewWithCols(0)
}

// NewWithCols creates a new Table with initial column capacity.
func NewWithCols(cols int) *Table {
	return &Table{
		g:         grid.New(0, cols),
		cells:     make(map[ID]*cell),
		cur:       cursor{row: -1, col: 0},
		colWidth:  make(map[int]int),
		rowHeight: make(map[int]int),
	}
}

// AddRow advances to the next row for cell placement.
func (t *Table) AddRow() *Table {
	t.addRow()
	return t
}

// AddCell adds a cell at the current cursor position.
// rowSpan and colSpan define cell dimensions (must be > 0).
// Returns *Table for chaining. Panics on invalid input.
func (t *Table) AddCell(ct CellType, rowSpan, colSpan int) *Table {
	if rowSpan <= 0 || colSpan <= 0 {
		panic(fmt.Errorf("tbl: invalid span rowSpan=%d colSpan=%d", rowSpan, colSpan))
	}

	if t.cur.row == 0 {
		t.addCols(colSpan)
	}

	t.addCell(ct, rowSpan, colSpan)

	return t
}
