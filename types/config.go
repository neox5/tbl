package types

// Config represents the public configuration for tables
type Config struct {
	Border      *TableBorder
	DefaultCell any // Interface type for flexibility
	Width       int
	MaxWidth    int
}
