package kong

import (
	"bytes"
	"net/http"
	"strings"
)

// String returns pointer to s.
func String(s string) *string {
	return &s
}

// Bool returns a pointer to b.
func Bool(b bool) *bool {
	return &b
}

// Int returns a pointer to i.
func Int(i int) *int {
	return &i
}

func isEmptyString(s *string) bool {
	return s == nil || strings.TrimSpace(*s) == ""
}

//Deduplicates a list a strings
func deduplicate(stringSlice []string) []string {
	existing := map[string]struct{}{}
	result := []string{}

	for _, s := range stringSlice {
		if _, exists := existing[s]; !exists {
			existing[s] = struct{}{}
			result = append(result, s)
		}
	}

	return result
}

// StringSlice converts a slice of string to a
// slice of *string
func StringSlice(elements ...string) []*string {
	var res []*string
	for _, element := range elements {
		e := element
		res = append(res, &e)
	}
	return res
}
func stringArrayToString(arr []*string) string {
	if arr == nil {
		return "nil"
	}

	var buf bytes.Buffer
	buf.WriteString("[ ")
	l := len(arr)
	for i, el := range arr {
		buf.WriteString(*el)
		if i != l-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(" ]")
	return buf.String()
}

// headerRoundTripper injects Headers into requests
// made via RT.
type headerRoundTripper struct {
	headers http.Header
	rt      http.RoundTripper
}

// RoundTrip satisfies the RoundTripper interface.
func (t headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = requestWithHeaders(req, t.headers)
	return t.rt.RoundTrip(req)
}

func requestWithHeaders(req *http.Request, headers http.Header) *http.Request {
	if req == nil {
		return nil
	}
	// fast-path
	if len(headers) == 0 {
		return req
	}
	newRequest := new(http.Request)
	*newRequest = *req
	newRequest.Header = req.Header.Clone()
	for k, values := range headers {
		for _, v := range values {
			newRequest.Header.Add(k, v)
		}
	}
	return newRequest
}

// RoundTripperWithHTTPHeaders returns a client which injects headers
// before sending any request.
func HTTPClientWithHeaders(client *http.Client,
	headers http.Header) *http.Client {
	var res *http.Client
	if client == nil {
		defaultTransport := http.DefaultTransport.(*http.Transport)
		res = &http.Client{}
		res.Transport = defaultTransport
	} else {
		res = client
	}
	res.Transport = headerRoundTripper{
		headers: headers,
		rt:      client.Transport,
	}
	return res
}
