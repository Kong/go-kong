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
	assert.NoError(err)
	assert.NotNil(client)

	basicAuth, err := client.BasicAuths.Create(defaultCtx,
		String("foo"), nil)
	assert.NotNil(err)
	assert.Nil(basicAuth)

	basicAuth = &BasicAuth{}
	basicAuth, err = client.BasicAuths.Create(defaultCtx, String(""),
		basicAuth)
	assert.NotNil(err)
	assert.Nil(basicAuth)

	basicAuth, err = client.BasicAuths.Create(defaultCtx,
		String("does-not-exist"), basicAuth)
	assert.NotNil(err)
	assert.Nil(basicAuth)

	// consumer for the basic-auth:
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	assert.NoError(err)
	require.NotNil(consumer)

	// no username is specified
	basicAuth = &BasicAuth{}
	basicAuth, err = client.BasicAuths.Create(defaultCtx,
		consumer.ID, basicAuth)
	assert.NotNil(err)
	assert.Nil(basicAuth)

	basicAuth = &BasicAuth{
		Username: String("foo"),
		Password: String("bar"),
	}
	basicAuth, err = client.BasicAuths.Create(defaultCtx,
		consumer.ID, basicAuth)
	assert.NoError(err)
	assert.NotNil(basicAuth)
	assert.NotEmpty(*basicAuth.Password)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestBasicAuthCreateWithID(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(consumer)

	createdBasicAuth, err := client.BasicAuths.Create(defaultCtx, consumer.ID,
		basicAuth)
	assert.NoError(err)
	assert.NotNil(createdBasicAuth)

	assert.Equal(uuid, *createdBasicAuth.ID)
	assert.Equal("my-username", *createdBasicAuth.Username)
	// password is hashed
	assert.NotEqual("my-password", *createdBasicAuth.Password)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestBasicAuthGet(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(consumer)

	createdBasicAuth, err := client.BasicAuths.Create(defaultCtx,
		consumer.ID, basicAuth)
	assert.NoError(err)
	assert.NotNil(createdBasicAuth)

	basicAuth, err = client.BasicAuths.Get(defaultCtx,
		consumer.ID, basicAuth.ID)
	assert.NoError(err)
	assert.Equal("my-username", *basicAuth.Username)

	basicAuth, err = client.BasicAuths.Get(defaultCtx, consumer.ID,
		basicAuth.Username)
	assert.NoError(err)
	assert.Equal("my-username", *basicAuth.Username)

	basicAuth, err = client.BasicAuths.Get(defaultCtx, consumer.ID,
		String("does-not-exists"))
	assert.Nil(basicAuth)
	assert.NotNil(err)

	basicAuth, err = client.BasicAuths.Get(defaultCtx, consumer.ID, String(""))
	assert.Nil(basicAuth)
	assert.NotNil(err)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestBasicAuthUpdate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(consumer)

	createdBasicAuth, err := client.BasicAuths.Create(defaultCtx,
		consumer.ID, basicAuth)
	assert.NoError(err)
	assert.NotNil(createdBasicAuth)

	basicAuth, err = client.BasicAuths.Get(defaultCtx,
		consumer.ID, basicAuth.ID)
	assert.NoError(err)
	assert.Equal("my-username", *basicAuth.Username)

	basicAuth.Username = String("my-new-username")
	basicAuth.Password = String("my-new-password")
	updatedBasicAuth, err := client.BasicAuths.Update(defaultCtx,
		consumer.ID, basicAuth)
	assert.NoError(err)
	assert.NotNil(updatedBasicAuth)
	assert.NotEqual("my-new-password", *updatedBasicAuth.Password)
	assert.Equal("my-new-username", *updatedBasicAuth.Username)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestBasicAuthDelete(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(consumer)

	createdBasicAuth, err := client.BasicAuths.Create(defaultCtx,
		consumer.ID, basicAuth)
	assert.NoError(err)
	assert.NotNil(createdBasicAuth)

	err = client.BasicAuths.Delete(defaultCtx, consumer.ID, basicAuth.Username)
	assert.NoError(err)

	basicAuth, err = client.BasicAuths.Get(defaultCtx,
		consumer.ID, basicAuth.Username)
	assert.NotNil(err)
	assert.Nil(basicAuth)

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
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
	assert.NoError(err)
	assert.NotNil(client)

	// consumer for the basic-auth:
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
		assert.NoError(err)
		assert.NotNil(basicAuth)
		basicAuths[i] = basicAuth
	}

	basicAuthsFromKong, next, err := client.BasicAuths.List(defaultCtx, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(basicAuthsFromKong)
	assert.Equal(4, len(basicAuthsFromKong))

	// first page
	page1, next, err := client.BasicAuths.List(defaultCtx, &ListOpt{Size: 1})
	assert.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))

	// last page
	next.Size = 4
	page2, next, err := client.BasicAuths.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Equal(3, len(page2))

	basicAuthsForConsumer, next, err := client.BasicAuths.ListForConsumer(defaultCtx, consumer1.ID, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(basicAuthsForConsumer)
	assert.Equal(2, len(basicAuthsForConsumer))

	basicAuths, err = client.BasicAuths.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(basicAuths)
	assert.Equal(4, len(basicAuths))

	assert.NoError(client.Consumers.Delete(defaultCtx, consumer1.ID))
	assert.NoError(client.Consumers.Delete(defaultCtx, consumer2.ID))
}
