package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var defaultConfigRedisEEPartial = Configuration{
	"cluster_max_redirections": float64(5),
	"cluster_nodes":            nil,
	"connect_timeout":          float64(2000),
	"connection_is_proxied":    bool(false),
	"database":                 float64(0),
	"host":                     "127.0.0.1",
	"keepalive_backlog":        nil,
	"keepalive_pool_size":      float64(256),
	"password":                 nil,
	"port":                     float64(6379),
	"read_timeout":             float64(2000),
	"send_timeout":             float64(2000),
	"sentinel_master":          nil,
	"sentinel_nodes":           nil,
	"sentinel_password":        nil,
	"sentinel_role":            nil,
	"sentinel_username":        nil,
	"server_name":              nil,
	"ssl":                      bool(false),
	"ssl_verify":               bool(false),
	"username":                 nil,
}

func TestPartialServiceCreateEndpoint(t *testing.T) {
	RunWhenEnterprise(t, ">=3.10.0", RequiredFeatures{})

	require := require.New(t)
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	t.Run("invalid partial creation -  nil partial", func(_ *testing.T) {
		notCreatedPartial, err := client.Partials.Create(defaultCtx, nil)

		assert.Error(err)
		assert.ErrorContains(err, "cannot create a nil partial")
		assert.Nil(notCreatedPartial)
	})

	t.Run("invalid partial creation -  empty partial", func(_ *testing.T) {
		notCreatedPartial, err := client.Partials.Create(defaultCtx, &Partial{})

		assert.Error(err)
		assert.ErrorContains(err, "partial type cannot be nil")
		assert.Nil(notCreatedPartial)
	})

	t.Run("invalid partial creation -  nil type", func(_ *testing.T) {
		notCreatedPartial, err := client.Partials.Create(defaultCtx, &Partial{
			Name: String("my-test-partial"),
		})

		assert.Error(err)
		assert.ErrorContains(err, "partial type cannot be nil")
		assert.Nil(notCreatedPartial)
	})

	t.Run("create valid partial", func(_ *testing.T) {
		validPartial := &Partial{
			Name: String("my-test-partial"),
			Type: String("redis-ee"),
		}

		createdPartial, err := client.Partials.Create(defaultCtx, validPartial)
		require.NoError(err)
		require.NotNil(createdPartial)
		assert.Equal("my-test-partial", *createdPartial.Name)
		assert.Equal("redis-ee", *createdPartial.Type)
		assert.Equal(defaultConfigRedisEEPartial, createdPartial.Config)

		t.Cleanup(func() {
			assert.NoError(client.Partials.Delete(defaultCtx, createdPartial.ID))
		})
	})
}

func TestPartialServiceGetEndpoint(t *testing.T) {
	RunWhenEnterprise(t, ">=3.10.0", RequiredFeatures{})

	require := require.New(t)
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	t.Run("invalid get -  empty id", func(_ *testing.T) {
		noPartial, err := client.Partials.Get(defaultCtx, String(""))

		assert.Error(err)
		assert.ErrorContains(err, "partialID cannot be nil for Get operation")
		assert.Nil(noPartial)
	})

	t.Run("get by id", func(_ *testing.T) {
		createdPartial, err := client.Partials.Create(defaultCtx, &Partial{
			Name: String("my-test-partial"),
			Type: String("redis-ee"),
		})
		require.NoError(err)
		require.NotNil(createdPartial)

		t.Cleanup(func() {
			assert.NoError(client.Partials.Delete(defaultCtx, createdPartial.ID))
		})

		fetchedPartial, err := client.Partials.Get(defaultCtx, createdPartial.ID)
		assert.NoError(err)
		assert.NotNil(fetchedPartial)
		assert.Equal("my-test-partial", *fetchedPartial.Name)
		assert.Equal("redis-ee", *fetchedPartial.Type)
	})
}

func TestPartialServiceUpdateEndpoint(t *testing.T) {
	RunWhenEnterprise(t, ">=3.10.0", RequiredFeatures{})

	require := require.New(t)
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	t.Run("invalid update -  nil partial", func(_ *testing.T) {
		notUpdatedPartial, err := client.Partials.Update(defaultCtx, nil)

		assert.Error(err)
		assert.ErrorContains(err, "cannot update a nil partial")
		assert.Nil(notUpdatedPartial)
	})

	t.Run("valid update", func(_ *testing.T) {
		createdPartial, err := client.Partials.Create(defaultCtx, &Partial{
			Name: String("my-test-partial"),
			Type: String("redis-ee"),
		})
		require.NoError(err)
		require.NotNil(createdPartial)

		t.Cleanup(func() {
			assert.NoError(client.Partials.Delete(defaultCtx, createdPartial.ID))
		})

		// initially created with default config
		assert.Equal(defaultConfigRedisEEPartial, createdPartial.Config)

		// config to update
		createdPartial.Config = Configuration{
			"send_timeout":    2001,
			"read_timeout":    3001,
			"connect_timeout": 4001,
		}

		// update partial
		updatedPartial, err := client.Partials.Update(defaultCtx, createdPartial)
		require.NoError(err)
		require.NotNil(updatedPartial)

		assert.Equal(float64(2001), updatedPartial.Config["send_timeout"])
		assert.Equal(float64(3001), updatedPartial.Config["read_timeout"])
		assert.Equal(float64(4001), updatedPartial.Config["connect_timeout"])
	})
}

func TestPartialServiceDeleteEndpoint(t *testing.T) {
	RunWhenEnterprise(t, ">=3.10.0", RequiredFeatures{})

	require := require.New(t)
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	t.Run("invalid delete -  empty id", func(_ *testing.T) {
		err := client.Partials.Delete(defaultCtx, String(""))

		assert.Error(err)
		assert.ErrorContains(err, "partialID cannot be nil for Delete operation")
	})

	t.Run("valid delete -  by id", func(_ *testing.T) {
		createdPartial, err := client.Partials.Create(defaultCtx, &Partial{
			Name: String("my-demo-partial"),
			Type: String("redis-ee"),
		})
		require.NoError(err)
		require.NotNil(createdPartial)

		err = client.Partials.Delete(defaultCtx, createdPartial.ID)
		assert.NoError(err)
	})
}

func TestPartialServiceListEndpoint(t *testing.T) {
	RunWhenEnterprise(t, ">=3.10.0", RequiredFeatures{})

	require := require.New(t)
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	populatePartials(t, client)

	// Testing pagination
	partialsFromKong := []*Partial{}

	// first page
	page1, next, err := client.Partials.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(err)
	require.NotNil(page1)
	require.NotNil(next)
	assert.Len(page1, 1)
	partialsFromKong = append(partialsFromKong, page1...)

	// last page
	next.Size = 3
	page2, next, err := client.Partials.List(defaultCtx, next)
	require.NoError(err)
	require.NotNil(page2)
	require.Nil(next)
	assert.Len(page2, 3)
	partialsFromKong = append(partialsFromKong, page2...)

	assert.Len(partialsFromKong, 4)
	for _, p := range partialsFromKong {
		assert.Equal("redis-ee", *p.Type)
		assert.Contains(*p.Name, "test-partial-")
	}

	t.Cleanup(func() {
		for _, p := range partialsFromKong {
			assert.NoError(client.Partials.Delete(defaultCtx, p.ID))
		}
	})
}

func TestPartialServiceListAllEndpoint(t *testing.T) {
	RunWhenEnterprise(t, ">=3.10.0", RequiredFeatures{})

	require := require.New(t)
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)
	populatePartials(t, client)

	partials, err := client.Partials.ListAll(defaultCtx)
	require.NoError(err)
	require.NotNil(partials)
	assert.Len(partials, 4)

	t.Cleanup(func() {
		for _, p := range partials {
			assert.NoError(client.Partials.Delete(defaultCtx, p.ID))
		}
	})
}

func populatePartials(t *testing.T, client *Client) {
	require := require.New(t)
	partials := []*Partial{
		{
			Name: String("test-partial-1"),
			Type: String("redis-ee"),
		},
		{
			Name: String("test-partial-2"),
			Type: String("redis-ee"),
		},
		{
			Name: String("test-partial-3"),
			Type: String("redis-ee"),
		},
		{
			Name: String("test-partial-4"),
			Type: String("redis-ee"),
		},
	}

	for _, p := range partials {
		createdPartial, err := client.Partials.Create(defaultCtx, p)
		require.NoError(err)
		require.NotNil(createdPartial)
	}
}

func TestPartialServiceGetFullSchema(t *testing.T) {
	RunWhenEnterprise(t, ">=3.10.0", RequiredFeatures{})

	require := require.New(t)
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	t.Run("successful schema retrieval for a partial", func(_ *testing.T) {
		schema, err := client.Partials.GetFullSchema(defaultCtx, String("redis-ee"))
		require.NoError(err)
		require.NotNil(schema)

		schema, err = client.Partials.GetFullSchema(defaultCtx, String("redis-ce"))
		require.NoError(err)
		require.NotNil(schema)
	})

	t.Run("invalid schema retrieval - empty string parameter", func(_ *testing.T) {
		schema, err := client.Partials.GetFullSchema(defaultCtx, String(""))
		require.Error(err)
		require.Nil(schema)
		assert.ErrorContains(err, "partialName cannot be nil for GetFullSchema operation")
	})

	t.Run("invalid schema retrieval - nil parameter", func(_ *testing.T) {
		schema, err := client.Partials.GetFullSchema(defaultCtx, nil)
		require.Error(err)
		require.Nil(schema)
		assert.ErrorContains(err, "partialName cannot be nil for GetFullSchema operation")
	})

	t.Run("invalid schema retrieval - fake partial", func(_ *testing.T) {
		schema, err := client.Partials.GetFullSchema(defaultCtx, String("fake-partial"))
		require.Error(err)
		require.Nil(schema)
		assert.ErrorContains(err, "No partial of type 'fake-partial'")
	})
}

func TestPartialServiceGetLinkedPlugins(t *testing.T) {
	RunWhenEnterprise(t, ">=3.10.0", RequiredFeatures{})

	require := require.New(t)
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	t.Run("successful links retrieval for a partial", func(_ *testing.T) {
		// Create redis-ee partial
		redisEEPartial := &Partial{
			Name: String("my-test-partial"),
			Type: String("redis-ee"),
		}

		redisEEPartial, err := client.Partials.Create(defaultCtx, redisEEPartial)
		require.NoError(err)
		require.NotNil(redisEEPartial)

		t.Cleanup(func() {
			assert.NoError(client.Partials.Delete(defaultCtx, redisEEPartial.ID))
		})

		// Create RLA plugin with partial
		rlaPlugin := &Plugin{
			Name: String("rate-limiting-advanced"),
			Config: Configuration{
				"limit":       []interface{}{5},
				"window_size": []interface{}{30},
			},
			Partials: []*PartialLink{
				{
					Partial: &Partial{
						ID: redisEEPartial.ID,
					},
				},
			},
			Enabled: Bool(true),
		}
		rlaPlugin, err = client.Plugins.Create(defaultCtx, rlaPlugin)
		require.NoError(err)
		require.NotNil(rlaPlugin)

		t.Cleanup(func() {
			assert.NoError(client.Plugins.Delete(defaultCtx, rlaPlugin.ID))
		})

		plugins, _, err := client.Partials.GetLinkedPlugins(defaultCtx, redisEEPartial.ID, nil)
		require.NoError(err)
		require.NotNil(plugins)
		assert.Len(plugins, 1)
		assert.Equal(rlaPlugin.Name, plugins[0].Name)
		assert.Equal(rlaPlugin.ID, plugins[0].ID)
	})

	t.Run("invalid links retrieval - empty string parameter", func(_ *testing.T) {
		plugins, next, err := client.Partials.GetLinkedPlugins(defaultCtx, String(""), nil)
		require.Error(err)
		require.Nil(plugins)
		require.Nil(next)
		assert.ErrorContains(err, "partialID cannot be nil for GetLinkedPlugins operation")
	})

	t.Run("invalid links retrieval - nil parameter", func(_ *testing.T) {
		plugins, next, err := client.Partials.GetLinkedPlugins(defaultCtx, nil, nil)
		require.Error(err)
		require.Nil(plugins)
		require.Nil(next)
		assert.ErrorContains(err, "partialID cannot be nil for GetLinkedPlugins operation")
	})

	t.Run("invalid links retrieval - fake partial", func(_ *testing.T) {
		randomID := uuid.NewString()
		plugins, next, err := client.Partials.GetLinkedPlugins(defaultCtx, String(randomID), nil)
		require.Error(err)
		require.Nil(plugins)
		require.Nil(next)
		assert.ErrorContains(err, "HTTP status 404")
	})
}
