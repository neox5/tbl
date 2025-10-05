package tbl

// cursor tracks current building position in the table.
type cursor struct {
	row int
	col int
}

// advance moves cursor forward by colSpan columns.
func (c *cursor) advance(colSpan int) {
	c.col += colSpan
}

// nextRow moves cursor to next row, resets column to 0.
func (c *cursor) nextRow() {
	c.row++
	c.col = 0
}

// position returns current row and column.
func (c *cursor) position() (row, col int) {
	return c.row, c.col
}
