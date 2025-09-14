package table

import (
	"github.com/neox5/tbl/internal/cell"
	"github.com/neox5/tbl/types"
)

// Default configuration values
var (
	defaultTableBorder = types.TableBorder{
		All:    false,
		Around: false,
		Style:  types.Single,
	}
)

// Default returns a configuration with sensible defaults
func Default() *types.Config {
	return &types.Config{
		Border:      &defaultTableBorder,
		CellDefault: cell.Default(),
		Width:       0,
		MaxWidth:    0,
	}
}

// GetBorder returns the table border configuration
func GetBorder(cfg *types.Config) types.TableBorder {
	if cfg.Border != nil {
		return *cfg.Border
	}
	return defaultTableBorder
}

// GetCellDefault returns the default cell configuration
func GetCellDefault(cfg *types.Config) *cell.Cell {
	if cfg.CellDefault != nil {
		return cell.NewFromValue(cfg.CellDefault)
	}
	return cell.Default()
}

// GetWidth returns the table width
func GetWidth(cfg *types.Config) int {
	return cfg.Width
}

// GetMaxWidth returns the maximum table width
func GetMaxWidth(cfg *types.Config) int {
	return cfg.MaxWidth
}
