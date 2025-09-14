package table

import (
	"github.com/neox5/tbl/internal/cell"
	"github.com/neox5/tbl/types"
)

// Table is the concrete table implementation
type Table struct {
	// Configuration
	config *types.Config

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
		config:        Default(),
		cells:         []*cell.Cell{},
		rowStarts:     []int{},
		colWidths:     []int{},
		colLevels:     []int{},
		openFlexCells: []int{},
		colIndex:      make(map[int][]int),
		rowIndex:      make(map[int][]int),
	}
}

// NewWithConfig creates a new table with the specified configuration
func NewWithConfig(cfg types.Config) *Table {
	t := New()
	t.config = &cfg
	return t
}

// AddRow adds a new row with the specified cells
func (t *Table) AddRow(cells ...any) {
	t.R(cells...)
}

// R is a short form of AddRow
func (t *Table) R(cells ...any) {
	t.newRow()
	t.addCells(cells...)
}

// C creates a new cell with the specified value
func (t *Table) C(value any) *cell.Cell {
	return t.newCell(value)
}

// newRow starts a new row by recording the current cell position
func (t *Table) newRow() {
	t.rowStarts = append(t.rowStarts, len(t.cells))
}

// newCell creates a new cell from any value type
func (t *Table) newCell(value any) *cell.Cell {
	switch v := value.(type) {
	case *cell.Cell:
		return v
	default:
		// Use default cell configuration and set content
		defaultCell := GetCellDefault(t.config)
		return defaultCell.WithContent(cell.NewFromValue(v).Content())
	}
}

// addCells appends multiple cells to the current row
func (t *Table) addCells(cells ...any) {
	for _, c := range cells {
		internalCell := t.newCell(c)
		t.cells = append(t.cells, internalCell)
	}
}
