// Package display handles terminal output formatting with colors and progress bars.
package display

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/TungCorn/antigravity-usage-checker/internal/api"
)

const (
	// ANSI color codes
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
)

// ShowUsage displays the usage data in the terminal.
func ShowUsage(data *api.UsageData, asJSON bool, isCached bool) {
	if asJSON {
		showJSON(data)
		return
	}
	
	showTable(data, isCached)
}

// showJSON outputs the usage data as formatted JSON.
func showJSON(data *api.UsageData) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error formatting JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

// showTable displays the usage data as a formatted table.
func showTable(data *api.UsageData, isCached bool) {
	// Header
	fmt.Println()
	fmt.Printf("%s%s🚀 Antigravity Usage Monitor%s\n", Bold, Cyan, Reset)
	fmt.Println(strings.Repeat("─", 55))
	
	// Cache indicator
	if isCached || data.IsCached {
		fmt.Printf("%s⚠️  Cached data from %s%s\n", Yellow, formatTime(data.FetchedAt), Reset)
		fmt.Println(strings.Repeat("─", 55))
	}
	
	// Table header
	fmt.Printf("%s%-20s %8s %8s %8s %s%s\n", Bold,
		"Model", "Used", "Limit", "Left", "Progress", Reset)
	fmt.Println(strings.Repeat("─", 55))
	
	// Model rows
	for _, model := range data.Models {
		color := getStatusColor(model.UsagePercent)
		progressBar := createProgressBar(model.UsagePercent, 15)
		
		fmt.Printf("%-20s %s%8d%s %8d %s%8d%s %s\n",
			truncateString(model.ModelName, 20),
			Magenta, model.Used, Reset,
			model.Limit,
			color, model.Remaining, Reset,
			progressBar,
		)
	}
	
	// Footer
	fmt.Println(strings.Repeat("─", 55))
	
	// Tier and credits
	footer := []string{}
	if data.Tier != "" {
		footer = append(footer, fmt.Sprintf("Tier: %s%s%s", Cyan, data.Tier, Reset))
	}
	if data.PromptCredit > 0 {
		footer = append(footer, fmt.Sprintf("Credits: %s%d%s", Green, data.PromptCredit, Reset))
	}
	
	if len(footer) > 0 {
		fmt.Printf("%s%s%s\n", Dim, strings.Join(footer, " | "), Reset)
	}
	
	fmt.Println()
}

// createProgressBar generates a visual progress bar.
func createProgressBar(percent, width int) string {
	if percent > 100 {
		percent = 100
	}
	if percent < 0 {
		percent = 0
	}
	
	filled := width * percent / 100
	empty := width - filled
	
	color := getStatusColor(percent)
	
	bar := fmt.Sprintf("%s%s%s%s %d%%",
		color,
		strings.Repeat("█", filled),
		strings.Repeat("░", empty),
		Reset,
		percent,
	)
	
	return bar
}

// getStatusColor returns the appropriate color based on usage percentage.
func getStatusColor(percent int) string {
	switch {
	case percent < 50:
		return Green
	case percent < 80:
		return Yellow
	default:
		return Red
	}
}

// truncateString limits string length with ellipsis.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// formatTime formats a time.Time for display.
func formatTime(t time.Time) string {
	return t.Format("15:04:05 02/01/2006")
}
