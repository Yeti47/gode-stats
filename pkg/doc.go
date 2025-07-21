/*
Package godestats provides a client library for the Code::Stats API.

Code::Stats is a programming statistics service that tracks your coding activity
across different languages and projects. This library provides a clean, idiomatic
Go interface for interacting with the Code::Stats API.

Key Features:
- Full support for Code::Stats API v1.0.1
- Floating-point XP calculations for precise level percentages
- Token-based authentication
- Comprehensive error handling
- Clean separation of concerns with well-defined interfaces

Basic Usage:

	package main

	import (
		"context"
		"fmt"
		"log"

		"github.com/user/gode-stats/pkg/client"
		"github.com/user/gode-stats/pkg/xp"
		godestats "github.com/user/gode-stats/pkg"
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
