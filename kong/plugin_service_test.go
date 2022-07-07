package kong

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPluginsServiceValidation(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	goodPlugin := &Plugin{
		Name: String("key-auth"),
		Config: Configuration{
			"anonymous": "true",
		},
	}

	badPlugin := &Plugin{
		Name: String("key-auth"),
		Config: Configuration{
			"garbage": true,
		},
	}

	valid, _, err := client.Plugins.Validate(defaultCtx, goodPlugin)
	assert.True(valid)
	assert.NoError(err)

	valid, msg, err := client.Plugins.Validate(defaultCtx, badPlugin)
	assert.False(valid)
	assert.NoError(err)
	assert.Equal("schema violation (config.garbage: unknown field)", msg)
}

func TestPluginsService(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	plugin := &Plugin{
		Name: String("key-auth"),
	}

	createdPlugin, err := client.Plugins.Create(defaultCtx, plugin)
	assert.NoError(err)
	require.NotNil(createdPlugin)

	plugin, err = client.Plugins.Get(defaultCtx, createdPlugin.ID)
	assert.NoError(err)
	assert.NotNil(plugin)

	plugin.Config["key_in_body"] = true
	plugin, err = client.Plugins.Update(defaultCtx, plugin)
	assert.NoError(err)
	assert.NotNil(plugin)
	assert.Equal(true, plugin.Config["key_in_body"])

	err = client.Plugins.Delete(defaultCtx, createdPlugin.ID)
	assert.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	plugin = &Plugin{
		Name: String("prometheus"),
		ID:   String(id),
	}

	createdPlugin, err = client.Plugins.Create(defaultCtx, plugin)
	assert.NoError(err)
	assert.NotNil(createdPlugin)
	assert.Equal(id, *createdPlugin.ID)

	err = client.Plugins.Delete(defaultCtx, createdPlugin.ID)
	assert.NoError(err)
}

func TestPluginWithTags(T *testing.T) {
	RunWhenKong(T, ">=1.1.0")
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	plugin := &Plugin{
		Name: String("key-auth"),
		Tags: StringSlice("tag1", "tag2"),
	}

	createdPlugin, err := client.Plugins.Create(defaultCtx, plugin)
	assert.NoError(err)
	assert.NotNil(createdPlugin)
	assert.Equal(StringSlice("tag1", "tag2"), createdPlugin.Tags)

	err = client.Plugins.Delete(defaultCtx, createdPlugin.ID)
	assert.NoError(err)
}

func TestPluginWithOrdering(T *testing.T) {
	RunWhenEnterprise(T, ">=3.0.0", RequiredFeatures{})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	plugin := &Plugin{
		Name: String("request-termination"),
		Ordering: &PluginOrdering{
			Before: PluginOrderingPhase{
				"access": []string{"key-auth", "basic-auth"},
			},
			After: PluginOrderingPhase{
				"access": []string{"correlation-id"},
			},
		},
	}

	createdPlugin, err := client.Plugins.Create(defaultCtx, plugin)
	assert.NoError(err)
	assert.NotNil(createdPlugin)
	assert.Equal(PluginOrdering{
		Before: PluginOrderingPhase{
			"access": []string{"key-auth", "basic-auth"},
		},
		After: PluginOrderingPhase{
			"access": []string{"correlation-id"},
		},
	}, *createdPlugin.Ordering)

	err = client.Plugins.Delete(defaultCtx, createdPlugin.ID)
	assert.NoError(err)

	plugin = &Plugin{
		Name: String("request-termination"),
		Ordering: &PluginOrdering{
			Before: PluginOrderingPhase{
				"not-a-phase": []string{"key-auth", "basic-auth"},
			},
		},
	}

	createdPlugin, err = client.Plugins.Create(defaultCtx, plugin)
	assert.Error(err)
	assert.Nil(createdPlugin)

	plugin = &Plugin{
		Name: String("request-termination"),
		Ordering: &PluginOrdering{
			Before: PluginOrderingPhase{
				"access": []string{"not-a-plugin"},
			},
		},
	}

	createdPlugin, err = client.Plugins.Create(defaultCtx, plugin)
	assert.Error(err)
	assert.Nil(createdPlugin)
}

func TestUnknownPlugin(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	plugin, err := client.Plugins.Create(defaultCtx, &Plugin{
		Name: String("plugin-not-present"),
	})
	assert.NotNil(err)
	assert.Nil(plugin)
}

func TestPluginListEndpoint(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	// fixtures
	plugins := []*Plugin{
		{
			Name: String("key-auth"),
		},
		{
			Name: String("basic-auth"),
		},
		{
			Name: String("jwt"),
		},
	}

	// create fixturs
	for i := 0; i < len(plugins); i++ {
		schema, err := client.Plugins.GetSchema(defaultCtx, plugins[i].Name)
		assert.NoError(err)
		assert.NotNil(schema)
		plugin, err := client.Plugins.Create(defaultCtx, plugins[i])
		assert.NoError(err)
		assert.NotNil(plugin)
		plugins[i] = plugin
	}

	pluginsFromKong, next, err := client.Plugins.List(defaultCtx, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(pluginsFromKong)
	assert.Equal(3, len(pluginsFromKong))

	// check if we see all plugins
	assert.True(comparePlugins(T, plugins, pluginsFromKong))

	// Test pagination
	pluginsFromKong = []*Plugin{}

	// first page
	page1, next, err := client.Plugins.List(defaultCtx, &ListOpt{Size: 1})
	assert.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	pluginsFromKong = append(pluginsFromKong, page1...)

	// second page
	page2, next, err := client.Plugins.List(defaultCtx, next)
	assert.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page2)
	assert.Equal(1, len(page2))
	pluginsFromKong = append(pluginsFromKong, page2...)

	// last page
	page3, next, err := client.Plugins.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(page3)
	assert.Equal(1, len(page3))
	pluginsFromKong = append(pluginsFromKong, page3...)

	assert.True(comparePlugins(T, plugins, pluginsFromKong))

	plugins, err = client.Plugins.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(plugins)
	assert.Equal(3, len(plugins))

	for i := 0; i < len(plugins); i++ {
		assert.NoError(client.Plugins.Delete(defaultCtx, plugins[i].ID))
	}
}

func TestPluginListAllForEntityEndpoint(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	// fixtures

	createdService, err := client.Services.Create(defaultCtx, &Service{
		Name: String("foo"),
		Host: String("upstream"),
		Port: Int(42),
		Path: String("/path"),
	})
	assert.NoError(err)
	assert.NotNil(createdService)

	createdRoute, err := client.Routes.Create(defaultCtx, &Route{
		Hosts:   StringSlice("host1.com", "host2.com"),
		Service: createdService,
	})
	assert.NoError(err)
	assert.NotNil(createdRoute)

	createdConsumer, err := client.Consumers.Create(defaultCtx, &Consumer{
		Username: String("foo"),
	})
	assert.NoError(err)
	require.NotNil(createdConsumer)

	plugins := []*Plugin{
		// global
		{
			Name: String("key-auth"),
		},
		{
			Name: String("basic-auth"),
		},
		{
			Name: String("jwt"),
		},
		// specific to route
		{
			Name:  String("key-auth"),
			Route: createdRoute,
		},
		{
			Name:  String("jwt"),
			Route: createdRoute,
		},
		// specific to service
		{
			Name:    String("key-auth"),
			Service: createdService,
		},
		{
			Name:    String("jwt"),
			Service: createdService,
		},
		// specific to consumer
		{
			Name:     String("rate-limiting"),
			Consumer: createdConsumer,
			Config: map[string]interface{}{
				"second": 1,
			},
		},
	}

	// create fixturs
	for i := 0; i < len(plugins); i++ {
		schema, err := client.Plugins.GetSchema(defaultCtx, plugins[i].Name)
		assert.NoError(err)
		assert.NotNil(schema)
		plugin, err := client.Plugins.Create(defaultCtx, plugins[i])
		assert.NoError(err)
		assert.NotNil(plugin)
		plugins[i] = plugin
	}

	pluginsFromKong, err := client.Plugins.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(pluginsFromKong)
	assert.Equal(len(plugins), len(pluginsFromKong))

	// check if we see all plugins
	assert.True(comparePlugins(T, plugins, pluginsFromKong))

	assert.True(comparePlugins(T, plugins, pluginsFromKong))

	pluginsFromKong, err = client.Plugins.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(pluginsFromKong)
	assert.Equal(8, len(pluginsFromKong))

	pluginsFromKong, err = client.Plugins.ListAllForConsumer(defaultCtx,
		createdConsumer.ID)
	assert.NoError(err)
	assert.NotNil(pluginsFromKong)
	assert.Equal(1, len(pluginsFromKong))

	pluginsFromKong, err = client.Plugins.ListAllForService(defaultCtx,
		createdService.ID)
	assert.NoError(err)
	assert.NotNil(pluginsFromKong)
	assert.Equal(2, len(pluginsFromKong))

	pluginsFromKong, err = client.Plugins.ListAllForRoute(defaultCtx,
		createdRoute.ID)
	assert.NoError(err)
	assert.NotNil(pluginsFromKong)
	assert.Equal(2, len(pluginsFromKong))

	for i := 0; i < len(plugins); i++ {
		assert.NoError(client.Plugins.Delete(defaultCtx, plugins[i].ID))
	}

	assert.NoError(client.Consumers.Delete(defaultCtx, createdConsumer.ID))
	assert.NoError(client.Routes.Delete(defaultCtx, createdRoute.ID))
	assert.NoError(client.Services.Delete(defaultCtx, createdService.ID))
}

func TestPluginGetFullSchema(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	schema, err := client.Plugins.GetFullSchema(defaultCtx, String("key-auth"))
	_, ok := schema["fields"]
	assert.True(ok)
	assert.NoError(err)

	schema, err = client.Plugins.GetFullSchema(defaultCtx, String("noexist"))
	assert.Nil(schema)
	assert.NotNil(err)
	assert.True(IsNotFoundErr(err))
}

func TestFillPluginDefaults(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	tests := []struct {
		name     string
		plugin   *Plugin
		expected *Plugin
	}{
		{
			name: "no config no protocols",
			plugin: &Plugin{
				Name:  String("basic-auth"),
				RunOn: String("test"),
			},
			expected: &Plugin{
				Name:  String("basic-auth"),
				RunOn: String("test"),
				Config: Configuration{
					"anonymous":        nil,
					"hide_credentials": false,
				},
				Protocols: []*string{String("grpc"), String("grpcs"), String("http"), String("https")},
				Enabled:   Bool(true),
			},
		},
		{
			name: "partial config no protocols",
			plugin: &Plugin{
				Name: String("basic-auth"),
				Consumer: &Consumer{
					ID: String("3bb9a73c-a467-11ec-b909-0242ac120002"),
				},
				Config: Configuration{
					"hide_credentials": true,
				},
			},
			expected: &Plugin{
				Name: String("basic-auth"),
				Consumer: &Consumer{
					ID: String("3bb9a73c-a467-11ec-b909-0242ac120002"),
				},
				Config: Configuration{
					"anonymous":        nil,
					"hide_credentials": true,
				},
				Protocols: []*string{String("grpc"), String("grpcs"), String("http"), String("https")},
				Enabled:   Bool(true),
			},
		},
		{
			name: "nested config partial protocols",
			plugin: &Plugin{
				Name: String("request-transformer"),
				Config: Configuration{
					"add": map[string]interface{}{
						"body":        []interface{}{},
						"headers":     "x-new-header:value",
						"querystring": []interface{}{},
					},
				},
				Enabled:   Bool(false),
				Protocols: []*string{String("grpc"), String("grpcs")},
			},
			expected: &Plugin{
				Name: String("request-transformer"),
				Config: Configuration{
					"http_method": nil,
					"add": map[string]interface{}{
						"body":        []interface{}{},
						"headers":     "x-new-header:value",
						"querystring": []interface{}{},
					},
					"append": map[string]interface{}{
						"body":        []interface{}{},
						"headers":     []interface{}{},
						"querystring": []interface{}{},
					},
					"remove": map[string]interface{}{
						"body":        []interface{}{},
						"headers":     []interface{}{},
						"querystring": []interface{}{},
					},
					"rename": map[string]interface{}{
						"body":        []interface{}{},
						"headers":     []interface{}{},
						"querystring": []interface{}{},
					},
					"replace": map[string]interface{}{
						"body":        []interface{}{},
						"headers":     []interface{}{},
						"querystring": []interface{}{},
						"uri":         nil,
					},
				},
				Protocols: []*string{String("grpc"), String("grpcs")},
				Enabled:   Bool(false),
			},
		},
	}

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			p := tc.plugin
			fullSchema, err := client.Plugins.GetFullSchema(defaultCtx, p.Name)
			assert.NoError(err)
			assert.NotNil(fullSchema)
			if err := FillPluginsDefaults(p, fullSchema); err != nil {
				t.Errorf(err.Error())
			}
			if diff := cmp.Diff(p, tc.expected); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func comparePlugins(T *testing.T, expected, actual []*Plugin) bool {
	var expectedNames, actualNames []string
	for _, plugin := range expected {
		if !assert.NotNil(T, plugin) {
			continue
		}
		expectedNames = append(expectedNames, *plugin.Name)
	}

	for _, plugin := range actual {
		actualNames = append(actualNames, *plugin.Name)
	}

	return (compareSlices(expectedNames, actualNames))
}
