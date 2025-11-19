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

// SetColConfig applies dimension constraints to column.
func (t *Table) SetColConfig(col int, cfg ColConfig) *Table {
	t.colConfigs[col] = cfg
	return t
}

// SetTableConfig applies global table constraints.
func (t *Table) SetTableConfig(cfg TableConfig) *Table {
	t.tableConfig = cfg
	return t
}

// SetDefaultStyle sets the base style for all cells.
// Can be overridden by column, row, or cell-specific styles.
func (t *Table) SetDefaultStyle(style CellStyle) *Table {
	t.defaultStyle = style
	return t
}

// SetColStyle sets style for all cells in column.
// Overrides default style, can be overridden by row or cell styles.
func (t *Table) SetColStyle(col int, style CellStyle) *Table {
	t.columnStyles[col] = style
	return t
}

// SetRowStyle sets style for all cells in row.
// Overrides default and column styles, can be overridden by cell styles.
func (t *Table) SetRowStyle(row int, style CellStyle) *Table {
	t.rowStyles[row] = style
	return t
}

// SetCellStyle sets style for specific cell.
// Highest priority, overrides all other styles.
func (t *Table) SetCellStyle(id ID, style CellStyle) *Table {
	t.cellStyles[id] = style
	return t
}

// BorderAll sets all borders on all cells as default style.
// Quick preset for fully bordered tables.
func (t *Table) BorderAll() *Table {
	t.defaultStyle.Border = Border{
		Sides: BorderAll,
	}
	return t
}

// BorderNone removes all borders from all cells as default style.
// Quick preset for borderless tables.
func (t *Table) BorderNone() *Table {
	t.defaultStyle.Border = Border{
		Sides: BorderNone,
	}
	return t
}

// BorderHeader sets borders on specified row only.
// Quick preset for header-bordered tables.
func (t *Table) BorderHeader(row int) *Table {
	t.rowStyles[row] = CellStyle{
		Border: Border{
			Sides: BorderAll,
		},
	}
	return t
}
