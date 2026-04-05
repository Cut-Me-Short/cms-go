package main

import (
	"context"
	"fmt"
	"log"

	openapi "github.com/cutmeshort/sdk-go"
	"github.com/cutmeshort/sdk-go/internal/config"
)

func main() {
	// Create configuration
	cfg := config.DefaultConfig().
		WithAPIKey("your-api-key-here").
		WithEnvironment(config.EnvironmentProduction)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Create API client
	apiCfg := &openapi.Configuration{
		UserAgent: "CutMeShort-SDK/1.0.0",
	}
	client := openapi.NewAPIClient(apiCfg)

	// Create a lead tracking request
	leadPayload := openapi.NewLeadPayload(
		"user_signup",           // eventName
		"customer_12345",        // customerExternalId
	)

	// Optional: Set additional fields
	leadPayload.SetClickId("click_550e8400e29b41d4a716446655440000")
	leadPayload.SetCustomerName("John Doe")
	leadPayload.SetCustomerEmail("john@example.com")

	// Execute the request with context
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, httpRes, err := client.TrackingAPI.
		TrackLead(ctx).
		LeadPayload(*leadPayload).
		Execute()

	if err != nil {
		log.Fatalf("API request failed: %v", err)
	}

	fmt.Printf("Response Status: %d\n", httpRes.StatusCode)
	fmt.Printf("Track Response: %+v\n", resp)
}
