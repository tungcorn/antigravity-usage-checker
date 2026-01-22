package api

import (
	"encoding/json"
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
								// No quotaInfo, struct will have zero values (RemainingFraction=0 -> 100% used)
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
			// Marshal map to JSON bytes to simulate API response
			respBytes, err := json.Marshal(tt.response)
			if err != nil {
				t.Fatalf("Failed to marshal test response: %v", err)
			}

			usage, err := client.parseUserStatusResponse(respBytes)

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

	respBytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal test response: %v", err)
	}

	usage, err := client.parseUserStatusResponse(respBytes)
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
	if model.Remaining != 25.0 {
		t.Errorf("Remaining = %.1f, want %.1f", model.Remaining, 25.0)
	}

	if model.Used != 75.0 {
		t.Errorf("Used = %.1f, want %.1f", model.Used, 75.0)
	}

	if model.UsagePercent != 75.0 {
		t.Errorf("UsagePercent = %.1f, want %.1f", model.UsagePercent, 75.0)
	}

	if model.ResetTime != "2024-12-31T23:59:59Z" {
		t.Errorf("ResetTime = %q, want %q", model.ResetTime, "2024-12-31T23:59:59Z")
	}
}
