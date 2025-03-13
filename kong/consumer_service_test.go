package kong

import (
	"reflect"
	"sort"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConsumersService(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	consumer := &Consumer{
		Username: String("foo"),
		CustomID: String("custom_id_foo"),
	}

	createdConsumer, err := client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(createdConsumer)

	consumer, err = client.Consumers.Get(defaultCtx, createdConsumer.ID)
	require.NoError(err)
	require.NotNil(consumer)

	consumer, err = client.Consumers.GetByCustomID(defaultCtx,
		String("does-not-exist"))
	require.Error(err)
	assert.Nil(consumer)

	consumer, err = client.Consumers.GetByCustomID(defaultCtx,
		String("custom_id_foo"))
	require.NoError(err)
	require.NotNil(consumer)

	consumer.Username = String("bar")
	consumer, err = client.Consumers.Update(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)
	assert.Equal("bar", *consumer.Username)

	err = client.Consumers.Delete(defaultCtx, createdConsumer.ID)
	require.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	consumer = &Consumer{
		Username: String("foo"),
		ID:       String(id),
	}

	createdConsumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	assert.NotNil(createdConsumer)
	assert.Equal(id, *createdConsumer.ID)

	err = client.Consumers.Delete(defaultCtx, createdConsumer.ID)
	require.NoError(err)
}

func TestConsumerWithTags(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=1.1.0")

	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	consumer := &Consumer{
		Username: String("foo"),
		Tags:     StringSlice("tag1", "tag2"),
	}

	createdConsumer, err := client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(createdConsumer)
	require.Equal(StringSlice("tag1", "tag2"), createdConsumer.Tags)

	err = client.Consumers.Delete(defaultCtx, createdConsumer.ID)
	require.NoError(err)
}

func TestConsumerListEndpoint(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	// Enterprise tests create an admin, which affects the list endpoints in peculiar ways. although the actual
	// consumer and credential entities are hidden from the API they still affect pagination. Tests that check
	// pagination behavior cannot check the same values on community and Enterprise. As such, we just don't run this
	// check against Enterprise, as the behavior is otherwise generally the same: if the endpoint works on community,
	// it will work on Enterprise, but you'll get different pagination for the same set of consumers
	SkipWhenEnterprise(T)
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	// fixtures
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

	// create fixturs
	for i := 0; i < len(consumers); i++ {
		consumer, err := client.Consumers.Create(defaultCtx, consumers[i])
		require.NoError(err)
		require.NotNil(consumer)
		consumers[i] = consumer
	}

	consumersFromKong, next, err := client.Consumers.List(defaultCtx, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(consumersFromKong)
	assert.Len(consumersFromKong, 3)

	// check if we see all consumers
	assert.True(compareConsumers(consumers, consumersFromKong))

	// Test pagination
	consumersFromKong = []*Consumer{}

	// first page
	page1, next, err := client.Consumers.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)
	consumersFromKong = append(consumersFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.Consumers.List(defaultCtx, next)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 2)
	consumersFromKong = append(consumersFromKong, page2...)

	assert.True(compareConsumers(consumers, consumersFromKong))

	consumers, err = client.Consumers.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(consumers)
	assert.Len(consumers, 3)

	for i := 0; i < len(consumers); i++ {
		require.NoError(client.Consumers.Delete(defaultCtx, consumers[i].ID))
	}
}

func TestConsumerListWithTags(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=1.1.0")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	// fixtures
	consumers := []*Consumer{
		{
			Username: String("user1"),
			Tags:     StringSlice("tag1", "tag2"),
		},
		{
			Username: String("user2"),
			Tags:     StringSlice("tag2", "tag3"),
		},
		{
			Username: String("user3"),
			Tags:     StringSlice("tag1", "tag3"),
		},
		{
			Username: String("user4"),
			Tags:     StringSlice("tag1", "tag2"),
		},
		{
			Username: String("user5"),
			Tags:     StringSlice("tag2", "tag3"),
		},
		{
			Username: String("user6"),
			Tags:     StringSlice("tag1", "tag3"),
		},
	}

	// create fixtures
	for i := 0; i < len(consumers); i++ {
		consumer, err := client.Consumers.Create(defaultCtx, consumers[i])
		require.NoError(err)
		require.NotNil(consumer)
		consumers[i] = consumer
	}

	consumersFromKong, next, err := client.Consumers.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag1"),
	})
	require.NoError(err)
	assert.Nil(next)
	assert.Len(consumersFromKong, 4)

	consumersFromKong, next, err = client.Consumers.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag2"),
	})
	require.NoError(err)
	assert.Nil(next)
	assert.Len(consumersFromKong, 4)

	consumersFromKong, next, err = client.Consumers.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag1", "tag2"),
	})
	require.NoError(err)
	assert.Nil(next)
	assert.Len(consumersFromKong, 6)

	consumersFromKong, next, err = client.Consumers.List(defaultCtx, &ListOpt{
		Tags:         StringSlice("tag1", "tag2"),
		MatchAllTags: true,
	})
	require.NoError(err)
	assert.Nil(next)
	assert.Len(consumersFromKong, 2)

	consumersFromKong, next, err = client.Consumers.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag1", "tag2"),
		Size: 3,
	})
	require.NoError(err)
	assert.NotNil(next)
	assert.Len(consumersFromKong, 3)

	consumersFromKong, next, err = client.Consumers.List(defaultCtx, next)
	require.NoError(err)
	assert.Nil(next)
	assert.Len(consumersFromKong, 3)

	consumersFromKong, next, err = client.Consumers.List(defaultCtx, &ListOpt{
		Tags:         StringSlice("tag1", "tag2"),
		MatchAllTags: true,
		Size:         1,
	})
	require.NoError(err)
	assert.NotNil(next)
	assert.Len(consumersFromKong, 1)

	consumersFromKong, next, err = client.Consumers.List(defaultCtx, next)
	require.NoError(err)
	assert.Nil(next)
	assert.Len(consumersFromKong, 1)

	for i := 0; i < len(consumers); i++ {
		require.NoError(client.Consumers.Delete(defaultCtx, consumers[i].Username))
	}
}

func compareConsumers(expected, actual []*Consumer) bool {
	var expectedUsernames, actualUsernames []string
	for _, consumer := range expected {
		expectedUsernames = append(expectedUsernames, *consumer.Username)
	}

	for _, consumer := range actual {
		actualUsernames = append(actualUsernames, *consumer.Username)
	}

	return (compareSlices(expectedUsernames, actualUsernames))
}

func compareSlices(expected, actual []string) bool {
	sort.Strings(expected)
	sort.Strings(actual)
	return (reflect.DeepEqual(expected, actual))
}
