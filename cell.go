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

// CellSpec specifies cell parameters for deferred construction.
type CellSpec struct {
	typ          CellType
	rSpan, cSpan int
	content      string
	stylers      []Freestyler
}

// C creates Static cell spec with content and optional styles.
// Default span: [1,1]
//
// Example:
//
//	t.AddRow(tbl.C("Name"), tbl.C("Age", tbl.Right()))
func C(content string, stylers ...Freestyler) CellSpec {
	return CellSpec{
		typ:     Static,
		rSpan:   1,
		cSpan:   1,
		content: content,
		stylers: stylers,
	}
}

// F creates Flex cell spec with content and optional styles.
// Default span: [1,1]
//
// Example:
//
//	t.AddRow(tbl.F("Bio"), tbl.C("Age"))
func F(content string, stylers ...Freestyler) CellSpec {
	return CellSpec{
		typ:     Flex,
		rSpan:   1,
		cSpan:   1,
		content: content,
		stylers: stylers,
	}
}

// Cx creates Static cell spec with custom span, content and optional styles.
//
// Example:
//
//	t.AddRow(tbl.Cx(2, 1, "Merged Cell"), tbl.C("Normal"))
func Cx(rSpan, cSpan int, content string, stylers ...Freestyler) CellSpec {
	return CellSpec{
		typ:     Static,
		rSpan:   rSpan,
		cSpan:   cSpan,
		content: content,
		stylers: stylers,
	}
}

// Fx creates Flex cell spec with custom span, content and optional styles.
//
// Example:
//
//	t.AddRow(tbl.Fx(1, 2, "Wide Flex"), tbl.C("Normal"))
func Fx(rSpan, cSpan int, content string, stylers ...Freestyler) CellSpec {
	return CellSpec{
		typ:     Flex,
		rSpan:   rSpan,
		cSpan:   cSpan,
		content: content,
		stylers: stylers,
	}
}

// Cell represents a table cell with position, span and content information.
type Cell struct {
	id           ID
	typ          CellType
	r, c         int
	rSpan, cSpan int
	initialSpan  int // original colSpan at creation
	content      string
	rawLines     []string // unconstrained content lines
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

// Width returns the required character width of the cell content (unconstrained).
func (c *Cell) Width() int {
	var width int
	for _, l := range c.rawLines {
		if len(l) > width {
			width = len(l)
		}
	}
	return width
}

// Height returns the required lines of the cell content (unconstrained).
func (c *Cell) Height() int {
	return len(c.rawLines)
}
