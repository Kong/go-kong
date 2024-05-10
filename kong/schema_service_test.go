package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchemaService(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	}
	for _, entity := range entities {
		schema, err := client.Schemas.Get(defaultCtx, entity)
		_, ok := schema["fields"]
		assert.True(ok)
		assert.NoError(err)
	}
}

func TestSchemaServiceValidate(t *testing.T) {
	client, err := NewTestClient(nil, nil)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	testCases := []struct {
		name       string
		entityType string
		entity     any
		valid      bool
	}{
		{
			name:       "valid service should pass the validation",
			entityType: "services",
			entity: &Service{
				Name: String("test.service"),
				Host: String("foo.com"),
			},
			valid: true,
		},
		{
			name:       "invalid service (invalid protocol) should fail the validation",
			entityType: "services",
			entity: &Service{
				Name:     String("test.service"),
				Host:     String("foo.com"),
				Protocol: String("not-exist-protocol"),
			},
			valid: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			valid, msg, err := client.Schemas.Validate(defaultCtx, tc.entityType, tc.entity)
			assert.NoError(t, err)
			if tc.valid {
				assert.True(t, valid)
				assert.Empty(t, msg)
			} else {
				assert.False(t, valid)
				assert.NotEmpty(t, msg)
			}
		})
	}
}
