package main

import (
	"context"
	"fmt"
	"log"
	"time"

	openapi "github.com/cutmeshort/sdk-go"
	"github.com/cutmeshort/sdk-go/internal/config"
)

func main() {
	// Create configuration for staging environment
	cfg := config.DefaultConfig().
		WithAPIKey("your-api-key-here").
		WithEnvironment(config.EnvironmentStaging)

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Create API client
	apiCfg := &openapi.Configuration{
		UserAgent: "CutMeShort-SDK/1.0.0",
	}
	client := openapi.NewAPIClient(apiCfg)

	// Create a sale tracking request
	salePayload := openapi.NewSalePayload(
		"click_550e8400e29b41d4a716446655440000", // clickId
		"purchase_completed",                      // eventName
		"customer_12345",                         // customerExternalId
		"INV_2024_001",                           // invoiceId
		9999,                                      // amount (in cents) = $99.99
		"USD",                                     // currency
	)

	// Optional: Set additional fields
	salePayload.Timestamp = &struct{*time.Time}{&time.Now()}.Time
	salePayload.SetCustomerName("John Doe")
	salePayload.SetCustomerEmail("john@example.com")

	// Execute the request
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, httpRes, err := client.TrackingAPI.
		TrackSale(ctx).
		SalePayload(*salePayload).
		Execute()

	if err != nil {
		log.Fatalf("API request failed: %v", err)
	}

	fmt.Printf("Response Status: %d\n", httpRes.StatusCode)
	fmt.Printf("Track Response: %+v\n", resp)
}
