package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicAuthCreate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	basicAuth, err := client.BasicAuths.Create(defaultCtx,
		String("foo"), nil)
	require.Error(err)
	assert.Nil(basicAuth)

	basicAuth = &BasicAuth{}
	basicAuth, err = client.BasicAuths.Create(defaultCtx, String(""),
		basicAuth)
	require.Error(err)
	assert.Nil(basicAuth)

	basicAuth, err = client.BasicAuths.Create(defaultCtx,
		String("does-not-exist"), basicAuth)
	require.Error(err)
	assert.Nil(basicAuth)

	// consumer for the basic-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	// no username is specified
	basicAuth = &BasicAuth{}
	basicAuth, err = client.BasicAuths.Create(defaultCtx,
		consumer.ID, basicAuth)
	require.Error(err)
	assert.Nil(basicAuth)

	basicAuth = &BasicAuth{
		Username: String("foo"),
		Password: String("bar"),
	}
	basicAuth, err = client.BasicAuths.Create(defaultCtx,
		consumer.ID, basicAuth)
	require.NoError(err)
	assert.NotNil(basicAuth)
	assert.NotEmpty(*basicAuth.Password)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestBasicAuthCreateWithID(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	basicAuth := &BasicAuth{
		ID:       String(uuid),
		Username: String("my-username"),
		Password: String("my-password"),
	}

	// consumer for the basic-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdBasicAuth, err := client.BasicAuths.Create(defaultCtx, consumer.ID,
		basicAuth)
	require.NoError(err)
	assert.NotNil(createdBasicAuth)

	assert.Equal(uuid, *createdBasicAuth.ID)
	assert.Equal("my-username", *createdBasicAuth.Username)
	// password is hashed
	assert.NotEqual("my-password", *createdBasicAuth.Password)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestBasicAuthGet(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	basicAuth := &BasicAuth{
		ID:       String(uuid),
		Username: String("my-username"),
		Password: String("my-password"),
	}

	// consumer for the basic-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdBasicAuth, err := client.BasicAuths.Create(defaultCtx,
		consumer.ID, basicAuth)
	require.NoError(err)
	assert.NotNil(createdBasicAuth)

	basicAuth, err = client.BasicAuths.Get(defaultCtx,
		consumer.ID, basicAuth.ID)
	require.NoError(err)
	assert.Equal("my-username", *basicAuth.Username)

	basicAuth, err = client.BasicAuths.Get(defaultCtx, consumer.ID,
		basicAuth.Username)
	require.NoError(err)
	assert.Equal("my-username", *basicAuth.Username)

	basicAuth, err = client.BasicAuths.Get(defaultCtx, consumer.ID,
		String("does-not-exists"))
	assert.Nil(basicAuth)
	require.Error(err)

	basicAuth, err = client.BasicAuths.Get(defaultCtx, consumer.ID, String(""))
	assert.Nil(basicAuth)
	require.Error(err)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestBasicAuthGetByID(t *testing.T) {
	RunWhenDBMode(t, "postgres")

	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	uuid := uuid.NewString()
	basicAuth := &BasicAuth{
		ID:       String(uuid),
		Username: String("my-username"),
		Password: String("my-password"),
	}

	// consumer for the basic-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(t, err)
	require.NotNil(t, consumer)

	createdBasicAuth, err := client.BasicAuths.Create(defaultCtx,
		consumer.ID, basicAuth)
	require.NoError(t, err)
	require.NotNil(t, createdBasicAuth)

	t.Cleanup(func() {
		require.NoError(t, client.Consumers.Delete(defaultCtx, consumer.ID))
	})

	t.Run("successful basic-auth retrieval by ID", func(_ *testing.T) {
		basicAuth, err = client.BasicAuths.GetByID(defaultCtx, basicAuth.ID)
		require.NoError(t, err)
		require.NotNil(t, basicAuth)
		require.Equal(t, "my-username", *basicAuth.Username)
	})

	t.Run("unsuccessful basic-auth retrieval using invalid ID", func(_ *testing.T) {
		basicAuth, err = client.BasicAuths.GetByID(defaultCtx, String("does-not-exist"))
		require.Nil(t, basicAuth)
		require.Error(t, err)
	})

	t.Run("unsuccessful basic-auth retrieval using empty string as ID", func(_ *testing.T) {
		basicAuth, err = client.BasicAuths.GetByID(defaultCtx, String(""))
		require.Nil(t, basicAuth)
		require.Error(t, err)
	})
}

func TestBasicAuthUpdate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	basicAuth := &BasicAuth{
		ID:       String(uuid),
		Username: String("my-username"),
		Password: String("my-password"),
	}

	// consumer for the basic-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdBasicAuth, err := client.BasicAuths.Create(defaultCtx,
		consumer.ID, basicAuth)
	require.NoError(err)
	assert.NotNil(createdBasicAuth)

	basicAuth, err = client.BasicAuths.Get(defaultCtx,
		consumer.ID, basicAuth.ID)
	require.NoError(err)
	assert.Equal("my-username", *basicAuth.Username)

	basicAuth.Username = String("my-new-username")
	basicAuth.Password = String("my-new-password")
	updatedBasicAuth, err := client.BasicAuths.Update(defaultCtx,
		consumer.ID, basicAuth)
	require.NoError(err)
	assert.NotNil(updatedBasicAuth)
	assert.NotEqual("my-new-password", *updatedBasicAuth.Password)
	assert.Equal("my-new-username", *updatedBasicAuth.Username)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestBasicAuthDelete(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	basicAuth := &BasicAuth{
		ID:       String(uuid),
		Username: String("my-username"),
		Password: String("my-password"),
	}

	// consumer for the basic-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdBasicAuth, err := client.BasicAuths.Create(defaultCtx,
		consumer.ID, basicAuth)
	require.NoError(err)
	assert.NotNil(createdBasicAuth)

	err = client.BasicAuths.Delete(defaultCtx, consumer.ID, basicAuth.Username)
	require.NoError(err)

	basicAuth, err = client.BasicAuths.Get(defaultCtx,
		consumer.ID, basicAuth.Username)
	require.Error(err)
	assert.Nil(basicAuth)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestBasicAuthListMethods(T *testing.T) {
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

	// consumer for the basic-auth:
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
	basicAuths := []*BasicAuth{
		{
			Username: String("username11"),
			Password: String("password11"),
			Consumer: consumer1,
		},
		{
			Username: String("username12"),
			Password: String("password12"),
			Consumer: consumer1,
		},
		{
			Username: String("username21"),
			Password: String("password21"),
			Consumer: consumer2,
		},
		{
			Username: String("username22"),
			Password: String("password22"),
			Consumer: consumer2,
		},
	}

	// create fixturs
	for i := 0; i < len(basicAuths); i++ {
		basicAuth, err := client.BasicAuths.Create(defaultCtx,
			basicAuths[i].Consumer.ID, basicAuths[i])
		require.NoError(err)
		assert.NotNil(basicAuth)
		basicAuths[i] = basicAuth
	}

	basicAuthsFromKong, next, err := client.BasicAuths.List(defaultCtx, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(basicAuthsFromKong)
	assert.Len(basicAuthsFromKong, 4)

	// first page
	page1, next, err := client.BasicAuths.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)

	// last page
	next.Size = 4
	page2, next, err := client.BasicAuths.List(defaultCtx, next)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 3)

	basicAuthsForConsumer, next, err := client.BasicAuths.ListForConsumer(defaultCtx, consumer1.ID, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(basicAuthsForConsumer)
	assert.Len(basicAuthsForConsumer, 2)

	basicAuths, err = client.BasicAuths.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(basicAuths)
	assert.Len(basicAuths, 4)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer1.ID))
	require.NoError(client.Consumers.Delete(defaultCtx, consumer2.ID))
}
