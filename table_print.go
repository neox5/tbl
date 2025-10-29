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
	processed := make(map[int]bool) // track processed columns

	for col := range cols {
		if processed[col] {
			continue
		}

		// Find cell at this position
		cellID := t.findCellAt(row, col)
		if cellID == ID(0) {
			// Empty cell (shouldn't happen with current design)
			b.WriteString("?")
			processed[col] = true
			if col < cols-1 {
				b.WriteString("|")
			}
			continue
		}

		cell := t.cells[cellID]

		// Get cell letter
		letter := t.getCellLetter(cellID, cell.typ)

		// Render span
		for i := range cell.cSpan {
			if i > 0 {
				b.WriteString(" ") // space instead of | for span
			}
			b.WriteString(letter)
			processed[col+i] = true
		}

		// Add separator if not last column
		if col+cell.cSpan < cols {
			b.WriteString("|")
		}

		col += cell.cSpan - 1 // advance by span (loop will +1)
	}

	// Cursor indicator: show / if row matches cursor and row is incomplete
	if row == t.c.Row() && t.c.Col() > 0 && t.c.Col() < cols {
		b.WriteString("/")
	} else {
		b.WriteString("]")
	}

	return b.String()
}

// findCellAt returns cell ID at grid position, or 0 if none.
func (t *Table) findCellAt(row, col int) ID {
	for id, cell := range t.cells {
		if cell.r <= row && row < cell.r+cell.rSpan &&
			cell.c <= col && col < cell.c+cell.cSpan {
			return id
		}
	}
	return ID(0)
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
