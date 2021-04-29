package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagExists(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	exists, err := client.Tags.Exists(defaultCtx)
	assert.Nil(err)
	assert.True(*exists)
}
