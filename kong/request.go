package kong

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/go-querystring/query"
)

// NewRequestRaw creates a request based on the inputs.
func (c *Client) NewRequestRaw(method, baseURL string, endpoint string, qs interface{},
	body interface{},
) (*http.Request, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint can't be nil")
	}
	// body to be sent in JSON
	var r io.Reader
	if body != nil {
		switch v := body.(type) {
		case string:
			r = bytes.NewBufferString(v)
		case []byte:
			r = bytes.NewBuffer(v)
		case io.Reader:
			r = v
		default:
			b, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			r = bytes.NewBuffer(b)
		}
	}

	// Create a new request
	req, err := http.NewRequest(method, baseURL+endpoint, r)
	if err != nil {
		return nil, err
	}

	// add body if needed
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	// add query string if any
	if qs != nil || len(c.QueryParams) > 0 {
		if qs == nil {
			req.URL.RawQuery = c.QueryParams.Encode()
		} else {
			values, err := query.Values(qs)
			if err != nil {
				return nil, err
			}

			// apply global query params
			for param, paramValues := range c.QueryParams {
				for _, paramValue := range paramValues {
					values.Add(param, paramValue)
				}
			}
			req.URL.RawQuery = values.Encode()
		}
	}
	return req, nil
}

// NewRequest creates a request based on the inputs.
// endpoint should be relative to the baseURL specified during
// client creation.
// body is always marshaled into JSON.
func (c *Client) NewRequest(method, endpoint string, qs interface{},
	body interface{},
) (*http.Request, error) {
	return c.NewRequestRaw(method, c.workspacedBaseURL(c.Workspace()), endpoint, qs, body)
}
