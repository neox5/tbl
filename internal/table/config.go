package table

import (
	"github.com/neox5/tbl/internal/cell"
	"github.com/neox5/tbl/types"
)

// defaultConfig returns configuration with sensible defaults
func defaultConfig() *types.Config {
	return &types.Config{
		Border:      &types.DefaultTableBorder,
		DefaultCell: cell.DefaultCell(),
		Width:       0,
		MaxWidth:    0,
	}
}

// applyDefaults applies default values to nil config fields
func applyDefaults(cfg *types.Config) {
	if cfg.Border == nil {
		cfg.Border = &types.DefaultTableBorder
	}
	if cfg.DefaultCell == nil {
		cfg.DefaultCell = cell.DefaultCell()
	}
	// Width and MaxWidth remain 0 if not set
}
