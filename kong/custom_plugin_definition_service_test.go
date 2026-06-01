package kong

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCustomPluginHandler = `return {
  VERSION = "1.0,0",
  PRIORITY = 500,
  access = function(self, config)
    kong.service.request.set_header(config.name, config.value)
  end
}`

func testCustomPluginSchema(name string) string {
	return fmt.Sprintf(`return {
  name = %q,
  fields = {
    { protocols = require("kong.db.schema.typedefs").protocols_http },
    {
      config = {
        type = "record",
        fields = {
          { name = { description = "The name of the header to set.", type = "string", required = true, }, },
          { value = { description = "The value for the header.", type = "string", required = true, }, },
        },
      },
    },
  },
}`, name)
}

func TestCustomPluginService(T *testing.T) {
	RunWhenEnterprise(T, ">=3.15.0", RequiredFeatures{CustomPlugins: true})

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	plugin := &CustomPluginDefintion{
		Name:    String("my-custom-plugin"),
		Handler: String(testCustomPluginHandler),
		Schema:  String(testCustomPluginSchema("my-custom-plugin")),
		Tags:    StringSlice("tag1", "tag2"),
	}

	createdPlugin, err := client.CustomPlugins.Create(defaultCtx, plugin)
	require.NoError(err)
	require.NotNil(createdPlugin)
	firstID := createdPlugin.ID
	T.Cleanup(func() {
		err := client.CustomPlugins.Delete(defaultCtx, firstID)
		require.NoError(err)
	})

	plugin, err = client.CustomPlugins.Get(defaultCtx, createdPlugin.ID)
	require.NoError(err)
	require.NotNil(plugin)

	plugin.Tags = StringSlice("tag1", "tag2", "tag3")
	plugin, err = client.CustomPlugins.Update(defaultCtx, plugin)
	require.NoError(err)
	require.NotNil(plugin)
	assert.Equal(StringSlice("tag1", "tag2", "tag3"), plugin.Tags)

	// ID can be specified
	id := uuid.NewString()
	plugin = &CustomPluginDefintion{
		Name:    String("my-custom-plugin-2"),
		ID:      String(id),
		Handler: String(testCustomPluginHandler),
		Schema:  String(testCustomPluginSchema("my-custom-plugin-2")),
	}

	createdPlugin, err = client.CustomPlugins.Create(defaultCtx, plugin)
	require.NoError(err)
	assert.NotNil(createdPlugin)
	assert.Equal(id, *createdPlugin.ID)
	secondID := createdPlugin.ID
	T.Cleanup(func() {
		err := client.CustomPlugins.Delete(defaultCtx, secondID)
		require.NoError(err)
	})
}

func TestCustomPluginServiceListAll(T *testing.T) {
	RunWhenEnterprise(T, ">=3.15.0", RequiredFeatures{CustomPlugins: true})

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	// fixtures
	plugins := []*CustomPluginDefintion{
		{
			Name:    String("custom-plugin-1"),
			Handler: String(testCustomPluginHandler),
			Schema:  String(testCustomPluginSchema("custom-plugin-1")),
			Tags:    StringSlice("tag1", "tag2"),
		},
		{
			Name:    String("custom-plugin-2"),
			Handler: String(testCustomPluginHandler),
			Schema:  String(testCustomPluginSchema("custom-plugin-2")),
			Tags:    StringSlice("tag2", "tag3"),
		},
		{
			Name:    String("custom-plugin-3"),
			Handler: String(testCustomPluginHandler),
			Schema:  String(testCustomPluginSchema("custom-plugin-3")),
			Tags:    StringSlice("tag1", "tag3"),
		},
	}

	// create fixtures
	for i := 0; i < len(plugins); i++ {
		plugin, err := client.CustomPlugins.Create(defaultCtx, plugins[i])
		require.NoError(err)
		require.NotNil(plugin)
		plugins[i] = plugin
		T.Cleanup(func() {
			err := client.CustomPlugins.Delete(defaultCtx, plugin.ID)
			require.NoError(err)
		})
	}

	pluginsFromKong, next, err := client.CustomPlugins.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag1"),
	})
	require.NoError(err)
	assert.Nil(next)
	assert.Len(pluginsFromKong, 2)

	allPlugins, err := client.CustomPlugins.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(allPlugins)
	assert.GreaterOrEqual(len(allPlugins), 3)
}
