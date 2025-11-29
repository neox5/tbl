package tbl

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// stripAnsiCodes removes ANSI escape sequences from string.
// Used for visual width calculation while preserving formatting in output.
func stripAnsiCodes(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// visualLength returns the visible character count of a string.
// Strips ANSI escape sequences and counts runes (not bytes).
func visualLength(s string) int {
	stripped := stripAnsiCodes(s)
	return utf8.RuneCountInString(stripped)
}

// layout formats cell content within given constraints.
// Returns complete lines with padding and alignment applied.
// width and height include padding space.
//
// Process:
//  1. Calculate content dimensions (subtract padding)
//  2. Build content lines with word wrapping based on WrapMode
//  3. Apply horizontal alignment to each line
//  4. Apply vertical alignment (padding/truncation to height)
//  5. Add horizontal padding to each line
func (c *Cell) layout(width, height int, style CellStyle) []string {
	if width <= 0 || height <= 0 {
		return []string{}
	}

	// Calculate content dimensions
	contentWidth := width - style.Padding.Left - style.Padding.Right
	contentHeight := height - style.Padding.Top - style.Padding.Bottom

	if contentWidth <= 0 || contentHeight <= 0 {
		// All padding, no content space
		emptyLine := strings.Repeat(" ", width)
		lines := make([]string, height)
		for i := range lines {
			lines[i] = emptyLine
		}
		return lines
	}

	// Build content lines with wrapping based on WrapMode
	var contentLines []string
	switch style.WrapMode {
	case WrapChar:
		contentLines = buildRawLinesChar(c.content, contentWidth)
	case WrapTruncate:
		contentLines = buildRawLinesTruncate(c.content, contentWidth)
	default: // WrapWord (default)
		contentLines = buildRawLines(c.content, contentWidth)
	}

	// Apply horizontal alignment
	alignedLines := applyHAlign(contentLines, contentWidth, style.HAlign)

	// Apply vertical alignment
	paddedLines := applyVAlign(alignedLines, contentWidth, contentHeight, style.VAlign)

	// Add horizontal padding to each line
	leftPad := strings.Repeat(" ", style.Padding.Left)
	rightPad := strings.Repeat(" ", style.Padding.Right)
	finalLines := make([]string, height)

	// Top padding
	emptyLine := strings.Repeat(" ", width)
	for i := range style.Padding.Top {
		finalLines[i] = emptyLine
	}

	// Content lines with horizontal padding
	for i, line := range paddedLines {
		finalLines[style.Padding.Top+i] = leftPad + line + rightPad
	}

	// Bottom padding
	for i := style.Padding.Top + len(paddedLines); i < height; i++ {
		finalLines[i] = emptyLine
	}

	return finalLines
}

// buildRawLines converts content into lines that fit width using word wrapping.
// Respects explicit line breaks (\n) in content.
// Words longer than width are truncated with ellipsis.
// Returns lines without padding (natural word wrap boundaries).
// Uses visual length for width calculation (strips ANSI codes, counts runes).
func buildRawLines(content string, width int) []string {
	if width <= 0 || content == "" {
		return []string{""}
	}

	segments := strings.Split(content, "\n")
	var lines []string

	for _, seg := range segments {
		seg = strings.TrimSpace(seg)

		// Empty segment → empty line
		if seg == "" {
			lines = append(lines, "")
			continue
		}

		// Process words
		words := strings.Fields(seg)
		var l strings.Builder

		for _, w := range words {
			// Long word → flush current, add truncated word
			if visualLength(w) > width {
				if l.Len() > 0 {
					lines = append(lines, l.String())
					l.Reset()
				}
				lines = append(lines, truncateWithEllipsis(w, width))
				continue
			}

			// Check if word fits on current line
			need := visualLength(l.String())
			if need > 0 {
				need++ // space
			}
			need += visualLength(w)

			// Doesn't fit → flush current line, start new
			if need > width {
				lines = append(lines, l.String())
				l.Reset()
			}

			// Add word to line
			if l.Len() > 0 {
				l.WriteByte(' ')
			}
			l.WriteString(w)
		}

		// Flush last line of segment
		if l.Len() > 0 {
			lines = append(lines, l.String())
		}
	}

	return lines
}

// buildRawLinesChar converts content into lines that fit width using character wrapping.
// Wraps at any character boundary, ignoring word boundaries.
// Respects explicit line breaks (\n) in content.
// Uses visual length for width calculation (strips ANSI codes, counts runes).
func buildRawLinesChar(content string, width int) []string {
	if width <= 0 || content == "" {
		return []string{""}
	}

	segments := strings.Split(content, "\n")
	var lines []string

	for _, seg := range segments {
		seg = strings.TrimSpace(seg)

		// Empty segment → empty line
		if seg == "" {
			lines = append(lines, "")
			continue
		}

		// Wrap at character boundaries
		// Note: Cannot use simple byte-slice indexing with ANSI codes
		// Must preserve ANSI sequences while measuring visual length
		for len(seg) > 0 {
			if visualLength(seg) <= width {
				lines = append(lines, seg)
				break
			}

			// Find break point that fits width visually
			breakPoint := findVisualBreakpoint(seg, width)
			lines = append(lines, seg[:breakPoint])
			seg = seg[breakPoint:]
		}
	}

	return lines
}

// findVisualBreakpoint finds byte index where visual length reaches target width.
// Preserves ANSI escape sequences and counts runes (not bytes).
func findVisualBreakpoint(s string, width int) int {
	visualLen := 0
	bytePos := 0

	for bytePos < len(s) {
		// Check for ANSI escape sequence
		if bytePos < len(s) && s[bytePos] == '\x1b' {
			// Find end of ANSI sequence
			endPos := bytePos + 1
			for endPos < len(s) && s[endPos] != 'm' {
				endPos++
			}
			if endPos < len(s) {
				endPos++ // include 'm'
			}
			bytePos = endPos
			continue
		}

		// Regular character - decode rune
		_, size := utf8.DecodeRuneInString(s[bytePos:])
		visualLen++
		if visualLen > width {
			return bytePos
		}
		bytePos += size
	}

	return bytePos
}

// buildRawLinesTruncate converts content into lines, truncating overflow with ellipsis.
// Each line from explicit breaks (\n) is truncated independently.
// No wrapping occurs - content exceeding width is cut off.
// Uses visual length for width calculation (strips ANSI codes, counts runes).
func buildRawLinesTruncate(content string, width int) []string {
	if width <= 0 || content == "" {
		return []string{""}
	}

	segments := strings.Split(content, "\n")
	var lines []string

	for _, seg := range segments {
		seg = strings.TrimSpace(seg)

		// Empty segment → empty line
		if seg == "" {
			lines = append(lines, "")
			continue
		}

		// Truncate if too long
		if visualLength(seg) > width {
			seg = truncateWithEllipsis(seg, width)
		}

		lines = append(lines, seg)
	}

	return lines
}

// applyHAlign applies horizontal alignment to each line.
// Pads lines to target width with spaces.
// Uses visual length for alignment calculation (strips ANSI codes, counts runes).
func applyHAlign(lines []string, width int, align HAlign) []string {
	result := make([]string, len(lines))

	for i, line := range lines {
		visualLen := visualLength(line)
		if visualLen >= width {
			result[i] = line
			continue
		}

		pad := width - visualLen
		var lPad, rPad int

		switch align {
		case HAlignLeft:
			rPad = pad
		case HAlignRight:
			lPad = pad
		case HAlignCenter:
			lPad = pad / 2
			rPad = pad - lPad
		default:
			rPad = pad
		}

		result[i] = strings.Repeat(" ", lPad) + line + strings.Repeat(" ", rPad)
	}

	return result
}

// applyVAlign applies vertical alignment across all lines.
// Truncates if lines exceed height.
// Pads with empty lines to reach height.
func applyVAlign(lines []string, width, height int, align VAlign) []string {
	// Already fits or exceeds
	if len(lines) >= height {
		return lines[:height]
	}

	emptyLines := height - len(lines)

	var tLines, bLines int

	switch align {
	case VAlignTop:
		bLines = emptyLines
	case VAlignMiddle:
		tLines = emptyLines / 2
		bLines = emptyLines - tLines
	case VAlignBottom:
		tLines = emptyLines
	default:
		bLines = emptyLines
	}

	result := make([]string, 0, height)
	emptyLine := strings.Repeat(" ", width)

	for range tLines {
		result = append(result, emptyLine)
	}
	result = append(result, lines...)
	for range bLines {
		result = append(result, emptyLine)
	}

	return result
}

// truncateWithEllipsis shortens text to fit width with ellipsis.
// Ellipsis length adapts to width: "..." (≥3), ".." (2), "." (1), "" (0).
// Preserves ANSI codes in truncated portion, adds ellipsis at visual breakpoint.
// Counts runes (not bytes) for visual width.
func truncateWithEllipsis(text string, width int) string {
	if width <= 0 {
		return ""
	}

	visualLen := visualLength(text)
	if visualLen <= width {
		return text
	}

	// Adapt ellipsis to width
	ellipsis := "..."
	switch width {
	case 1:
		ellipsis = "."
	case 2:
		ellipsis = ".."
	}

	// Find breakpoint for visual width minus ellipsis
	targetWidth := max(0, width-len(ellipsis))

	breakPoint := findVisualBreakpoint(text, targetWidth)
	return text[:breakPoint] + ellipsis
}
