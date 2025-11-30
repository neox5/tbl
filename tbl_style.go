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

// SetStyleFunc sets the programmable style resolver(s) for cells.
//
// Functions are composed left-to-right: later functions override earlier ones
// via CellStyle.merge(). Each function receives StyleContext with cell origin
// position, table dimensions, and cell content.
//
// Composition semantics:
//   - Functions evaluated in order: fn1, fn2, fn3, ...
//   - Results merged sequentially: base.merge(fn1).merge(fn2).merge(fn3)
//   - Non-zero values in later styles override earlier ones
//
// Multi-span cells:
//   - Only cell origin position is evaluated
//   - Spanning cells do NOT trigger functions for covered positions
//
// Resolution order:
//
//	defaultStyle < columnStyles < rowStyles < SetStyleFunc < cellStyles
//
// Example - Header with zebra striping:
//
//	t.SetStyleFunc(
//	    tbl.FirstRow(tbl.BBottom()),
//	    tbl.EvenRowsSkipN(1, tbl.Pad(1)),
//	)
//
// Example - Content-based styling:
//
//	negative := func(ctx tbl.StyleContext) bool {
//	    val, err := strconv.ParseFloat(ctx.Content, 64)
//	    return err == nil && val < 0
//	}
//	t.SetStyleFunc(
//	    tbl.FirstRow(tbl.BBottom()),
//	    tbl.Predicate(negative, tbl.Red()),
//	)
//
// Example - Custom predicate:
//
//	isSpecialRow := func(ctx tbl.StyleContext) bool {
//	    return ctx.Row%5 == 0 && ctx.Row > 0
//	}
//	t.SetStyleFunc(
//	    tbl.FirstRow(tbl.BBottom()),
//	    tbl.Predicate(isSpecialRow, tbl.BTop()),
//	)
func (t *Table) SetStyleFunc(fns ...Funcstyler) *Table {
	t.styleFunc = composeFuncstylers(fns...)
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

// composeFuncstylers combines multiple Funcstylers into one.
// Functions are evaluated left-to-right, with results merged sequentially.
// Returns nil if no functions provided.
func composeFuncstylers(fns ...Funcstyler) Funcstyler {
	if len(fns) == 0 {
		return nil
	}

	return func(ctx StyleContext) CellStyle {
		style := CellStyle{}
		for _, fn := range fns {
			if fn != nil {
				applied := fn(ctx)
				style = style.merge(applied)
			}
		}
		return style
	}
}

// resolveStyle returns effective style for cell using hierarchy.
// Resolution order: defaultStyle < columnStyles < rowStyles < styleFunc < cellStyles
// Only considers cell origin position for multi-span cells.
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
		ctx := StyleContext{
			Row:      cell.r,
			Col:      cell.c,
			RowCount: t.g.Rows(),
			ColCount: t.g.Cols(),
			Content:  cell.content,
		}
		fs := t.styleFunc(ctx)
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
