package grid

import (
	"fmt"
	"slices"
)

// indexAdd: add id once per spanned row, then sort that row by Area.Col().
func (g *Grid) indexAdd(id ID) {
	a := g.areas[id]
	if a == nil {
		panic(fmt.Errorf("grid: indexAdd unknown id=%d", id))
	}
	for i := range a.rowSpan { // Go 1.22
		r := a.row + i
		rowIDs := g.rowIndex[r]
		if !slices.Contains(rowIDs, id) {
			rowIDs = append(rowIDs, id)
			slices.SortFunc(rowIDs, g.cmpByCol)
			g.rowIndex[r] = rowIDs
		}
	}
}

// cmpByCol orders IDs by their Area.Col().
func (g *Grid) cmpByCol(a, b ID) int {
	aa, bb := g.areas[a], g.areas[b]
	if aa == nil || bb == nil {
		panic(fmt.Errorf("grid: cmpByCol missing area a=%d b=%d", a, b))
	}
	ca, cb := aa.Col(), bb.Col()
	switch {
	case ca < cb:
		return -1
	case ca > cb:
		return 1
	default:
		return 0
	}
}
