package tbl

const (
	FLEX     = -1
	NO_CAP   = -1
)

type CellAxis struct {
	Span   int // >0 or FLEX
	Weight int // for FLEX; default: 1
	Start  int // inclusive
	End    int // exclusive; OPEN_END until resolved
	EndCap int // nearest boundary hint; NO_CAP if none
}

func (a CellAxis) IsFixed() bool {
	return a.End != 0
}

func (a CellAxis) IsFlex() bool {
	return a.Span == FLEX
}
