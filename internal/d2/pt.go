package d2

// Pt is an immutable 2D integer point.
type Pt struct {
	x, y int
}

// NewPt constructs a new point.
func NewPt(x, y int) Pt {
	return Pt{x: x, y: y}
}

// X returns the x-coordinate.
func (p Pt) X() int { return p.x }

// Y returns the y-coordinate.
func (p Pt) Y() int { return p.y }

// Add returns a new point translated by v.
func (p Pt) Add(v Vec) Pt {
	return Pt{x: p.x + v.dx, y: p.y + v.dy}
}

// Dist returns the vector q - p (displacement from p to q).
func (p Pt) Dist(q Pt) Vec {
	return Vec{dx: q.x - p.x, dy: q.y - p.y}
}

// Eq reports whether p == q.
func (p Pt) Eq(q Pt) bool {
	return p.x == q.x && p.y == q.y
}
