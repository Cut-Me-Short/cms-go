/*
CutMeShort Go SDK - Comprehensive Test Suite

Tests cover validation, error handling, configuration, and logging systems.
*/

package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/cutmeshort/sdk-go/internal/config"
	sdkErrors "github.com/cutmeshort/sdk-go/internal/errors"
	"github.com/cutmeshort/sdk-go/internal/logger"
	"github.com/cutmeshort/sdk-go/internal/validation"
)

// TestValidationClickID tests click ID validation
func TestValidationClickID(t *testing.T) {
	tests := []struct {
		name    string
		clickID string
		isValid bool
		errMsg  string
	}{
		{"valid UUID", "550e8400-e29b-41d4-a716-446655440000", true, ""},
		{"valid alphanumeric", "click_123_abc", true, ""},
		{"empty string", "", false, "must not be empty"},
		{"too long", "a" + string(make([]byte, 256)), false, "must not exceed"},
		{"invalid special chars", "click@123#test", false, "must be a valid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateClickID(tt.clickID)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestValidationEventName tests event name validation
func TestValidationEventName(t *testing.T) {
	tests := []struct {
		name      string
		eventName string
		isValid   bool
	}{
		{"valid event", "user_signup", true},
		{"valid with hyphen", "user-signup", true},
		{"valid complex", "purchase_item_v2", true},
		{"empty string", "", false},
		{"invalid special char", "event@name", false},
		{"invalid space", "event name", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateEventName(tt.eventName)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestValidationCustomerExternalID tests customer ID validation
func TestValidationCustomerExternalID(t *testing.T) {
	tests := []struct {
		name    string
		custID  string
		isValid bool
	}{
		{"valid numeric", "12345", true},
		{"valid UUID", "550e8400-e29b-41d4-a716-446655440000", true},
		{"valid alphanumeric", "cust_XYZ_123", true},
		{"empty string", "", false},
		{"with spaces", "cust 123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateCustomerExternalID(tt.custID)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestValidationEmail tests email validation
func TestValidationEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		isValid bool
	}{
		{"valid email", "user@example.com", true},
		{"valid with plus", "user+tag@example.com", true},
		{"invalid format", "not-an-email", false},
		{"missing @", "usereexample.com", false},
		{"missing domain", "user@", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateEmail(tt.email)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestValidationAmount tests amount validation
func TestValidationAmount(t *testing.T) {
	tests := []struct {
		name    string
		amount  float64
		isValid bool
	}{
		{"valid positive", 99.99, true},
		{"valid large", 999999, true},
		{"zero", 0, false},
		{"negative", -10.50, false},
		{"exceeds max", 1000000000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateAmount(tt.amount)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestValidationCurrency tests currency validation
func TestValidationCurrency(t *testing.T) {
	tests := []struct {
		name     string
		currency string
		isValid  bool
	}{
		{"valid USD", "USD", true},
		{"valid lowercase", "eur", true},
		{"valid GBP", "GBP", true},
		{"invalid code", "INVALID", false},
		{"wrong length", "US", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateCurrency(tt.currency)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestValidationTimestamp tests timestamp validation
func TestValidationTimestamp(t *testing.T) {
	tests := []struct {
		name    string
		ts      time.Time
		isValid bool
	}{
		{"current time", time.Now(), true},
		{"1 minute in future", time.Now().Add(30 * time.Second), true},
		{"1 hour ago", time.Now().Add(-1 * time.Hour), true},
		{"24 hours ago", time.Now().Add(-23 * time.Hour), true},
		{"2 minutes in future", time.Now().Add(2 * time.Minute), false},
		{"25 hours ago", time.Now().Add(-25 * time.Hour), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateTimestamp(tt.ts)
			if tt.isValid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestErrorChaining tests error wrapping and type checking
func TestErrorChaining(t *testing.T) {
	t.Run("validation error", func(t *testing.T) {
		err := validation.ValidateEmail("invalid")
		sdkErr, ok := sdkErrors.IsValidationError(err)
		assert.True(t, ok)
		assert.Equal(t, sdkErrors.ErrorTypeValidation, sdkErr.Type)
		assert.Equal(t, "email", sdkErr.Field)
	})

	t.Run("error is retryable", func(t *testing.T) {
		rateLimitErr := sdkErrors.NewRateLimitError(60)
		assert.True(t, sdkErrors.IsRetryable(rateLimitErr))
	})

	t.Run("validation error not retryable", func(t *testing.T) {
		valErr := validation.ValidateEmail("invalid")
		assert.False(t, sdkErrors.IsRetryable(valErr))
	})
}

// TestConfigurationCreation tests config creation from defaults
func TestConfigurationCreation(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		cfg := config.DefaultConfig()
		assert.NotNil(t, cfg)
		assert.Equal(t, "https://www.cutmeshort.com", cfg.BaseURL)
		assert.Equal(t, config.EnvironmentProduction, cfg.Environment)
		assert.False(t, cfg.Debug)
		assert.NoError(t, cfg.Validate())
	})

	t.Run("staging config", func(t *testing.T) {
		cfg := config.DefaultConfig().WithEnvironment(config.EnvironmentStaging)
		assert.Equal(t, "https://staging.cutmeshort.com", cfg.BaseURL)
		assert.True(t, cfg.Debug)
	})

	t.Run("development config", func(t *testing.T) {
		cfg := config.DefaultConfig().WithEnvironment(config.EnvironmentDevelopment)
		assert.Equal(t, "http://localhost:8080", cfg.BaseURL)
		assert.True(t, cfg.Debug)
	})

	t.Run("with custom settings", func(t *testing.T) {
		cfg := config.DefaultConfig().
			WithAPIKey("secret-key").
			WithBaseURL("https://custom.example.com").
			WithDebug(true)

		assert.Equal(t, "secret-key", cfg.APIKey)
		assert.Equal(t, "https://custom.example.com", cfg.BaseURL)
		assert.True(t, cfg.Debug)
	})

	t.Run("config validation", func(t *testing.T) {
		cfg := config.DefaultConfig()
		cfg.BaseURL = ""
		err := cfg.Validate()
		assert.Error(t, err)
	})
}

// TestLoggerCreation tests creating logger instances
func TestLoggerCreation(t *testing.T) {
	t.Run("simple logger", func(t *testing.T) {
		log := logger.New("test", true)
		assert.NotNil(t, log)
		log.Debug("debug message", "key", "value")
		log.Info("info message")
		log.Warn("warn message", "code", 42)
		log.Error("error message", "error", "test error")
	})

	t.Run("noop logger", func(t *testing.T) {
		log := logger.NewNoOpLogger()
		assert.NotNil(t, log)
		log.Debug("debug")
		log.Info("info")
		log.Warn("warn")
		log.Error("error")
	})
}

// TestHTTPClientConfiguration tests that configuration creates proper HTTP client
func TestHTTPClientConfiguration(t *testing.T) {
	t.Run("production client", func(t *testing.T) {
		cfg := config.DefaultConfig()
		assert.NotNil(t, cfg.HTTPClient)
		assert.False(t, cfg.InsecureTLS)
	})

	t.Run("custom timeout", func(t *testing.T) {
		cfg := config.DefaultConfig().WithTimeout(5 * time.Second)
		assert.Equal(t, 5*time.Second, cfg.Timeout)
	})
}

// TestValidationIntegration tests validators working together
func TestValidationIntegration(t *testing.T) {
	t.Run("lead event validation", func(t *testing.T) {
		validations := map[string]error{
			"clickId":           validation.ValidateClickID("550e8400-e29b-41d4-a716-446655440000"),
			"eventName":         validation.ValidateEventName("lead_signup"),
			"customerId":        validation.ValidateCustomerExternalID("cust_12345"),
			"customerName":      validation.ValidateCustomerName("John Doe"),
			"customerEmail":     validation.ValidateEmail("john@example.com"),
		}

		for field, err := range validations {
			assert.NoError(t, err, "validation failed for %s", field)
		}
	})

	t.Run("sale event validation", func(t *testing.T) {
		validations := map[string]error{
			"clickId":           validation.ValidateClickID("550e8400-e29b-41d4-a716-446655440000"),
			"eventName":         validation.ValidateEventName("sale_completed"),
			"customerId":        validation.ValidateCustomerExternalID("cust_12345"),
			"invoiceId":         validation.ValidateInvoiceID("inv_001"),
			"amount":            validation.ValidateAmount(99.99),
			"currency":          validation.ValidateCurrency("USD"),
			"timestamp":         validation.ValidateTimestamp(time.Now()),
		}

		for field, err := range validations {
			assert.NoError(t, err, "validation failed for %s", field)
		}
	})
}
