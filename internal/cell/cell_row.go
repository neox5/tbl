package cell

// RowSpan returns the row span
func (c *Cell) RowSpan() int {
	return c.row.Span()
}

// IsRowFlex returns true if the row uses flexible spanning
func (c *Cell) IsRowFlex() bool {
	return c.row.IsFlex()
}

// IsRowPositioned returns true if the row has been positioned
func (c *Cell) IsRowPositioned() bool {
	return c.row.IsPositioned()
}

// RowStart returns the row starting position
func (c *Cell) RowStart() int {
	return c.row.Start()
}

// SetRowStart sets the row starting position
func (c *Cell) SetRowStart(start int) {
	c.row.SetStart(start)
}
