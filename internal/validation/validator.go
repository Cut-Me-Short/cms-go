/*
CutMeShort Go SDK - Input Validation

Package validation provides comprehensive input validation for all SDK payloads
following industry standards and security best practices.
*/

package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/cutmeshort/sdk-go/internal/errors"
)

const (
	// MaxStringLength is the maximum allowed string length for text fields
	MaxStringLength = 500

	// MaxEmailLength is the maximum length for email addresses
	MaxEmailLength = 254

	// MinAmount is the minimum transaction amount in cents
	MinAmount = 0.01

	// MaxAmount is the maximum transaction amount in cents
	MaxAmount = 999999999.99

	// AmountDecimals is the number of decimal places allowed for amounts
	AmountDecimals = 2

	// CurrencyCodeLength is the standard length for ISO 4217 currency codes
	CurrencyCodeLength = 3

	// UUIDPattern is a regex for validating UUIDs
	UUIDPattern = `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`

	// EventNamePattern allows alphanumeric, underscore, and hyphen
	EventNamePattern = `^[a-zA-Z0-9_-]+$`

	// MaxEventNameLength is the maximum length for event names
	MaxEventNameLength = 100

	// MaxIDLength is the maximum length for ID fields
	MaxIDLength = 255

	// MaxNameLength is the maximum length for customer name
	MaxNameLength = 200

	// MaxAvatarURLLength is the maximum length for avatar URL
	MaxAvatarURLLength = 2048
)

var (
	uuidRegex      = regexp.MustCompile(UUIDPattern)
	eventNameRegex = regexp.MustCompile(EventNamePattern)

	// ISO 4217 Currency Codes (common ones for validation)
	validCurrencies = map[string]bool{
		"USD": true, "EUR": true, "GBP": true, "JPY": true, "CHF": true,
		"CAD": true, "AUD": true, "NZD": true, "INR": true, "SGD": true,
		"HKD": true, "SEK": true, "NOK": true, "DKK": true, "AED": true,
		"SAR": true, "QAR": true, "BRL": true, "MXN": true, "ZAR": true,
		"TRY": true, "RUB": true, "KRW": true, "CNY": true, "TWD": true,
		"THB": true, "MYR": true, "PHY": true, "IDR": true, "VND": true,
		"PKR": true, "BDT": true, "LKR": true, "NPR": true, "BGN": true,
		"CZK": true, "HUF": true, "PLN": true, "RON": true, "HRK": true,
		"BGN": true, "RSD": true, "TND": true, "EGP": true, "ARS": true,
		"CLP": true, "COP": true, "PEN": true, "UYU": true,
	}
)

// ValidateClickID validates the click ID format
func ValidateClickID(clickID string) error {
	clickID = strings.TrimSpace(clickID)

	if clickID == "" {
		return errors.NewValidationError("clickId", "must not be empty")
	}

	if len(clickID) > MaxIDLength {
		return errors.NewValidationError("clickId", fmt.Sprintf("must not exceed %d characters", MaxIDLength))
	}

	// Check if it's a UUID format (most common)
	if !isValidUUID(clickID) && !isValidAlphanumeric(clickID) {
		return errors.NewValidationError("clickId", "must be a valid UUID or alphanumeric string")
	}

	return nil
}

// ValidateEventName validates the event name
func ValidateEventName(eventName string) error {
	eventName = strings.TrimSpace(eventName)

	if eventName == "" {
		return errors.NewValidationError("eventName", "must not be empty")
	}

	if len(eventName) > MaxEventNameLength {
		return errors.NewValidationError("eventName", fmt.Sprintf("must not exceed %d characters", MaxEventNameLength))
	}

	if !eventNameRegex.MatchString(eventName) {
		return errors.NewValidationError("eventName", "must contain only alphanumeric characters, underscores, and hyphens")
	}

	return nil
}

// ValidateCustomerExternalID validates the customer external ID
func ValidateCustomerExternalID(id string) error {
	id = strings.TrimSpace(id)

	if id == "" {
		return errors.NewValidationError("customerExternalId", "must not be empty")
	}

	if len(id) > MaxIDLength {
		return errors.NewValidationError("customerExternalId", fmt.Sprintf("must not exceed %d characters", MaxIDLength))
	}

	if !isValidAlphanumeric(id) {
		return errors.NewValidationError("customerExternalId", "must be alphanumeric (allowing underscores and hyphens)")
	}

	return nil
}

// ValidateEmail validates an email address format
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return errors.NewValidationError("email", "must not be empty")
	}

	if len(email) > MaxEmailLength {
		return errors.NewValidationError("email", fmt.Sprintf("must not exceed %d characters", MaxEmailLength))
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.NewValidationError("email", "must be a valid email address")
	}

	return nil
}

// ValidateCustomerName validates the customer name
func ValidateCustomerName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.NewValidationError("customerName", "must not be empty")
	}

	if len(name) > MaxNameLength {
		return errors.NewValidationError("customerName", fmt.Sprintf("must not exceed %d characters", MaxNameLength))
	}

	return nil
}

// ValidateAvatarURL validates the avatar URL format
func ValidateAvatarURL(url string) error {
	url = strings.TrimSpace(url)

	if url == "" {
		return errors.NewValidationError("customerAvatar", "must not be empty")
	}

	if len(url) > MaxAvatarURLLength {
		return errors.NewValidationError("customerAvatar", fmt.Sprintf("must not exceed %d characters", MaxAvatarURLLength))
	}

	// Basic URL validation
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "data:") {
		return errors.NewValidationError("customerAvatar", "must be a valid HTTP(S) URL or data URI")
	}

	return nil
}

// ValidateTimestamp validates the timestamp
func ValidateTimestamp(ts time.Time) error {
	now := time.Now()
	maxAge := 24 * time.Hour
	maxFuture := 1 * time.Minute

	// Check if timestamp is not in ancient past (more than 24 hours old)
	if now.Sub(ts) > maxAge {
		return errors.NewValidationError("timestamp", "must not be more than 24 hours in the past")
	}

	// Check if timestamp is not in future (more than 1 minute ahead)
	if ts.Sub(now) > maxFuture {
		return errors.NewValidationError("timestamp", "must not be more than 1 minute in the future")
	}

	return nil
}

// ValidateInvoiceID validates the invoice ID format
func ValidateInvoiceID(invoiceID string) error {
	invoiceID = strings.TrimSpace(invoiceID)

	if invoiceID == "" {
		return errors.NewValidationError("invoiceId", "must not be empty")
	}

	if len(invoiceID) > MaxIDLength {
		return errors.NewValidationError("invoiceId", fmt.Sprintf("must not exceed %d characters", MaxIDLength))
	}

	if !isValidAlphanumeric(invoiceID) {
		return errors.NewValidationError("invoiceId", "must be alphanumeric (allowing underscores and hyphens)")
	}

	return nil
}

// ValidateAmount validates the transaction amount
func ValidateAmount(amount float64) error {
	if amount <= 0 {
		return errors.NewValidationError("amount", fmt.Sprintf("must be greater than 0 (got %.2f)", amount))
	}

	if amount > MaxAmount {
		return errors.NewValidationError("amount", fmt.Sprintf("must not exceed %.2f", MaxAmount))
	}

	// Check decimal places (max 2 for cents)
	if hasMoreThanDecimalPlaces(amount, AmountDecimals) {
		return errors.NewValidationError("amount", fmt.Sprintf("must have at most %d decimal places", AmountDecimals))
	}

	return nil
}

// ValidateCurrency validates the currency code
func ValidateCurrency(currency string) error {
	currency = strings.TrimSpace(currency)
	currency = strings.ToUpper(currency)

	if currency == "" {
		return errors.NewValidationError("currency", "must not be empty")
	}

	if len(currency) != CurrencyCodeLength {
		return errors.NewValidationError("currency", fmt.Sprintf("must be exactly %d characters (ISO 4217 code)", CurrencyCodeLength))
	}

	if !validCurrencies[currency] {
		return errors.NewValidationError("currency", fmt.Sprintf("'%s' is not a valid ISO 4217 currency code", currency))
	}

	return nil
}

// Helper functions

// isValidUUID checks if a string is a valid UUID format
func isValidUUID(s string) bool {
	return uuidRegex.MatchString(strings.ToLower(s))
}

// isValidAlphanumeric checks if a string contains only alphanumeric characters, underscores, hyphens
func isValidAlphanumeric(s string) bool {
	for _, c := range s {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' || c == '-') {
			return false
		}
	}
	return true
}

// hasMoreThanDecimalPlaces checks if a float has more than the specified decimal places
func hasMoreThanDecimalPlaces(f float64, maxPlaces int) bool {
	multiplier := 1.0
	for i := 0; i < maxPlaces; i++ {
		multiplier *= 10
	}
	return f*multiplier != float64(int64(f*multiplier))
}
