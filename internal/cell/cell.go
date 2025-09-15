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
	col     Axis
	row     Axis
}

// New creates a new cell with default values
func New() *Cell {
	return &Cell{
		content: "",
		hAlign:  types.Left,
		vAlign:  types.Top,
		col:     NewAxis(1),
		row:     NewAxis(1),
	}
}

// NewColFlex creates a new cell with flexible column dimension
// opts: minSpan, maxSpan, weight
func NewColFlex(opts ...int) *Cell {
	if len(opts) > 3 {
		panic("NewColFlex: too many options (max 3: minSpan, maxSpan, weight)")
	}
	return New().WithColFlex(opts...)
}

// NewRowFlex creates a new cell with flexible row dimension
// opts: minSpan, maxSpan, weight
func NewRowFlex(opts ...int) *Cell {
	if len(opts) > 3 {
		panic("NewRowFlex: too many options (max 3: minSpan, maxSpan, weight)")
	}
	return New().WithRowFlex(opts...)
}

// NewFlex creates a new cell with flexible column and row dimensions
// opts: colMin, colMax, colWeight, rowMin, rowMax, rowWeight
func NewFlex(opts ...int) *Cell {
	if len(opts) > 6 {
		panic("NewFlex: too many options (max 6: colMin, colMax, colWeight, rowMin, rowMax, rowWeight)")
	}

	colOpts := opts[:3]
	rowOpts := []int{}
	if len(opts) > 3 {
		rowOpts = opts[3:]
	}

	return New().WithColFlex(colOpts...).WithRowFlex(rowOpts...)
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

// WithColSpan sets the column span
func (c *Cell) WithColSpan(span int) *Cell {
	c.col = NewAxis(span)
	return c
}

// WithRowSpan sets the row span
func (c *Cell) WithRowSpan(span int) *Cell {
	c.row = NewAxis(span)
	return c
}

// WithSpan sets the column and row span
func (c *Cell) WithSpan(colSpan, rowSpan int) *Cell {
	c.col = NewAxis(colSpan)
	c.row = NewAxis(rowSpan)
	return c
}

// WithColFlex sets flexible column span
// opts: minSpan, maxSpan, weight
func (c *Cell) WithColFlex(opts ...int) *Cell {
	minSpan, maxSpan, weight := 1, AxisNoCap, 1

	if len(opts) > 0 {
		minSpan = opts[0]
	}
	if len(opts) > 1 {
		maxSpan = opts[1]
	}
	if len(opts) > 2 {
		weight = opts[2]
	}

	c.col = NewFlexAxis(minSpan, maxSpan, weight)
	return c
}

// WithRowFlex sets flexible row span
// opts: minSpan, maxSpan, weight
func (c *Cell) WithRowFlex(opts ...int) *Cell {
	minSpan, maxSpan, weight := 1, AxisNoCap, 1

	if len(opts) > 0 {
		minSpan = opts[0]
	}
	if len(opts) > 1 {
		maxSpan = opts[1]
	}
	if len(opts) > 2 {
		weight = opts[2]
	}

	c.row = NewFlexAxis(minSpan, maxSpan, weight)
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

// Width returns the display width of the cell content
func (c *Cell) Width() int {
	return len(stripAnsiCodes(c.content))
}

// Height calculates the display height of the cell content
func (c *Cell) Height() int {
	return 1
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
