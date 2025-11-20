package tbl

// ID identifies a cell in the table.
type ID int64

// CellType indicates whether a cell is static or flexible.
type CellType int

const (
	Static CellType = iota
	Flex
)

// BorderSide defines which edges of a cell have borders.
type BorderSide uint8

const (
	BorderNone   BorderSide = 0
	BorderTop    BorderSide = 1 << 0 // 0001
	BorderRight  BorderSide = 1 << 1 // 0010
	BorderBottom BorderSide = 1 << 2 // 0100
	BorderLeft   BorderSide = 1 << 3 // 1000

	BorderAll = BorderTop | BorderRight | BorderBottom | BorderLeft
)

// HAlign specifies horizontal text alignment within a cell.
type HAlign int

const (
	HAlignLeft HAlign = iota
	HAlignCenter
	HAlignRight
)

// VAlign specifies vertical text alignment within a cell.
type VAlign int

const (
	VAlignTop VAlign = iota
	VAlignMiddle
	VAlignBottom
)

// Direction indicates traversal direction for flex cell scanning.
type Direction int

const (
	DirLeft Direction = iota
	DirRight
)

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
	Content struct{ Text string } // finalized line (includes padding, alignment)
	Space   struct{ Width int }
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
