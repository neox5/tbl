package tbl

import (
	"github.com/neox5/tbl/internal/table"
	"github.com/neox5/tbl/types"
)

// Tbl is the public interface for building tables
type Tbl interface {
	AddRow(cells ...any)
	R(cells ...any)
	C(value any) Cell
}

// Cell is the public interface for table cells
type Cell interface {
	WithAlign(h types.HorizontalAlignment, v types.VerticalAlignment) Cell
	WithContent(content string) Cell
	WithSpan(col, row int) Cell
	WithBorder(border types.CellBorder) Cell
	
	// Short form aliases
	A(h types.HorizontalAlignment, v types.VerticalAlignment) Cell
	C(content string) Cell
	S(col, row int) Cell
	B(border types.CellBorder) Cell
}

// New creates a new table with default configuration
func New() *table.Table {
	return table.New()
}

// NewWithConfig creates a new table with the specified configuration
func NewWithConfig(cfg types.Config) *table.Table {
	return table.NewWithConfig(cfg)
}
