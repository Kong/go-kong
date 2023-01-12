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
}

func NewAPIError(code int, msg string) *APIError {
	return &APIError{
		httpCode: code,
		message:  msg,
	}
}

func (e *APIError) Error() string {
	return fmt.Sprintf("HTTP status %d (message: %q)", e.httpCode, e.message)
}

// Code returns the HTTP status code for the error.
func (e *APIError) Code() int {
	return e.httpCode
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
