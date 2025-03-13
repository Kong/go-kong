package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagExists(T *testing.T) {
	RunWhenKong(T, ">=1.1.0")
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	exists, err := client.Tags.Exists(defaultCtx)
	require.NoError(T, err)
	assert.True(exists)
}

func TestTagDoesNotExists(T *testing.T) {
	RunWhenKong(T, "<1.1.0")
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	exists, err := client.Tags.Exists(defaultCtx)
	require.NoError(T, err)
	assert.False(exists)
}
