package tbl

import "strings"

// Layout formats cell content within given constraints.
// Returns one string per line, with alignment and truncation applied.
//
// Content processing:
//  1. Split content by '\n'
//  2. Word wrap each line at width boundary (full words only)
//  3. Truncate if total lines exceed height (ellipsis on last visible line)
//  4. Apply vertical alignment (distribute empty lines)
//  5. Apply horizontal alignment per line (distribute spaces)
//
// Word wrap behavior:
//   - Preserves whole words (breaks at spaces)
//   - Long words exceeding width are moved to next line
//   - If word cannot fit on any line, truncates with ellipsis
//
// Truncation:
//   - Ellipsis adapts to width: "..." (width≥3), ".." (width=2), "." (width=1)
//   - Applied to last visible line when height constraint exceeded
func (c *Cell) Layout(width, height int) []string {
	if width <= 0 || height <= 0 {
		return []string{}
	}

	// Step 1: Split by newlines
	rawLines := strings.Split(c.content, "\n")

	// Step 2: Word wrap each line
	var wrappedLines []string
	for _, line := range rawLines {
		wrappedLines = append(wrappedLines, wrapLine(line, width)...)
	}

	// Step 3: Height truncation
	if len(wrappedLines) > height {
		wrappedLines = wrappedLines[:height]
		lastIdx := height - 1
		wrappedLines[lastIdx] = truncateWithEllipsis(wrappedLines[lastIdx], width)
	}

	// Step 4: Vertical alignment
	totalLines := len(wrappedLines)
	emptyLines := height - totalLines

	var alignedLines []string

	switch c.vAlign {
	case VAlignTop:
		alignedLines = append(alignedLines, wrappedLines...)
		for range emptyLines {
			alignedLines = append(alignedLines, strings.Repeat(" ", width))
		}
	case VAlignMiddle:
		topPad := emptyLines / 2
		bottomPad := emptyLines - topPad
		for range topPad {
			alignedLines = append(alignedLines, strings.Repeat(" ", width))
		}
		alignedLines = append(alignedLines, wrappedLines...)
		for range bottomPad {
			alignedLines = append(alignedLines, strings.Repeat(" ", width))
		}
	case VAlignBottom:
		for range emptyLines {
			alignedLines = append(alignedLines, strings.Repeat(" ", width))
		}
		alignedLines = append(alignedLines, wrappedLines...)
	default:
		alignedLines = append(alignedLines, wrappedLines...)
		for range emptyLines {
			alignedLines = append(alignedLines, strings.Repeat(" ", width))
		}
	}

	// Step 5: Horizontal alignment
	for i, line := range alignedLines {
		alignedLines[i] = alignLine(line, width, c.hAlign)
	}

	return alignedLines
}

// wrapLine splits line into multiple lines at width boundary.
// Preserves whole words (breaks at spaces).
// Long words exceeding width are moved to next line and truncated if still too long.
func wrapLine(line string, width int) []string {
	if width <= 0 {
		return []string{""}
	}

	line = strings.TrimSpace(line)
	if len(line) == 0 {
		return []string{""}
	}

	words := strings.Fields(line)
	var result []string
	var currentLine strings.Builder

	for _, word := range words {
		// Word too long for any line
		if len(word) > width {
			// Flush current line if any
			if currentLine.Len() > 0 {
				result = append(result, currentLine.String())
				currentLine.Reset()
			}
			// Truncate word
			result = append(result, truncateWithEllipsis(word, width))
			continue
		}

		// Try adding word to current line
		testLen := currentLine.Len()
		if testLen > 0 {
			testLen++ // space
		}
		testLen += len(word)

		if testLen <= width {
			// Fits on current line
			if currentLine.Len() > 0 {
				currentLine.WriteByte(' ')
			}
			currentLine.WriteString(word)
		} else {
			// Start new line
			if currentLine.Len() > 0 {
				result = append(result, currentLine.String())
				currentLine.Reset()
			}
			currentLine.WriteString(word)
		}
	}

	// Flush last line
	if currentLine.Len() > 0 {
		result = append(result, currentLine.String())
	}

	if len(result) == 0 {
		return []string{""}
	}

	return result
}

// alignLine applies horizontal alignment to single line.
// Pads with spaces to reach target width.
func alignLine(line string, width int, align HAlign) string {
	lineLen := len(line)
	if lineLen >= width {
		return line
	}

	pad := width - lineLen

	switch align {
	case HAlignLeft:
		return line + strings.Repeat(" ", pad)
	case HAlignRight:
		return strings.Repeat(" ", pad) + line
	case HAlignCenter:
		leftPad := pad / 2
		rightPad := pad - leftPad
		return strings.Repeat(" ", leftPad) + line + strings.Repeat(" ", rightPad)
	default:
		return line + strings.Repeat(" ", pad)
	}
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
