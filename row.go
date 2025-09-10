package tbl

var DefaultRow = Row{
	Cells:  []Cell{},
	Border: RowBorder{},
	Height: 1,
}

type Row struct {
	Cells  []Cell
	Border RowBorder
	Height int // height in lines
}

func (r Row) WithBorder(border RowBorder) Row {
	r.Border = border
	return r
}

func (r Row) B(border RowBorder) Row {
	return r.WithBorder(border)
}
