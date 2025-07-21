package godestats

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Common error variables that consumers can check against
var (
	// ErrUserNotFound is returned when a user profile is not found or is private
	ErrUserNotFound = errors.New("user not found or profile is private")

	// ErrUnauthorized is returned when API token is missing or invalid
	ErrUnauthorized = errors.New("unauthorized: API token is missing or invalid")

	// ErrPulseTimestampTooOld is returned when a pulse timestamp is older than a week
	ErrPulseTimestampTooOld = errors.New("pulse timestamp is older than a week and will be rejected")

	// ErrEmptyUsername is returned when an empty username is provided
	ErrEmptyUsername = errors.New("username cannot be empty")

	// ErrNetworkError is returned when there are network connectivity issues
	ErrNetworkError = errors.New("network error")

	// ErrInvalidResponse is returned when the API response cannot be parsed
	ErrInvalidResponse = errors.New("invalid response from API")

	// ErrRateLimited is returned when the API rate limit is exceeded
	ErrRateLimited = errors.New("API rate limit exceeded")
)

// APIError represents an error response from the Code::Stats API
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Endpoint   string `json:"endpoint,omitempty"`
}

// Error implements the error interface for APIError
func (e *APIError) Error() string {
	if e.Endpoint != "" {
		return fmt.Sprintf("API error %d at %s: %s", e.StatusCode, e.Endpoint, e.Message)
	}
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// IsTemporary returns true if the error might be resolved by retrying
func (e *APIError) IsTemporary() bool {
	return e.StatusCode >= 500 || e.StatusCode == http.StatusTooManyRequests
}

// NetworkError wraps network-related errors with additional context
type NetworkError struct {
	Operation string `json:"operation"`
	URL       string `json:"url,omitempty"`
	Err       error  `json:"error"`
}

// Error implements the error interface for NetworkError
func (e *NetworkError) Error() string {
	if e.URL != "" {
		return fmt.Sprintf("network error during %s to %s: %v", e.Operation, e.URL, e.Err)
	}
	return fmt.Sprintf("network error during %s: %v", e.Operation, e.Err)
}

// Unwrap returns the underlying error for error unwrapping
func (e *NetworkError) Unwrap() error {
	return e.Err
}

// IsTemporary returns true if the network error might be resolved by retrying
func (e *NetworkError) IsTemporary() bool {
	if e.Err == nil {
		return false
	}

	// Check error message for common temporary conditions
	errMsg := e.Err.Error()
	return strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "connection refused") ||
		strings.Contains(errMsg, "no such host") ||
		strings.Contains(errMsg, "network is unreachable") ||
		strings.Contains(errMsg, "connection reset")
}

// Helper functions for creating specific errors

// NewAPIError creates a new APIError with the given status code and message
func NewAPIError(statusCode int, message, endpoint string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Endpoint:   endpoint,
	}
}

// NewNetworkError creates a new NetworkError with context
func NewNetworkError(operation, url string, err error) *NetworkError {
	return &NetworkError{
		Operation: operation,
		URL:       url,
		Err:       err,
	}
}

// Error classification helpers

// IsUserNotFound checks if an error indicates a user was not found
func IsUserNotFound(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrUserNotFound) {
		return true
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}

	return false
}

// IsUnauthorized checks if an error indicates unauthorized access
func IsUnauthorized(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrUnauthorized) {
		return true
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusUnauthorized
	}

	return false
}

// IsTemporary checks if an error is temporary and might be resolved by retrying
func IsTemporary(err error) bool {
	if err == nil {
		return false
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.IsTemporary()
	}

	var netErr *NetworkError
	if errors.As(err, &netErr) {
		return netErr.IsTemporary()
	}

	return false
}

// IsRateLimited checks if an error indicates rate limiting
func IsRateLimited(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrRateLimited) {
		return true
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}

	return false
}

// IsNetworkError checks if an error is a network-related error
func IsNetworkError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, ErrNetworkError) {
		return true
	}

	var netErr *NetworkError
	return errors.As(err, &netErr)
}
