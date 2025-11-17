package tbl

// BorderSide defines which edges of a cell have borders.
type BorderSide uint8

const (
	BorderNone   BorderSide = 0
	BorderTop    BorderSide = 1 << 0 // 0001
	BorderRight  BorderSide = 1 << 1 // 0010
	BorderBottom BorderSide = 1 << 2 // 0100
	BorderLeft   BorderSide = 1 << 3 // 1000

	BorderAll = BorderTop | BorderRight | BorderBottom | BorderLeft
)

// Border specifies which edges of a cell should render borders.
type Border struct {
	Sides    BorderSide // Which edges render visually (characters)
	Physical BorderSide // Which edges occupy physical space
}

// HasTop reports whether top edge has visual or physical border.
func (b Border) HasTop() bool {
	return b.IsVisualTop() || (b.Physical&BorderTop) != 0
}

// HasBottom reports whether bottom edge has visual or physical border.
func (b Border) HasBottom() bool {
	return b.IsVisualBottom() || (b.Physical&BorderBottom) != 0
}

// HasLeft reports whether left edge has visual or physical border.
func (b Border) HasLeft() bool {
	return b.IsVisualLeft() || (b.Physical&BorderLeft) != 0
}

// HasRight reports whether right edge has visual or physical border.
func (b Border) HasRight() bool {
	return b.IsVisualRight() || (b.Physical&BorderRight) != 0
}

// IsVisualTop reports whether top edge renders as character.
func (b Border) IsVisualTop() bool {
	return (b.Sides & BorderTop) != 0
}

// IsVisualBottom reports whether bottom edge renders as character.
func (b Border) IsVisualBottom() bool {
	return (b.Sides & BorderBottom) != 0
}

// IsVisualLeft reports whether left edge renders as character.
func (b Border) IsVisualLeft() bool {
	return (b.Sides & BorderLeft) != 0
}

// IsVisualRight reports whether right edge renders as character.
func (b Border) IsVisualRight() bool {
	return (b.Sides & BorderRight) != 0
}

// HAlign specifies horizontal text alignment within a cell.
type HAlign int

const (
	HAlignLeft HAlign = iota
	HAlignCenter
	HAlignRight
)

// VAlign specifies vertical text alignment within a cell.
type VAlign int

const (
	VAlignTop VAlign = iota
	VAlignMiddle
	VAlignBottom
)

// Padding specifies space around cell content.
type Padding struct {
	Top, Bottom, Left, Right int
}

// CellStyle contains presentation attributes for a cell.
type CellStyle struct {
	Padding Padding
	HAlign  HAlign
	VAlign  VAlign
	Border  Border
}

// merge applies non-zero overrides from other style to this style.
// Returns new CellStyle with merged values.
// Zero values in override are ignored (base value preserved).
func (s CellStyle) merge(other CellStyle) CellStyle {
	result := s

	// Merge padding (0 means use base value)
	if other.Padding.Top != 0 {
		result.Padding.Top = other.Padding.Top
	}
	if other.Padding.Bottom != 0 {
		result.Padding.Bottom = other.Padding.Bottom
	}
	if other.Padding.Left != 0 {
		result.Padding.Left = other.Padding.Left
	}
	if other.Padding.Right != 0 {
		result.Padding.Right = other.Padding.Right
	}

	// Alignment always override
	result.HAlign = other.HAlign
	result.VAlign = other.VAlign

	// Border always override
	result.Border = other.Border

	return result
}

// resolveStyle returns effective style for cell using hierarchy.
// Resolution order: defaultStyle < columnStyles < rowStyles < cellStyles
// Only considers cell origin position (cell.r, cell.c) for multi-span cells.
func (t *Table) resolveStyle(cell *Cell) CellStyle {
	style := t.defaultStyle

	// Column style at origin
	if cs, ok := t.columnStyles[cell.c]; ok {
		style = style.merge(cs)
	}

	// Row style at origin
	if rs, ok := t.rowStyles[cell.r]; ok {
		style = style.merge(rs)
	}

	// Cell-specific style (highest priority)
	if cs, ok := t.cellStyles[cell.id]; ok {
		style = style.merge(cs)
	}

	return style
}
