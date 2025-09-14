package cell

import (
	"regexp"

	"github.com/neox5/tbl/types"
)

const (
	COL_MIN_WIDTH = 1
)

// Cell is the concrete cell implementation
type Cell struct {
	content string
	border  types.CellBorder
	hAlign  types.HorizontalAlignment
	vAlign  types.VerticalAlignment
	col     Axis
	row     Axis
}

// New creates a new cell with default values
func New() *Cell {
	return &Cell{
		content: "",
		hAlign:  types.Left,
		vAlign:  types.Top,
		col:     Axis{Span: 1, Weight: 1},
		row:     Axis{Span: 1, Weight: 1},
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

// Content returns the cell content
func (c *Cell) Content() string {
	return c.content
}

// DisplayWidth calculates the display width of the cell content.
// When column is FLEX it needs to be fixed to report a width.
func (c *Cell) DisplayWidth() int {
	if c.col.IsFlex() && ! c.col.IsFixed() {
		return 0
	}
	return len(stripAnsiCodes(c.content))
}

// ColWidth calculates the column width of the cell.
// FLEX cells needs to be fixed/determined,
// otherwise it will report the COL_MIN_WIDTH
func (c *Cell) ColWidth() int {
	if c.col.IsFixed() {
		return c.col.End - c.col.Start
	}
	return COL_MIN_WIDTH
}

// Height calculates the display height of the cell content
func (c *Cell) Height() int {
	return 1
}

// Span return the column and row span
func (c *Cell) Span() (int, int) {
	return c.col.Span, c.row.Span
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

// stripAnsiCodes removes ANSI escape sequences for width calculation
func stripAnsiCodes(s string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(s, "")
}
