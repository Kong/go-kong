package kong

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// APIError is used for Kong Admin API errors.
type APIError struct {
	httpCode int
	message  string
	raw      []byte
	details  any
}

func NewAPIError(code int, msg string) *APIError {
	return &APIError{
		httpCode: code,
		message:  msg,
	}
}

func NewAPIErrorWithRaw(code int, msg string, raw []byte) *APIError {
	return &APIError{
		httpCode: code,
		message:  msg,
		raw:      raw,
	}
}

func (e *APIError) Error() string {
	if e.details == nil {
		return fmt.Sprintf("HTTP status %d (message: %q)", e.httpCode, e.message)
	}
	return fmt.Sprintf("HTTP status %d (message: %q; details: %v)", e.httpCode, e.message, e.details)
}

// Code returns the HTTP status code for the error.
func (e *APIError) Code() int {
	return e.httpCode
}

// Raw returns the raw HTTP error response body.
func (e *APIError) Raw() []byte {
	return e.raw
}

// Details returns optional details that might be relevant for proper
// handling of the APIError on the caller side.
func (e *APIError) Details() any {
	return e.details
}

// SetDetails allows setting optional details that might be relevant
// for proper handling of the APIError on the caller side.
func (e *APIError) SetDetails(details any) {
	e.details = details
}

// IsNotFoundErr returns true if the error or it's cause is
// a 404 response from Kong.
func IsNotFoundErr(e error) bool {
	var apiErr *APIError
	if errors.As(e, &apiErr) {
		return apiErr.httpCode == http.StatusNotFound
	}
	return false
}

// IsForbiddenErr returns true if the error or its cause is
// a 403 response from Kong.
func IsForbiddenErr(e error) bool {
	var apiErr *APIError
	if errors.As(e, &apiErr) {
		return apiErr.httpCode == http.StatusForbidden
	}
	return false
}

// ErrTooManyRequestsDetails is expected to be available under APIError.Details()
// when the API returns status code 429 (Too many requests) and a `Retry-After` header
// is set.
type ErrTooManyRequestsDetails struct {
	RetryAfter time.Duration
}
