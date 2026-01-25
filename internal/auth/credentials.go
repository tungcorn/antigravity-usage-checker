// Package auth handles OAuth credential loading from Antigravity config files.
package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Credentials holds OAuth token information from Antigravity.
type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiryDate   int64  `json:"expiry_date"` // Unix timestamp in milliseconds
}

// IsExpired checks if the access token has expired.
func (c *Credentials) IsExpired() bool {
	now := time.Now().UnixMilli()
	return now > c.ExpiryDate
}

// ExpiresInMinutes returns the number of minutes until token expiration.
func (c *Credentials) ExpiresInMinutes() int {
	remaining := c.ExpiryDate - time.Now().UnixMilli()
	return int(remaining / 60000)
}

// LoadCredentials loads OAuth credentials from the Antigravity config directory.
// Returns nil if credentials are not found or cannot be loaded.
func LoadCredentials() (*Credentials, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}
	credsFile := filepath.Join(configPath, "oauth_creds.json")

	data, err := os.ReadFile(credsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("failed to parse credentials: %w", err)
	}

	return &creds, nil
}

// getConfigPath returns the path to the .gemini config directory.
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".gemini"), nil
}
