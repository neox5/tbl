package cell

import (
	"regexp"

	"github.com/neox5/tbl/types"
)

// Cell is the concrete cell implementation
type Cell struct {
	content string
	border  types.CellBorder
	hAlign  types.HorizontalAlignment
	vAlign  types.VerticalAlignment
	col     CellAxis
	row     CellAxis
}

// New creates a new cell with default values
func New() *Cell {
	return &Cell{
		content: "",
		hAlign:  types.Left,
		vAlign:  types.Top,
		col:     CellAxis{Span: 1, Weight: 1},
		row:     CellAxis{Span: 1, Weight: 1},
	}
}

// WithContent sets the cell content
func (c *Cell) WithContent(content string) *Cell {
	c.content = content
	return c
}

// WithAlign sets the horizontal and vertical alignment
func (c *Cell) WithAlign(h types.HorizontalAlignment, v types.VerticalAlignment) *Cell {
	c.hAlign = h
	c.vAlign = v
	return c
}

// WithSpan sets the column and row span
func (c *Cell) WithSpan(col, row int) *Cell {
	c.col.Span = col
	c.row.Span = row
	return c
}

// WithBorder sets the cell border
func (c *Cell) WithBorder(border types.CellBorder) *Cell {
	c.border = border
	return c
}

// Short form aliases
func (c *Cell) C(content string) *Cell {
	return c.WithContent(content)
}

func (c *Cell) A(h types.HorizontalAlignment, v types.VerticalAlignment) *Cell {
	return c.WithAlign(h, v)
}

func (c *Cell) S(col, row int) *Cell {
	return c.WithSpan(col, row)
}

func (c *Cell) B(border types.CellBorder) *Cell {
	return c.WithBorder(border)
}

// Internal accessor methods for processing

// Width calculates the display width of the cell content
func (c *Cell) Width() int {
	return len(stripAnsiCodes(c.content))
}

// Height calculates the display height of the cell content
func (c *Cell) Height() int {
	return 1
}

// Content returns the cell content
func (c *Cell) Content() string {
	return c.content
}

// HAlign returns the horizontal alignment
func (c *Cell) HAlign() types.HorizontalAlignment {
	return c.hAlign
}

// VAlign returns the vertical alignment
func (c *Cell) VAlign() types.VerticalAlignment {
	return c.vAlign
}

// Border returns the cell border
func (c *Cell) Border() types.CellBorder {
	return c.border
}

// Col returns the column axis
func (c *Cell) Col() CellAxis {
	return c.col
}

// Row returns the row axis
func (c *Cell) Row() CellAxis {
	return c.row
}

// stripAnsiCodes removes ANSI escape sequences for width calculation
func stripAnsiCodes(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(s, "")
}
