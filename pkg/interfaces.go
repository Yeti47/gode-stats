package godestats

import (
	"context"
	"time"
)

// CodeStatsClient defines the interface for interacting with the Code::Stats API.
type CodeStatsClient interface {
	// GetUserProfile retrieves the public profile information for the specified user.
	// Returns an error if the user does not exist or their profile is private.
	GetUserProfile(ctx context.Context, username string) (*UserProfile, error)

	// SendPulse submits a pulse (collection of XPs for different languages) to the API.
	// The pulse must contain a coded_at timestamp and should be no older than a week.
	SendPulse(ctx context.Context, pulse Pulse) error
}

// XpCalculator defines the interface for calculating levels and percentages from XP.
type XpCalculator interface {
	// GetLevel calculates the level for the given XP amount.
	GetLevel(xp int) int

	// GetLevelPercentage calculates the percentage progress within the current level.
	// Returns a value between 0.0 and 1.0.
	GetLevelPercentage(xp int) float64

	// GetXpForLevel calculates the minimum XP required to reach the specified level.
	GetXpForLevel(level int) int

	// GetXpForNextLevel calculates the minimum XP required to reach the next level
	// from the current XP amount.
	GetXpForNextLevel(xp int) int
}

// UserProfile represents the public profile information of a user.
type UserProfile struct {
	User      string                  `json:"user"`
	TotalXP   int                     `json:"total_xp"`
	NewXP     int                     `json:"new_xp"`
	Machines  map[string]MachineInfo  `json:"machines"`
	Languages map[string]LanguageInfo `json:"languages"`
	Dates     map[string]int          `json:"dates"`
}

// MachineInfo represents XP information for a specific machine.
type MachineInfo struct {
	XPs    int `json:"xps"`
	NewXPs int `json:"new_xps"`
}

// LanguageInfo represents XP information for a specific language.
type LanguageInfo struct {
	XPs    int `json:"xps"`
	NewXPs int `json:"new_xps"`
}

// Pulse represents a collection of XPs for different languages at a specific time.
type Pulse struct {
	CodedAt time.Time    `json:"coded_at"`
	XPs     []LanguageXP `json:"xps"`
}

// LanguageXP represents the XP gained for a specific language.
type LanguageXP struct {
	Language string `json:"language"`
	XP       int    `json:"xp"`
}
