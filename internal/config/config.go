/*
CutMeShort Go SDK - Configuration

Package config provides configuration management for the SDK with support
for environment variables and custom settings.
*/

package config

import (
	"crypto/tls"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Environment represents the deployment environment
type Environment string

const (
	EnvironmentProduction  Environment = "production"
	EnvironmentStaging     Environment = "staging"
	EnvironmentDevelopment Environment = "development"
	EnvironmentTest        Environment = "test"
)

// Config holds the SDK configuration
type Config struct {
	// API Configuration
	BaseURL   string
	APIKey    string
	UserAgent string

	// Environment
	Environment Environment

	// HTTP Client Configuration
	HTTPClient         *http.Client
	Timeout            time.Duration
	MaxRetries         int
	BackoffMultiplier  float64
	InitialBackoff     time.Duration
	MaxBackoff         time.Duration

	// TLS Configuration
	InsecureTLS bool

	// Request Configuration
	Headers map[string]string

	// Logging Configuration
	Debug  bool
	Logger Logger

	// Validation Configuration
	StrictValidation bool
}

// Logger is an interface for logging
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// DefaultConfig returns a production-ready configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		BaseURL:           "https://www.cutmeshort.com",
		Environment:       EnvironmentProduction,
		Timeout:           30 * time.Second,
		MaxRetries:        3,
		BackoffMultiplier: 2.0,
		InitialBackoff:    100 * time.Millisecond,
		MaxBackoff:        30 * time.Second,
		InsecureTLS:       false,
		Headers:           make(map[string]string),
		Debug:             false,
		Logger:            NewNoOpLogger(),
		StrictValidation:  true,
		HTTPClient:        createDefaultHTTPClient(),
	}
}

// NewConfigFromEnv creates a configuration from environment variables
func NewConfigFromEnv() *Config {
	cfg := DefaultConfig()

	// Read from environment variables
	if apiKey := os.Getenv("CUTMESHORT_API_KEY"); apiKey != "" {
		cfg.APIKey = apiKey
	}

	if baseURL := os.Getenv("CUTMESHORT_BASE_URL"); baseURL != "" {
		cfg.BaseURL = baseURL
	}

	if env := os.Getenv("CUTMESHORT_ENV"); env != "" {
		cfg.Environment = Environment(env)
	}

	if debug := os.Getenv("CUTMESHORT_DEBUG"); debug != "" {
		cfg.Debug, _ = strconv.ParseBool(debug)
	}

	if timeout := os.Getenv("CUTMESHORT_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			cfg.Timeout = d
		}
	}

	if maxRetries := os.Getenv("CUTMESHORT_MAX_RETRIES"); maxRetries != "" {
		if mr, err := strconv.Atoi(maxRetries); err == nil {
			cfg.MaxRetries = mr
		}
	}

	if insecureTLS := os.Getenv("CUTMESHORT_INSECURE_TLS"); insecureTLS != "" {
		cfg.InsecureTLS, _ = strconv.ParseBool(insecureTLS)
	}

	return cfg
}

// WithAPIKey sets the API key
func (c *Config) WithAPIKey(key string) *Config {
	c.APIKey = key
	return c
}

// WithBaseURL sets the base URL
func (c *Config) WithBaseURL(url string) *Config {
	c.BaseURL = url
	return c
}

// WithEnvironment sets the environment
func (c *Config) WithEnvironment(env Environment) *Config {
	c.Environment = env
	switch env {
	case EnvironmentProduction:
		c.BaseURL = "https://www.cutmeshort.com"
		c.Debug = false
	case EnvironmentStaging:
		c.BaseURL = "https://staging.cutmeshort.com"
		c.Debug = true
	case EnvironmentDevelopment:
		c.BaseURL = "http://localhost:8080"
		c.Debug = true
	case EnvironmentTest:
		c.BaseURL = "http://localhost:3000"
		c.Debug = true
	}
	return c
}

// WithTimeout sets the request timeout
func (c *Config) WithTimeout(timeout time.Duration) *Config {
	c.Timeout = timeout
	if c.HTTPClient != nil {
		c.HTTPClient.Timeout = timeout
	}
	return c
}

// WithMaxRetries sets the maximum number of retries
func (c *Config) WithMaxRetries(maxRetries int) *Config {
	c.MaxRetries = maxRetries
	return c
}

// WithHTTPClient sets a custom HTTP client
func (c *Config) WithHTTPClient(client *http.Client) *Config {
	c.HTTPClient = client
	return c
}

// WithDebug enables or disables debug mode
func (c *Config) WithDebug(debug bool) *Config {
	c.Debug = debug
	return c
}

// WithLogger sets a custom logger
func (c *Config) WithLogger(logger Logger) *Config {
	c.Logger = logger
	return c
}

// AddHeader adds a custom header
func (c *Config) AddHeader(key, value string) *Config {
	c.Headers[key] = value
	return c
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.APIKey == "" {
		// API key might be optional for some endpoints, but warn if not set
		if c.Debug && c.Logger != nil {
			c.Logger.Warn("API key not configured")
		}
	}

	if c.BaseURL == "" {
		return NewConfigError("BaseURL must not be empty")
	}

	if c.Timeout < time.Second {
		return NewConfigError("Timeout must be at least 1 second")
	}

	if c.MaxRetries < 0 {
		return NewConfigError("MaxRetries must be non-negative")
	}

	return nil
}

// createDefaultHTTPClient creates a default HTTP client with TLS verification enabled
func createDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     100,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // Always verify TLS in production
			},
		},
	}
}

// ConfigError represents a configuration error
type ConfigError struct {
	message string
}

// NewConfigError creates a new configuration error
func NewConfigError(message string) *ConfigError {
	return &ConfigError{message: message}
}

// Error implements the error interface
func (e *ConfigError) Error() string {
	return "configuration error: " + e.message
}

// NoOpLogger is a logger that does nothing (default logger)
type NoOpLogger struct{}

// Debug logs a debug message
func (l *NoOpLogger) Debug(msg string, keysAndValues ...interface{}) {}

// Info logs an info message
func (l *NoOpLogger) Info(msg string, keysAndValues ...interface{}) {}

// Warn logs a warning message
func (l *NoOpLogger) Warn(msg string, keysAndValues ...interface{}) {}

// Error logs an error message
func (l *NoOpLogger) Error(msg string, keysAndValues ...interface{}) {}

// NewNoOpLogger creates a new no-op logger
func NewNoOpLogger() Logger {
	return &NoOpLogger{}
}
