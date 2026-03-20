package kong

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginConditionMarshal(t *testing.T) {
	t.Run("plugin with condition is marshaled correctly", func(t *testing.T) {
		plugin := &Plugin{
			Name: String("rate-limiting"),
			Config: Configuration{
				"minute": 10,
			},
			Condition: &PluginCondition{
				Type:       String("predicate"),
				Expression: String("http.path starts_with '/api'"),
			},
		}

		b, err := json.Marshal(plugin)
		require.NoError(t, err)

		assert.JSONEq(t, `{
			"name": "rate-limiting",
			"config": {"minute": 10},
			"condition": {
				"type": "predicate",
				"expression": "http.path starts_with '/api'"
			}
		}`, string(b))
	})

	t.Run("plugin without condition omits condition field", func(t *testing.T) {
		plugin := &Plugin{
			Name: String("rate-limiting"),
			Config: Configuration{
				"minute": 10,
			},
		}

		b, err := json.Marshal(plugin)
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, json.Unmarshal(b, &m))
		_, hasCondition := m["condition"]
		assert.False(t, hasCondition, "condition field should be omitted when nil")
	})

	t.Run("plugin condition with only expression", func(t *testing.T) {
		plugin := &Plugin{
			Name: String("key-auth"),
			Condition: &PluginCondition{
				Expression: String("http.method == 'GET'"),
			},
		}

		b, err := json.Marshal(plugin)
		require.NoError(t, err)

		assert.JSONEq(t, `{
			"name": "key-auth",
			"condition": {
				"expression": "http.method == 'GET'"
			}
		}`, string(b))
	})
}

func TestPluginConditionUnmarshal(t *testing.T) {
	t.Run("plugin JSON with condition is unmarshaled correctly", func(t *testing.T) {
		jsonData := `{
			"id": "abc123",
			"name": "rate-limiting",
			"config": {"minute": 10},
			"condition": {
				"type": "predicate",
				"expression": "http.path starts_with '/api'"
			}
		}`

		var plugin Plugin
		require.NoError(t, json.Unmarshal([]byte(jsonData), &plugin))

		require.NotNil(t, plugin.Condition)
		assert.Equal(t, "predicate", *plugin.Condition.Type)
		assert.Equal(t, "http.path starts_with '/api'", *plugin.Condition.Expression)
	})

	t.Run("plugin JSON without condition results in nil Condition", func(t *testing.T) {
		jsonData := `{
			"id": "abc123",
			"name": "rate-limiting",
			"config": {"minute": 10}
		}`

		var plugin Plugin
		require.NoError(t, json.Unmarshal([]byte(jsonData), &plugin))

		assert.Nil(t, plugin.Condition)
	})

	t.Run("plugin JSON with partial condition fields", func(t *testing.T) {
		jsonData := `{
			"name": "proxy-cache",
			"condition": {
				"expression": "http.method == 'POST'"
			}
		}`

		var plugin Plugin
		require.NoError(t, json.Unmarshal([]byte(jsonData), &plugin))

		require.NotNil(t, plugin.Condition)
		assert.Nil(t, plugin.Condition.Type)
		assert.Equal(t, "http.method == 'POST'", *plugin.Condition.Expression)
	})
}

func TestPluginConditionDeepCopy(t *testing.T) {
	t.Run("deep copy of plugin with condition", func(t *testing.T) {
		original := &Plugin{
			Name: String("rate-limiting"),
			Condition: &PluginCondition{
				Type:       String("predicate"),
				Expression: String("http.path starts_with '/api'"),
			},
		}

		copied := original.DeepCopy()
		require.NotNil(t, copied)
		require.NotNil(t, copied.Condition)

		assert.Equal(t, *original.Condition.Type, *copied.Condition.Type)
		assert.Equal(t, *original.Condition.Expression, *copied.Condition.Expression)

		// Ensure it's a deep copy (modifying one does not affect the other)
		*copied.Condition.Type = "modified"
		assert.NotEqual(t, *original.Condition.Type, *copied.Condition.Type)
	})

	t.Run("deep copy of plugin without condition", func(t *testing.T) {
		original := &Plugin{
			Name: String("rate-limiting"),
		}

		copied := original.DeepCopy()
		require.NotNil(t, copied)
		assert.Nil(t, copied.Condition)
	})

	t.Run("deep copy of PluginCondition", func(t *testing.T) {
		original := &PluginCondition{
			Type:       String("predicate"),
			Expression: String("http.path starts_with '/api'"),
		}

		copied := original.DeepCopy()
		require.NotNil(t, copied)
		assert.Equal(t, *original.Type, *copied.Type)
		assert.Equal(t, *original.Expression, *copied.Expression)
	})

	t.Run("deep copy of nil PluginCondition returns nil", func(t *testing.T) {
		var c *PluginCondition
		assert.Nil(t, c.DeepCopy())
	})
}
