package tbl

// VAlign implements Freestyler (direct field assignment).
type VAlign int

const (
	VAlignTop VAlign = iota
	VAlignMiddle
	VAlignBottom
)

// HAlign implements Freestyler (direct field assignment).
type HAlign int

const (
	HAlignLeft HAlign = iota
	HAlignCenter
	HAlignRight
)

// Style implements Freestyler (direct field assignment).
func (v VAlign) Style(base CellStyle) CellStyle {
	base.VAlign = v
	return base
}

// Style implements Freestyler (direct field assignment).
func (h HAlign) Style(base CellStyle) CellStyle {
	base.HAlign = h
	return base
}

// Vertical alignment constructors.
func Top() VAlign    { return VAlignTop }
func Middle() VAlign { return VAlignMiddle }
func Bottom() VAlign { return VAlignBottom }

// Horizontal alignment constructors.
func Left() HAlign   { return HAlignLeft }
func Center() HAlign { return HAlignCenter }
func Right() HAlign  { return HAlignRight }
