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

// Padding specifies space around cell content.
type Padding struct {
	Top, Bottom, Left, Right int
}

// Border specifies which edges of a cell should render borders.
type Border struct {
	Sides    BorderSide // Which edges render visually (characters)
	Physical BorderSide // Which edges occupy physical space
}

// Has reports whether border side occupies space (visual or physical).
func (b Border) Has(side BorderSide) bool {
	return b.IsVisual(side) || (b.Physical&side) != 0
}

// IsVisual reports whether border side renders as character.
func (b Border) IsVisual(side BorderSide) bool {
	return (b.Sides & side) != 0
}

// CellStyle contains presentation attributes for a cell.
type CellStyle struct {
	Padding Padding
	HAlign  HAlign
	VAlign  VAlign
	Border  Border
}

// merge applies non-zero overrides from other style to this style.
// Returns new CellStyle with merged values.
// Zero values in override are ignored (base value preserved).
func (s CellStyle) merge(other CellStyle) CellStyle {
	result := s

	// Merge padding (0 means use base value)
	if other.Padding.Top != 0 {
		result.Padding.Top = other.Padding.Top
	}
	if other.Padding.Bottom != 0 {
		result.Padding.Bottom = other.Padding.Bottom
	}
	if other.Padding.Left != 0 {
		result.Padding.Left = other.Padding.Left
	}
	if other.Padding.Right != 0 {
		result.Padding.Right = other.Padding.Right
	}

	// Alignment always override
	result.HAlign = other.HAlign
	result.VAlign = other.VAlign

	// Border always override
	result.Border = other.Border

	return result
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

// resolveStyle returns effective style for cell using hierarchy.
// Resolution order: defaultStyle < columnStyles < rowStyles < cellStyles
// Only considers cell origin position (cell.r, cell.c) for multi-span cells.
func (t *Table) resolveStyle(cell *Cell) CellStyle {
	style := t.defaultStyle

	// Column style at origin
	if cs, ok := t.columnStyles[cell.c]; ok {
		style = style.merge(cs)
	}

	// Row style at origin
	if rs, ok := t.rowStyles[cell.r]; ok {
		style = style.merge(rs)
	}

	// Cell-specific style (highest priority)
	if cs, ok := t.cellStyles[cell.id]; ok {
		style = style.merge(cs)
	}

	return style
}
