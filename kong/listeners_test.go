package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListeners(t *testing.T) {
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)

	t.Log("pulling the listener configurations from the kong root")
	proxyListeners, _, err := client.Listeners(defaultCtx)
	require.NoError(t, err)
	assert.NotEmpty(t, proxyListeners)

	t.Log("verifying that the standard http listener was found")
	foundHTTPListener := false
	for _, listener := range proxyListeners {
		if listener.SSL == false {
			foundHTTPListener = true
		}
	}
	assert.True(t, foundHTTPListener)
}
