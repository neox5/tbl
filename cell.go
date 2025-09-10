package tbl

const (
	FLEX = -1
)

var DefaultCell = Cell{
	Content: "",
	ColSpan: 1,
	RowSpan: 1,
	HAlign:  Left,
	VAlign:  Top,
}

type Cell struct {
	Content string
	ColSpan int
	RowSpan int
	Border  CellBorder
	HAlign  HorizontalAlignment
	VAlign  VerticalAlignment
}

func (c Cell) WithAlign(h HorizontalAlignment, v VerticalAlignment) Cell {
	c.HAlign = h
	c.VAlign = v
	return c
}

func (c Cell) WithBorder(border CellBorder) Cell {
	c.Border = border
	return c
}

func (c Cell) WithContent(content string) Cell {
	c.Content = content
	return c
}

func (c Cell) WithSpan(col, row int) Cell {
	c.ColSpan = col
	c.RowSpan = row
	return c
}

func (c Cell) A(h HorizontalAlignment, v VerticalAlignment) Cell {
	return c.WithAlign(h,v)
}

func (c Cell) B(border CellBorder) Cell {
	return c.WithBorder(border)
}

func (c Cell) C(content string) Cell {
	return c.WithContent(content)
}

func (c Cell) S(col, row int) Cell {
	return c.WithSpan(col,row)
}

