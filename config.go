package tbl

var DefaultConfig = Config{
	Border:      &DefaultTableBorder,
	CellDefault: &DefaultCell,
	Width:       0,
	MaxWidth:    0,
}

var DefaultTableBorder = TableBorder{
	All:    false,
	Around: false,
	Style:  Single,
}

type Config struct {
	Border      *TableBorder
	CellDefault *Cell
	Width       int
	MaxWidth    int
}

func (base Config) Merge(cfg Config) Config {
	result := Config{
		Border:      base.Border,
		CellDefault: base.CellDefault,
		Width:       base.Width,
		MaxWidth:    base.MaxWidth,
	}

	if cfg.Border != nil {
		result.Border = cfg.Border
	}
	if cfg.CellDefault != nil {
		result.CellDefault = cfg.CellDefault
	}
	if cfg.Width != 0 {
		result.Width = cfg.Width
	}
	if cfg.MaxWidth != 0 {
		result.MaxWidth = cfg.MaxWidth
	}

	return result
}
