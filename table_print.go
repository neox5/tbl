package tbl

import (
	"fmt"
	"strings"
)

// printDebug renders grid structure using TBL Grid Notation from README.
// Shows internal state: dimensions, cursor position, cell count.
// Grid display uses A-Z for Static cells, a-z for Flex cells.
// Cursor position indicated by / instead of closing ].
func (t *Table) printDebug() string {
	var b strings.Builder

	// Header with internal state
	b.WriteString("=== Table Debug ===\n")
	b.WriteString(fmt.Sprintf("Cols: %d, Rows: %d\n", t.g.Cols(), t.g.Rows()))
	b.WriteString(fmt.Sprintf("Cursor: (%d, %d)\n", t.c.Row(), t.c.Col()))
	b.WriteString(fmt.Sprintf("Cells: %d\n", len(t.cells)))
	b.WriteString(fmt.Sprintf("ColsFixed: %v\n", t.colsFixed))
	b.WriteString("\n")

	// Grid display
	b.WriteString("Grid:\n")
	if t.g.Rows() == 0 {
		b.WriteString("(empty)\n")
		return b.String()
	}

	for row := range t.g.Rows() {
		b.WriteString(t.renderRow(row))
		b.WriteString("\n")
	}

	return b.String()
}

// renderRow builds single row output in TBL Grid Notation format.
// Format: {row}: [{cells}] or {row}: [{cells}/ for cursor position.
func (t *Table) renderRow(row int) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("%d: [", row))

	cols := t.g.Cols()
	if cols == 0 {
		b.WriteString("]")
		return b.String()
	}

	// Track current column position for rendering
	col := 0

	for col < cols {
		cell := t.getCellAt(row, col)

		if cell == nil {
			// Empty cell (shouldn't happen with valid grid)
			b.WriteString("?")
			if col < cols-1 {
				b.WriteString("|")
			}
			col++
			continue
		}

		// Skip if cell doesn't start in this row (spanning from above)
		if cell.r != row {
			col++
			continue
		}

		// Get cell letter
		letter := t.getCellLetter(cell.id, cell.typ)

		// Render span
		for i := range cell.cSpan {
			if i > 0 {
				b.WriteString(" ") // space instead of | for span
			}
			b.WriteString(letter)
		}

		col += cell.cSpan

		// Add separator if not at end
		if col < cols {
			b.WriteString("|")
		}
	}

	// Cursor indicator: show / if row matches cursor and row is incomplete
	if row == t.c.Row() && t.c.Col() > 0 && t.c.Col() < cols {
		b.WriteString("/")
	} else {
		b.WriteString("]")
	}

	return b.String()
}

// getCellLetter returns display letter for cell.
// Static: A-Z (cycling), Flex: a-z (cycling).
func (t *Table) getCellLetter(id ID, typ CellType) string {
	// Use ID for letter assignment (cycling through alphabet)
	letterIdx := (int(id) - 1) % 26

	if typ == Static {
		return string(rune('A' + letterIdx))
	}
	return string(rune('a' + letterIdx))
}
