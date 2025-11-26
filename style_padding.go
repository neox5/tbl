package tbl

import "fmt"

// Padding specifies space around cell content.
type Padding struct {
	Top, Bottom, Left, Right int
}

// Style implements Freestyler (direct field assignment).
func (p Padding) Style(base CellStyle) CellStyle {
	base.Padding = p
	return base
}

// Pad creates a Padding from 1-4 values (CSS-like behavior).
// All values must be non-negative.
// Panics if count is 0, >4, or any value is negative.
//
// Usage:
//
//	Pad(a)          -> all sides = a
//	Pad(a, b)       -> top/bottom = a, left/right = b
//	Pad(a, b, c)    -> top = a, left/right = b, bottom = c
//	Pad(a, b, c, d) -> top = a, right = b, bottom = c, left = d
func Pad(values ...int) Padding {
	// Validate count
	if len(values) == 0 {
		panic("tbl: Pad requires at least 1 value")
	}
	if len(values) > 4 {
		panic(fmt.Sprintf("tbl: Pad accepts 1-4 values, got %d", len(values)))
	}

	// Validate all values are non-negative
	for i, v := range values {
		if v < 0 {
			panic(fmt.Sprintf("tbl: Pad value at index %d is negative: %d", i, v))
		}
	}

	switch len(values) {
	case 1:
		// All sides
		return Padding{
			Top:    values[0],
			Bottom: values[0],
			Left:   values[0],
			Right:  values[0],
		}
	case 2:
		// Vertical, Horizontal
		return Padding{
			Top:    values[0],
			Bottom: values[0],
			Left:   values[1],
			Right:  values[1],
		}
	case 3:
		// Top, Horizontal, Bottom
		return Padding{
			Top:    values[0],
			Bottom: values[2],
			Left:   values[1],
			Right:  values[1],
		}
	case 4:
		// Top, Right, Bottom, Left (clockwise)
		return Padding{
			Top:    values[0],
			Right:  values[1],
			Bottom: values[2],
			Left:   values[3],
		}
	default:
		// Unreachable due to validation above
		panic(fmt.Sprintf("tbl: Pad accepts 1-4 values, got %d", len(values)))
	}
}
