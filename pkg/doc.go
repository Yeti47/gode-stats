/*
Package godestats provides a client library for the Code::Stats API.

Repository: https://github.com/Yeti47/gode-stats
License: MIT

Code::Stats is a programming statistics service that tracks your coding activity
across different languages and projects. This library provides a clean, idiomatic
Go interface for interacting with the Code::Stats API.

Features:
- Code::Stats API client with full support for profile retrieval and pulse submission
- XP calculator with floating-point percentage calculations
- Token-based authentication
- Comprehensive error handling
- Clean interfaces following Go best practices

Basic Usage:

	package main

	import (
		"context"
		"fmt"
		"log"

		"github.com/Yeti47/gode-stats/pkg/client"
		"github.com/Yeti47/gode-stats/pkg/xp"
		godestats "github.com/Yeti47/gode-stats/pkg"
	)

	func main() {
		// Create API client
		c := client.New("your-api-token")

		// Get user profile
		profile, err := c.GetUserProfile(context.Background(), "username")
		if err != nil {
			log.Fatal(err)
		}

		// Use XP calculator
		calc := xp.NewCalculator()
		level := calc.GetLevel(profile.TotalXP)
		percentage := calc.GetLevelPercentage(profile.TotalXP)

		fmt.Printf("User %s is at level %d (%.2f%% progress)\n",
			profile.User, level, percentage*100)
	}

Interfaces:

The package defines two main interfaces:

1. CodeStatsClient: For interacting with the Code::Stats API
2. XpCalculator: For calculating levels and percentages from XP values

This design allows for easy testing and alternative implementations.

XP Calculation:

The XP calculator uses the official Code::Stats formula:

	level = floor(0.025 * sqrt(xp))

Unlike the reference implementation, this library returns floating-point
percentages, giving consumers the flexibility to round as needed.

Authentication:

API authentication uses tokens that can be generated from the Code::Stats
machine control panel. Tokens are sent via the X-API-Token HTTP header.

Error Handling:

The library provides detailed error messages and appropriate error types
for different failure scenarios, including network errors, authentication
failures, and API-specific errors.
*/
package godestats
