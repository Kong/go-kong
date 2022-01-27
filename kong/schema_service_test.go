package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaService(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	entities := []string{
		"services",
		"routes",
		"targets",
		"upstreams",
		"plugins",
		"ca_certificates",
		"certificates",
		"consumers",
		"snis",
		"tags",
		// not present in OSS
		// "workspaces",
		// not present with < 2.3
		// "clustering_data_planes",
		// "parameters",
	}
	for _, entity := range entities {
		schema, err := client.Schemas.Get(defaultCtx, entity)
		_, ok := schema["fields"]
		assert.True(ok)
		assert.Nil(err)
	}
}
