package kong

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInfoService(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	info, err := client.Info.Get(defaultCtx)
	require.NoError(err)
	require.NotNil(info)
	require.NotNil(info.Version)
	require.NotNil(info.Configuration)
	require.NotNil(info.Configuration.Database)
}

func TestConvert(T *testing.T) {
	assert := assert.New(T)
	information := map[string]interface{}{
		"version": "2.3.3.2-enterprise-edition",
		"configuration": map[string]interface{}{
			"portal":   true,
			"rbac":     "on",
			"database": "postgres",
		},
	}
	expected := &Info{
		Version: "2.3.3.2-enterprise-edition",
		Configuration: &RuntimeConfiguration{
			Portal:   true,
			RBAC:     "on",
			Database: "postgres",
		},
	}
	var actual Info
	err := convert(information, &actual)
	require.NoError(T, err)
	assert.True(reflect.DeepEqual(expected, &actual))
	assert.False(actual.Configuration.IsInMemory())
	assert.True(actual.Configuration.IsRBACEnabled())
}
