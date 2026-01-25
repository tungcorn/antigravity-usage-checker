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
	ModelName    string  `json:"model_name"`
	Used         float64 `json:"used"`
	Limit        float64 `json:"limit"`
	Remaining    float64 `json:"remaining"`
	ResetTime    string  `json:"reset_time"`
	UsagePercent float64 `json:"usage_percent"`
}

// UsageData contains the complete usage information.
type UsageData struct {
	Models       []QuotaInfo `json:"models"`
	Tier         string      `json:"tier"`
	Email        string      `json:"email,omitempty"`
	Name         string      `json:"name,omitempty"`
	PromptCredit int         `json:"prompt_credit,omitempty"`
	FetchedAt    time.Time   `json:"fetched_at"`
	IsCached     bool        `json:"is_cached,omitempty"`
}

// Internal structs for JSON parsing
type apiResponse struct {
	UserStatus userStatusRaw `json:"userStatus"`
}

type userStatusRaw struct {
	PlanName               string                 `json:"planName"`
	Email                  string                 `json:"email"`
	Name                   string                 `json:"name"`
	CascadeModelConfigData cascadeModelConfigData `json:"cascadeModelConfigData"`
	PromptCreditsInfo      promptCreditsInfo      `json:"promptCreditsInfo"`
}

type cascadeModelConfigData struct {
	ClientModelConfigs []clientModelConfig `json:"clientModelConfigs"`
}

type clientModelConfig struct {
	Label        string       `json:"label"`
	QuotaInfo    quotaInfoRaw `json:"quotaInfo"`
	ModelOrAlias modelOrAlias `json:"modelOrAlias"`
}

type modelOrAlias struct {
	Model string `json:"model"`
}

type quotaInfoRaw struct {
	RemainingFraction float64 `json:"remainingFraction"`
	ResetTime         string  `json:"resetTime"`
}

type promptCreditsInfo struct {
	RemainingCredits float64 `json:"remainingCredits"`
}

// NewClient creates a new API client with the given connection parameters.
func NewClient(connectPort int, csrfToken string, httpPort int) *Client {
	// SECURITY NOTE: TLS verification is intentionally disabled because:
	// 1. This ONLY connects to localhost (127.0.0.1) - never to external servers
	// 2. The Antigravity language server uses a self-signed certificate
	// 3. No sensitive data is transmitted over the network
	// This is safe for local-only connections and required for the tool to work.
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // #nosec G402
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

	respBytes, err := c.makeRequest(GetUserStatusPath, body)
	if err != nil {
		return nil, fmt.Errorf("GetUserStatus failed: %w", err)
	}

	return c.parseUserStatusResponse(respBytes)
}

// makeRequest performs an API request and returns raw bytes.
func (c *Client) makeRequest(path string, body interface{}) ([]byte, error) {
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
func (c *Client) doRequest(url string, body []byte) ([]byte, error) {
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

	return io.ReadAll(resp.Body)
}

// parseUserStatusResponse converts the API response to UsageData.
func (c *Client) parseUserStatusResponse(respBytes []byte) (*UsageData, error) {
	usage := &UsageData{
		FetchedAt: time.Now(),
		Models:    []QuotaInfo{},
	}

	var apiResp apiResponse
	if err := json.Unmarshal(respBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Extract plan name
	if apiResp.UserStatus.PlanName != "" {
		usage.Tier = apiResp.UserStatus.PlanName
	}

	// Extract user info
	usage.Email = apiResp.UserStatus.Email
	usage.Name = apiResp.UserStatus.Name

	// Extract prompt credits
	if apiResp.UserStatus.PromptCreditsInfo.RemainingCredits > 0 {
		usage.PromptCredit = int(apiResp.UserStatus.PromptCreditsInfo.RemainingCredits)
	}

	// Process models
	for _, config := range apiResp.UserStatus.CascadeModelConfigData.ClientModelConfigs {
		info := QuotaInfo{
			ModelName: config.Label,
		}

		// If label is missing, use model name as fallback
		if info.ModelName == "" {
			if config.ModelOrAlias.Model != "" {
				info.ModelName = config.ModelOrAlias.Model
			} else {
				info.ModelName = "Unknown Model"
			}
		}

		// Calculate remaining percentage from fraction (0.0 to 1.0)
		remaining := config.QuotaInfo.RemainingFraction
		info.Remaining = remaining * 100
		info.UsagePercent = 100 - info.Remaining
		info.Limit = 100
		info.Used = 100 - info.Remaining

		info.ResetTime = config.QuotaInfo.ResetTime

		usage.Models = append(usage.Models, info)
	}

	return usage, nil
}
