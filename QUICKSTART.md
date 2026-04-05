# Quick Start - CutMeShort SDK

Get started in 30 seconds with just your API token.

## Installation

```bash
go get github.com/cutmeshort/sdk-go
```

## Basic Usage - That's it!

```go
package main

import (
	"fmt"
	"log"
	"github.com/cutmeshort/sdk-go/sdk"
)

func main() {
	// 1. Initialize with your token (just one line!)
	cms := sdk.New("your-api-token-here")

	// 2. Track a lead (one simple function call)
	_, err := cms.TrackLead(
		"click-123",      // clickId
		"lead_captured",  // eventName  
		"customer-456",   // customerExternalId
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Lead tracked!")

	// 3. Track a sale (one simple function call)
	_, err = cms.TrackSale(
		"click-123",       // clickId
		"sale_completed",  // eventName
		"customer-456",    // customerExternalId
		"invoice-789",     // invoiceId
		5000,              // amount in cents ($50.00)
		"USD",             // currency
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sale tracked!")
}
```

## That's All!

No configuration, no builders, no complexity. Just:
1. **Initialize**: `cms := sdk.New("token")`
2. **Track Lead**: `cms.TrackLead(clickId, eventName, customerId)`
3. **Track Sale**: `cms.TrackSale(clickId, eventName, customerId, invoiceId, amount, currency)`

## Advanced Usage (Optional)

If you need to include optional fields like customer name, email, or avatar:

```go
leadData := sdk.LeadData{
	ClickId:            "click-123",
	EventName:          "lead_captured",
	CustomerExternalId: "customer-456",
	CustomerName:       sdk.StringPtr("John Doe"),
	CustomerEmail:      sdk.StringPtr("john@example.com"),
	CustomerAvatar:     sdk.StringPtr("https://example.com/avatar.jpg"),
}

cms.TrackLeadAdvanced(leadData)
```

## Error Handling

All SDK methods return errors for easy error handling:

```go
if err != nil {
	log.Printf("Failed to track lead: %v", err)
}
```

## Context & Timeouts

For custom context or timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

cms := sdk.NewWithContext(ctx, "your-token")
```

## Custom URL

Need to use a different API endpoint (staging, localhost, custom domain)?

```go
cms := sdk.New("your-token")

// Set custom base URL
cms.SetBaseURL("https://staging-api.example.com")
// or for local development:
cms.SetBaseURL("http://localhost:8080")

// Now all requests go to the custom URL
cms.TrackLead(...)
cms.TrackSale(...)
```

**Default URL**: `https://www.cutmeshort.com`

That's all you need to get started!
