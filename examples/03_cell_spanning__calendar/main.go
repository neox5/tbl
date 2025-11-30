package main

import (
	"os"
	"path/filepath"

	"github.com/neox5/tbl"
)

// main demonstrates cell spanning with a two-month calendar.
// Shows Flex cell stretching (single and multiple), Static cell spanning,
// and mixed cell types.
func main() {
	table := buildCalendar()

	// Print to stdout
	table.Print()

	// Write to output.txt
	outputPath := filepath.Join("examples", "03_cell_spanning__calendar", "output.txt")
	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := table.RenderTo(f); err != nil {
		panic(err)
	}
}

// buildCalendar constructs a calendar table for November and December 2025.
// Demonstrates:
//   - F("2025") for year header (stretches across full width)
//   - F("November"), F("December") for month headers (equal stretch)
//   - C() for individual day headers and day numbers
//   - Cx(1, 3) for multi-day events
//   - Empty C("") cells for completing the grid
//   - Simple() method for bulk static cell rows
func buildCalendar() *tbl.Table {
	t := tbl.New()

	// Year header - single Flex cell spans entire table
	t.AddRow(tbl.F("2025"))

	// Month headers - two Flex cells get equal width
	t.AddRow(tbl.F("November"), tbl.F("December"))

	// Day headers and date rows using Simple() method
	t.Simple(
		// Day headers for both months
		tbl.Row("Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"),

		// Week 1: Nov starts on Saturday
		tbl.Row("", "", "", "", "", "1", "2", "1", "2", "3", "4", "5", "6", "7"),

		// Week 2
		tbl.Row("3", "4", "5", "6", "7", "8", "9", "8", "9", "10", "11", "12", "13", "14"),

		// Week 3
		tbl.Row("10", "11", "12", "13", "14", "15", "16", "15", "16", "17", "18", "19", "20", "21"),
	)

	// Week 4: November has 3-day conference event (Wed-Fri) - requires manual AddRow
	t.AddRow(
		tbl.C("17"), tbl.C("18"), tbl.Cx(1, 3, "Conference"), tbl.C("22"), tbl.C("23"),
		tbl.C("22"), tbl.C("23"), tbl.C("24"), tbl.C("25"), tbl.C("26"), tbl.C("27"), tbl.C("28"),
	)

	// Week 5: November ends on 30th, December continues
	t.Simple(
		tbl.Row("24", "25", "26", "27", "28", "29", "30", "29", "30", "31", "", "", "", ""),
	)

	// Apply styling
	t.SetDefaultStyle(tbl.BAll(), tbl.Pad(0, 1), tbl.Center())
	t.SetRowStyle(0, tbl.Bold())    // Year
	t.SetRowStyle(1, tbl.Bold())    // Months
	t.SetRowStyle(2, tbl.BBottom()) // Day headers

	return t
}
