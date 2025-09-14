package tbl

func (t *Table) nextCol() int {
	for i, l := range t.colLevels {
		if l == 0 {
			return i
		}
	}
	return NO_COL
}

func (t *Table) reduceColLevels() {
	for i, l := range t.colLevels {
		if l == FLEX {
			continue
		}
		if l > 0 {
			t.colLevels[i]--
		}
	}
}

func (t *Table) advanceRow() {
	t.currIndex = 0
	t.reduceColLevels()

	for t.nextCol() == NO_COL {
		t.reduceColLevels()
	}
}

func (t *Table) availableSpan() int {
	max, span := 0, 0
	for _, l := range t.colLevels {
		if l == 0 {
			span++
		} else {
			if span > max {
				max = span
			}
			span = 0
		}
	}
	if span > max {
		max = span
	}
	return max
}
