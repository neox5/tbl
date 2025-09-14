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
	cells         []*cell.Cell
	rowStarts     []int
	colWidths     []int
	colLevels     []int
	hLines        []bool
	currIndex     int
	openFlexCells []int

	// Indices for optimization (future use)
	colIndex map[int][]int // cells overlapping a column
	rowIndex map[int][]int // cells overlapping a row
}

// New creates a new table with default configuration
func New() *Table {
	return &Table{
		border:        &types.DefaultTableBorder,
		width:         0,
		maxWidth:      0,
		cells:         []*cell.Cell{},
		rowStarts:     []int{},
		colWidths:     []int{},
		colLevels:     []int{},
		openFlexCells: []int{},
		colIndex:      make(map[int][]int),
		rowIndex:      make(map[int][]int),
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

// AddRow adds a new row with the specified cells
func (t *Table) AddRow(cells ...*cell.Cell) {
	t.startNewRow()
	t.addCells(cells)
}

// startNewRow starts a new row by recording the current cell position
func (t *Table) startNewRow() {
	t.rowStarts = append(t.rowStarts, len(t.cells))
}

// addCells appends multiple cells to the current row
func (t *Table) addCells(cells []*cell.Cell) {
	t.cells = append(t.cells, cells...)
}
