package kong

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type headersTestTransport func(*http.Request) (*http.Response, error)

func (f headersTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestHTTPClientWithHeadersCustomDefaultTransport(t *testing.T) {
	original := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = original })

	for _, tc := range []struct {
		name   string
		client *http.Client
	}{
		{name: "nil client"},
		{name: "nil transport", client: &http.Client{}},
		{name: "explicit transport", client: &http.Client{Transport: original}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			transport := headersTestTransport(func(req *http.Request) (*http.Response, error) {
				calls++
				require.Equal(t, "test-token", req.Header.Get("Kong-Admin-Token"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     make(http.Header),
					Body:       io.NopCloser(strings.NewReader("{}")),
					Request:    req,
				}, nil
			})
			http.DefaultTransport = transport
			if tc.client != nil && tc.client.Transport != nil {
				tc.client.Transport = transport
			}
			client := HTTPClientWithHeaders(tc.client, http.Header{
				"Kong-Admin-Token": {"test-token"},
			})
			req, err := http.NewRequestWithContext(t.Context(), http.MethodGet, "http://kong.invalid/status", nil)
			require.NoError(t, err)
			res, err := client.Do(req)
			require.NoError(t, err)
			defer res.Body.Close()
			require.Equal(t, http.StatusOK, res.StatusCode)
			require.Equal(t, 1, calls)
			require.Empty(t, req.Header.Get("Kong-Admin-Token"))
		})
	}
}
