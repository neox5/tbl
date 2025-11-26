package tbl

// Freestyler applies styling to a CellStyle.
//
// Tribute to the Bomfunk MC's Classic: https://youtu.be/ymNFyxvIdaM
type Freestyler interface {
	Style(base CellStyle) CellStyle
}

// Funcstyler computes a CellStyle for a given cell position.
// The returned style is merged into the resolved style for the cell.
type Funcstyler func(row, col int) CellStyle

// WrapMode controls how content overflow is handled.
type WrapMode int

const (
	WrapWord     WrapMode = iota // wrap at word boundaries (default)
	WrapChar                     // wrap at any character
	WrapTruncate                 // truncate with ellipsis
)

// Style implements Freestyler (direct field assignment).
func (w WrapMode) Style(base CellStyle) CellStyle {
	base.WrapMode = w
	return base
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
