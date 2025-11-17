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

// setColConfig applies dimension constraints to column.
// Internal use only.
func (t *Table) setColConfig(col int, cfg ColConfig) {
	t.colConfigs[col] = cfg
}

// setTableConfig applies global table constraints.
// Internal use only.
func (t *Table) setTableConfig(cfg TableConfig) {
	t.tableConfig = cfg
}
