package tbl

// Predicate applies styling when custom condition is met.
//
// The predicate function receives StyleContext with cell position,
// table dimensions, and cell content.
// Returns true to apply styling, false to skip.
//
// Example - Checkered pattern:
//
//	checkered := func(ctx tbl.StyleContext) bool {
//	    return (ctx.Row+ctx.Col)%2 == 0
//	}
//	t.SetStyleFunc(tbl.Predicate(checkered, tbl.Pad(1)))
//
// Example - Content-based styling:
//
//	negative := func(ctx tbl.StyleContext) bool {
//	    val, err := strconv.ParseFloat(ctx.Content, 64)
//	    return err == nil && val < 0
//	}
//	t.SetStyleFunc(tbl.Predicate(negative, tbl.Red()))
//
// Example - Column-specific content rules:
//
//	important := func(ctx tbl.StyleContext) bool {
//	    return ctx.Col == 2 && strings.Contains(ctx.Content, "URGENT")
//	}
//	t.SetStyleFunc(tbl.Predicate(important, tbl.Bold(), tbl.BgYellow()))
//
// Example - Complex conditions with combinators:
//
//	evenNotFirst := tbl.And(
//	    func(ctx tbl.StyleContext) bool { return ctx.Row%2 == 0 },
//	    tbl.Not(func(ctx tbl.StyleContext) bool { return ctx.Row == 0 }),
//	)
//	t.SetStyleFunc(tbl.Predicate(evenNotFirst, tbl.Pad(1)))
func Predicate(fn func(ctx StyleContext) bool, stylers ...Freestyler) Funcstyler {
	return func(ctx StyleContext) CellStyle {
		if fn(ctx) {
			return NewStyle(stylers...)
		}
		return CellStyle{}
	}
}

// And combines multiple predicates with logical AND.
// Returns true only if all predicates return true.
//
// Example - Even rows excluding first row:
//
//	evenNotFirst := tbl.And(
//	    func(ctx tbl.StyleContext) bool { return ctx.Row%2 == 0 },
//	    func(ctx tbl.StyleContext) bool { return ctx.Row > 0 },
//	)
//	t.SetStyleFunc(tbl.Predicate(evenNotFirst, tbl.Pad(1)))
//
// Example - Specific cell styling:
//
//	topRightCorner := tbl.And(
//	    func(ctx tbl.StyleContext) bool { return ctx.Row == 0 },
//	    func(ctx tbl.StyleContext) bool { return ctx.Col == ctx.ColCount-1 },
//	)
//	t.SetStyleFunc(tbl.Predicate(topRightCorner, tbl.Bold()))
func And(predicates ...func(StyleContext) bool) func(StyleContext) bool {
	return func(ctx StyleContext) bool {
		for _, p := range predicates {
			if !p(ctx) {
				return false
			}
		}
		return true
	}
}

// Or combines multiple predicates with logical OR.
// Returns true if any predicate returns true.
//
// Example - Style first or last row:
//
//	firstOrLast := tbl.Or(
//	    func(ctx tbl.StyleContext) bool { return ctx.Row == 0 },
//	    func(ctx tbl.StyleContext) bool { return ctx.Row == ctx.RowCount-1 },
//	)
//	t.SetStyleFunc(tbl.Predicate(firstOrLast, tbl.Bold()))
//
// Example - Multiple column styling:
//
//	col0or2 := tbl.Or(
//	    func(ctx tbl.StyleContext) bool { return ctx.Col == 0 },
//	    func(ctx tbl.StyleContext) bool { return ctx.Col == 2 },
//	)
//	t.SetStyleFunc(tbl.Predicate(col0or2, tbl.Right()))
func Or(predicates ...func(StyleContext) bool) func(StyleContext) bool {
	return func(ctx StyleContext) bool {
		for _, p := range predicates {
			if p(ctx) {
				return true
			}
		}
		return false
	}
}

// Not inverts a predicate's result.
// Returns true if the predicate returns false, and vice versa.
//
// Example - All rows except first:
//
//	notFirst := tbl.Not(func(ctx tbl.StyleContext) bool { return ctx.Row == 0 })
//	t.SetStyleFunc(tbl.Predicate(notFirst, tbl.Pad(1)))
//
// Example - Odd rows excluding first (alternative to OddRowsSkipN):
//
//	oddNotFirst := tbl.And(
//	    func(ctx tbl.StyleContext) bool { return ctx.Row%2 == 1 },
//	    tbl.Not(func(ctx tbl.StyleContext) bool { return ctx.Row == 0 }),
//	)
//	t.SetStyleFunc(tbl.Predicate(oddNotFirst, tbl.Pad(1)))
func Not(predicate func(StyleContext) bool) func(StyleContext) bool {
	return func(ctx StyleContext) bool {
		return !predicate(ctx)
	}
}

// FirstRow applies styling to the first row (row 0).
//
// Commonly used for header styling with bottom borders:
//
//	t.SetStyleFunc(tbl.FirstRow(tbl.BBottom(), tbl.Bold()))
func FirstRow(stylers ...Freestyler) Funcstyler {
	return Predicate(func(ctx StyleContext) bool { return ctx.Row == 0 }, stylers...)
}

// LastRow applies styling to the last row.
//
// Commonly used for footer styling with top borders:
//
//	t.SetStyleFunc(tbl.LastRow(tbl.BTop(), tbl.Bold()))
func LastRow(stylers ...Freestyler) Funcstyler {
	return Predicate(func(ctx StyleContext) bool { return ctx.Row == ctx.RowCount-1 }, stylers...)
}

// RowRange applies styling to rows in range [start, end) (exclusive end).
//
// Example - Style rows 2-5:
//
//	t.SetStyleFunc(tbl.RowRange(2, 5, tbl.BBottom()))
func RowRange(start, end int, stylers ...Freestyler) Funcstyler {
	return Predicate(func(ctx StyleContext) bool {
		return ctx.Row >= start && ctx.Row < end
	}, stylers...)
}

// EvenRows applies styling to even-numbered rows (0, 2, 4, ...).
//
// Example - Zebra striping:
//
//	t.SetStyleFunc(tbl.EvenRows(tbl.Pad(1)))
func EvenRows(stylers ...Freestyler) Funcstyler {
	return Predicate(func(ctx StyleContext) bool { return ctx.Row%2 == 0 }, stylers...)
}

// OddRows applies styling to odd-numbered rows (1, 3, 5, ...).
//
// Example - Inverse zebra striping:
//
//	t.SetStyleFunc(tbl.OddRows(tbl.Pad(1)))
func OddRows(stylers ...Freestyler) Funcstyler {
	return Predicate(func(ctx StyleContext) bool { return ctx.Row%2 == 1 }, stylers...)
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
	return Predicate(func(ctx StyleContext) bool {
		return ctx.Row >= n && ctx.Row%2 == 0
	}, stylers...)
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
	return Predicate(func(ctx StyleContext) bool {
		return ctx.Row >= n && ctx.Row%2 == 1
	}, stylers...)
}

// FirstCol applies styling to the first column (col 0).
//
// Commonly used for row label styling:
//
//	t.SetStyleFunc(tbl.FirstCol(tbl.BRight()))
func FirstCol(stylers ...Freestyler) Funcstyler {
	return Predicate(func(ctx StyleContext) bool { return ctx.Col == 0 }, stylers...)
}

// LastCol applies styling to the last column.
//
// Commonly used for summary column styling:
//
//	t.SetStyleFunc(tbl.LastCol(tbl.BLeft()))
func LastCol(stylers ...Freestyler) Funcstyler {
	return Predicate(func(ctx StyleContext) bool { return ctx.Col == ctx.ColCount-1 }, stylers...)
}

// ColRange applies styling to columns in range [start, end) (exclusive end).
//
// Example - Style columns 1-3:
//
//	t.SetStyleFunc(tbl.ColRange(1, 3, tbl.Right()))
func ColRange(start, end int, stylers ...Freestyler) Funcstyler {
	return Predicate(func(ctx StyleContext) bool {
		return ctx.Col >= start && ctx.Col < end
	}, stylers...)
}
