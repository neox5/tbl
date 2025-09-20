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

// indexRemove: find id in each spanned row and delete in-place.
func (g *Grid) indexRemove(id ID) {
	a := g.areas[id]
	if a == nil {
		panic(fmt.Errorf("grid: indexRemove unknown id=%d", id))
	}
	for i := range a.rowSpan {
		r := a.row + i
		rowIDs := g.rowIndex[r]
		if len(rowIDs) == 0 {
			continue
		}
		if idx := slices.Index(rowIDs, id); idx >= 0 {
			rowIDs = slices.Delete(rowIDs, idx, idx+1)
			if len(rowIDs) == 0 {
				delete(g.rowIndex, r)
			} else {
				g.rowIndex[r] = rowIDs
			}
		}
	}
}

// indexUpdate: remove then add using current Area.
func (g *Grid) indexUpdate(id ID) {
	g.indexRemove(id)
	g.indexAdd(id)
}
