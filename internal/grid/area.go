package grid

// Area is a rectangle of cells identified by index and span.
type Area struct {
	col     int
	row     int
	colSpan int
	rowSpan int
}

func NewArea(col, row, colSpan, rowSpan int) Area {
	if col < 0 || row < 0 || colSpan <= 0 || rowSpan <= 0 {
		panic("grid: invalid area")
	}
	return Area{col, row, colSpan, rowSpan}
}

func (a Area) Col() int     { return a.col }
func (a Area) Row() int     { return a.row }
func (a Area) ColSpan() int { return a.colSpan }
func (a Area) RowSpan() int { return a.rowSpan }
