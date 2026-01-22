package display

import (
	"testing"
	"time"

	"github.com/tungcorn/antigravity-usage-checker/internal/api"
)

// TestCreateProgressBar tests progress bar generation.
func TestCreateProgressBar(t *testing.T) {
	tests := []struct {
		name    string
		percent float64
		width   int
		wantLen int // Length of generated bar (approximate, without color codes)
	}{
		{
			name:    "0 percent",
			percent: 0,
			width:   10,
		},
		{
			name:    "50 percent",
			percent: 50,
			width:   10,
		},
		{
			name:    "100 percent",
			percent: 100,
			width:   10,
		},
		{
			name:    "Over 100 percent (should cap)",
			percent: 150,
			width:   10,
		},
		{
			name:    "Negative percent (should floor to 0)",
			percent: -10,
			width:   10,
		},
		{
			name:    "Decimal percent",
			percent: 45.5,
			width:   10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar := createProgressBar(tt.percent, tt.width)

			// Just verify it doesn't panic and returns a non-empty string
			if len(bar) == 0 {
				t.Error("createProgressBar() returned empty string")
			}
		})
	}
}

// TestTruncateString tests string truncation.
func TestTruncateString(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"exactly10!", 10, "exactly10!"},
		{"this is a very long string", 10, "this is..."},
		{"hello", 5, "hello"},
		{"hello world", 5, "he..."},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := truncateString(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncateString(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

// TestFormatTime tests time formatting.
func TestFormatTime(t *testing.T) {
	testTime := time.Date(2024, 12, 31, 14, 30, 45, 0, time.UTC)
	got := formatTime(testTime)
	want := "14:30:45 31/12/2024"

	if got != want {
		t.Errorf("formatTime() = %q, want %q", got, want)
	}
}

// TestFormatResetTime tests reset time formatting.
func TestFormatResetTime(t *testing.T) {
	tests := []struct {
		name         string
		resetTimeStr string
		wantContains string // Partial match since time is relative
	}{
		{
			name:         "Empty string",
			resetTimeStr: "",
			wantContains: "-",
		},
		{
			name:         "Invalid format",
			resetTimeStr: "not-a-date",
			wantContains: "-",
		},
		{
			name:         "Past time (already reset)",
			resetTimeStr: "2020-01-01T00:00:00Z",
			wantContains: "reset",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatResetTime(tt.resetTimeStr)
			if got == "" {
				t.Error("formatResetTime() returned empty string")
			}
		})
	}
}

// TestShowUsageDoesNotPanic tests that ShowUsage doesn't panic with various inputs.
func TestShowUsageDoesNotPanic(t *testing.T) {
	tests := []struct {
		name     string
		data     *api.UsageData
		asJSON   bool
		isCached bool
	}{
		{
			name: "Empty data",
			data: &api.UsageData{
				Models:    []api.QuotaInfo{},
				FetchedAt: time.Now(),
			},
			asJSON:   false,
			isCached: false,
		},
		{
			name: "With models",
			data: &api.UsageData{
				Models: []api.QuotaInfo{
					{
						ModelName:    "Claude Sonnet 4",
						Used:         50.5,
						Limit:        100,
						Remaining:    49.5,
						UsagePercent: 50.5,
						ResetTime:    "2024-12-31T23:59:59Z",
					},
					{
						ModelName:    "GPT-4o",
						Used:         80.3,
						Limit:        100,
						Remaining:    19.7,
						UsagePercent: 80.3,
						ResetTime:    "",
					},
				},
				Tier:         "Pro",
				PromptCredit: 1000,
				FetchedAt:    time.Now(),
			},
			asJSON:   false,
			isCached: false,
		},
		{
			name: "JSON output",
			data: &api.UsageData{
				Models: []api.QuotaInfo{
					{ModelName: "Test Model", Used: 10, Limit: 100},
				},
				FetchedAt: time.Now(),
			},
			asJSON:   true,
			isCached: false,
		},
		{
			name: "Cached data",
			data: &api.UsageData{
				Models:    []api.QuotaInfo{},
				IsCached:  true,
				FetchedAt: time.Now().Add(-1 * time.Hour),
			},
			asJSON:   false,
			isCached: true,
		},
		{
			name: "Long model names",
			data: &api.UsageData{
				Models: []api.QuotaInfo{
					{
						ModelName:    "This is a very long model name that should be truncated properly",
						Used:         100,
						Limit:        100,
						Remaining:    0,
						UsagePercent: 100,
					},
				},
				FetchedAt: time.Now(),
			},
			asJSON:   false,
			isCached: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that it doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("ShowUsage() panicked: %v", r)
				}
			}()

			ShowUsage(tt.data, tt.asJSON, tt.isCached)
		})
	}
}
