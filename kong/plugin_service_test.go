package kong

import (
	"fmt"
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
	RunWhenDBMode(T, "postgres")
	SkipWhenKongRouterFlavor(T, Expressions)

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
	require.Nil(createdPlugin.InstanceName)

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

	service := &Service{
		Name: String("fooWithPlugin"),
		Host: String("upstream"),
		Port: Int(42),
		Path: String("/path"),
	}
	// Clean Data
	err = client.Services.Delete(defaultCtx, service.Name)
	assert.NoError(err)
	// Test to create plugin from service endpoint
	createdService, err := client.Services.Create(defaultCtx, service)
	assert.NoError(err)

	id = uuid.NewString()
	pluginForService := &Plugin{
		Name: String("key-auth"),
		ID:   String(id),
		Config: Configuration{
			"anonymous": "true",
		},
	}

	createdPlugin, err = client.Plugins.CreateForService(defaultCtx, createdService.Name, pluginForService)
	assert.NoError(err)
	assert.NotNil(createdPlugin)
	assert.Equal(id, *createdPlugin.ID)
	assert.Equal("true", createdPlugin.Config["anonymous"])

	createdPlugin.Config["anonymous"] = "false"
	updatedPlugin, err := client.Plugins.UpdateForService(defaultCtx, createdService.Name, createdPlugin)
	assert.NoError(err)
	assert.NotNil(updatedPlugin)
	assert.Equal(id, *updatedPlugin.ID)
	assert.Equal("false", updatedPlugin.Config["anonymous"])

	err = client.Plugins.DeleteForService(defaultCtx, createdService.Name, updatedPlugin.ID)
	assert.NoError(err)

	// Create plugin without ID
	_, err = client.Plugins.CreateForService(defaultCtx, createdService.Name, &Plugin{Name: String("key-auth")})
	assert.NoError(err)
	assert.NotNil(createdPlugin)
	assert.NotNil(createdPlugin.ID)

	assert.NoError(client.Services.Delete(defaultCtx, createdService.ID))

	// Create Plugin for route
	route := &Route{
		Name:  String("route_plugin"),
		Paths: []*string{String("/route_plugin")},
	}
	// Clean Data
	err = client.Routes.Delete(defaultCtx, route.Name)
	assert.NoError(err)
	// Test to create plugin from route endpoint
	createdRoute, err := client.Routes.Create(defaultCtx, route)
	assert.NoError(err)

	id = uuid.NewString()
	pluginForRoute := &Plugin{
		Name: String("key-auth"),
		ID:   String(id),
		Config: Configuration{
			"anonymous": "true",
		},
	}

	createdPlugin, err = client.Plugins.CreateForRoute(defaultCtx, createdRoute.Name, pluginForRoute)
	assert.NoError(err)
	assert.NotNil(createdPlugin)
	assert.Equal(id, *createdPlugin.ID)
	assert.Equal("true", createdPlugin.Config["anonymous"])

	createdPlugin.Config["anonymous"] = "false"
	updatedPlugin, err = client.Plugins.UpdateForRoute(defaultCtx, createdRoute.Name, createdPlugin)
	assert.NoError(err)
	assert.NotNil(createdPlugin)
	assert.Equal(id, *createdPlugin.ID)
	assert.Equal("false", updatedPlugin.Config["anonymous"])

	err = client.Plugins.DeleteForRoute(defaultCtx, createdRoute.Name, updatedPlugin.ID)
	assert.NoError(err)

	// Create plugin without ID
	_, err = client.Plugins.CreateForRoute(defaultCtx, createdRoute.Name, &Plugin{Name: String("key-auth")})
	assert.NoError(err)
	assert.NotNil(createdPlugin)
	assert.NotNil(createdPlugin.ID)

	assert.NoError(client.Routes.Delete(defaultCtx, createdRoute.ID))
}

func TestPluginsWithInstanceNameService(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=3.2.0")

	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	plugin := &Plugin{
		Name:         String("key-auth"),
		InstanceName: String("my-plugin"),
	}

	// create a plugin with instance_name
	createdPlugin, err := client.Plugins.Create(defaultCtx, plugin)
	require.NoError(err)
	require.NotNil(createdPlugin)
	require.Equal(plugin.InstanceName, createdPlugin.InstanceName)

	// get a plugin by instance_name
	plugin, err = client.Plugins.Get(defaultCtx, createdPlugin.InstanceName)
	require.NoError(err)
	require.NotNil(plugin)
	require.Equal(plugin.ID, createdPlugin.ID)
	require.Equal(plugin.InstanceName, createdPlugin.InstanceName)

	// update a plugin with instance_name
	plugin.Config["key_in_body"] = true
	plugin, err = client.Plugins.Update(defaultCtx, plugin)
	require.NoError(err)
	require.NotNil(plugin)
	require.Equal(plugin.ID, createdPlugin.ID)
	require.Equal(plugin.InstanceName, createdPlugin.InstanceName)
	require.Equal(true, plugin.Config["key_in_body"])

	// delete a plugin by instance_name
	require.NoError(client.Plugins.Delete(defaultCtx, createdPlugin.InstanceName))
}

func TestPluginWithTags(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=1.1.0")

	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	plugin := &Plugin{
		Name: String("key-auth"),
		Tags: StringSlice("tag1", "tag2"),
	}

	createdPlugin, err := client.Plugins.Create(defaultCtx, plugin)
	require.NoError(err)
	require.NotNil(createdPlugin)
	require.Equal(StringSlice("tag1", "tag2"), createdPlugin.Tags)

	err = client.Plugins.Delete(defaultCtx, createdPlugin.ID)
	require.NoError(err)
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
	RunWhenDBMode(T, "postgres")

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

	// create fixtures
	for i := 0; i < len(plugins); i++ {
		schema, err := client.Plugins.GetFullSchema(defaultCtx, plugins[i].Name)
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
	RunWhenDBMode(T, "postgres")
	SkipWhenKongRouterFlavor(T, Expressions)

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

	// create fixtures
	for i := 0; i < len(plugins); i++ {
		schema, err := client.Plugins.GetFullSchema(defaultCtx, plugins[i].Name)
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
	// TODO https://github.com/Kong/go-kong/issues/214 this should only skip Enterprise 3.x (with a separate test)
	// not all Enterprise versions.
	SkipWhenEnterprise(T)
	RunWhenKong(T, ">=2.3.0")

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	require.NotNil(T, client)

	tests := []struct {
		name     string
		plugin   *Plugin
		expected *Plugin
		version  string
	}{
		{
			name:    "no config no protocols",
			version: ">=3.6.0",
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
					// NOTE: realm has been introduced in 3.6 basic auth schema
					// https://docs.konghq.com/hub/kong-inc/basic-auth/changelog/#kong-gateway-36x
					"realm": "service",
				},
				Protocols: []*string{String("grpc"), String("grpcs"), String("http"), String("https")},
				Enabled:   Bool(true),
			},
		},
		{
			name:    "no config no protocols",
			version: "<3.6.0",
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
			name:    "partial config no protocols",
			version: ">=3.6.0",
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
					// NOTE: realm has been introduced in 3.6 basic auth schema
					// https://docs.konghq.com/hub/kong-inc/basic-auth/changelog/#kong-gateway-36x
					"realm": "service",
				},
				Protocols: []*string{String("grpc"), String("grpcs"), String("http"), String("https")},
				Enabled:   Bool(true),
			},
		},
		{
			name:    "partial config no protocols",
			version: "<3.6.0",
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
					"append": map[string]interface{}{
						"headers": "x-append-header:value",
					},
					"remove": map[string]interface{}{
						"body":        []interface{}{},
						"querystring": "?query=val",
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
						"headers":     "x-append-header:value",
						"querystring": []interface{}{},
					},
					"remove": map[string]interface{}{
						"body":        []interface{}{},
						"headers":     []interface{}{},
						"querystring": "?query=val",
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
		name := tc.name
		if tc.version != "" {
			name = fmt.Sprintf("%s (kong %s)", name, tc.version)
		}

		T.Run(name, func(t *testing.T) {
			if tc.version != "" {
				RunWhenKong(t, tc.version)
			}

			p := tc.plugin
			fullSchema, err := client.Plugins.GetFullSchema(defaultCtx, p.Name)
			require.NoError(t, err)
			require.NotNil(t, fullSchema)
			require.NoError(t, FillPluginsDefaults(p, fullSchema))

			if diff := cmp.Diff(p, tc.expected); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

// TestFillPluginDefaultsArbitraryMap is split from TestFillPluginDefaults due to a version compatibility issue
func TestFillPluginDefaultsArbitraryMap(T *testing.T) {
	RunWhenKong(T, ">=2.3.0")

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	require.NotNil(T, client)

	tests := []struct {
		name     string
		plugin   *Plugin
		expected *Plugin
	}{
		{
			name: "nested config with arbitrary map field",
			plugin: &Plugin{
				Name: String("http-log"),
				Config: Configuration{
					"custom_fields_by_lua": map[string]interface{}{
						"foo": "bar",
					},
				},
				Enabled:   Bool(false),
				Protocols: []*string{String("grpc"), String("grpcs")},
			},
			expected: &Plugin{
				Name: String("http-log"),
				Config: Configuration{
					"content_type": string("application/json"),
					"custom_fields_by_lua": map[string]interface{}{
						"foo": "bar",
					},
					"flush_timeout": float64(2),
					"headers":       nil,
					"http_endpoint": nil,
					"keepalive":     float64(60000),
					"method":        string("POST"),
					"queue_size":    float64(1),
					"retry_count":   float64(10),
					"timeout":       float64(10000),
				},
				Enabled:   Bool(false),
				Protocols: []*string{String("grpc"), String("grpcs")},
			},
		},
	}

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			p := tc.plugin
			fullSchema, err := client.Plugins.GetFullSchema(defaultCtx, p.Name)
			require.NoError(t, err)
			require.NotNil(t, fullSchema)
			require.NoError(t, FillPluginsDefaults(p, fullSchema))

			// the log plugins are the only plugins that use the typedefs.lua_code type in their schema
			// https://github.com/Kong/kong/commit/9df893f6aff98cd51f27f1c27fa30fdcf13fcf48 changes a number of other
			// fields for 3.3, so this test only checks the relevant field to avoid needing a version split
			if diff := cmp.Diff(p.Config["custom_fields_by_lua"], tc.expected.Config["custom_fields_by_lua"]); diff != "" {
				t.Errorf("unexpected diff:\n%s", diff)
			}
		})
	}
}

func TestPluginsWithConsumerGroup(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenEnterprise(T, ">=3.4.0", RequiredFeatures{})

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	// create consumer group
	cg := &ConsumerGroup{
		Name: String("foo"),
	}
	createdCG, err := client.ConsumerGroups.Create(defaultCtx, cg)
	require.NoError(err)
	assert.NotNil(createdCG)

	plugin := &Plugin{
		Name: String("rate-limiting-advanced"),
		Config: Configuration{
			"limit":       []interface{}{5},
			"window_size": []interface{}{30},
		},
		ConsumerGroup: &ConsumerGroup{
			ID: createdCG.ID,
		},
	}

	createdPlugin, err := client.Plugins.Create(defaultCtx, plugin)
	assert.NoError(err)
	require.NotNil(createdPlugin)
	require.Nil(createdPlugin.InstanceName)

	plugin, err = client.Plugins.Get(defaultCtx, createdPlugin.ID)
	assert.NoError(err)
	assert.NotNil(plugin)
	assert.Equal(plugin.ConsumerGroup.ID, createdCG.ID)
	assert.Equal("sliding", plugin.Config["window_type"])

	createdPlugin.Config["window_type"] = "fixed"
	updatedPlugin, err := client.Plugins.UpdateForConsumerGroup(defaultCtx, createdCG.Name, createdPlugin)
	assert.NoError(err)
	assert.NotNil(createdPlugin)
	assert.Equal("fixed", updatedPlugin.Config["window_type"])

	assert.NoError(client.ConsumerGroups.Delete(defaultCtx, createdCG.ID))
	// assert the plugin was cascade deleted
	plugin, err = client.Plugins.Get(defaultCtx, createdPlugin.ID)
	assert.Nil(plugin)
	assert.True(IsNotFoundErr(err))

	// create another consumer group
	cg = &ConsumerGroup{
		Name: String("bar"),
	}
	createdCG, err = client.ConsumerGroups.Create(defaultCtx, cg)
	require.NoError(err)
	assert.NotNil(createdCG)

	id := uuid.NewString()
	pluginForCG := &Plugin{
		Name: String("request-transformer"),
		ID:   String(id),
	}

	createdPlugin, err = client.Plugins.CreateForConsumerGroup(defaultCtx, createdCG.Name, pluginForCG)
	assert.NoError(err)
	assert.NotNil(createdPlugin)
	assert.Equal(createdPlugin.ConsumerGroup.ID, createdCG.ID)

	assert.NoError(client.ConsumerGroups.Delete(defaultCtx, createdCG.ID))
	// assert the plugin was cascade deleted
	plugin, err = client.Plugins.Get(defaultCtx, createdPlugin.ID)
	assert.Nil(plugin)
	assert.True(IsNotFoundErr(err))

	// create another consumer group
	cg = &ConsumerGroup{
		Name: String("baz"),
	}
	createdCG, err = client.ConsumerGroups.Create(defaultCtx, cg)
	require.NoError(err)
	assert.NotNil(createdCG)

	plugins := []*Plugin{
		{
			Name:          String("request-transformer"),
			ConsumerGroup: createdCG,
		},
		{
			Name: String("rate-limiting-advanced"),
			Config: Configuration{
				"limit":       []interface{}{5},
				"window_size": []interface{}{30},
			},
			ConsumerGroup: createdCG,
		},
	}

	// create fixtures
	for i := 0; i < len(plugins); i++ {
		plugin, err := client.Plugins.Create(defaultCtx, plugins[i])
		assert.NoError(err)
		assert.NotNil(plugin)
		plugins[i] = plugin
	}

	pluginsFromKong, err := client.Plugins.ListAllForConsumerGroups(defaultCtx, createdCG.ID)
	assert.NoError(err)
	assert.NotNil(pluginsFromKong)
	assert.Len(pluginsFromKong, 2)

	assert.NoError(client.ConsumerGroups.Delete(defaultCtx, createdCG.ID))
	// assert the plugins were cascade deleted
	for _, plugin := range plugins {
		res, err := client.Plugins.Get(defaultCtx, plugin.ID)
		assert.Nil(res)
		assert.True(IsNotFoundErr(err))
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
