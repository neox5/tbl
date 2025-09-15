package cell

// ColSpan returns the column span
func (c *Cell) ColSpan() int {
	return c.col.Span()
}

// IsColFlex returns true if the column uses flexible spanning
func (c *Cell) IsColFlex() bool {
	return c.col.IsFlex()
}

// ColMaxSpan returns the maximum column span
func (c *Cell) ColMaxSpan() int {
	return c.col.MaxSpan()
}

// ColWeight returns the column weight for flex distribution
func (c *Cell) ColWeight() int {
	return c.col.Weight()
}

// IsColPositioned returns true if the column has been positioned
func (c *Cell) IsColPositioned() bool {
	return c.col.IsPositioned()
}

// ColStart returns the column starting position
func (c *Cell) ColStart() int {
	return c.col.Start()
}

// ColEnd returns the column ending position (start + span)
func (c *Cell) ColEnd() int {
	return c.col.End()
}

// SetColStart sets the column starting position
func (c *Cell) SetColStart(start int) {
	c.col.SetStart(start)
}

// ColCanGrowTo returns true if flex column can grow to the given span
func (c *Cell) ColCanGrowTo(span int) bool {
	return c.col.CanGrowTo(span)
}

// SetColSpan sets the current column span (used by table for flex resolution)
func (c *Cell) SetColSpan(span int) {
	c.col.SetSpan(span)
}
