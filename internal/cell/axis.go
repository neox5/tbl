package cell

// Constants for flexible cell dimensions
const (
	FLEX   = -1 // Indicates flexible span
	NO_CAP = -1 // Indicates no end cap boundary
)

// CellAxis represents a cell's dimension along one axis (column or row)
type CellAxis struct {
	Span   int // Number of columns/rows to span, or FLEX for dynamic
	Weight int // Weight for FLEX spans (default: 1)
	Start  int // Starting position (inclusive)
	End    int // Ending position (exclusive, 0 until resolved)
	EndCap int // Nearest boundary hint, NO_CAP if none
}

// IsFixed returns true if the axis has a fixed end position
func (a CellAxis) IsFixed() bool {
	return a.End != 0
}

// IsFlex returns true if the axis uses flexible spanning
func (a CellAxis) IsFlex() bool {
	return a.Span == FLEX
}
