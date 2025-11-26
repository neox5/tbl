package tbl

// TableConfig specifies global table constraints.
type TableConfig struct {
	MaxWidth int // maximum total table width (0 = no maximum)
}

// ColConfig specifies column dimension constraints.
type ColConfig struct {
	Width    int // fixed width (0 = auto)
	MinWidth int // minimum width (0 = no minimum)
	MaxWidth int // maximum width (0 = no maximum)
}

// SetTableConfig applies global table constraints.
func (t *Table) SetTableConfig(cfg TableConfig) *Table {
	t.tableConfig = cfg
	return t
}

// SetColConfig applies dimension constraints to column.
func (t *Table) SetColConfig(col int, cfg ColConfig) *Table {
	t.colConfigs[col] = cfg
	return t
}
