package grid

import "fmt"

// Area is a rectangle of cells. Fields are private and mutable via methods.
type Area struct {
	col     int
	row     int
	colSpan int
	rowSpan int
}

func NewArea(col, row, colSpan, rowSpan int) *Area {
	if col < 0 || row < 0 || colSpan <= 0 || rowSpan <= 0 {
		panic(fmt.Errorf("grid: invalid area (col=%d,row=%d,colSpan=%d,rowSpan=%d)", col, row, colSpan, rowSpan))
	}
	return &Area{col: col, row: row, colSpan: colSpan, rowSpan: rowSpan}
}

// Accessors
func (a Area) Col() int     { return a.col }
func (a Area) Row() int     { return a.row }
func (a Area) ColSpan() int { return a.colSpan }
func (a Area) RowSpan() int { return a.rowSpan }
func (a Area) ColEnd() int  { return a.col + a.colSpan } // exclusive
func (a Area) RowEnd() int  { return a.row + a.rowSpan } // exclusive

// Mutators
func (a *Area) MoveTo(col, row int) {
	if col < 0 || row < 0 {
		panic(fmt.Errorf("grid: MoveTo negative index col=%d row=%d", col, row))
	}
	a.col, a.row = col, row
}

func (a *Area) MoveBy(dCol, dRow int) {
	nc, nr := a.col+dCol, a.row+dRow
	if nc < 0 || nr < 0 {
		panic(fmt.Errorf("grid: MoveBy would go negative col=%d row=%d", nc, nr))
	}
	a.col, a.row = nc, nr
}

func (a *Area) Resize(colSpan, rowSpan int) {
	if colSpan <= 0 || rowSpan <= 0 {
		panic(fmt.Errorf("grid: Resize non-positive span colSpan=%d rowSpan=%d", colSpan, rowSpan))
	}
	a.colSpan, a.rowSpan = colSpan, rowSpan
}

// ForRows applies fn to each row in [Row, RowEnd).
func (a Area) ForRows(fn func(row int)) {
	for r := a.row; r < a.RowEnd(); r++ {
		fn(r)
	}
}

// ForRowsWithError applies fn and stops on first error.
func (a Area) ForRowsWithError(fn func(row int) error) error {
	for r := a.row; r < a.RowEnd(); r++ {
		if err := fn(r); err != nil {
			return err
		}
	}
	return nil
}

func (a Area) ForCols(do func(col int)) {
	for c := a.col; c < a.ColEnd(); c++ {
		do(c)
	}
}

func (a Area) ForEachCell(do func(col, row int)) {
	for r := a.row; r < a.RowEnd(); r++ {
		for c := a.col; c < a.ColEnd(); c++ {
			do(c, r)
		}
	}
}
