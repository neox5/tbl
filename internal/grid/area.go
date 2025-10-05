package grid

import "fmt"

// Area is a rectangle of cells. Fields are private and mutable via methods.
// The cells, rows, and cols slices are maintained automatically on mutation.
type Area struct {
	row     int
	col     int
	rowSpan int
	colSpan int
	cells   []Cell // cached cell positions
	rows    []int  // cached row indices
	cols    []int  // cached column indices
}

// NewArea creates a new area and initializes cache.
func NewArea(row, col, rowSpan, colSpan int) *Area {
	if row < 0 || col < 0 || rowSpan <= 0 || colSpan <= 0 {
		panic(fmt.Errorf("grid: invalid area (row=%d,col=%d,rowSpan=%d,colSpan=%d)", row, col, rowSpan, colSpan))
	}
	a := &Area{row: row, col: col, rowSpan: rowSpan, colSpan: colSpan}
	a.updateCache()
	return a
}

// Accessors
func (a Area) Row() int      { return a.row }
func (a Area) Col() int      { return a.col }
func (a Area) RowSpan() int  { return a.rowSpan }
func (a Area) ColSpan() int  { return a.colSpan }
func (a Area) RowEnd() int   { return a.row + a.rowSpan } // exclusive
func (a Area) ColEnd() int   { return a.col + a.colSpan } // exclusive
func (a Area) Cells() []Cell { return a.cells }
func (a Area) Rows() []int   { return a.rows }
func (a Area) Cols() []int   { return a.cols }

// Mutators
func (a *Area) MoveTo(row, col int) {
	if row < 0 || col < 0 {
		panic(fmt.Errorf("grid: MoveTo negative index row=%d col=%d", row, col))
	}
	a.row, a.col = row, col
	a.updateCache()
}

func (a *Area) MoveBy(dRow, dCol int) {
	nr, nc := a.row+dRow, a.col+dCol
	if nr < 0 || nc < 0 {
		panic(fmt.Errorf("grid: MoveBy would go negative row=%d col=%d", nr, nc))
	}
	a.row, a.col = nr, nc
	a.updateCache()
}

func (a *Area) Resize(rowSpan, colSpan int) {
	if rowSpan <= 0 || colSpan <= 0 {
		panic(fmt.Errorf("grid: Resize non-positive span rowSpan=%d colSpan=%d", rowSpan, colSpan))
	}
	a.rowSpan, a.colSpan = rowSpan, colSpan
	a.updateCache()
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

// ForCols applies fn to each col in [Col, ColEnd).
func (a Area) ForCols(fn func(col int)) {
	for c := a.col; c < a.ColEnd(); c++ {
		fn(c)
	}
}

// ForEachCell iterates row-major and passes (row, col).
func (a Area) ForEachCell(do func(row, col int)) {
	for r := a.row; r < a.RowEnd(); r++ {
		for c := a.col; c < a.ColEnd(); c++ {
			do(r, c)
		}
	}
}

// updateCache rebuilds all cached slices based on current area bounds.
// Called automatically after any mutation.
func (a *Area) updateCache() {
	// Update cells
	size := a.rowSpan * a.colSpan
	a.cells = make([]Cell, size)
	idx := 0
	for r := a.row; r < a.RowEnd(); r++ {
		for c := a.col; c < a.ColEnd(); c++ {
			a.cells[idx] = Cell{Col: c, Row: r}
			idx++
		}
	}

	// Update rows
	a.rows = make([]int, a.rowSpan)
	for i := range a.rowSpan {
		a.rows[i] = a.row + i
	}

	// Update cols
	a.cols = make([]int, a.colSpan)
	for i := range a.colSpan {
		a.cols[i] = a.col + i
	}
}
