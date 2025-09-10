package tbl

type HorizontalAlignment int

const (
	Left HorizontalAlignment = iota
	Center
	Right
)

type VerticalAlignment int

const (
	Top VerticalAlignment = iota
	Middle
	Bottom
)

type Alignment struct {
	Horizontal HorizontalAlignment
	Vertical   VerticalAlignment
}
