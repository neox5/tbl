package grid

// Cell is addressed row-major.
type Cell struct {
	Row int
	Col int
}

func (c Cell) Equal(o Cell) bool {
	return c.Row == o.Row && c.Col == o.Col
}
