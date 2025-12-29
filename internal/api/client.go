// Package api provides the client for communicating with the Antigravity
// language server API.
package api

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	// API endpoints
	GetUserStatusPath      = "/exa.language_server_pb.LanguageServerService/GetUserStatus"
	GetCommandModelConfigs = "/exa.language_server_pb.LanguageServerService/GetCommandModelConfigs"
	
	// Request timeout
	RequestTimeout = 10 * time.Second
)

// Client handles API communication with the Antigravity language server.
type Client struct {
	connectPort int
	httpPort    int
	csrfToken   string
	httpClient  *http.Client
}

// QuotaInfo contains quota information for a specific model.
type QuotaInfo struct {
	ModelName    string `json:"model_name"`
	Used         int    `json:"used"`
	Limit        int    `json:"limit"`
	Remaining    int    `json:"remaining"`
	ResetTime    string `json:"reset_time"`
	UsagePercent int    `json:"usage_percent"`
}

// UsageData contains the complete usage information.
type UsageData struct {
	Models       []QuotaInfo `json:"models"`
	Tier         string      `json:"tier"`
	PromptCredit int         `json:"prompt_credit,omitempty"`
	FetchedAt    time.Time   `json:"fetched_at"`
	IsCached     bool        `json:"is_cached,omitempty"`
}

// NewClient creates a new API client with the given connection parameters.
func NewClient(connectPort int, csrfToken string, httpPort int) *Client {
	// Create HTTP client that skips TLS verification (local connection)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	
	return &Client{
		connectPort: connectPort,
		httpPort:    httpPort,
		csrfToken:   csrfToken,
		httpClient: &http.Client{
			Timeout:   RequestTimeout,
			Transport: transport,
		},
	}
}

// GetUserStatus fetches the current user status and quota information.
func (c *Client) GetUserStatus() (*UsageData, error) {
	// Prepare request body (empty for GetUserStatus)
	body := map[string]interface{}{}
	
	resp, err := c.makeRequest(GetUserStatusPath, body)
	if err != nil {
		return nil, fmt.Errorf("GetUserStatus failed: %w", err)
	}
	
	return c.parseUserStatusResponse(resp)
}

// makeRequest performs an API request with proper headers.
func (c *Client) makeRequest(path string, body interface{}) (map[string]interface{}, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	
	// Try HTTPS first, then fall back to HTTP
	url := fmt.Sprintf("https://127.0.0.1:%d%s", c.connectPort, path)
	resp, err := c.doRequest(url, jsonBody)
	if err != nil {
		// Fallback to HTTP
		url = fmt.Sprintf("http://127.0.0.1:%d%s", c.httpPort, path)
		resp, err = c.doRequest(url, jsonBody)
		if err != nil {
			return nil, err
		}
	}
	
	return resp, nil
}

// doRequest performs the actual HTTP request.
func (c *Client) doRequest(url string, body []byte) (map[string]interface{}, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connect-Protocol-Version", "1")
	req.Header.Set("X-Codeium-Csrf-Token", c.csrfToken)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return result, nil
}

// parseUserStatusResponse converts the API response to UsageData.
func (c *Client) parseUserStatusResponse(resp map[string]interface{}) (*UsageData, error) {
	usage := &UsageData{
		FetchedAt: time.Now(),
		Models:    []QuotaInfo{},
	}
	
	// Navigate to userStatus.cascadeModelConfigData.clientModelConfigs
	userStatus, ok := resp["userStatus"].(map[string]interface{})
	if !ok {
		return usage, nil
	}
	
	// Try to get plan name from userStatus
	if planName, ok := userStatus["planName"].(string); ok {
		usage.Tier = planName
	}
	
	// Get cascadeModelConfigData
	cascadeData, ok := userStatus["cascadeModelConfigData"].(map[string]interface{})
	if !ok {
		return usage, nil
	}
	
	// Get clientModelConfigs array
	modelConfigs, ok := cascadeData["clientModelConfigs"].([]interface{})
	if !ok {
		return usage, nil
	}
	
	// Parse each model config
	for _, mc := range modelConfigs {
		config, ok := mc.(map[string]interface{})
		if !ok {
			continue
		}
		
		info := QuotaInfo{
			ModelName: getStringValue(config, "label", "Unknown"),
		}
		
		// Parse quotaInfo
		if quotaInfo, ok := config["quotaInfo"].(map[string]interface{}); ok {
			// remainingFraction is between 0 and 1
			if remaining, ok := quotaInfo["remainingFraction"].(float64); ok {
				info.UsagePercent = 100 - int(remaining*100)
				info.Remaining = int(remaining * 100) // As percentage
				info.Limit = 100
				info.Used = 100 - info.Remaining
			}
			
			// Reset time
			if resetTime, ok := quotaInfo["resetTime"].(string); ok {
				info.ResetTime = resetTime
			}
		}
		
		usage.Models = append(usage.Models, info)
	}
	
	// Get prompt credits from userStatus if available
	if promptCredits, ok := userStatus["promptCreditsInfo"].(map[string]interface{}); ok {
		if remaining, ok := promptCredits["remainingCredits"].(float64); ok {
			usage.PromptCredit = int(remaining)
		}
	}
	
	return usage, nil
}

// getStringValue safely extracts a string value from a map.
func getStringValue(m map[string]interface{}, key, defaultVal string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return defaultVal
}
