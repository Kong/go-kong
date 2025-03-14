package kong

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserAgentHeader(t *testing.T) {
	collectedUserAgents := make([]string, 0, 3)
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {
		collectedUserAgents = append(collectedUserAgents, req.Header.Get("User-Agent"))
	}))
	t.Cleanup(server.Close)

	const userAgent = "go-kong/v1.2.3"

	client, err := NewClient(&server.URL, nil)
	require.NoError(t, err)
	client.UserAgent = userAgent

	testingID := "test-id"
	_, _ = client.Licenses.Get(context.Background(), &testingID)
	require.Equal(t, userAgent, collectedUserAgents[0])

	_, _, _ = client.Plugins.List(context.Background(), nil) //nolint:dogsled
	require.Equal(t, userAgent, collectedUserAgents[1])

	_, _ = client.Vaults.Create(context.Background(), &Vault{})
	require.Equal(t, userAgent, collectedUserAgents[2])
}
