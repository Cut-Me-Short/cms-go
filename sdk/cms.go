/*
CutMeShort CMS SDK - Simplified Interface

This provides a simple, user-friendly interface for tracking leads and sales.
Users only need to initialize with a token and call TrackLead/TrackSale.
*/

package sdk

import (
	"context"
	"net/http"
	"time"
)

// CMS is the main SDK client for simplified usage
type CMS struct {
	token      string
	client     *APIClient
	ctx        context.Context
	httpClient *http.Client
}

// LeadData holds lead tracking data
type LeadData struct {
	ClickId            string
	EventName          string
	CustomerExternalId string
	CustomerName       *string
	CustomerEmail      *string
	CustomerAvatar     *string
	Timestamp          *time.Time
	Mode               *string // "deferred" for deferred attribution
}

// SaleData holds sale tracking data
type SaleData struct {
	ClickId            string
	EventName          string
	CustomerExternalId string
	InvoiceId          string
	Amount             float64 // in cents
	Currency           string  // 3-letter code
	CustomerName       *string
	CustomerEmail      *string
	CustomerAvatar     *string
	Timestamp          *time.Time
}

// New creates and initializes a new CMS SDK client with just a token
func New(token string) *CMS {
	cfg := NewConfiguration()
	cfg.DefaultHeader["Authorization"] = "Bearer " + token
	
	return &CMS{
		token:      token,
		client:     NewAPIClient(cfg),
		ctx:        context.Background(),
		httpClient: &http.Client{},
	}
}

// NewWithContext creates a new CMS SDK client with a custom context
func NewWithContext(ctx context.Context, token string) *CMS {
	cfg := NewConfiguration()
	cfg.DefaultHeader["Authorization"] = "Bearer " + token
	
	return &CMS{
		token:      token,
		client:     NewAPIClient(cfg),
		ctx:        ctx,
		httpClient: &http.Client{},
	}
}

// TrackLead tracks a lead event - simple one-liner
func (c *CMS) TrackLead(clickId, eventName, customerExternalId string) (*TrackResponse, error) {
	lead := &LeadPayload{
		ClickId:            &clickId,
		EventName:          eventName,
		CustomerExternalId: customerExternalId,
	}
	
	resp, _, err := c.client.TrackingAPI.TrackLead(c.ctx).LeadPayload(*lead).Execute()
	return resp, err
}

// TrackLeadAdvanced tracks a lead event with all optional fields
func (c *CMS) TrackLeadAdvanced(data LeadData) (*TrackResponse, error) {
	lead := &LeadPayload{
		ClickId:            &data.ClickId,
		EventName:          data.EventName,
		CustomerExternalId: data.CustomerExternalId,
		CustomerName:       data.CustomerName,
		CustomerEmail:      data.CustomerEmail,
		CustomerAvatar:     data.CustomerAvatar,
		Timestamp:          data.Timestamp,
		Mode:               data.Mode,
	}
	
	resp, _, err := c.client.TrackingAPI.TrackLead(c.ctx).LeadPayload(*lead).Execute()
	return resp, err
}

// TrackSale tracks a sale event - simple one-liner
func (c *CMS) TrackSale(clickId, eventName, customerExternalId, invoiceId, currency string, amount float64) (*TrackResponse, error) {
	sale := &SalePayload{
		ClickId:            clickId,
		EventName:          eventName,
		CustomerExternalId: customerExternalId,
		InvoiceId:          invoiceId,
		Amount:             amount,
		Currency:           currency,
	}
	
	resp, _, err := c.client.TrackingAPI.TrackSale(c.ctx).SalePayload(*sale).Execute()
	return resp, err
}

// TrackSaleAdvanced tracks a sale event with all optional fields
func (c *CMS) TrackSaleAdvanced(data SaleData) (*TrackResponse, error) {
	sale := &SalePayload{
		ClickId:            data.ClickId,
		EventName:          data.EventName,
		CustomerExternalId: data.CustomerExternalId,
		InvoiceId:          data.InvoiceId,
		Amount:             data.Amount,
		Currency:           data.Currency,
		CustomerName:       data.CustomerName,
		CustomerEmail:      data.CustomerEmail,
		CustomerAvatar:     data.CustomerAvatar,
		Timestamp:          data.Timestamp,
	}
	
	resp, _, err := c.client.TrackingAPI.TrackSale(c.ctx).SalePayload(*sale).Execute()
	return resp, err
}

// SetTimeout sets the HTTP client timeout
func (c *CMS) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// SetContext sets the context for API calls
func (c *CMS) SetContext(ctx context.Context) {
	c.ctx = ctx
}

// SetBaseURL sets a custom base URL for API requests
func (c *CMS) SetBaseURL(url string) {
	if len(c.client.cfg.Servers) > 0 {
		c.client.cfg.Servers[0].URL = url
	} else {
		c.client.cfg.Servers = ServerConfigurations{
			{
				URL:         url,
				Description: "Custom API Server",
			},
		}
	}
}

// StringPtr is a helper function to create a string pointer
func StringPtr(s string) *string {
	return &s
}

// TimePtr is a helper function to create a time pointer
func TimePtr(t time.Time) *time.Time {
	return &t
}

// Float64Ptr is a helper function to create a float64 pointer
func Float64Ptr(f float64) *float64 {
	return &f
}
