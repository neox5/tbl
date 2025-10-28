package tbl

// CellType indicates whether a cell is static or flexible.
type CellType int

const (
	// Static cells have fixed column spans.
	Static CellType = iota

	// Flex cells can expand to fill available space.
	Flex
)

type Cell struct {
	typ          CellType
	r, c         int
	rSpan, cSpan int

	// cache
	// rowCache []int // cached row indices
	// colCache []int // cached col indices
}

func NewCell(typ CellType, r, c, rSpan, cSpan int) *Cell {
	return &Cell{
		typ:   typ,
		r:     r,
		c:     c,
		rSpan: rSpan,
		cSpan: cSpan,
	}
}
