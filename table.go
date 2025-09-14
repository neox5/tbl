package tbl

type Table struct {
	// config fields
	border      TableBorder
	cellDefault Cell
	width       int
	maxWidth    int

	// table state
	cells         []Cell
	rowStarts     []int
	colWidths     []int
	colLevels     []int
	hLines        []bool
	currIndex     int
	openFlexCells []int

	// indices
	colIndex map[int][]int // index for cells overlapping a column
	rowIndex map[int][]int // index for cells overlapping a row
}

func New() *Table {
	return &Table{
		border:      DefaultTableBorder,
		cellDefault: DefaultCell,
		width:       0,
		maxWidth:    0,

		cells:         []Cell{},
		rowStarts:     []int{},
		colWidths:     []int{},
		colLevels:     []int{},
		openFlexCells: []int{},
	}
}

func NewWithConfig(cfg Config) *Table {
	t := New()

	if cfg.Border != nil {
		t.border = *cfg.Border
	}
	if cfg.CellDefault != nil {
		t.cellDefault = *cfg.CellDefault
	}
	if cfg.Width != 0 {
		t.width = cfg.Width
	}
	if cfg.MaxWidth != 0 {
		t.maxWidth = cfg.MaxWidth
	}

	return t
}
