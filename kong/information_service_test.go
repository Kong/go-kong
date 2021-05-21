package kong

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	err := Convert(information, &actual)
	assert.Nil(err)
	assert.True(reflect.DeepEqual(expected, &actual))
}
