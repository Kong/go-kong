package kong

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError is used for Kong Admin API errors.
type APIError struct {
	httpCode int
	message  string
	raw      []byte
}

func NewAPIError(code int, msg string, raw []byte) *APIError {
	return &APIError{
		httpCode: code,
		message:  msg,
		raw:      raw,
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
	return fmt.Sprintf("HTTP status %d (message: %q)", e.httpCode, e.message)
}

// Code returns the HTTP status code for the error.
func (e *APIError) Code() int {
	return e.httpCode
}

// Raw returns the raw HTTP error response body.
func (e *APIError) Raw() []byte {
	return e.raw
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
