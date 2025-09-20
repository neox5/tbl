package d2

// Vec is an immutable 2D integer displacement.
type Vec struct {
	dx, dy int
}

// NewVec constructs a new vector.
func NewVec(dx, dy int) Vec {
	return Vec{dx: dx, dy: dy}
}

// Dx returns the x-component.
func (v Vec) Dx() int { return v.dx }

// Dy returns the y-component.
func (v Vec) Dy() int { return v.dy }

// Neg returns -v.
func (v Vec) Neg() Vec {
	return Vec{dx: -v.dx, dy: -v.dy}
}

// Add returns v + u.
func (v Vec) Add(u Vec) Vec {
	return Vec{dx: v.dx + u.dx, dy: v.dy + u.dy}
}

// Sub returns v - u.
func (v Vec) Sub(u Vec) Vec {
	return Vec{dx: v.dx - u.dx, dy: v.dy - u.dy}
}

// Scale returns k * v.
func (v Vec) Scale(k int) Vec {
	return Vec{dx: k * v.dx, dy: k * v.dy}
}

// IsZero reports whether v == (0,0).
func (v Vec) IsZero() bool {
	return v.dx == 0 && v.dy == 0
}
