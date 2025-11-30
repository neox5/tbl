package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/neox5/tbl"
)

// ANSI color codes
const (
	Reset = "\x1b[0m"
	Green = "\x1b[32m"
	Red   = "\x1b[31m"
)

type FinnhubQuote struct {
	C  float64 `json:"c"`  // Current price
	D  float64 `json:"d"`  // Change
	DP float64 `json:"dp"` // Percent change
	H  float64 `json:"h"`  // High
	L  float64 `json:"l"`  // Low
	O  float64 `json:"o"`  // Open
	PC float64 `json:"pc"` // Previous close
	T  int64   `json:"t"`  // Timestamp
}

type StockQuote struct {
	Symbol        string  `json:"symbol"`
	CurrentPrice  float64 `json:"current_price"`
	Change        float64 `json:"change"`
	PercentChange float64 `json:"percent_change"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Open          float64 `json:"open"`
	PreviousClose float64 `json:"previous_close"`
	Timestamp     int64   `json:"timestamp"`
}

func main() {
	dataPath := filepath.Join("examples", "02_column_sizing__financial", "data.json")

	// Check if data.json exists
	if _, err := os.Stat(dataPath); os.IsNotExist(err) {
		fmt.Println("data.json not found, fetching quotes...")
		if err := fetchQuotes(dataPath); err != nil {
			fmt.Printf("Error fetching quotes: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()
	}

	// Load data
	data, err := os.ReadFile(dataPath)
	if err != nil {
		panic(err)
	}

	var quotes []StockQuote
	if err = json.Unmarshal(data, &quotes); err != nil {
		panic(err)
	}

	// Build table with column sizing and styling
	t := tbl.New()
	t.AddCol(10, 0, 0, tbl.Left())   // Symbol: fixed width 10, left-aligned
	t.AddCol(0, 12, 14, tbl.Right()) // Price: min 12, max 14, right-aligned
	t.AddCol(0, 12, 14, tbl.Right()) // Change%: min 12, max 14, right-aligned
	t.AddCol(0, 10, 12, tbl.Right()) // High: min 10, max 12, right-aligned
	t.AddCol(0, 10, 12, tbl.Right()) // Low: min 10, max 12, right-aligned

	// Header row
	t.AddRow(tbl.C("Symbol"), tbl.C("Price"), tbl.C("Change %"),
		tbl.C("High"), tbl.C("Low"))

	// Data rows
	for _, q := range quotes {
		t.AddRow(
			tbl.C(q.Symbol),
			tbl.C(formatPrice(q.CurrentPrice)),
			tbl.C(formatChangePercent(q.PercentChange)),
			tbl.C(formatPrice(q.High)),
			tbl.C(formatPrice(q.Low)),
		)
	}

	// Print to stdout
	fmt.Println()
	t.Print()

	// Write to output.txt
	outputPath := filepath.Join("examples", "02_column_sizing__financial", "output.txt")
	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := t.RenderTo(f); err != nil {
		panic(err)
	}
}

func fetchQuotes(outputPath string) error {
	// Get API key from environment
	key := os.Getenv("FINNHUB_API_KEY")
	if key == "" {
		return fmt.Errorf("FINNHUB_API_KEY environment variable not set\nGet free API key at: https://finnhub.io")
	}

	// Top 20 global companies by market cap
	symbols := []string{
		"AAPL", "MSFT", "GOOGL", "AMZN", "NVDA",
		"META", "TSLA", "BRK.B", "LLY", "AVGO",
		"TSM", "V", "WMT", "JPM", "XOM",
		"NVO", "UNH", "MA", "JNJ", "PG",
	}

	var quotes []StockQuote

	fmt.Printf("Fetching quotes for %d stocks...\n\n", len(symbols))

	for i, symbol := range symbols {
		fmt.Printf("[%2d/%d] Fetching %s... ", i+1, len(symbols), symbol)

		url := fmt.Sprintf("https://finnhub.io/api/v1/quote?symbol=%s&token=%s", symbol, key)

		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}

		var fq FinnhubQuote
		if err := json.Unmarshal(body, &fq); err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}

		quotes = append(quotes, StockQuote{
			Symbol:        symbol,
			CurrentPrice:  fq.C,
			Change:        fq.D,
			PercentChange: fq.DP,
			High:          fq.H,
			Low:           fq.L,
			Open:          fq.O,
			PreviousClose: fq.PC,
			Timestamp:     fq.T,
		})

		fmt.Printf("$%.2f (%+.2f%%)\n", fq.C, fq.DP)

		// Rate limiting (60 calls/min = 1 per second)
		if i < len(symbols)-1 {
			time.Sleep(500 * time.Millisecond)
		}
	}

	// Save to data.json (compact format)
	data, err := json.Marshal(quotes)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, data, 0o644); err != nil {
		return err
	}

	fmt.Printf("\nSuccess! Saved %d quotes to data.json\n", len(quotes))
	return nil
}

func formatPrice(p float64) string {
	return fmt.Sprintf("$%.2f", p)
}

func formatChangePercent(c float64) string {
	color := Green
	if c < 0 {
		color = Red
	}
	return fmt.Sprintf("%s%+.2f%%%s", color, c, Reset)
}
