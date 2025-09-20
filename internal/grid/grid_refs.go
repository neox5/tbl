package grid

// Ref addresses an Area stored in Grid. 0 means Nil.
type Ref uint32

const Nil Ref = 0

// RegisterArea stores a pointer to Area and returns its Ref.
// The Area must outlive the Grid usage or be managed by the caller.
func (g *Grid) RegisterArea(a *Area) Ref {
	g.areas = append(g.areas, a)
	return Ref(len(g.areas)) // 1..N
}

func (g *Grid) areaByRef(ref Ref) *Area {
	if ref == Nil {
		return nil
	}
	i := int(ref) - 1
	if i < 0 || i >= len(g.areas) {
		panic("grid: invalid ref")
	}
	return g.areas[i]
}
