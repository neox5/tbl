package config

import (
	"github.com/neox5/tbl/internal/cell"
	"github.com/neox5/tbl/types"
)

// Config is the internal configuration implementation
type Config struct {
	border      types.TableBorder
	cellDefault *cell.Cell
	width       int
	maxWidth    int
}

// Default configuration values
var (
	defaultTableBorder = types.TableBorder{
		All:    false,
		Around: false,
		Style:  types.Single,
	}
)

// Default returns a configuration with sensible defaults
func Default() *Config {
	return &Config{
		border:      defaultTableBorder,
		cellDefault: cell.Default(),
		width:       0,
		maxWidth:    0,
	}
}

// Border returns the table border configuration
func (c *Config) Border() types.TableBorder {
	return c.border
}

// CellDefault returns the default cell configuration
func (c *Config) CellDefault() *cell.Cell {
	return c.cellDefault
}

// Width returns the table width
func (c *Config) Width() int {
	return c.width
}

// MaxWidth returns the maximum table width
func (c *Config) MaxWidth() int {
	return c.maxWidth
}

// SetBorder sets the table border configuration
func (c *Config) SetBorder(border types.TableBorder) {
	c.border = border
}

// SetCellDefault sets the default cell configuration
func (c *Config) SetCellDefault(cellDefault *cell.Cell) {
	c.cellDefault = cellDefault
}

// SetWidth sets the table width
func (c *Config) SetWidth(width int) {
	c.width = width
}

// SetMaxWidth sets the maximum table width
func (c *Config) SetMaxWidth(maxWidth int) {
	c.maxWidth = maxWidth
}
