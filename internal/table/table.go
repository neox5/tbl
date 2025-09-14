package table

import (
	"fmt"

	"github.com/neox5/tbl/internal/cell"
	"github.com/neox5/tbl/types"
)

// Table is the concrete table implementation
type Table struct {
	// Configuration (flattened)
	border      *types.TableBorder
	width       int
	maxWidth    int
	newCellFunc func() *cell.Cell

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
		newCellFunc:   cell.New,
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
		if cfg.NewCellFunc != nil {
			// Convert public interface func to internal cell func
			t.newCellFunc = func() *cell.Cell {
				if c, ok := cfg.NewCellFunc().(*cell.Cell); ok {
					return c
				}
				return cell.New()
			}
		}
	}

	return t
}

// AddRow adds a new row with the specified cells
func (t *Table) AddRow(cells ...any) {
	t.newRow()
	t.addCells(cells...)
}

// R is a short form of AddRow
func (t *Table) R(cells ...any) {
	t.AddRow(cells...)
}

// NewCell creates a new cell with the specified value
func (t *Table) NewCell(value any) *cell.Cell {
	return t.newCell(value)
}

// C is a short form of NewCell
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
		return t.newCellFunc().WithContent(fmt.Sprintf("%v", v))
	}
}

// addCells appends multiple cells to the current row
func (t *Table) addCells(cells ...any) {
	for _, c := range cells {
		internalCell := t.newCell(c)
		t.cells = append(t.cells, internalCell)
	}
}
