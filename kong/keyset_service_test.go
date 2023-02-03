package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeySetService(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=3.1.0")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	keySet := &KeySet{
		Name: String("foo"),
	}

	createdKeySet, err := client.KeySets.Create(defaultCtx, keySet)
	assert.NoError(err)
	require.NotNil(createdKeySet)

	keySet, err = client.KeySets.Get(defaultCtx, createdKeySet.ID)
	assert.NoError(err)
	require.NotNil(keySet)

	keySet.Name = String("bar")
	keySet, err = client.KeySets.Update(defaultCtx, keySet)
	assert.NoError(err)
	require.NotNil(keySet)
	assert.Equal("bar", *keySet.Name)

	err = client.KeySets.Delete(defaultCtx, createdKeySet.ID)
	assert.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	keySet = &KeySet{
		Name: String("foo"),
		ID:   String(id),
	}

	createdKeySet, err = client.KeySets.Create(defaultCtx, keySet)
	assert.NoError(err)
	assert.NotNil(createdKeySet)
	assert.Equal(id, *createdKeySet.ID)

	err = client.KeySets.Delete(defaultCtx, createdKeySet.ID)
	assert.NoError(err)
}

func TestKeySetWithTags(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=3.1.0")

	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	keySet := &KeySet{
		Name: String("foo"),
		Tags: StringSlice("tag1", "tag2"),
	}

	createdKeySet, err := client.KeySets.Create(defaultCtx, keySet)
	require.NoError(err)
	require.NotNil(createdKeySet)
	require.Equal(StringSlice("tag1", "tag2"), createdKeySet.Tags)

	err = client.KeySets.Delete(defaultCtx, createdKeySet.ID)
	require.NoError(err)
}

func TestKeySetListWithTags(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=3.1.0")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	// fixtures
	keySets := []*KeySet{
		{
			Name: String("user1"),
			Tags: StringSlice("tag1", "tag2"),
		},
		{
			Name: String("user2"),
			Tags: StringSlice("tag2", "tag3"),
		},
		{
			Name: String("user3"),
			Tags: StringSlice("tag1", "tag3"),
		},
		{
			Name: String("user4"),
			Tags: StringSlice("tag1", "tag2"),
		},
		{
			Name: String("user5"),
			Tags: StringSlice("tag2", "tag3"),
		},
		{
			Name: String("user6"),
			Tags: StringSlice("tag1", "tag3"),
		},
	}

	// create fixtures
	for i := 0; i < len(keySets); i++ {
		keySet, err := client.KeySets.Create(defaultCtx, keySets[i])
		assert.NoError(err)
		require.NotNil(keySet)
		keySets[i] = keySet
	}

	keySetsFromKong, next, err := client.KeySets.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag1"),
	})
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(4, len(keySetsFromKong))

	keySetsFromKong, next, err = client.KeySets.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag2"),
	})
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(4, len(keySetsFromKong))

	keySetsFromKong, next, err = client.KeySets.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag1", "tag2"),
	})
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(6, len(keySetsFromKong))

	keySetsFromKong, next, err = client.KeySets.List(defaultCtx, &ListOpt{
		Tags:         StringSlice("tag1", "tag2"),
		MatchAllTags: true,
	})
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(2, len(keySetsFromKong))

	keySetsFromKong, next, err = client.KeySets.List(defaultCtx, &ListOpt{
		Tags: StringSlice("tag1", "tag2"),
		Size: 3,
	})
	assert.NoError(err)
	assert.NotNil(next)
	assert.Equal(3, len(keySetsFromKong))

	keySetsFromKong, next, err = client.KeySets.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(3, len(keySetsFromKong))

	keySetsFromKong, next, err = client.KeySets.List(defaultCtx, &ListOpt{
		Tags:         StringSlice("tag1", "tag2"),
		MatchAllTags: true,
		Size:         1,
	})
	assert.NoError(err)
	assert.NotNil(next)
	assert.Equal(1, len(keySetsFromKong))

	keySetsFromKong, next, err = client.KeySets.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.Equal(1, len(keySetsFromKong))

	for i := 0; i < len(keySets); i++ {
		assert.NoError(client.KeySets.Delete(defaultCtx, keySets[i].Name))
	}
}
