package table

import (
	"github.com/neox5/tbl/internal/cell"
	"github.com/neox5/tbl/types"
)

// Table is the concrete table implementation
type Table struct {
	// Configuration (flattened)
	border   *types.TableBorder
	width    int
	maxWidth int

	// Table state
	cells     []*cell.Cell
	rowStarts []int
	colWidths []int
	colLevels []int
	hLines    []bool
	flexRows  map[int]bool
	flexCols  map[int]bool
	colsFixed bool // true when column structure is finalized
	nextIndex int  // next cell index
	row, col  int  // next row/col index

	// Indices for optimization (future use)
	rowIndex map[int][]int // cells overlapping a row
}

// New creates a new table with default configuration
func New() *Table {
	return &Table{
		border:    &types.DefaultTableBorder,
		width:     0,
		maxWidth:  0,
		cells:     []*cell.Cell{},
		rowStarts: []int{},
		colWidths: []int{},
		colLevels: []int{},
		flexRows:  make(map[int]bool),
		flexCols:  make(map[int]bool),
		colsFixed: false,
		nextIndex: 0,
		row:       -1, // First advanceRow() brings us to row 0
		col:       0,
		rowIndex:  make(map[int][]int),
	}
}

// NewWithConfig creates a new table with merged configuration
func NewWithConfig(cfg *types.Config) *Table {
	t := New()

	if cfg != nil {
		if cfg.Border != nil {
			t.border = cfg.Border
		}
		if cfg.Width > 0 {
			t.width = cfg.Width
		}
		if cfg.MaxWidth > 0 {
			t.maxWidth = cfg.MaxWidth
		}
	}

	return t
}

// ColCount returns the number of columns in the table
func (t *Table) ColCount() int {
	return len(t.colWidths)
}

// RowCount returns the number of rows in the Table
func (t *Table) RowCount() int {
	return len(t.rowStarts)
}

// addFlexRow adds a row to the flex rows map
func (t *Table) addFlexRow(row int) {
	t.flexRows[row] = true
}

// removeFlexRow removes a row from the flex rows map
func (t *Table) removeFlexRow(row int) {
	delete(t.flexRows, row)
}

// isFlexRow returns true if the row is in the flex rows map
func (t *Table) isFlexRow(row int) bool {
	return t.flexRows[row]
}
