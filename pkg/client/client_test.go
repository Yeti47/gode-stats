package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	godestats "github.com/Yeti47/gode-stats/pkg"
)

func TestClient_GetUserProfile_Success(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedPath := "/api/users/testuser"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Send mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"user": "testuser",
			"total_xp": 1000,
			"new_xp": 50,
			"machines": {
				"laptop": {
					"xps": 800,
					"new_xps": 30
				}
			},
			"languages": {
				"Go": {
					"xps": 600,
					"new_xps": 25
				}
			},
			"dates": {
				"2023-01-01": 50
			}
		}`))
	}))
	defer server.Close()

	// Create client with test server URL
	client := NewWithBaseURL("test-token", server.URL)

	// Test the method
	profile, err := client.GetUserProfile(context.Background(), "testuser")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify response
	if profile.User != "testuser" {
		t.Errorf("Expected user 'testuser', got '%s'", profile.User)
	}
	if profile.TotalXP != 1000 {
		t.Errorf("Expected total XP 1000, got %d", profile.TotalXP)
	}
	if profile.NewXP != 50 {
		t.Errorf("Expected new XP 50, got %d", profile.NewXP)
	}
}

func TestClient_GetUserProfile_NotFound(t *testing.T) {
	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewWithBaseURL("test-token", server.URL)

	_, err := client.GetUserProfile(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("Expected error for non-existent user")
	}

	if !godestats.IsUserNotFound(err) {
		t.Errorf("Expected user not found error, got: %v", err)
	}
}

func TestClient_GetUserProfile_EmptyUsername(t *testing.T) {
	client := New("test-token")

	_, err := client.GetUserProfile(context.Background(), "")
	if err == nil {
		t.Fatal("Expected error for empty username")
	}

	if !errors.Is(err, godestats.ErrEmptyUsername) {
		t.Errorf("Expected ErrEmptyUsername, got: %v", err)
	}
}

func TestClient_SendPulse_Success(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		expectedPath := "/api/my/pulses"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify authentication header
		token := r.Header.Get("X-API-Token")
		if token != "test-token" {
			t.Errorf("Expected token 'test-token', got '%s'", token)
		}

		// Send success response
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"ok": "Great success!"}`))
	}))
	defer server.Close()

	client := NewWithBaseURL("test-token", server.URL)

	// Create a test pulse
	pulse := godestats.Pulse{
		CodedAt: time.Now(),
		XPs: []godestats.LanguageXP{
			{Language: "Go", XP: 15},
			{Language: "JavaScript", XP: 30},
		},
	}

	err := client.SendPulse(context.Background(), pulse)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestClient_SendPulse_NoToken(t *testing.T) {
	client := New("")

	pulse := godestats.Pulse{
		CodedAt: time.Now(),
		XPs: []godestats.LanguageXP{
			{Language: "Go", XP: 15},
		},
	}

	err := client.SendPulse(context.Background(), pulse)
	if err == nil {
		t.Fatal("Expected error for missing API token")
	}

	if !errors.Is(err, godestats.ErrUnauthorized) {
		t.Errorf("Expected ErrUnauthorized, got: %v", err)
	}
}

func TestClient_SendPulse_OldTimestamp(t *testing.T) {
	client := New("test-token")

	// Create a pulse with a timestamp older than a week
	oldTime := time.Now().AddDate(0, 0, -8)
	pulse := godestats.Pulse{
		CodedAt: oldTime,
		XPs: []godestats.LanguageXP{
			{Language: "Go", XP: 15},
		},
	}

	err := client.SendPulse(context.Background(), pulse)
	if err == nil {
		t.Fatal("Expected error for old timestamp")
	}

	if !errors.Is(err, godestats.ErrPulseTimestampTooOld) {
		t.Errorf("Expected ErrPulseTimestampTooOld, got: %v", err)
	}
}

// TestConstructors tests the various constructor functions
func TestConstructors(t *testing.T) {
	apiToken := "test-token"

	// Test New
	client := New(apiToken)
	if client == nil {
		t.Fatal("New() returned nil")
	}

	// Test NewAnonymous
	anonClient := NewAnonymous()
	if anonClient == nil {
		t.Fatal("NewAnonymous() returned nil")
	}
}

func TestAnonymousClient_SendPulse_ShouldFail(t *testing.T) {
	client := NewAnonymous()

	pulse := godestats.Pulse{
		CodedAt: time.Now(),
		XPs: []godestats.LanguageXP{
			{Language: "Go", XP: 15},
		},
	}

	err := client.SendPulse(context.Background(), pulse)
	if err == nil {
		t.Fatal("Expected error when sending pulse with anonymous client")
	}

	if !errors.Is(err, godestats.ErrUnauthorized) {
		t.Errorf("Expected ErrUnauthorized, got: %v", err)
	}
}
