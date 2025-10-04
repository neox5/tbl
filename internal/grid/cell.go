package grid

// Cell identifies a single discrete position in the grid.
type Cell struct {
	Col int
	Row int
}

// Equal returns true if other is the same Cell.
func (c Cell) Equal(o Cell) bool {
	return c.Col == o.Col && c.Row == o.Row
}
