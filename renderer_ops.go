package tbl

// RenderOp is a rendering instruction.
type RenderOp interface {
	renderOp() // unexported marker method
}

// Border instructions
type (
	CornerTL struct{}
	CornerTR struct{}
	CornerBL struct{}
	CornerBR struct{}
	CornerT  struct{} // top junction (─┬─)
	CornerB  struct{} // bottom junction (─┴─)
	CornerL  struct{} // left junction (├─)
	CornerR  struct{} // right junction (─┤)
	CornerX  struct{} // cross junction (─┼─)
	HLine    struct{ Width int }
)

// Content instructions
type (
	VLine   struct{}
	Content struct {
		Text   string
		Width  int
		HAlign HAlign
	}
	Space struct{ Width int }
)

// Implement marker
func (CornerTL) renderOp() {}
func (CornerTR) renderOp() {}
func (CornerBL) renderOp() {}
func (CornerBR) renderOp() {}
func (CornerT) renderOp()  {}
func (CornerB) renderOp()  {}
func (CornerL) renderOp()  {}
func (CornerR) renderOp()  {}
func (CornerX) renderOp()  {}
func (HLine) renderOp()    {}
func (VLine) renderOp()    {}
func (Content) renderOp()  {}
func (Space) renderOp()    {}
