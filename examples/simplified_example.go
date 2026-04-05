package main

import (
	"fmt"
	"log"

	"github.com/cutmeshort/sdk-go/sdk"
)

func ExampleSimplifiedUsage() {
	// That's it! Initialize with just your token
	cms := sdk.New("your-api-token-here")

	// Track a lead - one simple function call
	leadResp, err := cms.TrackLead(
		"click-123",                    // clickId
		"lead_captured",                // eventName
		"customer-456",                 // customerExternalId
	)
	if err != nil {
		log.Fatal("Error tracking lead:", err)
	}
	fmt.Printf("Lead tracked: %+v\n", leadResp)

	// Track a sale - one simple function call
	saleResp, err := cms.TrackSale(
		"click-123",                    // clickId
		"sale_completed",               // eventName
		"customer-456",                 // customerExternalId
		"invoice-789",                  // invoiceId
		5000,                           // amount in cents (e.g., $50.00)
		"USD",                          // currency code
	)
	if err != nil {
		log.Fatal("Error tracking sale:", err)
	}
	fmt.Printf("Sale tracked: %+v\n", saleResp)
}

func ExampleAdvancedUsage() {
	cms := sdk.New("your-api-token-here")

	// You can also use the advanced versions with all optional fields
	leadData := sdk.LeadData{
		ClickId:            "click-123",
		EventName:          "lead_captured",
		CustomerExternalId: "customer-456",
		CustomerName:       sdk.StringPtr("John Doe"),
		CustomerEmail:      sdk.StringPtr("john@example.com"),
		CustomerAvatar:     sdk.StringPtr("https://example.com/avatar.jpg"),
		Mode:               sdk.StringPtr("deferred"), // for deferred attribution
	}

	leadResp, err := cms.TrackLeadAdvanced(leadData)
	if err != nil {
		log.Fatal("Error tracking lead:", err)
	}
	fmt.Printf("Lead tracked: %+v\n", leadResp)

	saleData := sdk.SaleData{
		ClickId:            "click-123",
		EventName:          "sale_completed",
		CustomerExternalId: "customer-456",
		InvoiceId:          "invoice-789",
		Amount:             5000,
		Currency:           "USD",
		CustomerName:       sdk.StringPtr("John Doe"),
		CustomerEmail:      sdk.StringPtr("john@example.com"),
	}

	saleResp, err := cms.TrackSaleAdvanced(saleData)
	if err != nil {
		log.Fatal("Error tracking sale:", err)
	}
	fmt.Printf("Sale tracked: %+v\n", saleResp)
}
