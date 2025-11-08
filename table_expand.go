package tbl

import (
	"fmt"
	"sort"
)

// distributeAndExpand handles expansion for flex cells in one row.
// Distributes needed cols fairly: base amount to all, remainder to cells with least expansion.
// Tie-breaking: leftmost cells get priority.
// Directly expands cells and shifts adjacent content.
func (t *Table) distributeAndExpand(row int, flexCells []flexCell, needed int) {
	if len(flexCells) == 0 || needed <= 0 {
		return
	}

	// Sort by: 1) addedSpan ascending, 2) col ascending (left to right)
	sort.Slice(flexCells, func(i, j int) bool {
		if flexCells[i].addedSpan != flexCells[j].addedSpan {
			return flexCells[i].addedSpan < flexCells[j].addedSpan
		}
		return flexCells[i].cell.c < flexCells[j].cell.c
	})

	n := len(flexCells)

	// Calculate base distribution
	base := needed / n
	remainder := needed % n

	// Process each flex cell in sorted order (left to right)
	// This ensures we don't re-shift already expanded cells
	for i, fc := range flexCells {
		// Calculate expansion amount for this cell
		expandAmount := base
		if i < remainder {
			expandAmount++
		}

		if expandAmount == 0 {
			continue
		}

		// Find adjacent cells that need to shift right
		adjacentCol := fc.cell.c + fc.cell.cSpan
		shiftCells := t.findCellsToShift(row, adjacentCol)

		// Execute shifts from right to left
		t.shiftCellsRight(shiftCells, expandAmount)

		// Expand the flex cell
		t.expandCell(fc.cell, expandAmount)
	}
}

// expandCell expands a flex cell and updates grid state.
func (t *Table) expandCell(cell *Cell, cols int) {
	if cols <= 0 {
		return
	}
	if cell.typ != Flex {
		panic(fmt.Sprintf("cannot expand Static cell id=%d", cell.id))
	}

	// Clear old position
	t.g.ClearRect(cell.r, cell.c, cell.rSpan, cell.cSpan)

	// Update cell span
	cell.cSpan += cols

	// Set new position
	t.g.SetRect(cell.r, cell.c, cell.rSpan, cell.cSpan)
}

// findCellsToShift identifies cells that need to move for expansion.
func (t *Table) findCellsToShift(row, fromCol int) []*Cell {
	var cells []*Cell
	seen := make(map[ID]bool)

	// Find all cells in row starting at fromCol
	for col := fromCol; col < t.g.Cols(); col++ {
		cell := t.getCellAt(row, col)
		if cell != nil && !seen[cell.id] {
			seen[cell.id] = true
			cells = append(cells, cell)
			col = cell.c + cell.cSpan - 1 // Skip to end of cell
		}
	}

	return cells
}

// shiftCellsRight moves cells right by delta columns.
func (t *Table) shiftCellsRight(cells []*Cell, delta int) {
	if delta <= 0 || len(cells) == 0 {
		return
	}

	// Sort by column descending for right-to-left processing
	sort.Slice(cells, func(i, j int) bool {
		return cells[i].c > cells[j].c
	})

	// Clear all cells from grid
	for _, cell := range cells {
		t.g.ClearRect(cell.r, cell.c, cell.rSpan, cell.cSpan)
	}

	// Move cells and re-set in grid
	for _, cell := range cells {
		cell.c += delta
		t.g.SetRect(cell.r, cell.c, cell.rSpan, cell.cSpan)
	}
}
