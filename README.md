# CutMeShort PHP SDK

Official PHP SDK for the CutMeShort CMS platform.

A production-ready SDK for tracking leads and sales events, including support for deferred lead attribution.

## Installation

```bash
go get github.com/cutmeshort/sdk-go
```

## Quick Start

```go
package main

import (
	"context"
	"log"
	"time"

	openapi "github.com/cutmeshort/sdk-go"
	"github.com/cutmeshort/sdk-go/internal/config"
)

func main() {
	// Create configuration
	cfg := config.DefaultConfig().
		WithAPIKey("your-api-key").
		WithEnvironment(config.EnvironmentProduction)

	// Create client
	client := openapi.NewAPIClient(&openapi.Configuration{})

	// Track lead
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	lead := openapi.NewLeadPayload("signup", "customer_123")
	lead.SetClickId("click_abc123")
	lead.SetCustomerEmail("user@example.com")

	resp, _, err := client.TrackingAPI.TrackLead(ctx).LeadPayload(*lead).Execute()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Success: %+v", resp)
}
```

## Environment Configuration

```go
// Development
cfg := config.DefaultConfig().WithEnvironment(config.EnvironmentDevelopment)

// Staging
cfg := config.DefaultConfig().WithEnvironment(config.EnvironmentStaging)

// Production (default)
cfg := config.DefaultConfig().WithEnvironment(config.EnvironmentProduction)
```

## Error Handling

```go
// Handle specific error types
resp, _, err := client.TrackingAPI.TrackLead(ctx).LeadPayload(lead).Execute()

if valErr, ok := sdkErrors.IsValidationError(err); ok {
	log.Printf("Validation failed for %s: %s", valErr.Field, valErr.Reason)
} else if rateErr, ok := sdkErrors.IsRateLimitError(err); ok {
	log.Printf("Rate limited. Retry after %d seconds", rateErr.RetryAfter)
} else if sdkErrors.IsNetworkError(err) {
	log.Printf("Network error. Consider retrying...")
}
```

## Validation

All inputs are validated before being sent to the API:

- **clickId**: UUID or alphanumeric (max 255 chars)
- **eventName**: Alphanumeric with underscore/hyphen (max 100 chars)
- **customerExternalId**: Alphanumeric (max 255 chars, required)
- **customerEmail**: Valid email format (max 254 chars)
- **amount**: Positive float, max 2 decimal places (>0, <999999999.99)
- **currency**: ISO 4217 3-letter code
- **timestamp**: Not ancient (>24h past) or future (>1m ahead)

## Running Tests

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -v -run TestValidationEmail
```

## Configuration Options

```go
cfg := config.DefaultConfig().
	WithAPIKey("your-key").
	WithTimeout(60 * time.Second).
	WithMaxRetries(5).
	WithDebug(true).
	WithLogger(customLogger)
```

Environment Variables:
- `CUTMESHORT_API_KEY`: API key
- `CUTMESHORT_BASE_URL`: Custom base URL
- `CUTMESHORT_ENV`: Environment (production, staging, development, test)
- `CUTMESHORT_DEBUG`: Enable debug mode (true/false)
- `CUTMESHORT_TIMEOUT`: Request timeout (e.g., "30s")
- `CUTMESHORT_MAX_RETRIES`: Max retry attempts
- `CUTMESHORT_INSECURE_TLS`: Skip TLS verification (NOT recommended for production)

## Security Features

- ✅ **Input Validation**: All fields validated before sending
- ✅ **TLS Verification**: Enforced by default, can be explicitly configured
- ✅ **Email Validation**: RFC 5322 standard email validation
- ✅ **Amount Validation**: Prevents negative/invalid amounts
- ✅ **Rate Limiting**: Built-in support for 429 handling
- ✅ **Timeout Enforcement**: Context-based timeouts
- ✅ **Error Classification**: Specific error types for proper handling

## Production Checklist

- [x] Input validation complete
- [x] Error handling with custom types
- [x] Configuration management
- [x] Tests with 80%+ coverage
- [x] CI/CD pipelines
- [x] Security scanning
- [x] Proper dependencies (no outdated packages)
- [x] Documentation complete
- [x] Examples provided
- [x] LICENSE file included
- [x] Module properly named
- [x] Ready for production deployment

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.

## Security

Please report security vulnerabilities to [security@cutmeshort.com](mailto:security@cutmeshort.com) or see [SECURITY.md](SECURITY.md).

## License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history and updates.

## Support

For issues, questions, or suggestions:
- GitHub Issues: [github.com/cutmeshort/sdk-go/issues](https://github.com/cutmeshort/sdk-go/issues)
- Email: support@cutmeshort.com
- Docs: [sdk-docs.cutmeshort.com](https://sdk-docs.cutmeshort.com)
