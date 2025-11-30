package tbl

// FontStyle represents text styling attributes.
type FontStyle int

const (
	FontNormal FontStyle = 0
	FontBold   FontStyle = 1 << iota
	FontDim
	FontItalic
	FontUnderline
	FontBlink
	FontReverse
	FontStrikethrough
)

// Style implements Freestyler for FontStyle.
func (f FontStyle) Style(base CellStyle) CellStyle {
	base.FontStyle = f
	return base
}

// Font style constructors

func Bold() FontStyle          { return FontBold }
func Dim() FontStyle           { return FontDim }
func Italic() FontStyle        { return FontItalic }
func Underline() FontStyle     { return FontUnderline }
func Blink() FontStyle         { return FontBlink }
func Reverse() FontStyle       { return FontReverse }
func Strikethrough() FontStyle { return FontStrikethrough }
