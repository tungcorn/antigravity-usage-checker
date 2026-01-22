// Package display handles terminal output formatting with colors and progress bars.
package display

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/tungcorn/antigravity-usage-checker/internal/api"
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
	printHeader(data, isCached)

	// Calculate totals
	totalUsed, totalLimit, totalRemaining := calculateTotals(data.Models)

	// Print rows
	printRows(data.Models)

	printFooter(totalUsed, totalLimit, totalRemaining, data.Tier, data.PromptCredit)
	fmt.Println()
}

func printHeader(data *api.UsageData, isCached bool) {
	// Header with box drawing
	fmt.Println()
	fmt.Println("┌" + strings.Repeat("─", 68) + "┐")
	fmt.Printf("│ %s%s🚀 Antigravity Usage Monitor%s%-35s│\n", Bold, Cyan, Reset, "")

	// Cache indicator
	if isCached || data.IsCached {
		fmt.Println("├" + strings.Repeat("─", 68) + "┤")
		fmt.Printf("│ %s⚠️  Cached data from %s%s%-28s│\n", Yellow, formatTime(data.FetchedAt), Reset, "")
	}

	fmt.Println("├" + strings.Repeat("─", 68) + "┤")

	// Table header
	fmt.Printf("│ %-30s %-7s %-14s %-12s│\n",
		"Model", "Used", "Progress", "Reset")
	fmt.Println("├" + strings.Repeat("─", 68) + "┤")
}

type quotaKey struct {
	Used      float64
	Limit     float64
	Remaining float64
	ResetTime string
}

func calculateTotals(models []api.QuotaInfo) (float64, float64, float64) {
	uniqueQuotas := make(map[quotaKey]bool)
	var totalUsed, totalLimit, totalRemaining float64

	for _, model := range models {
		key := quotaKey{
			Used:      model.Used,
			Limit:     model.Limit,
			Remaining: model.Remaining,
			ResetTime: model.ResetTime,
		}

		if !uniqueQuotas[key] {
			uniqueQuotas[key] = true
			totalUsed += model.Used
			totalLimit += model.Limit
			totalRemaining += model.Remaining
		}
	}
	return totalUsed, totalLimit, totalRemaining
}

// formatPercent formats a percentage, showing decimals only when needed.
func formatPercent(percent float64) string {
	if percent == float64(int(percent)) {
		return fmt.Sprintf("%.0f%%", percent)
	}
	return fmt.Sprintf("%.1f%%", percent)
}

func printRows(models []api.QuotaInfo) {
	for _, model := range models {
		remainingPercent := 100 - model.UsagePercent
		color := getRemainingColor(remainingPercent)
		progressBar := createProgressBar(model.UsagePercent, 14)
		resetStr := formatResetTime(model.ResetTime)

		usedStr := formatPercent(model.Used)

		fmt.Printf("│ %-30s %s%-7s%s %s %-12s│\n",
			truncateString(model.ModelName, 28),
			color, usedStr, Reset,
			progressBar,
			truncateString(resetStr, 11),
		)
	}
}

func printFooter(used, limit, remaining float64, tier string, credits int) {
	fmt.Println(strings.Repeat("─", 68))

	// Total usage summary
	var totalUsagePercent float64
	if limit > 0 {
		totalUsagePercent = (used * 100) / limit
	}
	totalRemainingPercent := 100 - totalUsagePercent
	summaryColor := getRemainingColor(totalRemainingPercent)

	fmt.Printf("%s📊 Total: %.1f%% used (%.1f%% remaining)%s\n",
		summaryColor, totalUsagePercent, totalRemainingPercent, Reset)
	fmt.Println(strings.Repeat("─", 68))

	// Tier and credits
	var footer []string
	if tier != "" {
		footer = append(footer, fmt.Sprintf("Tier: %s%s%s", Cyan, tier, Reset))
	}
	if credits > 0 {
		footer = append(footer, fmt.Sprintf("Credits: %s%d%s", Green, credits, Reset))
	}

	if len(footer) > 0 {
		fmt.Printf("%s%s%s\n", Dim, strings.Join(footer, " | "), Reset)
	}
}

// createProgressBar generates a visual progress bar.
func createProgressBar(percent float64, width int) string {
	if percent > 100 {
		percent = 100
	}
	if percent < 0 {
		percent = 0
	}

	filled := int(float64(width) * percent / 100)
	empty := width - filled

	// Color based on remaining percentage
	remainingPercent := 100 - percent
	color := getRemainingColor(remainingPercent)

	bar := fmt.Sprintf("%s%s%s%s",
		color,
		strings.Repeat("█", filled),
		strings.Repeat("░", empty),
		Reset,
	)

	return bar
}

// getRemainingColor returns the appropriate color based on remaining percentage.
func getRemainingColor(remainingPercent float64) string {
	switch {
	case remainingPercent > 50:
		return Green
	case remainingPercent > 20:
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

// formatResetTime converts an ISO timestamp to relative time with exact time (e.g., "1h 30m (14:30)").
func formatResetTime(resetTimeStr string) string {
	if resetTimeStr == "" {
		return "-"
	}

	// Parse ISO 8601 timestamp
	resetTime, err := time.Parse(time.RFC3339, resetTimeStr)
	if err != nil {
		return "-"
	}

	now := time.Now()
	diff := resetTime.Sub(now)

	// If already reset
	if diff <= 0 {
		return Dim + "reset" + Reset
	}

	// Convert reset time to local timezone for display
	localResetTime := resetTime.Local()
	exactTime := localResetTime.Format("15:04")

	// Format as relative time with exact time
	hours := int(diff.Hours())
	minutes := int(diff.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%s%dh %dm (%s)%s", Dim, hours, minutes, exactTime, Reset)
	}
	return fmt.Sprintf("%s%dm (%s)%s", Dim, minutes, exactTime, Reset)
}
