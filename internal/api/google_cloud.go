// Package api provides clients for fetching quota data from various sources.
package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// Google Cloud Code API endpoints
	GoogleCloudCodeBaseURL   = "https://cloudcode-pa.googleapis.com"
	FetchAvailableModelsPath = "/v1internal:fetchAvailableModels"
	LoadCodeAssistPath       = "/v1internal:loadCodeAssist"

	// User agent to mimic official client
	AntigravityUserAgent = "antigravity/1.11.3"
)

// GoogleCloudClient handles communication with Google Cloud Code API.
type GoogleCloudClient struct {
	accessToken string
	httpClient  *http.Client
	projectID   string
}

// GoogleCloudQuotaInfo represents quota info from Google Cloud API.
type GoogleCloudQuotaInfo struct {
	RemainingFraction float64 `json:"remainingFraction"`
	ResetTime         string  `json:"resetTime"`
}

// GoogleCloudModelInfo represents a model from Google Cloud API.
type GoogleCloudModelInfo struct {
	Model     string               `json:"model"`
	QuotaInfo GoogleCloudQuotaInfo `json:"quotaInfo"`
}

// fetchAvailableModelsResponse represents the API response structure.
type fetchAvailableModelsResponse struct {
	Models []struct {
		ModelID   string `json:"modelId"`
		Label     string `json:"label"`
		QuotaInfo struct {
			RemainingFraction float64 `json:"remainingFraction"`
			ResetTime         string  `json:"resetTime"`
		} `json:"quotaInfo"`
	} `json:"models"`
	// Alternative structure where models is a map
	ModelsMap map[string]struct {
		QuotaInfo struct {
			RemainingFraction float64 `json:"remainingFraction"`
			ResetTime         string  `json:"resetTime"`
		} `json:"quotaInfo"`
	} `json:"modelsMap"`
}

// loadCodeAssistResponse represents the loadCodeAssist API response.
type loadCodeAssistResponse struct {
	CloudAICompanionProject string `json:"cloudaicompanionProject"`
	Project                 string `json:"project"`
}

// NewGoogleCloudClient creates a new Google Cloud API client.
func NewGoogleCloudClient(accessToken string) *GoogleCloudClient {
	return &GoogleCloudClient{
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: RequestTimeout,
		},
	}
}

// GetProjectID fetches the project ID from loadCodeAssist endpoint.
func (c *GoogleCloudClient) GetProjectID() (string, error) {
	if c.projectID != "" {
		return c.projectID, nil
	}

	url := GoogleCloudCodeBaseURL + LoadCodeAssistPath
	body := map[string]interface{}{}

	respBytes, err := c.doRequest(url, body)
	if err != nil {
		return "", fmt.Errorf("loadCodeAssist failed: %w", err)
	}

	var resp loadCodeAssistResponse
	if err := json.Unmarshal(respBytes, &resp); err != nil {
		return "", fmt.Errorf("failed to parse loadCodeAssist response: %w", err)
	}

	// Try cloudaicompanionProject first, then project
	if resp.CloudAICompanionProject != "" {
		c.projectID = resp.CloudAICompanionProject
	} else if resp.Project != "" {
		c.projectID = resp.Project
	}

	return c.projectID, nil
}

// GetUsageData fetches quota data from Google Cloud API with exact percentages.
func (c *GoogleCloudClient) GetUsageData() (*UsageData, error) {
	// First, try to get project ID (optional, may not be required)
	projectID, _ := c.GetProjectID()

	url := GoogleCloudCodeBaseURL + FetchAvailableModelsPath
	body := map[string]interface{}{}
	if projectID != "" {
		body["project"] = projectID
	}

	respBytes, err := c.doRequest(url, body)
	if err != nil {
		return nil, fmt.Errorf("fetchAvailableModels failed: %w", err)
	}

	return c.parseModelsResponse(respBytes)
}

// doRequest performs an authenticated request to Google Cloud API.
func (c *GoogleCloudClient) doRequest(url string, body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("User-Agent", AntigravityUserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBytes))
	}

	return respBytes, nil
}

// parseModelsResponse converts the API response to UsageData.
func (c *GoogleCloudClient) parseModelsResponse(respBytes []byte) (*UsageData, error) {
	usage := &UsageData{
		FetchedAt: time.Now(),
		Models:    []QuotaInfo{},
	}

	// Try to parse as generic JSON first to handle different response formats
	var rawResp map[string]interface{}
	if err := json.Unmarshal(respBytes, &rawResp); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Try to extract models from various possible response structures
	if models, ok := rawResp["models"].([]interface{}); ok {
		// Array format
		for _, m := range models {
			if modelMap, ok := m.(map[string]interface{}); ok {
				info := c.parseModelMap(modelMap)
				if info.ModelName != "" {
					usage.Models = append(usage.Models, info)
				}
			}
		}
	} else if modelsMap, ok := rawResp["models"].(map[string]interface{}); ok {
		// Map format (keyed by model name)
		for modelName, m := range modelsMap {
			if modelMap, ok := m.(map[string]interface{}); ok {
				info := c.parseModelMap(modelMap)
				if info.ModelName == "" {
					info.ModelName = modelName
				}
				usage.Models = append(usage.Models, info)
			}
		}
	}

	// Try to extract tier/plan info
	if planName, ok := rawResp["planName"].(string); ok {
		usage.Tier = planName
	}

	return usage, nil
}

// parseModelMap extracts QuotaInfo from a model map.
func (c *GoogleCloudClient) parseModelMap(modelMap map[string]interface{}) QuotaInfo {
	info := QuotaInfo{}

	// Get model name from various possible fields
	if label, ok := modelMap["label"].(string); ok && label != "" {
		info.ModelName = label
	} else if modelID, ok := modelMap["modelId"].(string); ok && modelID != "" {
		info.ModelName = modelID
	} else if model, ok := modelMap["model"].(string); ok && model != "" {
		info.ModelName = model
	}

	// Get quota info
	if quotaInfo, ok := modelMap["quotaInfo"].(map[string]interface{}); ok {
		if remainingFraction, ok := quotaInfo["remainingFraction"].(float64); ok {
			// Calculate exact percentages from fraction
			info.Remaining = remainingFraction * 100
			info.Used = 100 - info.Remaining
			info.UsagePercent = info.Used
			info.Limit = 100
		}
		if resetTime, ok := quotaInfo["resetTime"].(string); ok {
			info.ResetTime = resetTime
		}
	}

	return info
}
