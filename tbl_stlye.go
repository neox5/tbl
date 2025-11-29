package tbl

// containsTemplate checks if any styler is a CharTemplate.
func containsTemplate(stylers []Freestyler) bool {
	for _, s := range stylers {
		if _, ok := s.(CharTemplate); ok {
			return true
		}
	}
	return false
}

// SetDefaultStyle sets the base style for all cells.
// Can be overridden by column, row, or cell-specific styles.
func (t *Table) SetDefaultStyle(stylers ...Freestyler) *Table {
	t.defaultStyle = t.defaultStyle.Apply(stylers...)
	return t
}

// SetColStyle sets style for all cells in column.
// Overrides default style, can be overridden by row or cell styles.
// Panics if stylers contain CharTemplate (templates are table-level only).
func (t *Table) SetColStyle(col int, stylers ...Freestyler) *Table {
	if containsTemplate(stylers) {
		panic("tbl: CharTemplate only supported via SetDefaultStyle")
	}
	t.columnStyles[col] = t.columnStyles[col].Apply(stylers...)
	return t
}

// SetRowStyle sets style for all cells in row.
// Overrides default and column styles, can be overridden by cell styles.
// Panics if stylers contain CharTemplate (templates are table-level only).
func (t *Table) SetRowStyle(row int, stylers ...Freestyler) *Table {
	if containsTemplate(stylers) {
		panic("tbl: CharTemplate only supported via SetDefaultStyle")
	}
	t.rowStyles[row] = t.rowStyles[row].Apply(stylers...)
	return t
}

// SetStyleFunc sets the programmable style resolver for all cells.
//
// The Funcstyler is evaluated after default, column, and row styles and
// before cell-specific styles in resolveStyle.
func (t *Table) SetStyleFunc(fn Funcstyler) *Table {
	t.styleFunc = fn
	return t
}

// SetCellStyle sets style for specific cell.
// Highest priority, overrides all other styles.
// Panics if stylers contain CharTemplate (templates are table-level only).
func (t *Table) SetCellStyle(id ID, stylers ...Freestyler) *Table {
	if containsTemplate(stylers) {
		panic("tbl: CharTemplate only supported via SetDefaultStyle")
	}
	t.cellStyles[id] = t.cellStyles[id].Apply(stylers...)
	return t
}

// resolveStyle returns effective style for cell using hierarchy.
// Resolution order: defaultStyle < columnStyles < rowStyles < styleFunc < cellStyles
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

	// Programmable style at origin
	if t.styleFunc != nil {
		fs := t.styleFunc(cell.r, cell.c)
		if fs.Template != (CharTemplate{}) {
			panic("tbl: CharTemplate only supported via SetDefaultStyle")
		}
		style = style.merge(fs)
	}

	// Cell-specific style (highest priority)
	if cs, ok := t.cellStyles[cell.id]; ok {
		style = style.merge(cs)
	}

	return style
}
