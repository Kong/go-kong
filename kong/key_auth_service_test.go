package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyAuthCreate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	keyAuth, err := client.KeyAuths.Create(defaultCtx,
		String("foo"), nil)
	require.Error(err)
	assert.Nil(keyAuth)

	keyAuth = &KeyAuth{}
	keyAuth, err = client.KeyAuths.Create(defaultCtx, String(""),
		keyAuth)
	require.Error(err)
	assert.Nil(keyAuth)

	keyAuth, err = client.KeyAuths.Create(defaultCtx,
		String("does-not-exist"), keyAuth)
	require.Error(err)
	assert.Nil(keyAuth)

	// consumer for the key-auth
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	keyAuth = &KeyAuth{}
	createdKeyAuth, err := client.KeyAuths.Create(defaultCtx,
		consumer.ID, keyAuth)
	require.NoError(err)
	assert.NotNil(createdKeyAuth)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestKeyAuthCreateWithID(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	keyAuth := &KeyAuth{
		ID:  String(uuid),
		Key: String("my-apikey"),
	}

	// consumer for the key-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdKeyAuth, err := client.KeyAuths.Create(defaultCtx,
		consumer.ID, keyAuth)
	require.NoError(err)
	assert.NotNil(createdKeyAuth)

	assert.Equal(uuid, *createdKeyAuth.ID)
	assert.Equal("my-apikey", *createdKeyAuth.Key)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestKeyAuthGet(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	keyAuth := &KeyAuth{
		ID:  String(uuid),
		Key: String("my-apikey"),
	}

	// consumer for the key-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdKeyAuth, err := client.KeyAuths.Create(defaultCtx,
		consumer.ID, keyAuth)
	require.NoError(err)
	assert.NotNil(createdKeyAuth)

	searchKeyAuth, err := client.KeyAuths.Get(defaultCtx,
		consumer.ID, keyAuth.ID)
	require.NoError(err)
	assert.Equal("my-apikey", *searchKeyAuth.Key)

	searchKeyAuth, err = client.KeyAuths.Get(defaultCtx,
		consumer.ID, keyAuth.Key)
	require.NoError(err)
	assert.Equal("my-apikey", *searchKeyAuth.Key)

	searchKeyAuth, err = client.KeyAuths.Get(defaultCtx, consumer.ID,
		String("does-not-exists"))
	assert.Nil(searchKeyAuth)
	require.Error(err)

	searchKeyAuth, err = client.KeyAuths.Get(defaultCtx,
		consumer.ID, String(""))
	assert.Nil(searchKeyAuth)
	require.Error(err)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestKeyAuthGetByID(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	uuid := uuid.NewString()
	keyAuth := &KeyAuth{
		ID:  String(uuid),
		Key: String("my-apikey"),
	}

	// consumer for the key-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdKeyAuth, err := client.KeyAuths.Create(defaultCtx,
		consumer.ID, keyAuth)
	require.NoError(err)
	require.NotNil(createdKeyAuth)

	T.Cleanup(func() {
		require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
	})

	T.Run("successful key-auth retrieval by ID", func(_ *testing.T) {
		searchKeyAuth, err := client.KeyAuths.GetByID(defaultCtx, keyAuth.ID)
		require.NoError(err)
		require.NotNil(searchKeyAuth)
		require.Equal("my-apikey", *searchKeyAuth.Key)
	})

	T.Run("unsuccessful key-auth retrieval by ID", func(_ *testing.T) {
		searchKeyAuth, err := client.KeyAuths.GetByID(defaultCtx, String("does-not-exist"))
		require.Nil(searchKeyAuth)
		require.Error(err)
	})

	T.Run("unsuccessful key-auth retrieval using empty string", func(_ *testing.T) {
		searchKeyAuth, err := client.KeyAuths.GetByID(defaultCtx, String(""))
		require.Nil(searchKeyAuth)
		require.Error(err)
	})
}

func TestKeyAuthUpdate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	keyAuth := &KeyAuth{
		ID:  String(uuid),
		Key: String("my-apikey"),
	}

	// consumer for the key-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdKeyAuth, err := client.KeyAuths.Create(defaultCtx,
		consumer.ID, keyAuth)
	require.NoError(err)
	assert.NotNil(createdKeyAuth)

	searchKeyAuth, err := client.KeyAuths.Get(defaultCtx,
		consumer.ID, keyAuth.ID)
	require.NoError(err)
	assert.Equal("my-apikey", *searchKeyAuth.Key)

	keyAuth.Key = String("my-new-apikey")
	updatedKeyAuth, err := client.KeyAuths.Update(defaultCtx,
		consumer.ID, keyAuth)
	require.NoError(err)
	assert.NotNil(updatedKeyAuth)
	assert.Equal("my-new-apikey", *updatedKeyAuth.Key)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestKeyAuthDelete(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	keyAuth := &KeyAuth{
		ID:  String(uuid),
		Key: String("my-apikey"),
	}

	// consumer for the key-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdKeyAuth, err := client.KeyAuths.Create(defaultCtx,
		consumer.ID, keyAuth)
	require.NoError(err)
	assert.NotNil(createdKeyAuth)

	err = client.KeyAuths.Delete(defaultCtx, consumer.ID, keyAuth.Key)
	require.NoError(err)

	searchKeyAuth, err := client.KeyAuths.Get(defaultCtx,
		consumer.ID, keyAuth.ID)
	require.Error(err)
	assert.Nil(searchKeyAuth)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestKeyAuthListMethods(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	// consumer for the key-auth:
	consumer1 := &Consumer{
		Username: String("foo"),
	}

	consumer1, err = client.Consumers.Create(defaultCtx, consumer1)
	require.NoError(err)
	require.NotNil(consumer1)

	consumer2 := &Consumer{
		Username: String("bar"),
	}

	consumer2, err = client.Consumers.Create(defaultCtx, consumer2)
	require.NoError(err)
	require.NotNil(consumer2)

	// fixtures
	keyAuths := []*KeyAuth{
		{
			Key:      String("key11"),
			Consumer: consumer1,
		},
		{
			Key:      String("key12"),
			Consumer: consumer1,
		},
		{
			Key:      String("key21"),
			Consumer: consumer2,
		},
		{
			Key:      String("key22"),
			Consumer: consumer2,
		},
	}

	// create fixturs
	for i := 0; i < len(keyAuths); i++ {
		keyAuth, err := client.KeyAuths.Create(defaultCtx,
			keyAuths[i].Consumer.ID, keyAuths[i])
		require.NoError(err)
		assert.NotNil(keyAuth)
		keyAuths[i] = keyAuth
	}

	keyAuthsFromKong, next, err := client.KeyAuths.List(defaultCtx, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(keyAuthsFromKong)
	assert.Len(keyAuthsFromKong, 4)

	// first page
	page1, next, err := client.KeyAuths.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)

	// last page
	next.Size = 3
	page2, next, err := client.KeyAuths.List(defaultCtx, next)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 3)

	keyAuthsForConsumer, next, err := client.KeyAuths.ListForConsumer(defaultCtx, consumer1.ID, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(keyAuthsForConsumer)
	assert.Len(keyAuthsForConsumer, 2)

	keyAuths, err = client.KeyAuths.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(keyAuths)
	assert.Len(keyAuths, 4)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer1.ID))
	require.NoError(client.Consumers.Delete(defaultCtx, consumer2.ID))
}

func TestKeyAuthCreateWithTTL(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=1.4.0")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	keyAuth := &KeyAuth{
		TTL: Int(10),
		Key: String("my-apikey"),
	}

	// consumer for the key-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdKeyAuth, err := client.KeyAuths.Create(defaultCtx,
		consumer.ID, keyAuth)
	require.NoError(err)
	assert.NotNil(createdKeyAuth)

	assert.Less(*createdKeyAuth.TTL, 10)
	assert.Equal("my-apikey", *createdKeyAuth.Key)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}
