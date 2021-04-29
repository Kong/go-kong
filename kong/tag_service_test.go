package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagExists(T *testing.T) {
	runWhenKong(T, ">=1.1.0")
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	exists, err := client.Tags.Exists(defaultCtx)
	assert.Nil(err)
	assert.True(exists)
}

func TestTagDoesNotExists(T *testing.T) {
	runWhenKong(T, "<1.1.0")
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	exists, err := client.Tags.Exists(defaultCtx)
	assert.Nil(err)
	assert.False(exists)
}
