package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHMACAuthCreate(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	hmacAuth, err := client.HMACAuths.Create(defaultCtx,
		String("foo"), nil)
	assert.NotNil(err)
	assert.Nil(hmacAuth)

	hmacAuth = &HMACAuth{}
	hmacAuth, err = client.HMACAuths.Create(defaultCtx, String(""),
		hmacAuth)
	assert.NotNil(err)
	assert.Nil(hmacAuth)

	// consumer for the hmac-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	assert.NoError(err)
	require.NotNil(consumer)

	// no username is specified
	hmacAuth = &HMACAuth{}
	hmacAuth, err = client.HMACAuths.Create(defaultCtx, consumer.ID, hmacAuth)
	assert.NotNil(err)
	assert.Nil(hmacAuth)

	hmacAuth = &HMACAuth{
		Username: String("foo"),
	}
	hmacAuth, err = client.HMACAuths.Create(defaultCtx, consumer.ID, hmacAuth)
	assert.NoError(err)
	assert.NotNil(hmacAuth)
	assert.NotNil(hmacAuth.Secret)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestHMACAuthCreateWithID(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(consumer)

	createdHMACAuth, err := client.HMACAuths.Create(defaultCtx, consumer.ID,
		hmacAuth)
	assert.NoError(err)
	assert.NotNil(createdHMACAuth)

	assert.Equal(uuid, *createdHMACAuth.ID)
	assert.Equal("my-username", *createdHMACAuth.Username)
	assert.Equal("my-secret", *createdHMACAuth.Secret)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestHMACAuthGet(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(consumer)

	createdHMACAuth, err := client.HMACAuths.Create(defaultCtx,
		consumer.ID, hmacAuth)
	assert.NoError(err)
	assert.NotNil(createdHMACAuth)

	hmacAuth, err = client.HMACAuths.Get(defaultCtx, consumer.ID, hmacAuth.ID)
	assert.NoError(err)
	assert.Equal("my-username", *hmacAuth.Username)

	hmacAuth, err = client.HMACAuths.Get(defaultCtx, consumer.ID,
		hmacAuth.Username)
	assert.NoError(err)
	assert.Equal("my-username", *hmacAuth.Username)

	hmacAuth, err = client.HMACAuths.Get(defaultCtx, consumer.ID,
		String("does-not-exists"))
	assert.Nil(hmacAuth)
	assert.NotNil(err)

	hmacAuth, err = client.HMACAuths.Get(defaultCtx, consumer.ID, String(""))
	assert.Nil(hmacAuth)
	assert.NotNil(err)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestHMACAuthUpdate(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(consumer)

	createdHMACAuth, err := client.HMACAuths.Create(defaultCtx,
		consumer.ID, hmacAuth)
	assert.NoError(err)
	assert.NotNil(createdHMACAuth)

	hmacAuth, err = client.HMACAuths.Get(defaultCtx, consumer.ID, hmacAuth.ID)
	assert.NoError(err)
	assert.Equal("my-username", *hmacAuth.Username)

	hmacAuth.Username = String("my-new-username")
	hmacAuth.Secret = String("my-new-secret")
	updatedHMACAuth, err := client.HMACAuths.Update(defaultCtx,
		consumer.ID, hmacAuth)
	assert.NoError(err)
	assert.NotNil(updatedHMACAuth)
	assert.Equal("my-new-secret", *updatedHMACAuth.Secret)
	assert.Equal("my-new-username", *updatedHMACAuth.Username)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestHMACAuthDelete(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(consumer)

	createdHMACAuth, err := client.HMACAuths.Create(defaultCtx,
		consumer.ID, hmacAuth)
	assert.NoError(err)
	assert.NotNil(createdHMACAuth)

	err = client.HMACAuths.Delete(defaultCtx, consumer.ID, hmacAuth.Username)
	assert.NoError(err)

	hmacAuth, err = client.HMACAuths.Get(defaultCtx,
		consumer.ID, hmacAuth.Username)
	assert.NotNil(err)
	assert.Nil(hmacAuth)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestHMACAuthListMethods(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	// consumer for the hmac-auth:
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
		assert.NoError(err)
		assert.NotNil(hmacAuth)
		hmacAuths[i] = hmacAuth
	}

	hmacAuthsFromKong, next, err := client.HMACAuths.List(defaultCtx, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(hmacAuthsFromKong)
	assert.Equal(4, len(hmacAuthsFromKong))

	// first page
	page1, next, err := client.HMACAuths.List(defaultCtx, &ListOpt{Size: 1})
	assert.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))

	// last page
	next.Size = 3
	page2, next, err := client.HMACAuths.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Equal(3, len(page2))

	hmacAuthsForConsumer, next, err := client.HMACAuths.ListForConsumer(defaultCtx, consumer1.ID, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(hmacAuthsForConsumer)
	assert.Equal(2, len(hmacAuthsForConsumer))

	hmacAuths, err = client.HMACAuths.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(hmacAuths)
	assert.Equal(4, len(hmacAuths))

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer1.ID))
	assert.NoError(client.Consumers.Delete(defaultCtx, consumer2.ID))
}
