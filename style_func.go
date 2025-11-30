package tbl

// FirstRow applies styling to the first row (row 0).
//
// Commonly used for header styling with bottom borders:
//
//	t.SetStyleFunc(tbl.FirstRow(tbl.BBottom(), tbl.Bold()))
func FirstRow(stylers ...Freestyler) Funcstyler {
	return func(row, col, rowCount, colCount int) CellStyle {
		if row == 0 {
			return NewStyle(stylers...)
		}
		return CellStyle{}
	}
}

// LastRow applies styling to the last row.
//
// Commonly used for footer styling with top borders:
//
//	t.SetStyleFunc(tbl.LastRow(tbl.BTop(), tbl.Bold()))
func LastRow(stylers ...Freestyler) Funcstyler {
	return func(row, col, rowCount, colCount int) CellStyle {
		if row == rowCount-1 {
			return NewStyle(stylers...)
		}
		return CellStyle{}
	}
}

// RowRange applies styling to rows in range [start, end) (exclusive end).
//
// Example - Style rows 2-5:
//
//	t.SetStyleFunc(tbl.RowRange(2, 5, tbl.BBottom()))
func RowRange(start, end int, stylers ...Freestyler) Funcstyler {
	return func(row, col, rowCount, colCount int) CellStyle {
		if row >= start && row < end {
			return NewStyle(stylers...)
		}
		return CellStyle{}
	}
}

// EvenRows applies styling to even-numbered rows (0, 2, 4, ...).
//
// Example - Zebra striping:
//
//	t.SetStyleFunc(tbl.EvenRows(tbl.Pad(1)))
func EvenRows(stylers ...Freestyler) Funcstyler {
	return func(row, col, rowCount, colCount int) CellStyle {
		if row%2 == 0 {
			return NewStyle(stylers...)
		}
		return CellStyle{}
	}
}

// OddRows applies styling to odd-numbered rows (1, 3, 5, ...).
//
// Example - Inverse zebra striping:
//
//	t.SetStyleFunc(tbl.OddRows(tbl.Pad(1)))
func OddRows(stylers ...Freestyler) Funcstyler {
	return func(row, col, rowCount, colCount int) CellStyle {
		if row%2 == 1 {
			return NewStyle(stylers...)
		}
		return CellStyle{}
	}
}

// EvenRowsSkipN applies styling to even rows, skipping the first n rows.
//
// Commonly used for zebra striping with header exclusion:
//
//	// Skip first row (header)
//	t.SetStyleFunc(
//	    tbl.FirstRow(tbl.BBottom()),
//	    tbl.EvenRowsSkipN(1, tbl.Pad(1)),
//	)
//
//	// Skip first two rows (header + subheader)
//	t.SetStyleFunc(
//	    tbl.RowRange(0, 2, tbl.BBottom()),
//	    tbl.EvenRowsSkipN(2, tbl.Pad(1)),
//	)
func EvenRowsSkipN(n int, stylers ...Freestyler) Funcstyler {
	return func(row, col, rowCount, colCount int) CellStyle {
		if row >= n && row%2 == 0 {
			return NewStyle(stylers...)
		}
		return CellStyle{}
	}
}

// OddRowsSkipN applies styling to odd rows, skipping the first n rows.
//
// Useful for inverse zebra striping with header exclusion:
//
//	// Skip first row (header)
//	t.SetStyleFunc(
//	    tbl.FirstRow(tbl.BBottom()),
//	    tbl.OddRowsSkipN(1, tbl.Pad(1)),
//	)
//
//	// Skip first two rows (header + subheader)
//	t.SetStyleFunc(
//	    tbl.RowRange(0, 2, tbl.BBottom()),
//	    tbl.OddRowsSkipN(2, tbl.Pad(1)),
//	)
func OddRowsSkipN(n int, stylers ...Freestyler) Funcstyler {
	return func(row, col, rowCount, colCount int) CellStyle {
		if row >= n && row%2 == 1 {
			return NewStyle(stylers...)
		}
		return CellStyle{}
	}
}

// FirstCol applies styling to the first column (col 0).
//
// Commonly used for row label styling:
//
//	t.SetStyleFunc(tbl.FirstCol(tbl.BRight()))
func FirstCol(stylers ...Freestyler) Funcstyler {
	return func(row, col, rowCount, colCount int) CellStyle {
		if col == 0 {
			return NewStyle(stylers...)
		}
		return CellStyle{}
	}
}

// LastCol applies styling to the last column.
//
// Commonly used for summary column styling:
//
//	t.SetStyleFunc(tbl.LastCol(tbl.BLeft()))
func LastCol(stylers ...Freestyler) Funcstyler {
	return func(row, col, rowCount, colCount int) CellStyle {
		if col == colCount-1 {
			return NewStyle(stylers...)
		}
		return CellStyle{}
	}
}

// ColRange applies styling to columns in range [start, end) (exclusive end).
//
// Example - Style columns 1-3:
//
//	t.SetStyleFunc(tbl.ColRange(1, 3, tbl.Right()))
func ColRange(start, end int, stylers ...Freestyler) Funcstyler {
	return func(row, col, rowCount, colCount int) CellStyle {
		if col >= start && col < end {
			return NewStyle(stylers...)
		}
		return CellStyle{}
	}
}

// Predicate applies styling when custom condition is met.
//
// The predicate function receives cell position and table dimensions.
// Returns true to apply styling, false to skip.
//
// Example - Checkered pattern:
//
//	checkered := func(row, col, _, _ int) bool {
//	    return (row+col)%2 == 0
//	}
//	t.SetStyleFunc(tbl.Predicate(checkered, tbl.Pad(1)))
//
// Example - Custom row condition:
//
//	isSpecialRow := func(row, col, _, _ int) bool {
//	    return row == 5 || row == 10
//	}
//	t.SetStyleFunc(tbl.Predicate(isSpecialRow, tbl.BTop()))
func Predicate(fn func(row, col, rowCount, colCount int) bool, stylers ...Freestyler) Funcstyler {
	return func(row, col, rowCount, colCount int) CellStyle {
		if fn(row, col, rowCount, colCount) {
			return NewStyle(stylers...)
		}
		return CellStyle{}
	}
}
