package main

import (
	"fmt"

	"github.com/neox5/tbl"
)

// ANSI color codes
const (
	Reset  = "\x1b[0m"
	Red    = "\x1b[31m"
	Green  = "\x1b[32m"
	Yellow = "\x1b[33m"
	Blue   = "\x1b[34m"
	Bold   = "\x1b[1m"
)

func main() {
	fmt.Println("=== ANSI Color Support Example ===")
	fmt.Println()

	// Example 1: Simple colored content
	fmt.Println("Example 1: Basic Colorization")
	tbl.Simple(
		tbl.Row("Status", "Message"),
		tbl.Row(Green+"Success"+Reset, "Operation completed"),
		tbl.Row(Yellow+"Warning"+Reset, "Check configuration"),
		tbl.Row(Red+"Error"+Reset, "Connection failed"),
	).SetDefaultStyle(tbl.BAll()).Print()

	fmt.Println()

	// Example 2: Comparison table (dynamite use case)
	fmt.Println("Example 2: Configuration Comparison")
	t := tbl.New()

	t.AddRow(tbl.C("Field"), tbl.C("Local"), tbl.C("Remote"))
	t.AddRow(
		tbl.C("namingPattern"),
		tbl.C(Green+"/api/{service}"+Reset),
		tbl.C(Green+"/api/{service}"+Reset),
	)
	t.AddRow(
		tbl.C("enabled"),
		tbl.C(Red+"true"+Reset),
		tbl.C(Red+"false"+Reset),
	)
	t.AddRow(
		tbl.C("order"),
		tbl.C(Yellow+"100"+Reset),
		tbl.C(Yellow+"200"+Reset),
	)

	t.SetDefaultStyle(tbl.BAll())
	t.SetRowStyle(0, tbl.BBottom(), tbl.Center())
	t.SetColStyle(0, tbl.Left())
	t.SetColStyle(1, tbl.Right())
	t.SetColStyle(2, tbl.Right())

	t.Print()

	fmt.Println()

	// Example 3: Word wrapping with ANSI codes
	fmt.Println("Example 3: Word Wrapping with Colors")
	t2 := tbl.New()
	t2.AddCol(0, 15, 20) // Description
	t2.AddCol(0, 30, 40) // Details

	t2.AddRow(tbl.C("Type"), tbl.C("Description"))
	t2.AddRow(
		tbl.C(Bold+"Success"+Reset),
		tbl.C(Green+"This is a very long success message that demonstrates word wrapping with ANSI color codes preserved correctly"+Reset),
	)
	t2.AddRow(
		tbl.C(Bold+"Error"+Reset),
		tbl.C(Red+"This is a very long error message that shows how ANSI escape sequences are handled during text wrapping operations"+Reset),
	)

	t2.SetDefaultStyle(tbl.BAll())
	t2.SetRowStyle(0, tbl.BBottom(), tbl.Center())

	t2.Print()

	fmt.Println()

	// Example 4: Mixed formatting
	fmt.Println("Example 4: Mixed Formatting (Bold + Color)")
	tbl.Simple(
		tbl.Row("Item", "Status", "Notes"),
		tbl.Row("Config A", Bold+Green+"✓ Match"+Reset, "No changes"),
		tbl.Row("Config B", Bold+Yellow+"~ Differ"+Reset, "Value mismatch"),
		tbl.Row("Config C", Bold+Red+"✗ Missing"+Reset, "Not in remote"),
	).SetDefaultStyle(tbl.BAll(), tbl.Pad(0, 1)).Print()

	fmt.Println()

	// Example 5: Truncation with ANSI codes
	fmt.Println("Example 5: Truncation Preserves ANSI Codes")
	t3 := tbl.New()
	t3.AddCol(15, 0, 0) // Fixed width column

	t3.AddRow(tbl.C("Truncated"))
	t3.AddRow(tbl.C(Red + "This text is much too long and will be truncated with ellipsis" + Reset))
	t3.AddRow(tbl.C(Green + "Short text" + Reset))

	t3.SetDefaultStyle(tbl.BAll())
	t3.SetRowStyle(0, tbl.BBottom())
	t3.SetRowStyle(1, tbl.WrapTruncate)

	t3.Print()
}
