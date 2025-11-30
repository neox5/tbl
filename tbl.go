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

	// Programmable style hook
	styleFunc Funcstyler
}

// New creates a new Table with zero columns (dynamic sizing).
func New() *Table {
	return NewWithCols(0)
}

// NewWithCols creates a new Table with initial column capacity.
// If cols > 0, column count is fixed and cannot expand.
func NewWithCols(cols int) *Table {
	if cols < 0 {
		panic("tbl: invalid cols value")
	}

	t := &Table{
		g:          btmp.NewGridWithSize(1, cols),
		cells:      make(map[ID]*Cell),
		row:        -1,
		col:        0,
		colsFixed:  false,
		nextCellID: 1,

		// Initialize configuration
		tableConfig: TableConfig{},
		colConfigs:  make(map[int]ColConfig),

		// Initialize style registry
		defaultStyle: NewStyle(
			Pad(0, 1),
			Left(),
			Top(),
			BNone(),
			Thin(),
		),
		columnStyles: make(map[int]CellStyle),
		rowStyles:    make(map[int]CellStyle),
		cellStyles:   make(map[ID]CellStyle),
	}

	if cols > 0 {
		t.colsFixed = true
	}

	return t
}

// AddCol adds a column with width constraints and optional styling.
// Returns the column index.
// Must be called before first AddRow().
// Sets column count as fixed, but first row can still expand via flex cells.
// Subsequent rows cannot expand beyond established column count.
func (t *Table) AddCol(width, minWidth, maxWidth int, stylers ...Freestyler) int {
	// Validate: no rows yet
	if t.row >= 0 {
		panic("tbl: cannot add columns after AddRow")
	}

	// Validate: non-negative values
	if width < 0 || minWidth < 0 || maxWidth < 0 {
		panic(fmt.Sprintf("tbl: invalid column config width=%d minWidth=%d maxWidth=%d", width, minWidth, maxWidth))
	}

	// Validate: logical constraints
	if width > 0 && minWidth > 0 && width < minWidth {
		panic(fmt.Sprintf("tbl: width %d less than minWidth %d", width, minWidth))
	}
	if width > 0 && maxWidth > 0 && width > maxWidth {
		panic(fmt.Sprintf("tbl: width %d greater than maxWidth %d", width, maxWidth))
	}
	if minWidth > 0 && maxWidth > 0 && minWidth > maxWidth {
		panic(fmt.Sprintf("tbl: minWidth %d greater than maxWidth %d", minWidth, maxWidth))
	}

	col := t.g.Cols()
	t.addCol(width, minWidth, maxWidth)

	// Apply style if provided
	if len(stylers) > 0 {
		if containsTemplate(stylers) {
			panic("tbl: CharTemplate only supported via SetDefaultStyle")
		}
		t.columnStyles[col] = t.columnStyles[col].Apply(stylers...)
	}

	return col
}

// AddRow advances to next row and optionally adds cells.
// Returns the row index.
// No args: advances cursor only (compatible with explicit AddCell calls).
// With args: advances cursor and adds cells left-to-right.
//
// Example with CellSpec:
//
//	t.AddRow(tbl.C("Name"), tbl.C("Age"), tbl.F("Bio"))
//	t.AddRow(tbl.Cx(2, 1, "Merged"), tbl.C("30"), tbl.C("Engineer"))
//
// Example without args (current behavior):
//
//	t.AddRow()
//	t.AddCell(tbl.Static, 1, 1, "Data")
func (t *Table) AddRow(specs ...CellSpec) int {
	t.addRow(specs...)
	return t.row
}

// AddCell adds a cell at cursor position with specified type and span.
// Returns the cell ID.
// Expands columns if needed (when not fixed). Validates span fits in grid.
// Panics if: no row active, span invalid, insufficient columns (when fixed),
// or space occupied.
func (t *Table) AddCell(ct CellType, rowSpan, colSpan int, content string) ID {
	if rowSpan <= 0 || colSpan <= 0 {
		panic(fmt.Sprintf("tbl: invalid span rowSpan=%d colSpan=%d at cursor (%d,%d)", rowSpan, colSpan, t.row, t.col))
	}

	if t.row == -1 {
		panic(fmt.Sprintf("tbl: no row to add cell at cursor (%d,%d)", t.row, t.col))
	}

	id := t.nextCellID
	t.addCell(ct, rowSpan, colSpan, content)
	return id
}

// Row creates a row slice for use with Simple.
// Syntactic sugar for cleaner table construction.
//
// Example:
//
//	tbl.New().Simple(
//	    tbl.Row("Name", "Age", "City"),
//	    tbl.Row("Alice", "30", "NYC"),
//	    tbl.Row("Bob", "25", "LA"),
//	)
func Row(cells ...string) []string {
	return cells
}

// Simple adds rows of static [1,1] cells from string slices.
// Returns the table for method chaining.
//
// Example - basic usage:
//
//	tbl.New().Simple(
//	    tbl.Row("Name", "Age"),
//	    tbl.Row("Alice", "30"),
//	).Print()
//
// Example - mixed with other APIs:
//
//	t := tbl.New()
//	t.AddRow(tbl.F("Header"))
//	t.Simple(
//	    tbl.Row("Name", "Age"),
//	    tbl.Row("Alice", "30"),
//	)
func (t *Table) Simple(rows ...[]string) *Table {
	for _, row := range rows {
		t.AddRow()
		for _, cell := range row {
			t.AddCell(Static, 1, 1, cell)
		}
	}
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
