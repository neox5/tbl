package tbl

import (
	"github.com/neox5/tbl/internal/cell"
	"github.com/neox5/tbl/internal/table"
	"github.com/neox5/tbl/types"
)

// Tbl is the public table wrapper
type Tbl struct {
	table *table.Table
}

// New creates a new table with default configuration
func New() *Tbl {
	return &Tbl{
		table: table.New(),
	}
}

// NewWithConfig creates a new table with the specified configuration
func NewWithConfig(cfg *types.Config) *Tbl {
	return &Tbl{
		table: table.NewWithConfig(cfg),
	}
}

// AddRow adds a new row with the specified cells
func (t *Tbl) AddRow(cells ...*Cell) *Tbl {
	// Convert public cells to internal cells
	internalCells := make([]*cell.Cell, len(cells))
	for i, c := range cells {
		internalCells[i] = c.cell
	}

	t.table.AddRow(internalCells...)
	return t
}

// R is a short form of AddRow
func (t *Tbl) R(cells ...*Cell) *Tbl {
	return t.AddRow(cells...)
}

// NewCell creates a new cell with the specified value
func (t *Tbl) NewCell(value any) *Cell {
	return &Cell{
		cell: t.table.NewCell(value),
	}
}

// C is a short form of NewCell
func (t *Tbl) C(value any) *Cell {
	return t.NewCell(value)
}

// Render renders the table to string
func (t *Tbl) Render() string {
	return t.table.Render()
}
