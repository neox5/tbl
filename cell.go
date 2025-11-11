package tbl

import (
	"math"
	"strings"
)

// CellType indicates whether a cell is static or flexible.
type CellType int

const (
	Static CellType = iota
	Flex
)

// HAlign specifies horizontal text alignment within a cell.
type HAlign int

const (
	HAlignLeft HAlign = iota
	HAlignCenter
	HAlignRight
)

// VAlign specifies vertical text alignment within a cell.
type VAlign int

const (
	VAlignTop VAlign = iota
	VAlignMiddle
	VAlignBottom
)

// Cell represents a table cell with position, span and content information.
type Cell struct {
	id           ID
	typ          CellType
	r, c         int
	rSpan, cSpan int
	initialSpan  int // original colSpan at creation
	content      string
	hAlign       HAlign   // horizontal alignment (not used yet)
	vAlign       VAlign   // vertical alignment (not used yet)
	rawLines     []string // unconstraint content lines
}

// NewCell creates a new cell.
func NewCell(id ID, typ CellType, r, c, rSpan, cSpan int, content string) *Cell {
	cell := &Cell{
		id:          id,
		typ:         typ,
		r:           r,
		c:           c,
		rSpan:       rSpan,
		cSpan:       cSpan,
		initialSpan: cSpan,
		content:     strings.TrimSpace(content),
		hAlign:      HAlignLeft, // default horizontal alignment
		vAlign:      VAlignTop,  // default vertical alignment
	}

	cell.rawLines = buildRawLines(cell.content, math.MaxInt)

	return cell
}

// Contains reports whether the cell covers the given grid position.
func (c *Cell) Contains(row, col int) bool {
	return row >= c.r && row < c.r+c.rSpan &&
		col >= c.c && col < c.c+c.cSpan
}

// TouchesRow reports whether the cell overlaps the given row.
func (c *Cell) TouchesRow(row int) bool {
	return row >= c.r && row < c.r+c.rSpan
}

// AddedSpan returns how many columns were added by flex expansion.
func (c *Cell) AddedSpan() int {
	return c.cSpan - c.initialSpan
}

// Content returns the cell text.
func (c *Cell) Content() string { return c.content }

// Width returns the required character width of the cell content (unconstraint).
func (c *Cell) Width() int {
	var width int
	for _, l := range c.rawLines {
		if len(l) > width {
			width = len(l)
		}
	}
	return width
}

// Height returns the required lines of the cell content (unconstraint).
func (c *Cell) Height() int {
	return len(c.rawLines)
}

// HAlign returns the horizontal alignment of the cell.
func (c *Cell) HAlign() HAlign { return c.hAlign }

// VAlign returns the vertical alignment of the cell.
func (c *Cell) VAlign() VAlign { return c.vAlign }
