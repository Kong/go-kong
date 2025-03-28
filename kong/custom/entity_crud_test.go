package custom

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestRender(t *testing.T) {
	assert := assert.New(t)

	entity := NewEntityObject("key-auth")
	entity.AddRelation("consumer_id", "bob")
	result, err := render("/consumers/${consumer_id}/key-auths", entity)
	require.NoError(t, err)
	assert.Equal("/consumers/bob/key-auths", result)

	result, err = render("/consumers/${random_id}/key-auths", entity)
	require.Error(t, err)
	assert.Empty(result)
}

func TestEntityCRUDDefinition(t *testing.T) {
	assert := assert.New(t)

	var typ Type = "foo"
	e := EntityCRUDDefinition{
		Name:       typ,
		CRUDPath:   "/consumers/${consumer_id}/foo",
		PrimaryKey: "id",
	}

	var object Object = map[string]interface{}{
		"id":   "unique-id",
		"key1": "value1",
		"key2": "value2",
	}
	entity := NewEntityObject(typ)
	entity.SetObject(object)
	entity.AddRelation("consumer_id", "gopher")

	assert.Equal(typ, e.Type())
	url, err := e.GetEndpoint(entity)
	require.NoError(t, err)
	assert.Equal("/consumers/gopher/foo/unique-id", url)

	url, err = e.PatchEndpoint(entity)
	require.NoError(t, err)
	assert.Equal("/consumers/gopher/foo/unique-id", url)

	url, err = e.DeleteEndpoint(entity)
	require.NoError(t, err)
	assert.Equal("/consumers/gopher/foo/unique-id", url)

	url, err = e.PostEndpoint(entity)
	require.NoError(t, err)
	assert.Equal("/consumers/gopher/foo", url)

	url, err = e.ListEndpoint(entity)
	require.NoError(t, err)
	assert.Equal("/consumers/gopher/foo", url)

	entity = NewEntityObject(typ)
	entity.SetObject(object)
	url, err = e.GetEndpoint(entity)
	require.Error(t, err)
	assert.Empty(url)

	object = map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	}
	entity = NewEntityObject(typ)
	entity.AddRelation("consumer_id", "gopher")
	entity.SetObject(object)
	url, err = e.GetEndpoint(entity)
	require.Error(t, err)
	assert.Empty(url)

	object = map[string]interface{}{
		"id":   -1,
		"key1": "value1",
		"key2": "value2",
	}
	entity = NewEntityObject(typ)
	entity.SetObject(object)
	entity.AddRelation("consumer_id", "gopher")
	url, err = e.GetEndpoint(entity)
	require.Error(t, err)
	assert.Empty(url)
}

func TestEntityCRUDUnmarshal(t *testing.T) {
	assert := assert.New(t)

	t.Run("unmarshal JSON into EntityCRUDDefinition", func(t *testing.T) {
		bytes := []byte(`{
			"name": "name",
			"crud": "crud-path",
			"primary_key": "primary-key"
		}`)
		var def EntityCRUDDefinition
		err := json.Unmarshal(bytes, &def)
		require.NoError(t, err)
		assert.Equal(Type("name"), def.Name)
		assert.Equal("crud-path", def.CRUDPath)
		assert.Equal("primary-key", def.PrimaryKey)
	})

	t.Run("unmarshal YAML into EntityCRUDDefinition", func(t *testing.T) {
		var def EntityCRUDDefinition
		bytes := []byte(`
name: "name"
crud: "crud-path"
primary_key: "primary-key"`)
		err := yaml.Unmarshal(bytes, &def)
		require.NoError(t, err)
		assert.Equal(Type("name"), def.Name)
		assert.Equal("crud-path", def.CRUDPath)
		assert.Equal("primary-key", def.PrimaryKey)
	})
}
