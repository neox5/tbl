package tbl

// Direction indicates traversal direction for flex cell scanning.
type Direction int

const (
	DirLeft Direction = iota
	DirRight
)

// flexCell represents a flex cell found during traversal.
type flexCell struct {
	cell      *Cell
	addedSpan int
}

// traverseFlex scans row for flex cells in both directions from col position.
// Recursively processes flex cells found in row above (row-1).
// Returns success status and flex cells organized by row.
//
// Return structure: map[int][]flexCell where key = row
func (t *Table) traverseFlex(row, col int) (bool, map[int][]flexCell) {
	if row < 0 {
		return true, nil
	}

	var canTraverse bool
	result := make(map[int][]flexCell)
	seen := make(map[ID]bool)

	// Scan left
	if t.traverseFlexDir(row, col, DirLeft, result, seen) {
		canTraverse = true
	}

	// Scan right (includes origin col)
	if t.traverseFlexDir(row, col, DirRight, result, seen) {
		canTraverse = true
	}

	return canTraverse, result
}

// traverseFlexDir scans single direction from col position in row.
// Stops at walls or grid boundaries.
// For each flex cell found, recursively calls traverseFlex(row-1, flexCol).
func (t *Table) traverseFlexDir(row, col int, dir Direction, result map[int][]flexCell, seen map[ID]bool) bool {
	if row == 0 {
		return true
	}

	var canTraverse bool

	// Determine iteration bounds
	start, end, step := col, t.g.Cols(), 1
	if dir == DirLeft {
		start, end, step = col-1, -1, -1
	}

	// Iterate in direction
	for c := start; c != end; c += step {
		// Check for wall
		if t.isWall(row, c) {
			break
		}

		// Check above if flex cell
		if t.isFlex(row-1, c) {
			cell := t.getCellAt(row-1, c)

			// Skip if already seen
			if seen[cell.id] {
				continue
			}
			seen[cell.id] = true

			// Add to current row result
			result[row-1] = append(result[row-1], flexCell{
				cell:      cell,
				addedSpan: cell.AddedSpan(),
			})

			// Recurse to row above
			ok, aboveResult := t.traverseFlex(row-1, c)
			if ok {
				canTraverse = true
			}

			// Merge results
			for r, flexCells := range aboveResult {
				for _, fc := range flexCells {
					if !seen[fc.cell.id] {
						seen[fc.cell.id] = true
						result[r] = append(result[r], fc)
					}
				}
			}
		}
	}

	return canTraverse
}
