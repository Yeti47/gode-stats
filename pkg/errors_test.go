package godestats

import (
	"errors"
	"fmt"
	"testing"
)

func TestAPIError(t *testing.T) {
	err := NewAPIError(404, "User not found", "/api/users/test")

	expected := "API error 404 at /api/users/test: User not found"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}

	// Test without endpoint
	err2 := NewAPIError(500, "Internal server error", "")
	expected2 := "API error 500: Internal server error"
	if err2.Error() != expected2 {
		t.Errorf("Expected '%s', got '%s'", expected2, err2.Error())
	}
}

func TestAPIError_IsTemporary(t *testing.T) {
	tests := []struct {
		statusCode int
		expected   bool
	}{
		{400, false},
		{401, false},
		{404, false},
		{429, true},
		{500, true},
		{502, true},
		{503, true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("status_%d", tt.statusCode), func(t *testing.T) {
			err := NewAPIError(tt.statusCode, "test error", "")
			if err.IsTemporary() != tt.expected {
				t.Errorf("Expected IsTemporary() = %v for status %d", tt.expected, tt.statusCode)
			}
		})
	}
}

func TestNetworkError(t *testing.T) {
	originalErr := errors.New("connection refused")
	err := NewNetworkError("GET", "https://api.example.com", originalErr)

	expected := "network error during GET to https://api.example.com: connection refused"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}

	// Test unwrapping
	if !errors.Is(err, originalErr) {
		t.Error("Expected NetworkError to wrap original error")
	}

	// Test without URL
	err2 := NewNetworkError("POST", "", originalErr)
	expected2 := "network error during POST: connection refused"
	if err2.Error() != expected2 {
		t.Errorf("Expected '%s', got '%s'", expected2, err2.Error())
	}
}

func TestIsUserNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"ErrUserNotFound", ErrUserNotFound, true},
		{"404 API error", NewAPIError(404, "Not found", ""), true},
		{"500 API error", NewAPIError(500, "Server error", ""), false},
		{"other error", errors.New("random error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUserNotFound(tt.err)
			if result != tt.expected {
				t.Errorf("Expected IsUserNotFound() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"ErrUnauthorized", ErrUnauthorized, true},
		{"401 API error", NewAPIError(401, "Unauthorized", ""), true},
		{"404 API error", NewAPIError(404, "Not found", ""), false},
		{"other error", errors.New("random error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUnauthorized(tt.err)
			if result != tt.expected {
				t.Errorf("Expected IsUnauthorized() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsTemporary(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"temporary API error", NewAPIError(500, "Server error", ""), true},
		{"non-temporary API error", NewAPIError(400, "Bad request", ""), false},
		{"temporary network error", NewNetworkError("GET", "", errors.New("timeout")), true},
		{"other error", errors.New("random error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTemporary(tt.err)
			if result != tt.expected {
				t.Errorf("Expected IsTemporary() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsRateLimited(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"ErrRateLimited", ErrRateLimited, true},
		{"429 API error", NewAPIError(429, "Too many requests", ""), true},
		{"500 API error", NewAPIError(500, "Server error", ""), false},
		{"other error", errors.New("random error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRateLimited(tt.err)
			if result != tt.expected {
				t.Errorf("Expected IsRateLimited() = %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestIsNetworkError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"ErrNetworkError", ErrNetworkError, true},
		{"network error", NewNetworkError("GET", "", errors.New("timeout")), true},
		{"API error", NewAPIError(500, "Server error", ""), false},
		{"other error", errors.New("random error"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNetworkError(tt.err)
			if result != tt.expected {
				t.Errorf("Expected IsNetworkError() = %v, got %v", tt.expected, result)
			}
		})
	}
}
