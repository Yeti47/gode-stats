package client

import (
	"context"
	"testing"

	godestats "github.com/Yeti47/gode-stats/pkg"
)

// TestAnonymousIntegration tests anonymous read-only access against the real API
// This test is safe as it only performs read operations
func TestAnonymousIntegration(t *testing.T) {
	client := NewAnonymous()

	// Test getting a public profile (using a well-known public user)
	// This is safe as it's read-only and doesn't modify any data
	profile, err := client.GetUserProfile(context.Background(), "Nicd")
	if err != nil {
		// Check if it's a network error
		if godestats.IsNetworkError(err) {
			t.Skipf("Network error connecting to API (expected in some environments): %v", err)
		}
		// Check if user not found (maybe the user changed or profile became private)
		if godestats.IsUserNotFound(err) {
			t.Skipf("Test user not found or profile is private: %v", err)
		}
		t.Fatalf("Unexpected error retrieving public profile: %v", err)
	}

	if profile.User != "Nicd" {
		t.Errorf("Expected user 'Nicd', got '%s'", profile.User)
	}

	if profile.TotalXP <= 0 {
		t.Errorf("Expected positive total XP, got %d", profile.TotalXP)
	}

	t.Logf("Successfully retrieved public profile for %s with %d total XP", profile.User, profile.TotalXP)
}

// TestAuthenticatedClientCreation tests that authenticated clients can be created
// This test doesn't perform any API calls, so it's safe
func TestAuthenticatedClientCreation(t *testing.T) {
	// Test creating authenticated client
	client := New("fake-token")
	if client == nil {
		t.Error("Expected non-nil client")
	}

	// Test creating anonymous client
	anonClient := NewAnonymous()
	if anonClient == nil {
		t.Error("Expected non-nil anonymous client")
	}

	// Test creating custom base URL client
	customClient := NewWithBaseURL("fake-token", "https://example.com")
	if customClient == nil {
		t.Error("Expected non-nil custom client")
	}

	t.Log("All client constructors work correctly")
}
