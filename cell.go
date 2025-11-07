package tbl

// CellType indicates whether a cell is static or flexible.
type CellType int

const (
	// Static cells have fixed column spans.
	Static CellType = iota

	// Flex cells can expand to fill available space.
	Flex
)

// Cell represents a table cell with position and span information.
type Cell struct {
	id           ID
	typ          CellType
	r, c         int
	rSpan, cSpan int
	initialSpan  int // Original colSpan at creation
}

// NewCell creates a new cell with specified properties.
func NewCell(id ID, typ CellType, r, c, rSpan, cSpan int) *Cell {
	return &Cell{
		id:          id,
		typ:         typ,
		r:           r,
		c:           c,
		rSpan:       rSpan,
		cSpan:       cSpan,
		initialSpan: cSpan,
	}
}

// Contains checks if cell contains the position (row, col).
func (c *Cell) Contains(row, col int) bool {
	return row >= c.r && row < c.r+c.rSpan &&
		col >= c.c && col < c.c+c.cSpan
}

// TouchesRow checks if cell spans through the specified row.
func (c *Cell) TouchesRow(row int) bool {
	return row >= c.r && row < c.r+c.rSpan
}

// AddedSpan returns total expansion applied to this cell.
func (c *Cell) AddedSpan() int {
	return c.cSpan - c.initialSpan
}
