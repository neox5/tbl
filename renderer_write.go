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
	cornerX  rune
	hLine    rune
	vLine    rune
}

// defaultTemplate is the standard ASCII box drawing character set.
var defaultTemplate = charTemplate{
	cornerTL: '+',
	cornerTR: '+',
	cornerBL: '+',
	cornerBR: '+',
	cornerT:  '+',
	cornerB:  '+',
	cornerX:  '+',
	hLine:    '-',
	vLine:    '|',
}

// writeBorder interprets border instruction sequence using template.
func (r *renderer) writeBorder(b *strings.Builder, ops []RenderOp) {
	tpl := defaultTemplate

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
		case CornerX:
			b.WriteRune(tpl.cornerX)
		case HLine:
			for range v.Width {
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
			writeAlignedContent(b, v.Text, v.Width, v.HAlign)
		case Space:
			for range v.Width {
				b.WriteByte(' ')
			}
		}
	}
	b.WriteByte('\n')
}

// writeAlignedContent writes text with alignment and padding.
func writeAlignedContent(b *strings.Builder, text string, width int, align HAlign) {
	pad := width - len(text)
	if pad < 0 {
		pad = 0
	}

	// Simple left align with padding for now
	b.WriteByte(' ')
	b.WriteString(text)
	for range pad + 1 {
		b.WriteByte(' ')
	}
}
