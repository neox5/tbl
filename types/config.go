package types

// Config represents the public configuration for tables
type Config struct {
	Border      *TableBorder
	NewCellFunc func() any // Function to create default cells
	Width       int
	MaxWidth    int
}
