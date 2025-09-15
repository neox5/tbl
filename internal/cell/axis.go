package cell

// Constants for flexible cell dimensions
const (
	AxisNoCap = -1 // Indicates no end cap boundary
)

// Axis represents a cell's dimension along one axis (column or row)
type Axis struct {
	flex    bool // True if axis uses flexible spanning
	span    int  // Current span (minSpan for flex, actual span for fixed)
	maxSpan int  // Maximum span, or NO_CAP for unlimited
	weight  int  // Weight for flex distribution (default: 1)
	start   int  // Starting position (0 until positioned)
}

// NewAxis creates a new axis with fixed span
func NewAxis(span int) Axis {
	if span <= 0 {
		panic("Axis: span must be greater than 0")
	}
	return Axis{
		flex:    false,
		span:    span,
		maxSpan: span,
		weight:  1,
		start:   0,
	}
}

// NewFlexAxis creates a new axis with flexible span
func NewFlexAxis(minSpan, maxSpan, weight int) Axis {
	if minSpan <= 0 {
		panic("Axis: minSpan must be greater than 0")
	}
	if maxSpan != AxisNoCap && minSpan > maxSpan {
		panic("Axis: minSpan cannot be greater than maxSpan")
	}
	if weight <= 0 {
		panic("Axis: weight must be greater than 0")
	}
	return Axis{
		flex:    true,
		span:    minSpan,
		maxSpan: maxSpan,
		weight:  weight,
		start:   0,
	}
}

// IsFlex returns true if the axis uses flexible spanning
func (a Axis) IsFlex() bool {
	return a.flex
}

// IsPositioned returns true if the axis has been positioned
func (a Axis) IsPositioned() bool {
	return a.start > 0
}

// Span returns the current span
func (a Axis) Span() int {
	return a.span
}

// MaxSpan returns the maximum span
func (a Axis) MaxSpan() int {
	return a.maxSpan
}

// Weight returns the weight for flex distribution
func (a Axis) Weight() int {
	return a.weight
}

// Start returns the starting position
func (a Axis) Start() int {
	return a.start
}

// End returns the ending position (start + span)
func (a Axis) End() int {
	return a.start + a.span
}

// SetStart sets the starting position
func (a *Axis) SetStart(start int) {
	a.start = start
}

// CanGrow returns true if flex axis can grow
func (a Axis) CanGrow() bool {
	if !a.IsFlex() {
		return false
	}
	return a.maxSpan == AxisNoCap || a.span < a.maxSpan
}

// AddSpan adds to the current span
func (a *Axis) AddSpan(add int) {
	if !a.IsFlex() {
		panic("Axis: cannot add span to fixed axis")
	}
	if add <= 0 {
		panic("Axis: add must be greater than 0")
	}
	
	newSpan := a.span + add
	if a.maxSpan != AxisNoCap && newSpan > a.maxSpan {
		panic("Axis: cannot add span above maximum")
	}
	a.span = newSpan
}

