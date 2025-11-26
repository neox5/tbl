package tbl

// BorderSide controls which sides of a border are active.
type BorderSide int

const (
	BorderNone BorderSide = 0
	BorderTop  BorderSide = 1 << iota
	BorderBottom
	BorderLeft
	BorderRight

	BorderAll = BorderTop | BorderBottom | BorderLeft | BorderRight
)

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

// Style implements Freestyler (direct field assignment).
func (b Border) Style(base CellStyle) CellStyle {
	base.Border = b
	return base
}

// Border constructors.
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

// Common border combinations.
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
