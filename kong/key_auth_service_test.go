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
	assert.NoError(err)
	assert.NotNil(client)

	keyAuth, err := client.KeyAuths.Create(defaultCtx,
		String("foo"), nil)
	assert.NotNil(err)
	assert.Nil(keyAuth)

	keyAuth = &KeyAuth{}
	keyAuth, err = client.KeyAuths.Create(defaultCtx, String(""),
		keyAuth)
	assert.NotNil(err)
	assert.Nil(keyAuth)

	keyAuth, err = client.KeyAuths.Create(defaultCtx,
		String("does-not-exist"), keyAuth)
	assert.NotNil(err)
	assert.Nil(keyAuth)

	// consumer for the key-auth
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	assert.NoError(err)
	require.NotNil(consumer)

	keyAuth = &KeyAuth{}
	createdKeyAuth, err := client.KeyAuths.Create(defaultCtx,
		consumer.ID, keyAuth)
	assert.NoError(err)
	assert.NotNil(createdKeyAuth)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestKeyAuthCreateWithID(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(consumer)

	createdKeyAuth, err := client.KeyAuths.Create(defaultCtx,
		consumer.ID, keyAuth)
	assert.NoError(err)
	assert.NotNil(createdKeyAuth)

	assert.Equal(uuid, *createdKeyAuth.ID)
	assert.Equal("my-apikey", *createdKeyAuth.Key)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestKeyAuthGet(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(consumer)

	createdKeyAuth, err := client.KeyAuths.Create(defaultCtx,
		consumer.ID, keyAuth)
	assert.NoError(err)
	assert.NotNil(createdKeyAuth)

	searchKeyAuth, err := client.KeyAuths.Get(defaultCtx,
		consumer.ID, keyAuth.ID)
	assert.NoError(err)
	assert.Equal("my-apikey", *searchKeyAuth.Key)

	searchKeyAuth, err = client.KeyAuths.Get(defaultCtx,
		consumer.ID, keyAuth.Key)
	assert.NoError(err)
	assert.Equal("my-apikey", *searchKeyAuth.Key)

	searchKeyAuth, err = client.KeyAuths.Get(defaultCtx, consumer.ID,
		String("does-not-exists"))
	assert.Nil(searchKeyAuth)
	assert.NotNil(err)

	searchKeyAuth, err = client.KeyAuths.Get(defaultCtx,
		consumer.ID, String(""))
	assert.Nil(searchKeyAuth)
	assert.NotNil(err)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestKeyAuthUpdate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(consumer)

	createdKeyAuth, err := client.KeyAuths.Create(defaultCtx,
		consumer.ID, keyAuth)
	assert.NoError(err)
	assert.NotNil(createdKeyAuth)

	searchKeyAuth, err := client.KeyAuths.Get(defaultCtx,
		consumer.ID, keyAuth.ID)
	assert.NoError(err)
	assert.Equal("my-apikey", *searchKeyAuth.Key)

	keyAuth.Key = String("my-new-apikey")
	updatedKeyAuth, err := client.KeyAuths.Update(defaultCtx,
		consumer.ID, keyAuth)
	assert.NoError(err)
	assert.NotNil(updatedKeyAuth)
	assert.Equal("my-new-apikey", *updatedKeyAuth.Key)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestKeyAuthDelete(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(consumer)

	createdKeyAuth, err := client.KeyAuths.Create(defaultCtx,
		consumer.ID, keyAuth)
	assert.NoError(err)
	assert.NotNil(createdKeyAuth)

	err = client.KeyAuths.Delete(defaultCtx, consumer.ID, keyAuth.Key)
	assert.NoError(err)

	searchKeyAuth, err := client.KeyAuths.Get(defaultCtx,
		consumer.ID, keyAuth.ID)
	assert.NotNil(err)
	assert.Nil(searchKeyAuth)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestKeyAuthListMethods(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	// consumer for the key-auth:
	consumer1 := &Consumer{
		Username: String("foo"),
	}

	consumer1, err = client.Consumers.Create(defaultCtx, consumer1)
	assert.NoError(err)
	require.NotNil(consumer1)

	consumer2 := &Consumer{
		Username: String("bar"),
	}

	consumer2, err = client.Consumers.Create(defaultCtx, consumer2)
	assert.NoError(err)
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
		assert.NoError(err)
		assert.NotNil(keyAuth)
		keyAuths[i] = keyAuth
	}

	keyAuthsFromKong, next, err := client.KeyAuths.List(defaultCtx, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(keyAuthsFromKong)
	assert.Equal(4, len(keyAuthsFromKong))

	// first page
	page1, next, err := client.KeyAuths.List(defaultCtx, &ListOpt{Size: 1})
	assert.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))

	// last page
	next.Size = 3
	page2, next, err := client.KeyAuths.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Equal(3, len(page2))

	keyAuthsForConsumer, next, err := client.KeyAuths.ListForConsumer(defaultCtx, consumer1.ID, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(keyAuthsForConsumer)
	assert.Equal(2, len(keyAuthsForConsumer))

	keyAuths, err = client.KeyAuths.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(keyAuths)
	assert.Equal(4, len(keyAuths))

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer1.ID))
	assert.NoError(client.Consumers.Delete(defaultCtx, consumer2.ID))
}

func TestKeyAuthCreateWithTTL(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=1.4.0")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(consumer)

	createdKeyAuth, err := client.KeyAuths.Create(defaultCtx,
		consumer.ID, keyAuth)
	assert.NoError(err)
	assert.NotNil(createdKeyAuth)

	assert.True(*createdKeyAuth.TTL < 10)
	assert.Equal("my-apikey", *createdKeyAuth.Key)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}
