package kong

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// NewRequestRaw creates a request based on the inputs.
func (c *Client) NewRequestRaw(method, baseURL string, endpoint string, qs interface{},
	body interface{}) (*http.Request, error) {
	if endpoint == "" {
		return nil, fmt.Errorf("endpoint can't be nil")
	}
	// body to be sent in JSON
	var buf []byte
	if body != nil {
		var err error
		buf, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	// Create a new request
	fmt.Printf("<< " + baseURL + endpoint + ">>>")
	req, err := http.NewRequest(method, baseURL+endpoint, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	// add body if needed
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	// add query string if any
	if qs != nil {
		values, err := query.Values(qs)
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = values.Encode()
	}
	return req, nil
}

// NewRequest creates a request based on the inputs.
// endpoint should be relative to the baseURL specified during
// client creation.
// body is always marshaled into JSON.
func (c *Client) NewRequest(method, endpoint string, qs interface{},
	body interface{}) (*http.Request, error) {
	return c.NewRequestRaw(method, c.workspacedBaseURL(c.Workspace()), endpoint, qs, body)
}
