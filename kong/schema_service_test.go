package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	entities := []string{"services", "routes", "targets", "upstreams", "plugins"}
	for _, entity := range entities {
		schema, err := client.Schemas.Get(defaultCtx, entity)
		_, ok := schema["fields"]
		assert.True(ok)
		assert.Nil(err)
	}
}
