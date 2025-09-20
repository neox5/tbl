package grid

import "fmt"

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
