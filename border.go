package tbl

// Border style definitions
type BorderStyle int

const (
	Single BorderStyle = iota
	Double
	Heavy
	Rounded
	ASCII
)

// Default border character sets
var DefaultBorderChars = map[BorderStyle]BorderChars{
	Single: {
		Horizontal: '─', Vertical: '│',
		TopLeft: '┌', TopRight: '┐', BottomLeft: '└', BottomRight: '┘',
		Cross: '┼', TUp: '┴', TDown: '┬', TLeft: '┤', TRight: '├',
	},
	Double: {
		Horizontal: '═', Vertical: '║',
		TopLeft: '╔', TopRight: '╗', BottomLeft: '╚', BottomRight: '╝',
		Cross: '╬', TUp: '╩', TDown: '╦', TLeft: '╣', TRight: '╠',
	},
	Heavy: {
		Horizontal: '━', Vertical: '┃',
		TopLeft: '┏', TopRight: '┓', BottomLeft: '┗', BottomRight: '┛',
		Cross: '╋', TUp: '┻', TDown: '┳', TLeft: '┫', TRight: '┣',
	},
	Rounded: {
		Horizontal: '─', Vertical: '│',
		TopLeft: '╭', TopRight: '╮', BottomLeft: '╰', BottomRight: '╯',
		Cross: '┼', TUp: '┴', TDown: '┬', TLeft: '┤', TRight: '├',
	},
	ASCII: {
		Horizontal: '-', Vertical: '|',
		TopLeft: '+', TopRight: '+', BottomLeft: '+', BottomRight: '+',
		Cross: '+', TUp: '+', TDown: '+', TLeft: '+', TRight: '+',
	},
}

type BorderChars struct {
	// Edges
	Horizontal, Vertical rune
	// Vertices (intersections)
	TopLeft, TopRight, BottomLeft, BottomRight rune
	Cross, TUp, TDown, TLeft, TRight           rune
}

type TableBorder struct {
	All    bool
	Around bool
	Style  BorderStyle
}

type RowBorder struct {
	ColSeparation bool
	RowSeparation bool
	Style         BorderStyle
}

type CellBorder struct {
	Top, Right, Bottom, Left bool
	Style                    BorderStyle
}
