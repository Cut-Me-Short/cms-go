# CutMeShort Go SDK - Getting Started Guide

## Installation

```bash
go get github.com/cutmeshort/sdk-go
```

## Basic Setup

### 1. Create a Configuration

```go
package main

import (
    "github.com/cutmeshort/sdk-go/internal/config"
)

func main() {
    // Default production configuration
    cfg := config.DefaultConfig().
        WithAPIKey("your-api-key-here")
    
    // Validate configuration
    if err := cfg.Validate(); err != nil {
        panic(err)
    }
}
```

### 2. Create an API Client

```go
import openapi "github.com/cutmeshort/sdk-go"

func main() {
    // ... config setup ...
    
    // Create client
    client := openapi.NewAPIClient(&openapi.Configuration{
        UserAgent: "MyApp/1.0.0",
    })
}
```

### 3. Make Your First Request

```go
import (
    "context"
    "time"
    openapi "github.com/cutmeshort/sdk-go"
)

func main() {
    // ... config and client setup ...
    
    // Create lead payload
    lead := openapi.NewLeadPayload(
        "user_signup",      // eventName
        "customer_12345",   // customerExternalId
    )
    
    // Set optional fields
    lead.SetClickId("click_abc123")
    lead.SetCustomerEmail("user@example.com")
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Send request
    response, httpResp, err := client.TrackingAPI.
        TrackLead(ctx).
        LeadPayload(*lead).
        Execute()
    
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Success: %+v\n", response)
}
```

## Environment Configuration

### Development

```go
cfg := config.DefaultConfig().
    WithEnvironment(config.EnvironmentDevelopment).
    WithDebug(true)

// This sets:
// - BaseURL: http://localhost:8080
// - Debug: true
```

### Staging

```go
cfg := config.DefaultConfig().
    WithEnvironment(config.EnvironmentStaging).
    WithAPIKey(os.Getenv("STAGING_API_KEY"))

// This sets:
// - BaseURL: https://staging.cutmeshort.com
// - Debug: true
```

### Production

```go
cfg := config.DefaultConfig().
    WithEnvironment(config.EnvironmentProduction).
    WithAPIKey(os.Getenv("CUTMESHORT_API_KEY"))

// This sets:
// - BaseURL: https://www.cutmeshort.com
// - Debug: false
// - TLS verification: enabled
```

## Error Handling

### Basic Error Handling

```go
resp, _, err := client.TrackingAPI.TrackLead(ctx).LeadPayload(lead).Execute()

if err != nil {
    log.Printf("Error: %v", err)
}
```

### Specific Error Types

```go
import (
    sdkErrors "github.com/cutmeshort/sdk-go/internal/errors"
    "github.com/cutmeshort/sdk-go/internal/validation"
)

resp, _, err := client.TrackingAPI.TrackLead(ctx).LeadPayload(lead).Execute()

if err != nil {
    // Check for validation error
    if valErr, ok := sdkErrors.IsValidationError(err); ok {
        log.Printf("Validation failed for %s: %s", valErr.Field, valErr.Reason)
        return
    }
    
    // Check for rate limiting
    if rateErr, ok := sdkErrors.IsRateLimitError(err); ok {
        log.Printf("Rate limited. Retry after %d seconds", rateErr.RetryAfter)
        return
    }
    
    // Check if retryable
    if sdkErrors.IsRetryable(err) {
        log.Printf("Transient error, consider retrying")
        return
    }
    
    // Generic error
    log.Printf("API error: %v", err)
}
```

## Input Validation

### Validating Fields Manually

```go
import "github.com/cutmeshort/sdk-go/internal/validation"

// Validate email
if err := validation.ValidateEmail(email); err != nil {
    log.Printf("Invalid email: %v", err)
}

// Validate amount
if err := validation.ValidateAmount(amount); err != nil {
    log.Printf("Invalid amount: %v", err)
}

// Validate currency
if err := validation.ValidateCurrency("USD"); err != nil {
    log.Printf("Invalid currency: %v", err)
}
```

### Validation Rules

| Field | Rule | Example |
|-------|------|---------|
| `clickId` | UUID or alphanumeric (max 255) | `click_abc123` or UUID |
| `eventName` | Alphanumeric + underscore/hyphen (max 100) | `user_signup` |
| `customerExternalId` | Alphanumeric (max 255) | `cust_12345` |
| `customerEmail` | Valid email (max 254) | `user@example.com` |
| `amount` | Positive, 2 decimals (max 999999999.99) | `99.99` |
| `currency` | ISO 4217 code (3 letters) | `USD`, `EUR` |
| `timestamp` | Not ancient (>24h) or future (>1m) | Now ± 1 minute OK |

## Logging

### Using a Custom Logger

```go
type ConsoleLogger struct{}

func (l *ConsoleLogger) Debug(msg string, keysAndValues ...interface{}) {
    fmt.Printf("[DEBUG] %s %v\n", msg, keysAndValues)
}

func (l *ConsoleLogger) Info(msg string, keysAndValues ...interface{}) {
    fmt.Printf("[INFO] %s %v\n", msg, keysAndValues)
}

func (l *ConsoleLogger) Warn(msg string, keysAndValues ...interface{}) {
    fmt.Printf("[WARN] %s %v\n", msg, keysAndValues)
}

func (l *ConsoleLogger) Error(msg string, keysAndValues ...interface{}) {
    fmt.Printf("[ERROR] %s %v\n", msg, keysAndValues)
}

// Use with config
cfg := config.DefaultConfig().WithLogger(&ConsoleLogger{})
```

### Structured Logging with zerolog (Recommended)

```go
import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

// Adapt zerolog to SDK logger interface
type ZerologyAdapter struct {
    logger zerolog.Logger
}

func (z *ZerologyAdapter) Debug(msg string, keysAndValues ...interface{}) {
    z.logger.Debug().Fields(parseKeyValues(keysAndValues)).Msg(msg)
}

// ... implement other methods ...

// Use
cfg := config.DefaultConfig().WithLogger(&ZerologyAdapter{logger: log.Logger})
```

## Advanced Configuration

### Custom HTTP Client

```go
import "net/http"

customClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        50,
        MaxIdleConnsPerHost: 5,
    },
}

cfg := config.DefaultConfig().WithHTTPClient(customClient)
```

### Custom Headers

```go
cfg.AddHeader("X-Custom-Header", "custom-value")
cfg.AddHeader("X-Request-ID", generateRequestID())
```

### Loading from Environment Variables

```go
// Set environment variables
// CUTMESHORT_API_KEY=your-key
// CUTMESHORT_ENV=staging
// CUTMESHORT_DEBUG=true

cfg := config.NewConfigFromEnv()
```

## Patterns and Best Practices

### Resource Cleanup

```go
// Always use context.WithTimeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Use the context
resp, _, err := client.TrackingAPI.TrackLead(ctx).LeadPayload(lead).Execute()
```

### Retry Logic

```go
import (
    "time"
    sdkErrors "github.com/cutmeshort/sdk-go/internal/errors"
)

func retryWithBackoff(fn func() error, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        
        if !sdkErrors.IsRetryable(err) {
            return err
        }
        
        backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
        time.Sleep(backoff)
    }
    return fmt.Errorf("max retries exceeded")
}

// Usage
err := retryWithBackoff(func() error {
    _, _, err := client.TrackingAPI.TrackLead(ctx).LeadPayload(lead).Execute()
    return err
}, 3)
```

### Concurrent Requests

```go
// SDK is safe for concurrent use
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(index int) {
        defer wg.Done()
        // Safe to use same client concurrently
        _, _, _ = client.TrackingAPI.TrackLead(ctx).LeadPayload(lead).Execute()
    }(i)
}
wg.Wait()
```

## Testing

### Unit Testing

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestTrackingFlow(t *testing.T) {
    cfg := config.DefaultConfig()
    client := openapi.NewAPIClient(&openapi.Configuration{})
    
    assert.NotNil(t, client)
    assert.NotNil(t, client.TrackingAPI)
}
```

### Integration Testing with Mock Server

See [examples/track_lead.go](examples/track_lead.go) for full examples.

## Troubleshooting

### Common Issues

**Issue**: "API key not configured"
- **Solution**: Set `CUTMESHORT_API_KEY` environment variable or use `WithAPIKey()`

**Issue**: "Invalid email" validation error
- **Solution**: Ensure email is in standard format (user@domain.com)

**Issue**: "Amount must be greater than 0"
- **Solution**: Amounts must be positive floats (e.g., 99.99, not -10 or 0)

**Issue**: "Timeout exceeded"
- **Solution**: Increase timeout using `WithTimeout()` or check network connectivity

**Issue**: "Rate limit exceeded"
- **Solution**: Implement backoff retry logic or reduce request frequency

## Next Steps

1. Check [examples/](examples/) for complete working examples
2. Read [ARCHITECTURE.md](ARCHITECTURE.md) for detailed design
3. Review [SECURITY.md](SECURITY.md) for security practices
4. See [CONTRIBUTING.md](CONTRIBUTING.md) if you want to contribute

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/cutmeshort/sdk-go/issues)
- **Email**: support@cutmeshort.com
- **Documentation**: [docs/](docs/)
