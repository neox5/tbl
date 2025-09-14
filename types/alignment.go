package types

// HorizontalAlignment represents horizontal text alignment within a cell
type HorizontalAlignment int

const (
	Left HorizontalAlignment = iota
	Center
	Right
)

// VerticalAlignment represents vertical text alignment within a cell
type VerticalAlignment int

const (
	Top VerticalAlignment = iota
	Middle
	Bottom
)

// Alignment combines horizontal and vertical alignment
type Alignment struct {
	Horizontal HorizontalAlignment
	Vertical   VerticalAlignment
}
