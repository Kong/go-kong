package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumerGroupConsumersService(t *testing.T) {
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

	// delete consumer
	err = client.Consumers.Delete(defaultCtx, createdConsumer.ID)
	require.NoError(t, err)

	// make sure adding a not-existing consumer to a group fails
	_, err = client.ConsumerGroupConsumers.Create(defaultCtx, cg.Name, consumer.Username)
	assert.NotNil(err)

	// delete consumer group
	err = client.ConsumerGroups.Delete(defaultCtx, createdConsumerGroup.ID)
	require.NoError(t, err)
}

func TestConsumerGroupConsumersListEndpoint(t *testing.T) {
	RunWhenEnterprise(t, ">=2.7.0", RequiredFeatures{})
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	assert.NotNil(client)

	// create consumer group
	cg := &ConsumerGroup{
		Name: String("bar"),
	}
	createdConsumerGroup, err := client.ConsumerGroups.Create(defaultCtx, cg)
	require.NoError(t, err)
	assert.NotNil(createdConsumerGroup)

	// create consumers
	consumers := []*Consumer{
		{
			Username: String("foo1"),
		},
		{
			Username: String("foo2"),
		},
		{
			Username: String("foo3"),
		},
	}
	for i := 0; i < len(consumers); i++ {
		consumer, err := client.Consumers.Create(defaultCtx, consumers[i])
		require.NoError(t, err)
		assert.NotNil(consumer)
		consumers[i] = consumer

		// add consumers to consumer group
		var response *ConsumerGroupObject
		response, err = client.ConsumerGroupConsumers.Create(defaultCtx, cg.Name, consumer.Username)
		require.NoError(t, err)
		assert.NotNil(response)
		assert.Equal(response.Consumers[0].ID, consumer.ID)
		assert.Equal(response.Consumers[0].Username, consumer.Username)
	}

	consumerGroupConsumersFromKong, err := client.ConsumerGroupConsumers.ListAll(defaultCtx, cg.Name)
	require.NoError(t, err)
	assert.NotNil(consumerGroupConsumersFromKong)
	assert.Equal(3, len(consumerGroupConsumersFromKong))

	// check if we see all consumer groups
	assert.True(compareConsumers(consumers, consumerGroupConsumersFromKong))

	for i := 0; i < len(consumerGroupConsumersFromKong); i++ {
		assert.NoError(client.Consumers.Delete(defaultCtx, consumers[i].ID))
	}

	assert.NoError(client.ConsumerGroups.Delete(defaultCtx, cg.Name))
}
