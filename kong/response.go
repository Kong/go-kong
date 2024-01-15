package kong

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Response is a Kong Admin API response.
// It contains the response headers, status and status code.
type Response struct {
	Header     http.Header
	Status     string
	StatusCode int
}

func newResponse(res *http.Response) *Response {
	return &Response{
		Header:     res.Header,
		Status:     res.Status,
		StatusCode: res.StatusCode,
	}
}

func messageFromBody(b []byte) string {
	s := struct {
		Message string
	}{}

	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Sprintf("<failed to parse response body: %v>", err)
	}

	return s.Message
}

// detailsFromBodyDetailsField extract details from body if the response body contains a "details" field.
// Used for extracting details from response from Konnect APIs when error happens.
func detailsFromBodyDetailsField(b []byte) any {
	s := struct {
		Details any `json:"details"`
	}{}
	if err := json.Unmarshal(b, &s); err != nil {
		return nil
	}
	return s.Details
}

func hasError(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 399 {
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read error body: %w", err)
	}

	apiErr := NewAPIError(res.StatusCode, messageFromBody(body))
	if details, ok := extractErrDetails(res, body); ok {
		apiErr.SetDetails(details)
	}

	return apiErr
}

func extractErrDetails(res *http.Response, body []byte) (any, bool) {
	// firstly deal with certain status code.
	switch res.StatusCode {
	case http.StatusTooManyRequests:
		return extractErrTooManyRequestsDetails(res)
	}
	// Then extract details from "details" field in the response body.
	if detailsFromRespBody := detailsFromBodyDetailsField(body); detailsFromRespBody != nil {
		return detailsFromRespBody, true
	}

	return nil, false
}

func extractErrTooManyRequestsDetails(res *http.Response) (ErrTooManyRequestsDetails, bool) {
	const (
		base    = 10
		bitSize = 64
	)
	if retryAfter := res.Header.Get("Retry-After"); retryAfter != "" {
		if sleep, err := strconv.ParseInt(retryAfter, base, bitSize); err == nil {
			return ErrTooManyRequestsDetails{
				RetryAfter: time.Second * time.Duration(sleep),
			}, true
		}
	}

	return ErrTooManyRequestsDetails{}, false
}
