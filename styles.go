package tbl

import "fmt"

// Freestyler applies styling to a CellStyle.
//
// Tribute to the Bomfunk MC's Classic: https://youtu.be/ymNFyxvIdaM
type Freestyler interface {
	Style(base CellStyle) CellStyle
}

// WrapMode controls how content overflow is handled.
type WrapMode int

const (
	WrapWord     WrapMode = iota // wrap at word boundaries (default)
	WrapChar                     // wrap at any character
	WrapTruncate                 // truncate with ellipsis
)

// Padding specifies space around cell content.
type Padding struct {
	Top, Bottom, Left, Right int
}

// Border specifies which edges of a cell should render borders.
type Border struct {
	Sides    BorderSide // Which edges render visually (characters)
	Physical BorderSide // Which edges occupy physical space
}

// Has reports whether border side occupies space (visual or physical).
func (b Border) Has(side BorderSide) bool {
	return b.IsVisual(side) || (b.Physical&side) != 0
}

// IsVisual reports whether border side renders as character.
func (b Border) IsVisual(side BorderSide) bool {
	return (b.Sides & side) != 0
}

// CellStyle contains presentation attributes for a cell.
type CellStyle struct {
	Padding  Padding
	HAlign   HAlign
	VAlign   VAlign
	Border   Border
	WrapMode WrapMode
	Template CharTemplate
}

// NewStyle creates a CellStyle from stylers.
func NewStyle(stylers ...Freestyler) CellStyle {
	result := CellStyle{}
	for _, s := range stylers {
		if s != nil {
			result = s.Style(result)
		}
	}
	return result
}

// Apply applies stylers to this CellStyle, returning a new CellStyle.
func (s CellStyle) Apply(stylers ...Freestyler) CellStyle {
	result := s
	for _, styler := range stylers {
		if styler != nil {
			result = styler.Style(result)
		}
	}
	return result
}

// Style implements Freestyler for CellStyle (uses merge).
func (s CellStyle) Style(base CellStyle) CellStyle {
	return base.merge(s)
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

	// Alignment: only override if explicitly set (non-zero)
	if other.HAlign != 0 {
		result.HAlign = other.HAlign
	}
	if other.VAlign != 0 {
		result.VAlign = other.VAlign
	}

	// Border: only override if explicitly set (non-zero)
	if other.Border.Sides != 0 || other.Border.Physical != 0 {
		result.Border = other.Border
	}

	// WrapMode: only override if explicitly set (non-zero)
	if other.WrapMode != 0 {
		result.WrapMode = other.WrapMode
	}

	// Template: only override if explicitly set (non-zero runes)
	if other.Template != (CharTemplate{}) {
		result.Template = other.Template
	}

	return result
}

// HAlign implements Freestyler (direct field assignment).
func (h HAlign) Style(base CellStyle) CellStyle {
	base.HAlign = h
	return base
}

// VAlign implements Freestyler (direct field assignment).
func (v VAlign) Style(base CellStyle) CellStyle {
	base.VAlign = v
	return base
}

// Padding implements Freestyler (direct field assignment).
func (p Padding) Style(base CellStyle) CellStyle {
	base.Padding = p
	return base
}

// Border implements Freestyler (direct field assignment).
func (b Border) Style(base CellStyle) CellStyle {
	base.Border = b
	return base
}

// WrapMode implements Freestyler (direct field assignment).
func (w WrapMode) Style(base CellStyle) CellStyle {
	base.WrapMode = w
	return base
}

// Horizontal alignment constructors
func Left() HAlign   { return HAlignLeft }
func Center() HAlign { return HAlignCenter }
func Right() HAlign  { return HAlignRight }

// Vertical alignment constructors
func Top() VAlign    { return VAlignTop }
func Middle() VAlign { return VAlignMiddle }
func Bottom() VAlign { return VAlignBottom }

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

// Border constructors
func BLeft() Border {
	return Border{Sides: BorderLeft}
}

func BRight() Border {
	return Border{Sides: BorderRight}
}

func BTop() Border {
	return Border{Sides: BorderTop}
}

func BBottom() Border {
	return Border{Sides: BorderBottom}
}

func BAll() Border {
	return Border{Sides: BorderAll}
}

func BNone() Border {
	return Border{Sides: BorderNone}
}

// Common border combinations
func BTopBottom() Border {
	return Border{Sides: BorderTop | BorderBottom}
}

func BLeftRight() Border {
	return Border{Sides: BorderLeft | BorderRight}
}

// Borders creates a Border with custom BorderSide combination.
func Borders(sides BorderSide) Border {
	return Border{Sides: sides}
}

// containsTemplate checks if any styler is a CharTemplate.
func containsTemplate(stylers []Freestyler) bool {
	for _, s := range stylers {
		if _, ok := s.(CharTemplate); ok {
			return true
		}
	}
	return false
}

// SetDefaultStyle sets the base style for all cells.
// Can be overridden by column, row, or cell-specific styles.
func (t *Table) SetDefaultStyle(stylers ...Freestyler) *Table {
	t.defaultStyle = t.defaultStyle.Apply(stylers...)
	return t
}

// SetColStyle sets style for all cells in column.
// Overrides default style, can be overridden by row or cell styles.
// Panics if stylers contain CharTemplate (templates are table-level only).
func (t *Table) SetColStyle(col int, stylers ...Freestyler) *Table {
	if containsTemplate(stylers) {
		panic("tbl: CharTemplate only supported via SetDefaultStyle")
	}
	t.columnStyles[col] = t.columnStyles[col].Apply(stylers...)
	return t
}

// SetRowStyle sets style for all cells in row.
// Overrides default and column styles, can be overridden by cell styles.
// Panics if stylers contain CharTemplate (templates are table-level only).
func (t *Table) SetRowStyle(row int, stylers ...Freestyler) *Table {
	if containsTemplate(stylers) {
		panic("tbl: CharTemplate only supported via SetDefaultStyle")
	}
	t.rowStyles[row] = t.rowStyles[row].Apply(stylers...)
	return t
}

// SetCellStyle sets style for specific cell.
// Highest priority, overrides all other styles.
// Panics if stylers contain CharTemplate (templates are table-level only).
func (t *Table) SetCellStyle(id ID, stylers ...Freestyler) *Table {
	if containsTemplate(stylers) {
		panic("tbl: CharTemplate only supported via SetDefaultStyle")
	}
	t.cellStyles[id] = t.cellStyles[id].Apply(stylers...)
	return t
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
