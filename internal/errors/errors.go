/*
CutMeShort Go SDK - Error Types and Handling

Package errors defines custom error types for SDK operations with proper
error classification and context.
*/

package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorType classifies the type of error that occurred
type ErrorType string

const (
	ErrorTypeValidation    ErrorType = "validation"
	ErrorTypeAuthentication ErrorType = "authentication"
	ErrorTypeAuthorization  ErrorType = "authorization"
	ErrorTypeNotFound       ErrorType = "not_found"
	ErrorTypeRateLimit      ErrorType = "rate_limit"
	ErrorTypeConflict       ErrorType = "conflict"
	ErrorTypeServer         ErrorType = "server_error"
	ErrorTypeNetwork        ErrorType = "network_error"
	ErrorTypeTimeout        ErrorType = "timeout"
	ErrorTypeUnknown        ErrorType = "unknown"
)

// SDKError is the base error type for all SDK errors
type SDKError struct {
	Type       ErrorType
	StatusCode int
	Message    string
	Details    map[string]interface{}
	Err        error
	Operation  string
	RequestID  string
}

// Error implements the error interface
func (e *SDKError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %s (underlying: %v)", e.Type, e.Message, e.Operation, e.Err)
	}
	return fmt.Sprintf("[%s] %s: %s", e.Type, e.Message, e.Operation)
}

// Unwrap provides support for error chaining
func (e *SDKError) Unwrap() error {
	return e.Err
}

// NewSDKError creates a new SDK error
func NewSDKError(errorType ErrorType, statusCode int, message string, operation string, err error) *SDKError {
	return &SDKError{
		Type:       errorType,
		StatusCode: statusCode,
		Message:    message,
		Operation:  operation,
		Err:        err,
		Details:    make(map[string]interface{}),
	}
}

// AddDetail adds a detail to the error for debugging
func (e *SDKError) AddDetail(key string, value interface{}) *SDKError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithRequestID adds a request ID for tracing
func (e *SDKError) WithRequestID(id string) *SDKError {
	e.RequestID = id
	return e
}

// ValidationError is returned when input validation fails
type ValidationError struct {
	*SDKError
	Field  string
	Reason string
}

// NewValidationError creates a validation error
func NewValidationError(field string, reason string) *ValidationError {
	return &ValidationError{
		SDKError: NewSDKError(
			ErrorTypeValidation,
			http.StatusBadRequest,
			fmt.Sprintf("validation failed for field '%s': %s", field, reason),
			"validation",
			nil,
		),
		Field:  field,
		Reason: reason,
	}
}

// RateLimitError indicates the API rate limit was exceeded
type RateLimitError struct {
	*SDKError
	RetryAfter int
}

// NewRateLimitError creates a rate limit error
func NewRateLimitError(retryAfter int) *RateLimitError {
	return &RateLimitError{
		SDKError: NewSDKError(
			ErrorTypeRateLimit,
			http.StatusTooManyRequests,
			"rate limit exceeded",
			"api_call",
			nil,
		),
		RetryAfter: retryAfter,
	}
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) (*ValidationError, bool) {
	var ve *ValidationError
	return ve, errors.As(err, &ve)
}

// IsRateLimitError checks if an error is a rate limit error
func IsRateLimitError(err error) (*RateLimitError, bool) {
	var rle *RateLimitError
	return rle, errors.As(err, &rle)
}

// IsSDKError checks if an error is an SDK error
func IsSDKError(err error) (*SDKError, bool) {
	var se *SDKError
	return se, errors.As(err, &se)
}

// IsNetworkError checks if an error was caused by network issues
func IsNetworkError(err error) bool {
	if se, ok := IsSDKError(err); ok {
		return se.Type == ErrorTypeNetwork || se.Type == ErrorTypeTimeout
	}
	return false
}

// IsRetryable determines if an operation should be retried
func IsRetryable(err error) bool {
	if se, ok := IsSDKError(err); ok {
		switch se.Type {
		case ErrorTypeRateLimit, ErrorTypeNetwork, ErrorTypeTimeout, ErrorTypeServer:
			// These are transient and should be retried
			return true
		case ErrorTypeAuthentication, ErrorTypeAuthorization, ErrorTypeValidation:
			// These are permanent and shouldn't be retried
			return false
		default:
			return false
		}
	}
	return false
}
