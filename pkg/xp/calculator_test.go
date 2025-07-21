package xp

import (
	"math"
	"testing"
)

func TestCalculator_GetLevel(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		xp       int
		expected int
	}{
		{"Zero XP", 0, 0},
		{"Negative XP", -100, 0},
		{"Small XP", 100, 0},
		{"Level 1", 1600, 1},
		{"Level 2", 6400, 2},
		{"Level 5", 40000, 5},
		{"Level 10", 160000, 10},
		{"Large XP", 1000000, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.GetLevel(tt.xp)
			if result != tt.expected {
				t.Errorf("GetLevel(%d) = %d, expected %d", tt.xp, result, tt.expected)
			}
		})
	}
}

func TestCalculator_GetXpForLevel(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		level    int
		expected int
	}{
		{"Level 0", 0, 0},
		{"Level 1", 1, 1600},
		{"Level 2", 2, 6400},
		{"Level 5", 5, 40000},
		{"Level 10", 10, 160000},
		{"Level 25", 25, 1000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.GetXpForLevel(tt.level)
			if result != tt.expected {
				t.Errorf("GetXpForLevel(%d) = %d, expected %d", tt.level, result, tt.expected)
			}
		})
	}
}

func TestCalculator_GetLevelPercentage(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		xp       int
		expected float64
		delta    float64
	}{
		{"Zero XP", 0, 0.0, 0.01},
		{"Negative XP", -100, 0.0, 0.01},
		{"Start of level 1", 1600, 0.0, 0.01},
		{"Middle of level 1", 4000, 0.5, 0.01},
		{"Almost level 2", 6300, 0.979, 0.01},
		{"Start of level 2", 6400, 0.0, 0.01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.GetLevelPercentage(tt.xp)
			if math.Abs(result-tt.expected) > tt.delta {
				t.Errorf("GetLevelPercentage(%d) = %f, expected %f ± %f",
					tt.xp, result, tt.expected, tt.delta)
			}
		})
	}
}

func TestCalculator_GetXpForNextLevel(t *testing.T) {
	calc := NewCalculator()

	tests := []struct {
		name     string
		xp       int
		expected int
	}{
		{"Zero XP", 0, 1600},
		{"XP at level 1", 1600, 6400},
		{"XP in middle of level 1", 4000, 6400},
		{"XP at level 2", 6400, 14400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calc.GetXpForNextLevel(tt.xp)
			if result != tt.expected {
				t.Errorf("GetXpForNextLevel(%d) = %d, expected %d", tt.xp, result, tt.expected)
			}
		})
	}
}

// TestLevelCalculationConsistency ensures that level calculations are consistent
// between GetLevel and GetXpForLevel functions.
func TestLevelCalculationConsistency(t *testing.T) {
	calc := NewCalculator()

	for level := 1; level <= 50; level++ {
		minXpForLevel := calc.GetXpForLevel(level)
		calculatedLevel := calc.GetLevel(minXpForLevel)

		if calculatedLevel != level {
			t.Errorf("Inconsistency at level %d: GetXpForLevel(%d) = %d, but GetLevel(%d) = %d",
				level, level, minXpForLevel, minXpForLevel, calculatedLevel)
		}
	}
}

// BenchmarkGetLevel benchmarks the GetLevel function.
func BenchmarkGetLevel(b *testing.B) {
	calc := NewCalculator()

	for i := 0; i < b.N; i++ {
		calc.GetLevel(100000)
	}
}

// BenchmarkGetLevelPercentage benchmarks the GetLevelPercentage function.
func BenchmarkGetLevelPercentage(b *testing.B) {
	calc := NewCalculator()

	for i := 0; i < b.N; i++ {
		calc.GetLevelPercentage(100000)
	}
}
