package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTCreate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	jwt, err := client.JWTAuths.Create(defaultCtx, String("foo"), nil)
	require.Error(err)
	assert.Nil(jwt)

	jwt = &JWTAuth{}
	jwt, err = client.JWTAuths.Create(defaultCtx, String(""), jwt)
	require.Error(err)
	assert.Nil(jwt)

	jwt, err = client.JWTAuths.Create(defaultCtx,
		String("does-not-exist"), jwt)
	require.Error(err)
	assert.Nil(jwt)

	// consumer for the JWT
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	jwt = &JWTAuth{
		Key:          String("foo"),
		RSAPublicKey: String("bar"),
	}
	jwt, err = client.JWTAuths.Create(defaultCtx, consumer.ID, jwt)
	require.NoError(err)
	assert.NotNil(jwt)
	assert.NotEmpty(*jwt.Secret)
	assert.Equal("bar", *jwt.RSAPublicKey)
	assert.NotEmpty(*jwt.Algorithm)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestJWTCreateWithID(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	jwt := &JWTAuth{
		ID:     String(uuid),
		Key:    String("my-key"),
		Secret: String("my-secret"),
	}

	// consumer for the jwt
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdJWT, err := client.JWTAuths.Create(defaultCtx, consumer.ID,
		jwt)
	require.NoError(err)
	assert.NotNil(createdJWT)

	assert.Equal(uuid, *createdJWT.ID)
	assert.Equal("my-key", *createdJWT.Key)
	assert.Equal("my-secret", *createdJWT.Secret)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestJWTGet(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	jwt := &JWTAuth{
		ID:  String(uuid),
		Key: String("my-key"),
	}

	// consumer for the jwt
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdJWT, err := client.JWTAuths.Create(defaultCtx, consumer.ID, jwt)
	require.NoError(err)
	assert.NotNil(createdJWT)

	jwt, err = client.JWTAuths.Get(defaultCtx, consumer.ID, jwt.ID)
	require.NoError(err)
	assert.Equal("my-key", *jwt.Key)

	jwt, err = client.JWTAuths.Get(defaultCtx, consumer.ID,
		jwt.Key)
	require.NoError(err)
	assert.Equal("my-key", *jwt.Key)

	jwt, err = client.JWTAuths.Get(defaultCtx, consumer.ID,
		String("does-not-exists"))
	assert.Nil(jwt)
	require.Error(err)

	jwt, err = client.JWTAuths.Get(defaultCtx, consumer.ID, String(""))
	assert.Nil(jwt)
	require.Error(err)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestJWTGetByID(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	jwt := &JWTAuth{
		ID:  String(uuid),
		Key: String("my-key"),
	}

	// consumer for the jwt
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdJWT, err := client.JWTAuths.Create(defaultCtx, consumer.ID, jwt)
	require.NoError(err)
	assert.NotNil(createdJWT)

	jwt, err = client.JWTAuths.GetByID(defaultCtx, jwt.ID)
	require.NoError(err)
	assert.Equal("my-key", *jwt.Key)

	jwt, err = client.JWTAuths.GetByID(defaultCtx, String("does-not-exist"))
	assert.Nil(jwt)
	require.Error(err)

	jwt, err = client.JWTAuths.GetByID(defaultCtx, String(""))
	assert.Nil(jwt)
	require.Error(err)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestJWTUpdate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	jwt := &JWTAuth{
		ID:  String(uuid),
		Key: String("my-key"),
	}

	// consumer for the jwt
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdJWT, err := client.JWTAuths.Create(defaultCtx, consumer.ID, jwt)
	require.NoError(err)
	assert.NotNil(createdJWT)

	jwt, err = client.JWTAuths.Get(defaultCtx, consumer.ID, jwt.ID)
	require.NoError(err)
	assert.Equal("my-key", *jwt.Key)

	jwt.Key = String("my-new-key")
	jwt.Secret = String("my-new-secret")
	updatedJWT, err := client.JWTAuths.Update(defaultCtx, consumer.ID, jwt)
	require.NoError(err)
	assert.NotNil(updatedJWT)
	assert.Equal("my-new-secret", *updatedJWT.Secret)
	assert.Equal("my-new-key", *updatedJWT.Key)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestJWTDelete(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	jwt := &JWTAuth{
		ID:  String(uuid),
		Key: String("my-key"),
	}

	// consumer for the jwt
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdJWT, err := client.JWTAuths.Create(defaultCtx, consumer.ID, jwt)
	require.NoError(err)
	assert.NotNil(createdJWT)

	err = client.JWTAuths.Delete(defaultCtx, consumer.ID, jwt.Key)
	require.NoError(err)

	jwt, err = client.JWTAuths.Get(defaultCtx, consumer.ID, jwt.Key)
	require.Error(err)
	assert.Nil(jwt)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestJWTListMethods(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	// consumer for the JWT
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
	jwts := []*JWTAuth{
		{
			Key:      String("username11"),
			Consumer: consumer1,
		},
		{
			Key:      String("username12"),
			Consumer: consumer1,
		},
		{
			Key:      String("username21"),
			Consumer: consumer2,
		},
		{
			Key:      String("username22"),
			Consumer: consumer2,
		},
	}

	// create fixturs
	for i := 0; i < len(jwts); i++ {
		jwt, err := client.JWTAuths.Create(defaultCtx,
			jwts[i].Consumer.ID, jwts[i])
		require.NoError(err)
		assert.NotNil(jwt)
		jwts[i] = jwt
	}

	jwtsFromKong, next, err := client.JWTAuths.List(defaultCtx, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(jwtsFromKong)
	assert.Len(jwtsFromKong, 4)

	// first page
	page1, next, err := client.JWTAuths.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)

	// last page
	next.Size = 3
	page2, next, err := client.JWTAuths.List(defaultCtx, next)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 3)

	jwtsForConsumer, next, err := client.JWTAuths.ListForConsumer(defaultCtx,
		consumer1.ID, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(jwtsForConsumer)
	assert.Len(jwtsForConsumer, 2)

	jwts, err = client.JWTAuths.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(jwts)
	assert.Len(jwts, 4)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer1.ID))
	require.NoError(client.Consumers.Delete(defaultCtx, consumer2.ID))
}
