package main

import (
	"fmt"
	"log"

	"github.com/cutmeshort/sdk-go/internal/config"
	"github.com/cutmeshort/sdk-go/internal/validation"
)

// SimpleLogger is a basic logger implementation
type SimpleLogger struct{}

func (l *SimpleLogger) Debug(msg string, keysAndValues ...interface{}) {
	fmt.Printf("[DEBUG] %s %v\n", msg, keysAndValues)
}

func (l *SimpleLogger) Info(msg string, keysAndValues ...interface{}) {
	fmt.Printf("[INFO] %s %v\n", msg, keysAndValues)
}

func (l *SimpleLogger) Warn(msg string, keysAndValues ...interface{}) {
	fmt.Printf("[WARN] %s %v\n", msg, keysAndValues)
}

func (l *SimpleLogger) Error(msg string, keysAndValues ...interface{}) {
	fmt.Printf("[ERROR] %s %v\n", msg, keysAndValues)
}

func main() {
	// Example 1: Configuration with environment variables
	fmt.Println("=== Configuration Example ===")

	cfg := config.NewConfigFromEnv().
		WithAPIKey("your-secret-key").
		WithLogger(&SimpleLogger{})

	cfg.Logger.Info("Configuration loaded", "environment", cfg.Environment, "baseURL", cfg.BaseURL)

	// Example 2: Input validation
	fmt.Println("\n=== Validation Example ===")

	// Valid email
	if err := validation.ValidateEmail("user@example.com"); err != nil {
		log.Printf("Invalid email: %v", err)
	} else {
		fmt.Println("Email is valid")
	}

	// Invalid email
	if err := validation.ValidateEmail("notanemail"); err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	// Valid amount
	if err := validation.ValidateAmount(99.99); err != nil {
		log.Printf("Invalid amount: %v", err)
	} else {
		fmt.Println("Amount is valid")
	}

	// Invalid currency
	if err := validation.ValidateCurrency("INVALID"); err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}

	// Example 3: Custom headers
	fmt.Println("\n=== Custom Headers Example ===")

	cfg.AddHeader("X-Custom-Header", "custom-value").
		AddHeader("X-Request-ID", "req-12345")

	for key, value := range cfg.Headers {
		fmt.Printf("Header: %s = %s\n", key, value)
	}

	// Example 4: Configuration for different environments
	fmt.Println("\n=== Environment Configuration Example ===")

	environments := []config.Environment{
		config.EnvironmentDevelopment,
		config.EnvironmentStaging,
		config.EnvironmentProduction,
	}

	for _, env := range environments {
		envCfg := config.DefaultConfig().WithEnvironment(env)
		fmt.Printf("Environment: %s -> BaseURL: %s (Debug: %v)\n", env, envCfg.BaseURL, envCfg.Debug)
	}
}
