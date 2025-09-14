package tbl

import (
	"github.com/neox5/tbl/internal/cell"
	"github.com/neox5/tbl/types"
)

// Cell is the public cell wrapper
type Cell struct {
	cell *cell.Cell
}

// NewCell creates a standalone cell with the specified value
func NewCell(value any) *Cell {
	return &Cell{
		cell: cell.NewFromValue(value),
	}
}

// C is a short form of NewCell
func C(value any) *Cell {
	return NewCell(value)
}

// WithContent sets the cell content
func (c *Cell) WithContent(content string) *Cell {
	c.cell.WithContent(content)
	return c
}

// WithAlign sets the horizontal and vertical alignment
func (c *Cell) WithAlign(h types.HorizontalAlignment, v types.VerticalAlignment) *Cell {
	c.cell.WithAlign(h, v)
	return c
}

// WithSpan sets the column and row span
func (c *Cell) WithSpan(col, row int) *Cell {
	c.cell.WithSpan(col, row)
	return c
}

// WithBorder sets the cell border
func (c *Cell) WithBorder(border types.CellBorder) *Cell {
	c.cell.WithBorder(border)
	return c
}

// Short form aliases

// C sets the cell content (short form)
func (c *Cell) C(content string) *Cell {
	return c.WithContent(content)
}

// A sets the horizontal and vertical alignment (short form)
func (c *Cell) A(h types.HorizontalAlignment, v types.VerticalAlignment) *Cell {
	return c.WithAlign(h, v)
}

// S sets the column and row span (short form)
func (c *Cell) S(col, row int) *Cell {
	return c.WithSpan(col, row)
}

// B sets the cell border (short form)
func (c *Cell) B(border types.CellBorder) *Cell {
	return c.WithBorder(border)
}
