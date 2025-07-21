package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	godestats "github.com/Yeti47/gode-stats/pkg"
)

const (
	// DefaultBaseURL is the default base URL for the Code::Stats API.
	DefaultBaseURL = "https://codestats.net"
	// APIPrefix is the prefix for all API endpoints.
	APIPrefix = "/api"
	// AuthHeader is the HTTP header used for API token authentication.
	AuthHeader = "X-API-Token"
	// UserAgent is the User-Agent header sent with requests.
	UserAgent = "gode-stats/1.0.0"
)

// Client implements the CodeStatsClient interface for interacting with the Code::Stats API.
type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

// New creates a new Code::Stats API client with the provided API token.
func New(apiToken string) godestats.CodeStatsClient {
	return NewWithBaseURL(apiToken, DefaultBaseURL)
}

// NewAnonymous creates a new anonymous Code::Stats API client for read-only operations.
// This client can only retrieve public user profiles and cannot send pulses.
func NewAnonymous() godestats.CodeStatsClient {
	return NewWithBaseURL("", DefaultBaseURL)
}

// NewWithBaseURL creates a new Code::Stats API client with a custom base URL.
// This is useful for testing against custom instances or local development servers.
func NewWithBaseURL(apiToken, baseURL string) godestats.CodeStatsClient {
	return &Client{
		baseURL:  baseURL,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetUserProfile retrieves the public profile information for the specified user.
func (c *Client) GetUserProfile(ctx context.Context, username string) (*godestats.UserProfile, error) {
	if username == "" {
		return nil, godestats.ErrEmptyUsername
	}

	// Construct the API URL
	endpoint := fmt.Sprintf("%s%s/users/%s", c.baseURL, APIPrefix, url.PathEscape(username))

	// Create the request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "application/json")

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, godestats.NewNetworkError("GET request", endpoint, err)
	}
	defer resp.Body.Close()

	// Handle HTTP errors
	if resp.StatusCode == http.StatusNotFound {
		return nil, godestats.ErrUserNotFound
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, godestats.ErrUnauthorized
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, godestats.ErrRateLimited
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		// Try to parse error message from JSON
		var errorResp struct {
			Error string `json:"error"`
		}
		message := string(body)
		if json.Unmarshal(body, &errorResp) == nil && errorResp.Error != "" {
			message = errorResp.Error
		}

		return nil, godestats.NewAPIError(resp.StatusCode, message, endpoint)
	}

	// Parse the response
	var profile godestats.UserProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("%w: %v", godestats.ErrInvalidResponse, err)
	}

	return &profile, nil
}

// SendPulse submits a pulse (collection of XPs for different languages) to the API.
func (c *Client) SendPulse(ctx context.Context, pulse godestats.Pulse) error {
	if c.apiToken == "" {
		return godestats.ErrUnauthorized
	}

	// Validate pulse timestamp (must not be older than a week)
	weekAgo := time.Now().AddDate(0, 0, -7)
	if pulse.CodedAt.Before(weekAgo) {
		return godestats.ErrPulseTimestampTooOld
	}

	// Construct the API URL
	endpoint := fmt.Sprintf("%s%s/my/pulses", c.baseURL, APIPrefix)

	// Serialize the pulse to JSON
	pulseData, err := json.Marshal(pulse)
	if err != nil {
		return fmt.Errorf("failed to serialize pulse: %w", err)
	}

	// Create the request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(pulseData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set(AuthHeader, c.apiToken)

	// Execute the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return godestats.NewNetworkError("POST request", endpoint, err)
	}
	defer resp.Body.Close()

	// Handle HTTP errors
	if resp.StatusCode == http.StatusCreated {
		return nil // Success
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return godestats.ErrUnauthorized
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return godestats.ErrRateLimited
	}

	// Read error response
	body, _ := io.ReadAll(resp.Body)

	// Try to parse error message from JSON
	var errorResp struct {
		Error string `json:"error"`
	}
	message := string(body)
	if json.Unmarshal(body, &errorResp) == nil && errorResp.Error != "" {
		message = errorResp.Error
	}

	return godestats.NewAPIError(resp.StatusCode, message, endpoint)
}
