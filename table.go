package tbl

type Table struct {
	config Config
	state  state
	rows   []Row
}

type state struct {
	colLevels         []int
	colWidths         []int
	rowHeights        []int
	currentRow        int
	stillFlexibleCols bool
}

func New(cfg Config) *Table {
	return &Table{
		config: cfg,
		state:  state{colLevels: []int{}, colWidths: []int{}},
		rows:   []Row{},
	}
}
