package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagExists(T *testing.T) {
	RunWhenKong(T, ">=1.1.0")
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	exists, err := client.Tags.Exists(defaultCtx)
	assert.NoError(err)
	assert.True(exists)
}

func TestTagDoesNotExists(T *testing.T) {
	RunWhenKong(T, "<1.1.0")
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	exists, err := client.Tags.Exists(defaultCtx)
	assert.NoError(err)
	assert.False(exists)
}
