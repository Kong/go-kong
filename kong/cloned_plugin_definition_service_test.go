package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClonedPluginService(T *testing.T) {
	RunWhenEnterprise(T, ">=3.15.0", RequiredFeatures{})

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	plugin := &ClonedPluginDefintion{
		Name:     String("my-cloned-plugin"),
		Priority: Uint64(100),
		Ref:      String("file-log"),
		Tags:     StringSlice("tag1", "tag2"),
	}

	createdPlugin, err := client.ClonedPlugins.Create(defaultCtx, plugin)
	require.NoError(err)
	require.NotNil(createdPlugin)
	firstID := createdPlugin.ID
	T.Cleanup(func() {
		err := client.ClonedPlugins.Delete(defaultCtx, firstID)
		require.NoError(err)
	})

	// Get by ID
	plugin, err = client.ClonedPlugins.Get(defaultCtx, createdPlugin.ID)
	require.NoError(err)
	require.NotNil(plugin)
	assert.Equal(String("my-cloned-plugin"), plugin.Name)
	assert.Equal(String("file-log"), plugin.Ref)
	assert.Equal(StringSlice("tag1", "tag2"), plugin.Tags)

	// Get by Name
	plugin, err = client.ClonedPlugins.Get(defaultCtx, createdPlugin.Name)
	require.NoError(err)
	require.NotNil(plugin)
	assert.Equal(createdPlugin.ID, plugin.ID)
	assert.Equal(String("file-log"), plugin.Ref)
	assert.Equal(StringSlice("tag1", "tag2"), plugin.Tags)

	// Update priority
	plugin.Priority = Uint64(200)
	plugin, err = client.ClonedPlugins.Update(defaultCtx, plugin)
	require.NoError(err)
	require.NotNil(plugin)
	assert.Equal(Uint64(200), plugin.Priority)

	// ID can be specified
	id := uuid.NewString()
	plugin = &ClonedPluginDefintion{
		Name:     String("my-new-cloned-plugin"),
		ID:       String(id),
		Priority: Uint64(100),
		Ref:      String("file-log"),
	}

	createdPlugin, err = client.ClonedPlugins.Create(defaultCtx, plugin)
	require.NoError(err)
	assert.NotNil(createdPlugin)
	assert.Equal(id, *createdPlugin.ID)
	secondID := createdPlugin.ID
	T.Cleanup(func() {
		err := client.ClonedPlugins.Delete(defaultCtx, secondID)
		require.NoError(err)
	})
}

func TestClonedPluginServiceListAll(T *testing.T) {
	RunWhenEnterprise(T, ">=3.15.0", RequiredFeatures{})

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	// fixtures
	plugins := []*ClonedPluginDefintion{
		{
			Name:     String("cloned-plugin-1"),
			Priority: Uint64(100),
			Ref:      String("file-log"),
			Tags:     StringSlice("tag1", "tag2"),
		},
		{
			Name:     String("cloned-plugin-2"),
			Priority: Uint64(200),
			Ref:      String("key-auth"),
			Tags:     StringSlice("tag2", "tag3"),
		},
		{
			Name:     String("cloned-plugin-3"),
			Priority: Uint64(300),
			Ref:      String("acl"),
			Tags:     StringSlice("tag1", "tag3"),
		},
	}

	// create fixtures
	for i := 0; i < len(plugins); i++ {
		plugin, err := client.ClonedPlugins.Create(defaultCtx, plugins[i])
		require.NoError(err)
		require.NotNil(plugin)
		plugins[i] = plugin
		T.Cleanup(func() {
			err := client.ClonedPlugins.Delete(defaultCtx, plugin.ID)
			require.NoError(err)
		})
	}

	pluginsFromKong, next, err := client.ClonedPlugins.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag1"),
	})
	require.NoError(err)
	assert.Nil(next)
	assert.Len(pluginsFromKong, 2)

	allPlugins, err := client.ClonedPlugins.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(allPlugins)
	assert.GreaterOrEqual(len(allPlugins), 3)
}
