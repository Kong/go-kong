package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHMACAuthCreate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	hmacAuth, err := client.HMACAuths.Create(defaultCtx,
		String("foo"), nil)
	require.Error(err)
	assert.Nil(hmacAuth)

	hmacAuth = &HMACAuth{}
	hmacAuth, err = client.HMACAuths.Create(defaultCtx, String(""),
		hmacAuth)
	require.Error(err)
	assert.Nil(hmacAuth)

	// consumer for the hmac-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	// no username is specified
	hmacAuth = &HMACAuth{}
	hmacAuth, err = client.HMACAuths.Create(defaultCtx, consumer.ID, hmacAuth)
	require.Error(err)
	assert.Nil(hmacAuth)

	hmacAuth = &HMACAuth{
		Username: String("foo"),
	}
	hmacAuth, err = client.HMACAuths.Create(defaultCtx, consumer.ID, hmacAuth)
	require.NoError(err)
	assert.NotNil(hmacAuth)
	assert.NotNil(hmacAuth.Secret)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestHMACAuthCreateWithID(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	hmacAuth := &HMACAuth{
		ID:       String(uuid),
		Username: String("my-username"),
		Secret:   String("my-secret"),
	}

	// consumer for the hmac-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdHMACAuth, err := client.HMACAuths.Create(defaultCtx, consumer.ID,
		hmacAuth)
	require.NoError(err)
	assert.NotNil(createdHMACAuth)

	assert.Equal(uuid, *createdHMACAuth.ID)
	assert.Equal("my-username", *createdHMACAuth.Username)
	assert.Equal("my-secret", *createdHMACAuth.Secret)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestHMACAuthGet(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	hmacAuth := &HMACAuth{
		ID:       String(uuid),
		Username: String("my-username"),
	}

	// consumer for the hmac-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdHMACAuth, err := client.HMACAuths.Create(defaultCtx,
		consumer.ID, hmacAuth)
	require.NoError(err)
	assert.NotNil(createdHMACAuth)

	hmacAuth, err = client.HMACAuths.Get(defaultCtx, consumer.ID, hmacAuth.ID)
	require.NoError(err)
	assert.Equal("my-username", *hmacAuth.Username)

	hmacAuth, err = client.HMACAuths.Get(defaultCtx, consumer.ID,
		hmacAuth.Username)
	require.NoError(err)
	assert.Equal("my-username", *hmacAuth.Username)

	hmacAuth, err = client.HMACAuths.Get(defaultCtx, consumer.ID,
		String("does-not-exists"))
	assert.Nil(hmacAuth)
	require.Error(err)

	hmacAuth, err = client.HMACAuths.Get(defaultCtx, consumer.ID, String(""))
	assert.Nil(hmacAuth)
	require.Error(err)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestHMACAuthGetByID(t *testing.T) {
	RunWhenDBMode(t, "postgres")

	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	uuid := uuid.NewString()
	hmacAuth := &HMACAuth{
		ID:       String(uuid),
		Username: String("my-username"),
	}

	// consumer for the hmac-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(t, err)
	require.NotNil(t, consumer)

	createdHMACAuth, err := client.HMACAuths.Create(defaultCtx,
		consumer.ID, hmacAuth)
	require.NoError(t, err)
	require.NotNil(t, createdHMACAuth)

	t.Cleanup(func() {
		require.NoError(t, client.Consumers.Delete(defaultCtx, consumer.ID))
	})

	t.Run("successful hmac-auth retrieval by ID", func(_ *testing.T) {
		hmacAuth, err = client.HMACAuths.GetByID(defaultCtx, hmacAuth.ID)
		require.NoError(t, err)
		require.NotNil(t, hmacAuth)
		require.Equal(t, "my-username", *hmacAuth.Username)
	})

	t.Run("unsuccessful hmac-auth retrieval by ID", func(_ *testing.T) {
		hmacAuth, err = client.HMACAuths.GetByID(defaultCtx, String("does-not-exist"))
		require.Nil(t, hmacAuth)
		require.Error(t, err)
	})

	t.Run("unsuccessful hmac-auth retrieval when empty string is passed", func(_ *testing.T) {
		hmacAuth, err = client.HMACAuths.GetByID(defaultCtx, String(""))
		require.Nil(t, hmacAuth)
		require.Error(t, err)
	})
}

func TestHMACAuthUpdate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	hmacAuth := &HMACAuth{
		ID:       String(uuid),
		Username: String("my-username"),
		Secret:   String("secret"),
	}

	// consumer for the hmac-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdHMACAuth, err := client.HMACAuths.Create(defaultCtx,
		consumer.ID, hmacAuth)
	require.NoError(err)
	assert.NotNil(createdHMACAuth)

	hmacAuth, err = client.HMACAuths.Get(defaultCtx, consumer.ID, hmacAuth.ID)
	require.NoError(err)
	assert.Equal("my-username", *hmacAuth.Username)

	hmacAuth.Username = String("my-new-username")
	hmacAuth.Secret = String("my-new-secret")
	updatedHMACAuth, err := client.HMACAuths.Update(defaultCtx,
		consumer.ID, hmacAuth)
	require.NoError(err)
	assert.NotNil(updatedHMACAuth)
	assert.Equal("my-new-secret", *updatedHMACAuth.Secret)
	assert.Equal("my-new-username", *updatedHMACAuth.Username)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestHMACAuthDelete(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	hmacAuth := &HMACAuth{
		ID:       String(uuid),
		Username: String("my-username"),
	}

	// consumer for the hmac-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdHMACAuth, err := client.HMACAuths.Create(defaultCtx,
		consumer.ID, hmacAuth)
	require.NoError(err)
	assert.NotNil(createdHMACAuth)

	err = client.HMACAuths.Delete(defaultCtx, consumer.ID, hmacAuth.Username)
	require.NoError(err)

	hmacAuth, err = client.HMACAuths.Get(defaultCtx,
		consumer.ID, hmacAuth.Username)
	require.Error(err)
	assert.Nil(hmacAuth)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestHMACAuthListMethods(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	// consumer for the hmac-auth:
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
	hmacAuths := []*HMACAuth{
		{
			Username: String("username11"),
			Consumer: consumer1,
		},
		{
			Username: String("username12"),
			Consumer: consumer1,
		},
		{
			Username: String("username21"),
			Consumer: consumer2,
		},
		{
			Username: String("username22"),
			Consumer: consumer2,
		},
	}

	// create fixturs
	for i := 0; i < len(hmacAuths); i++ {
		hmacAuth, err := client.HMACAuths.Create(defaultCtx,
			hmacAuths[i].Consumer.ID, hmacAuths[i])
		require.NoError(err)
		assert.NotNil(hmacAuth)
		hmacAuths[i] = hmacAuth
	}

	hmacAuthsFromKong, next, err := client.HMACAuths.List(defaultCtx, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(hmacAuthsFromKong)
	assert.Len(hmacAuthsFromKong, 4)

	// first page
	page1, next, err := client.HMACAuths.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)

	// last page
	next.Size = 3
	page2, next, err := client.HMACAuths.List(defaultCtx, next)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 3)

	hmacAuthsForConsumer, next, err := client.HMACAuths.ListForConsumer(defaultCtx, consumer1.ID, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(hmacAuthsForConsumer)
	assert.Len(hmacAuthsForConsumer, 2)

	hmacAuths, err = client.HMACAuths.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(hmacAuths)
	assert.Len(hmacAuths, 4)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer1.ID))
	require.NoError(client.Consumers.Delete(defaultCtx, consumer2.ID))
}
