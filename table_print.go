package tbl

import (
	"fmt"
	"strings"
)

// printDebug renders grid structure using TBL Grid Notation from README.
// Shows internal state: dimensions, cursor position, cell count.
// Grid display uses A-Z for Static cells, a-z for Flex cells.
// Cursor position indicated by / instead of closing ].
// Includes btmp.Grid bitmap visualization and underlying Bitmap representation.
func (t *Table) printDebug() string {
	var b strings.Builder

	// Header with internal state
	b.WriteString("=== Table Debug ===\n")
	b.WriteString(fmt.Sprintf("Cols: %d, Rows: %d\n", t.g.Cols(), t.g.Rows()))
	b.WriteString(fmt.Sprintf("Cursor: (%d, %d)\n", t.row, t.col))
	b.WriteString(fmt.Sprintf("Cells: %d\n", len(t.cells)))
	b.WriteString(fmt.Sprintf("ColsFixed: %v\n", t.colsFixed))
	b.WriteString("\n")

	// TBL Grid Notation
	b.WriteString("Grid (TBL Notation):\n")
	if t.g.Rows() == 0 {
		b.WriteString("(empty)\n")
	} else {
		for row := range t.g.Rows() {
			b.WriteString(t.renderRow(row))
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")

	// btmp Grid bitmap visualization
	b.WriteString("Grid (btmp):\n")
	gridPrint := t.g.Print()
	if gridPrint == "" {
		b.WriteString("(empty)\n")
	} else {
		b.WriteString(gridPrint)
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Underlying Bitmap
	b.WriteString("Bitmap:\n")
	bitmapPrint := t.g.B.Print()
	if bitmapPrint == "" {
		b.WriteString("(empty)\n")
	} else {
		b.WriteString(bitmapPrint)
		b.WriteString("\n")
	}

	return b.String()
}

// renderRow builds single row output in TBL Grid Notation format.
// Format: {row}: [{cells}] or {row}: [{cells}/ for cursor position.
// Cells spanning from above rows are rendered with their letter repeated.
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

		// Get cell letter
		letter := t.getCellLetter(cell.id, cell.typ)

		// Render cell with its colSpan (whether starting here or spanning from above)
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
	if row == t.row && t.col > 0 && t.col < cols {
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
