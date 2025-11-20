package tbl

import (
	"fmt"
	"io"

	"github.com/neox5/btmp"
)

// ID identifies a cell in the table.
type ID int64

// Table manages incremental table construction with flex/static cells.
type Table struct {
	g     *btmp.Grid
	cells map[ID]*Cell

	// Cursor state
	row int // current row position, -1 before first AddRow
	col int // current column position within row

	colsFixed  bool
	nextCellID ID

	// Configuration
	tableConfig TableConfig
	colConfigs  map[int]ColConfig

	// Style registry
	defaultStyle CellStyle
	columnStyles map[int]CellStyle
	rowStyles    map[int]CellStyle
	cellStyles   map[ID]CellStyle
}

// New creates a new Table with zero columns.
func New() *Table {
	return NewWithCols(0)
}

// NewWithCols creates a new Table with initial column capacity.
func NewWithCols(cols int) *Table {
	if cols < 0 {
		panic("tbl: invalid cols value")
	}

	t := &Table{
		g:          btmp.NewGridWithSize(1, max(cols, 1)),
		cells:      make(map[ID]*Cell),
		row:        -1,
		col:        0,
		nextCellID: 1,

		// Initialize configuration
		tableConfig: TableConfig{},
		colConfigs:  make(map[int]ColConfig),

		// Initialize style registry
		defaultStyle: CellStyle{
			Padding: Padding{Top: 0, Bottom: 0, Left: 1, Right: 1},
			HAlign:  HAlignLeft,
			VAlign:  VAlignTop,
			Border:  Border{Sides: BorderNone},
		},
		columnStyles: make(map[int]CellStyle),
		rowStyles:    make(map[int]CellStyle),
		cellStyles:   make(map[ID]CellStyle),
	}

	if cols > 0 {
		t.colsFixed = true
	}

	return t
}

// AddRow advances to next row with validation and cursor positioning.
func (t *Table) AddRow() *Table {
	t.addRow()
	return t
}

// AddCell adds a cell at cursor position with specified type and span.
// Expands columns if needed (when not fixed). Validates span fits in grid.
// Panics if: no row active, span invalid, insufficient columns (when fixed),
// or space occupied.
func (t *Table) AddCell(ct CellType, rowSpan, colSpan int, content string) *Table {
	if rowSpan <= 0 || colSpan <= 0 {
		panic(fmt.Sprintf("tbl: invalid span rowSpan=%d colSpan=%d at cursor (%d,%d)", rowSpan, colSpan, t.row, t.col))
	}

	if t.row == -1 {
		panic(fmt.Sprintf("tbl: no row to add cell at cursor (%d,%d)", t.row, t.col))
	}

	t.addCell(ct, rowSpan, colSpan, content)
	return t
}

// Render returns the ASCII table as a string.
func (t *Table) Render() string {
	t.finalize()
	return newRenderer(t).render()
}

// RenderTo writes the table to w.
func (t *Table) RenderTo(w io.Writer) error {
	s := t.Render()
	_, err := io.WriteString(w, s)
	return err
}

// Print prints the rendered output to stdout.
func (t *Table) Print() {
	fmt.Print(t.Render())
}

// PrintDebug renders table structure in TBL Grid Notation format.
// Shows grid layout with cell types and current cursor position.
// Returns empty string if table has no rows.
// For debug/development purposes.
func (t *Table) PrintDebug() string {
	return t.printDebug()
}
