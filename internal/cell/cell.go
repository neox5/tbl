package cell

import (
	"fmt"
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

// DefaultCell creates a new cell with default values
func DefaultCell() *Cell {
	return &Cell{
		content: "",
		hAlign:  types.Left,
		vAlign:  types.Top,
		col:     CellAxis{Span: 1, Weight: 1},
		row:     CellAxis{Span: 1, Weight: 1},
	}
}

// WithContent sets the cell content and returns a new cell
func (c *Cell) WithContent(content string) *Cell {
	newCell := *c
	newCell.content = content
	return &newCell
}

// WithAlign sets the horizontal and vertical alignment
func (c *Cell) WithAlign(h types.HorizontalAlignment, v types.VerticalAlignment) *Cell {
	newCell := *c
	newCell.hAlign = h
	newCell.vAlign = v
	return &newCell
}

// WithSpan sets the column and row span
func (c *Cell) WithSpan(col, row int) *Cell {
	newCell := *c
	newCell.col.Span = col
	newCell.row.Span = row
	return &newCell
}

// WithBorder sets the cell border
func (c *Cell) WithBorder(border types.CellBorder) *Cell {
	newCell := *c
	newCell.border = border
	return &newCell
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

// NewFromValue creates a new cell from any value type
func NewFromValue(value any) *Cell {
	switch v := value.(type) {
	case string:
		return DefaultCell().WithContent(v)
	case *Cell:
		return v
	default:
		return DefaultCell().WithContent(fmt.Sprintf("%v", v))
	}
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
