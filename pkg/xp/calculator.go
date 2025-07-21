package xp

import (
	"math"

	godestats "github.com/Yeti47/gode-stats/pkg"
)

const (
	// LevelFactor is the constant used in level calculation as per Code::Stats documentation.
	LevelFactor = 0.025
)

// Calculator implements the XpCalculator interface for calculating levels and percentages from XP.
type Calculator struct{}

// NewCalculator creates a new XP calculator instance.
func NewCalculator() godestats.XpCalculator {
	return &Calculator{}
}

// GetLevel calculates the level for the given XP amount.
// Formula: floor(LEVEL_FACTOR * sqrt(xp))
func (c *Calculator) GetLevel(xp int) int {
	if xp < 0 {
		return 0
	}
	return int(math.Floor(LevelFactor * math.Sqrt(float64(xp))))
}

// GetLevelPercentage calculates the percentage progress within the current level.
// Returns a value between 0.0 and 1.0 representing the progress to the next level.
func (c *Calculator) GetLevelPercentage(xp int) float64 {
	if xp < 0 {
		return 0.0
	}

	currentLevel := c.GetLevel(xp)
	currentLevelXP := c.GetXpForLevel(currentLevel)
	nextLevelXP := c.GetXpForLevel(currentLevel + 1)

	// If we're already at the maximum calculable level, return 1.0
	if nextLevelXP <= currentLevelXP {
		return 1.0
	}

	// Calculate percentage within the current level
	xpInCurrentLevel := xp - currentLevelXP
	xpNeededForNextLevel := nextLevelXP - currentLevelXP

	if xpNeededForNextLevel == 0 {
		return 1.0
	}

	percentage := float64(xpInCurrentLevel) / float64(xpNeededForNextLevel)

	// Ensure the percentage is within bounds
	if percentage < 0.0 {
		return 0.0
	}
	if percentage > 1.0 {
		return 1.0
	}

	return percentage
}

// GetXpForLevel calculates the minimum XP required to reach the specified level.
// This is the inverse of the GetLevel function.
// Formula: (level / LEVEL_FACTOR)^2
func (c *Calculator) GetXpForLevel(level int) int {
	if level <= 0 {
		return 0
	}

	// Calculate the minimum XP needed for this level
	levelFloat := float64(level)
	xpFloat := math.Pow(levelFloat/LevelFactor, 2)

	return int(math.Ceil(xpFloat))
}

// GetXpForNextLevel calculates the minimum XP required to reach the next level
// from the current XP amount.
func (c *Calculator) GetXpForNextLevel(xp int) int {
	currentLevel := c.GetLevel(xp)
	return c.GetXpForLevel(currentLevel + 1)
}
