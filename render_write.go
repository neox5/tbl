package tbl

import "strings"

// RenderOp is a rendering instruction.
type RenderOp interface {
	renderOp() // unexported marker method
}

// Border instructions (implement RenderOp)
type (
	CornerTL struct{}
	CornerTR struct{}
	CornerBL struct{}
	CornerBR struct{}
	CornerT  struct{} // top junction (─┬─)
	CornerB  struct{} // bottom junction (─┴─)
	CornerL  struct{} // left junction (├─)
	CornerR  struct{} // right junction (─┤)
	CornerX  struct{} // cross junction (─┼─)
	HLine    struct{ Width int }
)

// Content instructions
type (
	VLine   struct{}
	Content struct{ Text string } // finalized line (includes padding, alignment)
	Space   struct{ Width int }
)

// Implement marker
func (CornerTL) renderOp() {}
func (CornerTR) renderOp() {}
func (CornerBL) renderOp() {}
func (CornerBR) renderOp() {}
func (CornerT) renderOp()  {}
func (CornerB) renderOp()  {}
func (CornerL) renderOp()  {}
func (CornerR) renderOp()  {}
func (CornerX) renderOp()  {}
func (HLine) renderOp()    {}
func (VLine) renderOp()    {}
func (Content) renderOp()  {}
func (Space) renderOp()    {}

// CharTemplate defines character set for table rendering.
// Template is a table-level property applied via SetDefaultStyle.
// Attempting to use templates with SetRowStyle, SetColStyle, or
// SetCellStyle will panic.
type CharTemplate struct {
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

// Thin returns thin Unicode box-drawing template.
func Thin() CharTemplate {
	return CharTemplate{
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
}

// Thick returns thick Unicode box-drawing template.
func Thick() CharTemplate {
	return CharTemplate{
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
}

// Double returns double-line Unicode box template.
func Double() CharTemplate {
	return CharTemplate{
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
}

// ASCII returns ASCII-only template using +|- characters.
func ASCII() CharTemplate {
	return CharTemplate{
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
}

// Style implements Styler for CharTemplate.
func (ct CharTemplate) Style(base CellStyle) CellStyle {
	base.Template = ct
	return base
}

// writeLine writes single line (border or content) using ops and configured template.
func (r *renderer) writeLine(b *strings.Builder, ops []RenderOp) {
	tpl := r.tpl

	for _, op := range ops {
		switch v := op.(type) {
		// Border ops
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
			for range v.Width {
				b.WriteRune(tpl.hLine)
			}

		// Content ops
		case VLine:
			b.WriteRune(tpl.vLine)
		case Content:
			b.WriteString(v.Text)
		case Space:
			for range v.Width {
				b.WriteByte(' ')
			}
		}
	}
	b.WriteByte('\n')
}
