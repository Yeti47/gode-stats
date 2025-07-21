# Gode Stats

A Go client library for the Code::Stats API.

## Features

- Code::Stats API client with full support for profile retrieval and pulse submission
- XP calculator with floating-point percentage calculations
- Supports both authenticated and anonymous access (if endpoint allows it)
- Clean interfaces following Go best practices
- Token-based authentication
- Comprehensive structured error handling with specific error types

## Installation

```bash
go get github.com/Yeti47/gode-stats
```

## Authentication

Get your API token from the [Code::Stats machine control panel](https://codestats.net/my/machines).

## Usage

### Creating a Client

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/Yeti47/gode-stats/pkg/client"
)

func main() {
    // Authenticated client for production
    apiToken := "your-api-token-here"
    c := client.New(apiToken)
    
    // Anonymous client (read-only, no token required)
    anonClient := client.NewAnonymous()
    
    // Get user profile
    profile, err := c.GetUserProfile(context.Background(), "username")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("User: %s, Total XP: %d\n", profile.User, profile.TotalXP)
}
```

### Sending Pulses

```go
package main

import (
    "context"
    "time"
    
    "github.com/Yeti47/gode-stats/pkg/client"
    "github.com/Yeti47/gode-stats/pkg"
)

func main() {
    c := client.New("your-api-token")
    
    // Create a pulse
    pulse := godestats.Pulse{
        CodedAt: time.Now(),
        XPs: []godestats.LanguageXP{
            {Language: "Go", XP: 25},
            {Language: "JavaScript", XP: 15},
        },
    }
    
    err := c.SendPulse(context.Background(), pulse)
    if err != nil {
        panic(err)
    }
}
```

### Calculating XP and Levels

```go
package main

import (
    "fmt"
    
    "github.com/Yeti47/gode-stats/pkg/xp"
)

func main() {
    calc := xp.NewCalculator()
    
    xpAmount := 10000
    level := calc.GetLevel(xpAmount)
    percentage := calc.GetLevelPercentage(xpAmount)
    
    fmt.Printf("XP: %d, Level: %d, Progress: %.2f%%\n", 
               xpAmount, level, percentage*100)
}
```

### Error Handling

The library provides comprehensive error handling with specific error types that you can check for:

```go
package main

import (
    "context"
    "errors"
    "fmt"
    
    "github.com/Yeti47/gode-stats/pkg/client"
    "github.com/Yeti47/gode-stats/pkg"
)

func main() {
    c := client.NewAnonymous()
    
    _, err := c.GetUserProfile(context.Background(), "nonexistent-user")
    if err != nil {
        // Check for specific error types
        switch {
        case godestats.IsUserNotFound(err):
            fmt.Println("User not found - try a different username")
        case godestats.IsUnauthorized(err):
            fmt.Println("Need to provide an API token")
        case godestats.IsRateLimited(err):
            fmt.Println("Rate limited - wait before retrying")
        case godestats.IsTemporary(err):
            fmt.Println("Temporary error - safe to retry")
        case godestats.IsNetworkError(err):
            fmt.Println("Network error - check connectivity")
        default:
            fmt.Printf("Other error: %v\n", err)
        }
        
        // Or check for specific error variables
        if errors.Is(err, godestats.ErrUserNotFound) {
            fmt.Println("Definitely a user not found error")
        }
        
        // Or extract detailed error information
        var apiErr *godestats.APIError
        if errors.As(err, &apiErr) {
            fmt.Printf("API error: status %d, message: %s\n", 
                       apiErr.StatusCode, apiErr.Message)
        }
    }
}
```

## API Reference

See the [Code::Stats API documentation](https://codestats.net/api-docs) for more information about the API endpoints.
