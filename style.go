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
	Sides BorderSide
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
