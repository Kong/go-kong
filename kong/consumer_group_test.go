package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumerGroupsService(t *testing.T) {
	RunWhenEnterprise(t, ">=2.7.0", RequiredFeatures{})
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	assert.NotNil(client)

	cg := &ConsumerGroup{
		Name: String("foo"),
	}

	createdConsumerGroup, err := client.ConsumerGroups.Create(defaultCtx, cg)
	require.NoError(t, err)
	assert.NotNil(createdConsumerGroup)

	var cgc *ConsumerGroupObject
	cgc, err = client.ConsumerGroups.Get(defaultCtx, createdConsumerGroup.ID)
	require.NoError(t, err)
	assert.NotNil(cgc)

	cg.Name = String("bar")
	cg.ID = cgc.ConsumerGroup.ID
	cg, err = client.ConsumerGroups.Update(defaultCtx, cg)
	require.NoError(t, err)
	assert.NotNil(cg)
	assert.Equal("bar", *cg.Name)

	err = client.ConsumerGroups.Delete(defaultCtx, cg.ID)
	require.NoError(t, err)

	// ID can be specified
	id := uuid.NewString()
	cg = &ConsumerGroup{
		Name: String("foo"),
		ID:   String(id),
	}

	createdConsumerGroup, err = client.ConsumerGroups.Create(defaultCtx, cg)
	require.NoError(t, err)
	assert.NotNil(createdConsumerGroup)
	assert.Equal(id, *createdConsumerGroup.ID)

	err = client.ConsumerGroups.Delete(defaultCtx, createdConsumerGroup.ID)
	require.NoError(t, err)
}

func TestConsumerGroupWithTags(t *testing.T) {
	RunWhenEnterprise(t, ">=3.1.1", RequiredFeatures{})
	require := require.New(t)
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	cg := &ConsumerGroup{
		Name: String("foo"),
		Tags: StringSlice("tag1", "tag2"),
	}

	createdConsumerGroup, err := client.ConsumerGroups.Create(defaultCtx, cg)
	require.NoError(err)
	t.Cleanup(func() {
		err = client.ConsumerGroups.Delete(defaultCtx, createdConsumerGroup.ID)
		require.NoError(err)
	})
	assert.NotNil(createdConsumerGroup)
	require.Equal(cg.Tags, createdConsumerGroup.Tags)

	createdConsumerGroup.Tags = StringSlice("tag1", "tag2", "tag3")
	updatedConsumerGroup, err := client.ConsumerGroups.Update(defaultCtx, createdConsumerGroup)
	require.NoError(err)
	assert.NotNil(updatedConsumerGroup)
	assert.Equal(updatedConsumerGroup.Tags, createdConsumerGroup.Tags)
}

func TestConsumerGroupListEndpoint(t *testing.T) {
	RunWhenEnterprise(t, ">=2.7.0", RequiredFeatures{})
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	assert.NotNil(client)

	// fixtures
	consumerGroups := []*ConsumerGroup{
		{
			Name: String("foo1"),
		},
		{
			Name: String("foo2"),
		},
		{
			Name: String("foo3"),
		},
	}

	// create fixturs
	for i := 0; i < len(consumerGroups); i++ {
		cg, err := client.ConsumerGroups.Create(defaultCtx, consumerGroups[i])
		require.NoError(t, err)
		assert.NotNil(cg)
		consumerGroups[i] = cg
	}

	consumerGroupsFromKong, next, err := client.ConsumerGroups.List(defaultCtx, nil)
	require.NoError(t, err)
	assert.Nil(next)
	assert.NotNil(consumerGroupsFromKong)
	assert.Len(consumerGroupsFromKong, 3)

	// check if we see all consumer groups
	assert.True(compareConsumerGroups(consumerGroups, consumerGroupsFromKong))

	// Test pagination
	consumerGroupsFromKong = []*ConsumerGroup{}

	// first page
	page1, next, err := client.ConsumerGroups.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(t, err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)
	consumerGroupsFromKong = append(consumerGroupsFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.ConsumerGroups.List(defaultCtx, next)
	require.NoError(t, err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 2)
	consumerGroupsFromKong = append(consumerGroupsFromKong, page2...)

	assert.True(compareConsumerGroups(consumerGroups, consumerGroupsFromKong))

	consumerGroupsFromKong, err = client.ConsumerGroups.ListAll(defaultCtx)
	require.NoError(t, err)
	assert.NotNil(consumerGroupsFromKong)
	assert.Len(consumerGroupsFromKong, 3)

	for i := 0; i < len(consumerGroups); i++ {
		require.NoError(t, client.ConsumerGroups.Delete(defaultCtx, consumerGroups[i].ID))
	}
}

func TestConsumerGroupsRLAService(t *testing.T) {
	RunWhenEnterprise(t, ">=2.7.0", RequiredFeatures{})
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	assert.NotNil(client)

	// create consumer
	consumer := &Consumer{
		Username: String("foo"),
	}
	createdConsumer, err := client.Consumers.Create(defaultCtx, consumer)
	require.NoError(t, err)
	assert.NotNil(createdConsumer)

	// create consumer group
	cg := &ConsumerGroup{
		Name: String("bar"),
	}
	createdConsumerGroup, err := client.ConsumerGroups.Create(defaultCtx, cg)
	require.NoError(t, err)
	assert.NotNil(createdConsumerGroup)

	// add consumer to consumer group
	var response *ConsumerGroupObject
	response, err = client.ConsumerGroupConsumers.Create(defaultCtx, cg.Name, consumer.Username)
	require.NoError(t, err)
	assert.NotNil(response)
	assert.Equal(response.Consumers[0].ID, createdConsumer.ID)
	assert.Equal(response.Consumers[0].Username, createdConsumer.Username)

	// create RLA plugin with consumer groups support
	rlaPlugin := &Plugin{
		Name: String("rate-limiting-advanced"),
		Config: Configuration{
			"limit":                   []interface{}{5},
			"window_size":             []interface{}{30},
			"window_type":             "sliding",
			"retry_after_jitter_max":  float64(0),
			"enforce_consumer_groups": true,
			"consumer_groups":         []string{"bar"},
			"sync_rate":               float64(1),
			"strategy":                "cluster",
		},
	}

	createdPlugin, err := client.Plugins.Create(defaultCtx, rlaPlugin)
	require.NoError(t, err)
	assert.NotNil(createdPlugin)

	// add override to bar consumer group
	var override *ConsumerGroupRLA
	config := map[string]Configuration{
		"config": {
			"limit":                  []interface{}{float64(10)},
			"window_size":            []interface{}{float64(10)},
			"retry_after_jitter_max": float64(1),
		},
	}
	override, err = client.ConsumerGroups.UpdateRateLimitingAdvancedPlugin(
		defaultCtx, cg.Name, config)
	require.NoError(t, err)
	assert.NotNil(override)

	assert.Equal(override.Config, config["config"])
	assert.Equal("bar", *override.ConsumerGroup)
	assert.Equal("rate-limiting-advanced", *override.Plugin)

	// clean up
	require.NoError(t, client.ConsumerGroups.Delete(defaultCtx, createdConsumerGroup.ID))
	require.NoError(t, client.Consumers.Delete(defaultCtx, createdConsumer.ID))
	require.NoError(t, client.Plugins.Delete(defaultCtx, createdPlugin.ID))
}

func compareConsumerGroups(expected, actual []*ConsumerGroup) bool {
	var expectedNames, actualNames []string
	for _, cg := range expected {
		expectedNames = append(expectedNames, *cg.Name)
	}

	for _, cg := range actual {
		actualNames = append(actualNames, *cg.Name)
	}

	return (compareSlices(expectedNames, actualNames))
}

func TestConsumerGroupGetEndpointPreGW39(t *testing.T) {
	RunWhenEnterprise(t, ">=2.7.0 <3.9.0", RequiredFeatures{})
	assert := assert.New(t)
	require := require.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	// create consumer
	consumer := &Consumer{
		Username: String("foo"),
	}
	createdConsumer, err := client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	assert.NotNil(createdConsumer)
	t.Cleanup(func() {
		if createdConsumer != nil {
			require.NoError(client.Consumers.Delete(defaultCtx, createdConsumer.ID))
		}
	})

	// create a consumer-group
	cg := &ConsumerGroup{
		Name: String("foo-group"),
	}

	createdConsumerGroup, err := client.ConsumerGroups.Create(defaultCtx, cg)
	require.NoError(err)
	assert.NotNil(createdConsumerGroup)
	t.Cleanup(func() {
		if createdConsumerGroup != nil {
			require.NoError(client.ConsumerGroups.Delete(defaultCtx, createdConsumerGroup.ID))
		}
	})

	// add consumer to consumer group
	var response *ConsumerGroupObject
	response, err = client.ConsumerGroupConsumers.Create(defaultCtx, cg.Name, consumer.Username)
	require.NoError(err)
	assert.NotNil(response)
	assert.Equal(response.Consumers[0].ID, createdConsumer.ID)
	assert.Equal(response.Consumers[0].Username, createdConsumer.Username)
	assert.Equal(response.ConsumerGroup.ID, createdConsumerGroup.ID)

	// Check get endpoint
	consumerGroupFromKong, err := client.ConsumerGroups.Get(defaultCtx, createdConsumerGroup.ID)
	require.NoError(err)
	assert.NotNil(consumerGroupFromKong)
	assert.Equal(consumerGroupFromKong.ConsumerGroup.ID, createdConsumerGroup.ID)
	assert.NotNil(consumerGroupFromKong.Consumers)
	assert.Len(consumerGroupFromKong.Consumers, 1)
	assert.Equal(consumerGroupFromKong.Consumers[0].ID, createdConsumer.ID)
}

func TestConsumerGroupGetEndpointPostGW39(t *testing.T) {
	RunWhenEnterprise(t, ">=3.9.0", RequiredFeatures{})
	assert := assert.New(t)
	require := require.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	// create consumer
	consumer := &Consumer{
		Username: String("foo"),
	}
	createdConsumer, err := client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	assert.NotNil(createdConsumer)
	t.Cleanup(func() {
		if createdConsumer != nil {
			require.NoError(client.Consumers.Delete(defaultCtx, createdConsumer.ID))
		}
	})

	// create a consumer-group
	cg := &ConsumerGroup{
		Name: String("foo-group"),
	}

	createdConsumerGroup, err := client.ConsumerGroups.Create(defaultCtx, cg)
	require.NoError(err)
	assert.NotNil(createdConsumerGroup)
	t.Cleanup(func() {
		if createdConsumerGroup != nil {
			require.NoError(client.ConsumerGroups.Delete(defaultCtx, createdConsumerGroup.ID))
		}
	})

	// add consumer to consumer group
	var response *ConsumerGroupObject
	response, err = client.ConsumerGroupConsumers.Create(defaultCtx, cg.Name, consumer.Username)
	require.NoError(err)
	assert.NotNil(response)
	assert.Equal(response.Consumers[0].ID, createdConsumer.ID)
	assert.Equal(response.Consumers[0].Username, createdConsumer.Username)
	assert.Equal(response.ConsumerGroup.ID, createdConsumerGroup.ID)

	t.Run("Get", func(_ *testing.T) {
		consumerGroupFromKong, err := client.ConsumerGroups.Get(defaultCtx, createdConsumerGroup.ID)
		require.NoError(err)
		assert.NotNil(consumerGroupFromKong)
		assert.Equal(consumerGroupFromKong.ConsumerGroup.ID, createdConsumerGroup.ID)
		assert.NotNil(consumerGroupFromKong.Consumers)
		assert.Len(consumerGroupFromKong.Consumers, 1, "consumers are listed")
	})

	t.Run("GetWithNoConsumers", func(_ *testing.T) {
		consumerGroupFromKong, err := client.ConsumerGroups.GetWithNoConsumers(defaultCtx, createdConsumerGroup.ID)
		require.NoError(err)
		assert.NotNil(consumerGroupFromKong)
		assert.Equal(consumerGroupFromKong.ConsumerGroup.ID, createdConsumerGroup.ID)
		assert.Nil(consumerGroupFromKong.Consumers, "consumers should not be listed")
	})
}

func TestConsumerGroupPluginWithPartials(t *testing.T) {
	RunWhenEnterprise(t, ">=3.10.0", RequiredFeatures{})
	require := require.New(t)
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	// create a consumer group
	cg := &ConsumerGroup{
		Name: String("test-group-with-partials"),
	}
	createdConsumerGroup, err := client.ConsumerGroups.Create(defaultCtx, cg)
	require.NoError(err)
	require.NotNil(createdConsumerGroup)
	t.Cleanup(func() {
		if createdConsumerGroup != nil {
			require.NoError(client.ConsumerGroups.Delete(defaultCtx, createdConsumerGroup.ID))
		}
	})

	// create a partial
	partial := &Partial{
		Name: String("test-redis-partial"),
		Type: String("redis-ee"),
		Config: Configuration{
			"send_timeout":    2001,
			"read_timeout":    3001,
			"connect_timeout": 4001,
		},
	}
	createdPartial, err := client.Partials.Create(defaultCtx, partial)
	require.NoError(err)
	require.NotNil(createdPartial)
	t.Cleanup(func() {
		require.NoError(client.Partials.Delete(defaultCtx, createdPartial.ID))
	})

	// create a consumer group plugin with partial
	plugin := &Plugin{
		Name: String("rate-limiting-advanced"),
		Config: Configuration{
			"limit":       []interface{}{50},
			"window_size": []interface{}{30},
		},
		Partials: []*PartialLink{
			{
				Partial: &Partial{
					ID: createdPartial.ID,
				},
			},
		},
	}

	createdPlugin, err := client.Plugins.CreateForConsumerGroup(defaultCtx, createdConsumerGroup.ID, plugin)
	require.NoError(err)
	require.NotNil(createdPlugin)
	t.Cleanup(func() {
		if createdPlugin != nil {
			require.NoError(client.Plugins.Delete(defaultCtx, createdPlugin.ID))
		}
	})

	// verify the plugin has the partial link
	assert.NotNil(createdPlugin.Partials)
	assert.Len(createdPlugin.Partials, 1)
	assert.Equal(createdPartial.ID, createdPlugin.Partials[0].ID)
	assert.Equal(String("config.redis"), createdPlugin.Partials[0].Path)

	// verify the partial config was applied
	redisConfig, ok := createdPlugin.Config["redis"].(map[string]interface{})
	require.True(ok, "redis config should be present")
	assert.InEpsilon(2001, redisConfig["send_timeout"], 0.1)
	assert.InEpsilon(3001, redisConfig["read_timeout"], 0.1)
	assert.InEpsilon(4001, redisConfig["connect_timeout"], 0.1)

	// get the plugin and verify partials are returned
	fetchedPlugin, err := client.Plugins.Get(defaultCtx, createdPlugin.ID)
	require.NoError(err)
	assert.NotNil(fetchedPlugin)
	assert.NotNil(fetchedPlugin.Partials)
	assert.Len(fetchedPlugin.Partials, 1)
	assert.Equal(createdPartial.ID, fetchedPlugin.Partials[0].ID)
}
