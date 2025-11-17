package tbl

import "strings"

// charTemplate defines the character set for rendering.
type charTemplate struct {
	cornerTL rune
	cornerTR rune
	cornerBL rune
	cornerBR rune
	cornerT  rune
	cornerB  rune
	cornerL  rune
	cornerR  rune
	cornerX  rune
	hLine    rune
	vLine    rune
}

// ASCII template (simple)
var asciiTemplate = charTemplate{
	cornerTL: '+',
	cornerTR: '+',
	cornerBL: '+',
	cornerBR: '+',
	cornerT:  '+',
	cornerB:  '+',
	cornerL:  '+',
	cornerR:  '+',
	cornerX:  '+',
	hLine:    '-',
	vLine:    '|',
}

// Unicode thin box-drawing
var thinTemplate = charTemplate{
	cornerTL: '┌',
	cornerTR: '┐',
	cornerBL: '└',
	cornerBR: '┘',
	cornerT:  '┬',
	cornerB:  '┴',
	cornerL:  '├',
	cornerR:  '┤',
	cornerX:  '┼',
	hLine:    '─',
	vLine:    '│',
}

// Unicode thick box-drawing
var thickTemplate = charTemplate{
	cornerTL: '┏',
	cornerTR: '┓',
	cornerBL: '┗',
	cornerBR: '┛',
	cornerT:  '┳',
	cornerB:  '┻',
	cornerL:  '┣',
	cornerR:  '┫',
	cornerX:  '╋',
	hLine:    '━',
	vLine:    '┃',
}

// Unicode double-line box
var doubleTemplate = charTemplate{
	cornerTL: '╔',
	cornerTR: '╗',
	cornerBL: '╚',
	cornerBR: '╝',
	cornerT:  '╦',
	cornerB:  '╩',
	cornerL:  '╠',
	cornerR:  '╣',
	cornerX:  '╬',
	hLine:    '═',
	vLine:    '║',
}

// Select the default
var defaultTemplate = thinTemplate

// writeBorder interprets border instruction sequence using template.
func (r *renderer) writeBorder(b *strings.Builder, ops []RenderOp) {
	tpl := asciiTemplate

	for _, op := range ops {
		switch v := op.(type) {
		case CornerTL:
			b.WriteRune(tpl.cornerTL)
		case CornerTR:
			b.WriteRune(tpl.cornerTR)
		case CornerBL:
			b.WriteRune(tpl.cornerBL)
		case CornerBR:
			b.WriteRune(tpl.cornerBR)
		case CornerT:
			b.WriteRune(tpl.cornerT)
		case CornerB:
			b.WriteRune(tpl.cornerB)
		case CornerL:
			b.WriteRune(tpl.cornerL)
		case CornerR:
			b.WriteRune(tpl.cornerR)
		case CornerX:
			b.WriteRune(tpl.cornerX)
		case HLine:
			for i := 0; i < v.Width; i++ {
				b.WriteRune(tpl.hLine)
			}
		}
	}
	b.WriteByte('\n')
}

// writeContent interprets content instruction sequence using template.
func (r *renderer) writeContent(b *strings.Builder, ops []RenderOp) {
	tpl := defaultTemplate

	for _, op := range ops {
		switch v := op.(type) {
		case VLine:
			b.WriteRune(tpl.vLine)
		case Content:
			b.WriteString(v.Text) // just write finalized text
		case Space:
			for i := 0; i < v.Width; i++ {
				b.WriteByte(' ')
			}
		}
	}
	b.WriteByte('\n')
}
