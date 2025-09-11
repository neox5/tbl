package tbl

type Table struct {
	config       Config
	rows         []Row
	virtualRows  []int
	colLevels    []int
	colWidths    []int
	hLines       []bool
	flexibleCols bool
	width        int
	col, row     int
}

func New() *Table {
	return NewWithConfig(DefaultConfig)
}

func NewWithConfig(cfg Config) *Table {
	mergedCfg := DefaultConfig.Merge(cfg)
	return &Table{
		config:      mergedCfg,
		rows:        []Row{},
		virtualRows: []int{},
		colLevels:   []int{},
		colWidths:   []int{},
	}
}
