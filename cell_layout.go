package tbl

import "strings"

// Layout formats cell content within given constraints.
// Returns complete lines with padding and alignment applied.
// width and height include padding space.
//
// Process:
//  1. Calculate content dimensions (subtract padding)
//  2. Build content lines with word wrapping
//  3. Apply horizontal alignment to each line
//  4. Apply vertical alignment (padding/truncation to height)
//  5. Add horizontal padding to each line
func (c *Cell) Layout(width, height int, style CellStyle) []string {
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

	// Build content lines with wrapping
	contentLines := buildRawLines(c.content, contentWidth)

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

// buildRawLines converts content into lines that fit width.
// Respects explicit line breaks (\n) in content.
// Words longer than width are truncated with ellipsis.
// Returns lines without padding (natural word wrap boundaries).
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
			if len(w) > width {
				if l.Len() > 0 {
					lines = append(lines, l.String())
					l.Reset()
				}
				lines = append(lines, truncateWithEllipsis(w, width))
				continue
			}

			// Check if word fits on current line
			need := l.Len()
			if need > 0 {
				need++ // space
			}
			need += len(w)

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

// applyHAlign applies horizontal alignment to each line.
// Pads lines to target width with spaces.
func applyHAlign(lines []string, width int, align HAlign) []string {
	result := make([]string, len(lines))

	for i, line := range lines {
		if len(line) >= width {
			result[i] = line
			continue
		}

		pad := width - len(line)
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
func truncateWithEllipsis(text string, width int) string {
	if width <= 0 {
		return ""
	}
	if len(text) <= width {
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

	// Truncate to fit
	cutLen := max(0, width-len(ellipsis))
	return text[:cutLen] + ellipsis
}
