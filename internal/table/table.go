package table

import (
	"fmt"

	"github.com/neox5/tbl/internal/cell"
	"github.com/neox5/tbl/types"
)

// Table is the concrete table implementation
type Table struct {
	// Configuration (flattened)
	border   *types.TableBorder
	width    int
	maxWidth int

	// Table state
	cells         []*cell.Cell
	rowStarts     []int
	colWidths     []int
	colLevels     []int
	hLines        []bool
	openFlexCells []int
	nextIndex     int // next cell index
	row, col      int // next row/col index

	// Indices for optimization (future use)
	colIndex map[int][]int // cells overlapping a column
	rowIndex map[int][]int // cells overlapping a row
}

// New creates a new table with default configuration
func New() *Table {
	return &Table{
		border:        &types.DefaultTableBorder,
		width:         0,
		maxWidth:      0,
		cells:         []*cell.Cell{},
		rowStarts:     []int{},
		colWidths:     []int{},
		colLevels:     []int{},
		openFlexCells: []int{},
		nextIndex:     0,
		row:           0,
		col:           0,
		colIndex:      make(map[int][]int),
		rowIndex:      make(map[int][]int),
	}
}

// NewWithConfig creates a new table with merged configuration
func NewWithConfig(cfg *types.Config) *Table {
	t := New()

	if cfg != nil {
		if cfg.Border != nil {
			t.border = cfg.Border
		}
		if cfg.Width > 0 {
			t.width = cfg.Width
		}
		if cfg.MaxWidth > 0 {
			t.maxWidth = cfg.MaxWidth
		}
	}

	return t
}

func (t *Table) ColCount() int {
	return len(t.colWidths)
}

// AddRow adds a new row with the specified cells
func (t *Table) AddRow(cells ...*cell.Cell) {
	t.startNewRow()
	if t.row == 0 {
		t.initializeWithFirstRow(cells)
	}
	t.addCells(cells)
	t.validateRow()
	t.advanceRow()
}

// initializeWithFirstRow initializes colWidths and colLevels so that addCells
// works also with the first row
func (t *Table) initializeWithFirstRow(cells []*cell.Cell) {
	totalCols := 0
	for _, c := range cells {
		cSpan, _ := c.Span()
		if cSpan == cell.FLEX {
			cSpan = 1
		}
		totalCols += cSpan
	}

	t.colWidths = make([]int, totalCols)
	t.colLevels = make([]int, totalCols)
}

// advanceCol advances the column index to the next free column (colLevel = 0)
func (t *Table) advanceCol() {
	for t.col < t.ColCount() && t.colLevels[t.col] > 0 {
		t.col++
	}

	// If we've reached the end of the row, reset to column 0 for next row
	if t.col >= t.ColCount() {
		t.col = 0
	}
}

// validateRow validates if all 0 colLevels were increased.
// If not we need to panic as the row is incomplete.
func (t *Table) validateRow() {
	for i, level := range t.colLevels {
		if level == 0 {
			panic(fmt.Sprintf("incomplete row %d: column %d was not filled", t.row, i))
		}
	}
}

// advanceRow advances to the next row and decrements colLevels
func (t *Table) advanceRow() {
	t.row++

	hasFreeCol := false
	for i := range t.colLevels {
		t.colLevels[i]--
		if t.colLevels[i] == 0 {
			hasFreeCol = true
		}
	}

	if !hasFreeCol {
		t.advanceRow() // Recursive call if no columns are free
	}
}

// startNewRow starts a new row by recording the current cell position
func (t *Table) startNewRow() {
	t.rowStarts = append(t.rowStarts, len(t.cells))
}

// addCells appends multiple cells to the current row
func (t *Table) addCells(cells []*cell.Cell) {
	for _, c := range cells {
		cSpan, rSpan := c.Span()
		if !t.spanFit(cSpan) {
			panic(fmt.Sprintf("cell span %d does not fit at row %d col %d", cSpan, t.row, t.col))
		}

		w := c.ColWidth()
		for i := range cSpan {
			// update colWidths when cell width is wider then the current value
			if t.colWidths[t.col+i] < w {
				t.colWidths[t.col+i] = w
			}
			t.colLevels[t.col+i] = rSpan // colLevels will always be set; >0 or -1 (FLEX)
		}

		t.advanceCol()

		t.cells = append(t.cells, c)
		t.nextIndex++
	}
}

// spanFit returns true if the span can fit at the current row position
func (t *Table) spanFit(span int) bool {
	r := 0
	for i := t.col; i < t.ColCount(); i++ {
		if t.colLevels[i] == 0 {
			r++
			continue
		}
		break
	}
	return r >= span
}
