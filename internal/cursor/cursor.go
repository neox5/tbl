// Package cursor implements a simple row col reference
package cursor

// Cursor tracks current building position in the table.
type Cursor struct {
	row int
	col int
}

func New() *Cursor {
	return &Cursor{row: -1, col: 0}
}

func (c *Cursor) Row() int {
	return c.row
}

func (c *Cursor) Col() int {
	return c.col
}

// Pos returns current row and column.
func (c *Cursor) Pos() (row, col int) {
	return c.row, c.col
}

// Advance moves cursor forward by colSpan columns.
func (c *Cursor) Advance(colSpan int) {
	if colSpan <= 0 {
		return
	}
	c.col += colSpan
}

// NextRow moves cursor to next row, resets column to 0, returns new row.
func (c *Cursor) NextRow() int {
	c.row++
	c.col = 0
	return c.row
}
