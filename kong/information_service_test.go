package kong

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfoService(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	info, err := client.Info.Get(defaultCtx)
	assert.Nil(err)
	assert.NotNil(info)
	assert.NotNil(info.Version)
	assert.NotNil(info.Configuration)
	assert.NotNil(info.Configuration.Database)
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
	assert.NoError(err)
	assert.True(reflect.DeepEqual(expected, &actual))
}
