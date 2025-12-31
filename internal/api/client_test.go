package api

import (
	"testing"
)

// TestParseUserStatusResponse tests API response parsing.
func TestParseUserStatusResponse(t *testing.T) {
	client := NewClient(42100, "test-token", 8080)

	tests := []struct {
		name       string
		response   map[string]interface{}
		wantTier   string
		wantModels int
		wantErr    bool
	}{
		{
			name:       "Empty response",
			response:   map[string]interface{}{},
			wantTier:   "",
			wantModels: 0,
			wantErr:    false,
		},
		{
			name: "Response with userStatus only",
			response: map[string]interface{}{
				"userStatus": map[string]interface{}{
					"planName": "Pro",
				},
			},
			wantTier:   "Pro",
			wantModels: 0,
			wantErr:    false,
		},
		{
			name: "Full response with models",
			response: map[string]interface{}{
				"userStatus": map[string]interface{}{
					"planName": "Enterprise",
					"cascadeModelConfigData": map[string]interface{}{
						"clientModelConfigs": []interface{}{
							map[string]interface{}{
								"label": "Claude Sonnet 4",
								"quotaInfo": map[string]interface{}{
									"remainingFraction": 0.75,
									"resetTime":         "2024-12-31T23:59:59Z",
								},
							},
							map[string]interface{}{
								"label": "GPT-4o",
								"quotaInfo": map[string]interface{}{
									"remainingFraction": 0.50,
									"resetTime":         "2024-12-31T12:00:00Z",
								},
							},
						},
					},
					"promptCreditsInfo": map[string]interface{}{
						"remainingCredits": float64(500),
					},
				},
			},
			wantTier:   "Enterprise",
			wantModels: 2,
			wantErr:    false,
		},
		{
			name: "Response with missing quotaInfo",
			response: map[string]interface{}{
				"userStatus": map[string]interface{}{
					"planName": "Free",
					"cascadeModelConfigData": map[string]interface{}{
						"clientModelConfigs": []interface{}{
							map[string]interface{}{
								"label": "Basic Model",
								// No quotaInfo
							},
						},
					},
				},
			},
			wantTier:   "Free",
			wantModels: 1,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usage, err := client.parseUserStatusResponse(tt.response)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseUserStatusResponse() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("parseUserStatusResponse() unexpected error: %v", err)
				return
			}

			if usage.Tier != tt.wantTier {
				t.Errorf("Tier = %q, want %q", usage.Tier, tt.wantTier)
			}

			if len(usage.Models) != tt.wantModels {
				t.Errorf("got %d models, want %d", len(usage.Models), tt.wantModels)
			}
		})
	}
}

// TestParseUserStatusResponseModelDetails tests detailed model parsing.
func TestParseUserStatusResponseModelDetails(t *testing.T) {
	client := NewClient(42100, "test-token", 8080)

	response := map[string]interface{}{
		"userStatus": map[string]interface{}{
			"planName": "Pro",
			"cascadeModelConfigData": map[string]interface{}{
				"clientModelConfigs": []interface{}{
					map[string]interface{}{
						"label": "Test Model",
						"quotaInfo": map[string]interface{}{
							"remainingFraction": 0.25, // 25% remaining = 75% used
							"resetTime":         "2024-12-31T23:59:59Z",
						},
					},
				},
			},
		},
	}

	usage, err := client.parseUserStatusResponse(response)
	if err != nil {
		t.Fatalf("parseUserStatusResponse() error: %v", err)
	}

	if len(usage.Models) != 1 {
		t.Fatalf("expected 1 model, got %d", len(usage.Models))
	}

	model := usage.Models[0]

	if model.ModelName != "Test Model" {
		t.Errorf("ModelName = %q, want %q", model.ModelName, "Test Model")
	}

	// remainingFraction = 0.25 means 25% remaining
	// So: Remaining = 25, Used = 75, Limit = 100, UsagePercent = 75
	if model.Remaining != 25 {
		t.Errorf("Remaining = %d, want %d", model.Remaining, 25)
	}

	if model.Used != 75 {
		t.Errorf("Used = %d, want %d", model.Used, 75)
	}

	if model.UsagePercent != 75 {
		t.Errorf("UsagePercent = %d, want %d", model.UsagePercent, 75)
	}

	if model.ResetTime != "2024-12-31T23:59:59Z" {
		t.Errorf("ResetTime = %q, want %q", model.ResetTime, "2024-12-31T23:59:59Z")
	}
}

// TestGetStringValue tests safe string extraction.
func TestGetStringValue(t *testing.T) {
	tests := []struct {
		name       string
		m          map[string]interface{}
		key        string
		defaultVal string
		want       string
	}{
		{
			name:       "Key exists with string value",
			m:          map[string]interface{}{"name": "Claude"},
			key:        "name",
			defaultVal: "Unknown",
			want:       "Claude",
		},
		{
			name:       "Key does not exist",
			m:          map[string]interface{}{"other": "value"},
			key:        "name",
			defaultVal: "Unknown",
			want:       "Unknown",
		},
		{
			name:       "Key exists but not string",
			m:          map[string]interface{}{"count": 42},
			key:        "count",
			defaultVal: "N/A",
			want:       "N/A",
		},
		{
			name:       "Empty map",
			m:          map[string]interface{}{},
			key:        "name",
			defaultVal: "Default",
			want:       "Default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getStringValue(tt.m, tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("getStringValue() = %q, want %q", got, tt.want)
			}
		})
	}
}
