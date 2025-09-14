package types

// BorderStyle represents different border drawing styles
type BorderStyle int

const (
	None BorderStyle = iota
	Single
	Double
	Heavy
	Rounded
	ASCII
)

// BorderChars defines the characters used for drawing borders
type BorderChars struct {
	// Edges
	Horizontal, Vertical rune
	// Vertices (intersections)
	TopLeft, TopRight, BottomLeft, BottomRight rune
	Cross, TUp, TDown, TLeft, TRight           rune
}

// DefaultBorderChars provides default character sets for each border style
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

// DefaultTableBorder provides sensible table border defaults
var DefaultTableBorder = TableBorder{
	All:    false,
	Around: false,
	Style:  Single,
}

// TableBorder configures borders for the entire table
type TableBorder struct {
	All    bool        // Show all borders (internal and external)
	Around bool        // Show only outer border
	Style  BorderStyle // Border style to use
}

// RowBorder configures borders for table rows
type RowBorder struct {
	Top           bool        // Show border above row
	ColSeparation bool        // Show vertical separators between columns
	Bottom        bool        // Show border below row
	Style         BorderStyle // Border style to use
}

// CellBorder configures borders for individual cells
type CellBorder struct {
	Top, Right, Bottom, Left bool        // Which sides to show borders on
	Style                    BorderStyle // Border style to use
}
