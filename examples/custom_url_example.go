package main

import (
	"fmt"
	"log"

	"github.com/cutmeshort/sdk-go/sdk"
)

// Example 1: Using default production URL
func ExampleDefaultURL() {
	cms := sdk.New("your-api-token")
	
	_, err := cms.TrackLead("click-123", "lead_event", "customer-456")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tracked with default URL: https://www.cutmeshort.com")
}

// Example 2: Using custom URL (e.g., staging/testing)
func ExampleCustomURL() {
	cms := sdk.New("your-api-token")
	
	// Set custom base URL
	cms.SetBaseURL("https://staging-api.example.com")
	
	_, err := cms.TrackLead("click-123", "lead_event", "customer-456")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tracked with custom URL: https://staging-api.example.com")
}

// Example 3: Using local/development URL
func ExampleLocalURL() {
	cms := sdk.New("your-api-token")
	
	// Set localhost for development
	cms.SetBaseURL("http://localhost:8080")
	
	_, err := cms.TrackLead("click-123", "lead_event", "customer-456")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tracked with local URL: http://localhost:8080")
}

// Example 4: Switching URLs on the fly
func ExampleSwitchingURLs() {
	cms := sdk.New("your-api-token")
	
	// Track with production
	cms.SetBaseURL("https://api.example.com")
	_, err := cms.TrackLead("click-1", "lead", "cust-1")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("1️⃣ Tracked with: https://api.example.com")
	
	// Switch to staging
	cms.SetBaseURL("https://staging.example.com")
	_, err = cms.TrackSale("click-2", "sale", "cust-2", "inv-2", 10000, "USD")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("2️⃣ Tracked with: https://staging.example.com")
	
	// Switch back to production
	cms.SetBaseURL("https://api.example.com")
	_, err = cms.TrackLead("click-3", "lead", "cust-3")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("3️⃣ Tracked with: https://api.example.com")
}
