package kong

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Response is a Kong Admin API response. It wraps http.Response.
type Response struct {
	*http.Response
	// other Kong specific fields
}

func newResponse(res *http.Response) *Response {
	return &Response{Response: res}
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

func hasError(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 399 {
		return nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read error body: %w", err)
	}

	apiErr := NewAPIError(res.StatusCode, messageFromBody(body))
	if details, ok := extractErrDetails(res); ok {
		apiErr.SetDetails(details)
	}

	return apiErr
}

func extractErrDetails(res *http.Response) (any, bool) {
	switch res.StatusCode {
	case http.StatusTooManyRequests:
		return extractErrTooManyRequestsDetails(res)
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
