package tbl

const (
	FLEX = -1
)

type Cell struct {
	Content string
	ColSpan int
	RowSpan int
	Border  CellBorder
	HAlign  HorizontalAlignment
	VAlign  VerticalAlignment
}
