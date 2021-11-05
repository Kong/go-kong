package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListeners(t *testing.T) {
	client, err := NewTestClient(nil, nil)
	assert.NoError(t, err)

	t.Log("pulling the listener configurations from the kong root")
	proxyListeners, _, err := client.Listeners(defaultCtx)
	assert.NoError(t, err)
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
