package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	publicRSA = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAruemrMy2Q0wZMuaREt7v
/cs5UY1JCYo+3S/F3JZhsZB2urQ8yDh9f8s2CM1N+88b05k3CTcwRagQXDxWvw+J
WUIcD6sERmh61SCUNuJ8A2YT4z4Uii3spVB/osJVkvVxR6fJtJHovNVGjtlKIioD
x/KM0tpeF4D+6G5vAFlmBUgyDecv3ijXtlP/hrwRcj3gzQUIfLhzYlUaibGG86B9
BLC3Yquw9Xg+MyTeTYCH7OuDZSGvTe7iqVBT1TsrrEjc3zNWyIJRHaHcb/BJs5eL
2ToXo1HGySDKA3gxuMweq22n8s++1tyQpMzvjhCFjPTYKviREg9CpB3d+p2E+8Iw
rwIDAQAB
-----END PUBLIC KEY-----`
	key = "cFt2eJLrFBUhCfgrShwnMKku&ZZaabca"
)

func TestJWTCreate(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	jwt, err := client.JWTAuths.Create(defaultCtx, String("foo"), nil)
	assert.NotNil(err)
	assert.Nil(jwt)

	jwt = &JWTAuth{}
	jwt, err = client.JWTAuths.Create(defaultCtx, String(""), jwt)
	assert.NotNil(err)
	assert.Nil(jwt)

	jwt, err = client.JWTAuths.Create(defaultCtx,
		String("does-not-exist"), jwt)
	assert.NotNil(err)
	assert.Nil(jwt)

	// consumer for the JWT
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	assert.Nil(err)
	assert.NotNil(consumer)

	jwt = &JWTAuth{
		Key:          String(key),
		RSAPublicKey: String(publicRSA),
	}
	jwt, err = client.JWTAuths.Create(defaultCtx, consumer.ID, jwt)
	assert.Nil(err)
	assert.NotNil(jwt)
	assert.NotEmpty(*jwt.Secret)
	assert.Equal(publicRSA, *jwt.RSAPublicKey)
	assert.NotEmpty(*jwt.Algorithm)

	assert.Nil(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestJWTCreateWithID(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	jwt := &JWTAuth{
		ID:     String(uuid),
		Key:    String(key),
		Secret: String("my-secret"),
	}

	// consumer for the jwt
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	assert.Nil(err)
	assert.NotNil(consumer)

	createdJWT, err := client.JWTAuths.Create(defaultCtx, consumer.ID,
		jwt)
	assert.Nil(err)
	assert.NotNil(createdJWT)

	assert.Equal(uuid, *createdJWT.ID)
	assert.Equal(key, *createdJWT.Key)
	assert.Equal("my-secret", *createdJWT.Secret)

	assert.Nil(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestJWTGet(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	jwt := &JWTAuth{
		ID:  String(uuid),
		Key: String(key),
	}

	// consumer for the jwt
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	assert.Nil(err)
	assert.NotNil(consumer)

	createdJWT, err := client.JWTAuths.Create(defaultCtx, consumer.ID, jwt)
	assert.Nil(err)
	assert.NotNil(createdJWT)

	jwt, err = client.JWTAuths.Get(defaultCtx, consumer.ID, jwt.ID)
	assert.Nil(err)
	assert.Equal(key, *jwt.Key)

	// jwt, err = client.JWTAuths.Get(defaultCtx, consumer.ID,
	// 	jwt.Key)
	// assert.Nil(err)
	// assert.Equal("my-key", *jwt.Key)

	jwt, err = client.JWTAuths.Get(defaultCtx, consumer.ID,
		String("does-not-exists"))
	assert.Nil(jwt)
	assert.NotNil(err)

	jwt, err = client.JWTAuths.Get(defaultCtx, consumer.ID, String(""))
	assert.Nil(jwt)
	assert.NotNil(err)

	assert.Nil(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestJWTUpdate(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	jwt := &JWTAuth{
		ID:  String(uuid),
		Key: String(key),
	}

	// consumer for the jwt
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	assert.Nil(err)
	assert.NotNil(consumer)

	createdJWT, err := client.JWTAuths.Create(defaultCtx, consumer.ID, jwt)
	assert.Nil(err)
	assert.NotNil(createdJWT)

	jwt, err = client.JWTAuths.Get(defaultCtx, consumer.ID, jwt.ID)
	assert.Nil(err)
	assert.Equal(key, *jwt.Key)

	newKey := "NpTDhujWV48RxkSP9qWE8xxewxy79PNF"
	jwt.Key = String(newKey)
	jwt.Secret = String("my-new-secret")
	updatedJWT, err := client.JWTAuths.Create(defaultCtx, consumer.ID, jwt)
	assert.Nil(err)
	assert.NotNil(updatedJWT)
	assert.Equal("my-new-secret", *updatedJWT.Secret)
	assert.Equal(newKey, *updatedJWT.Key)

	assert.Nil(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestJWTDelete(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	jwt := &JWTAuth{
		ID:  String(uuid),
		Key: String(key),
	}

	// consumer for the jwt
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	assert.Nil(err)
	assert.NotNil(consumer)

	createdJWT, err := client.JWTAuths.Create(defaultCtx, consumer.ID, jwt)
	assert.Nil(err)
	assert.NotNil(createdJWT)

	err = client.JWTAuths.Delete(defaultCtx, consumer.ID, jwt.Key)
	assert.Nil(err)

	// jwt, err = client.JWTAuths.Get(defaultCtx, consumer.ID, jwt.Key)
	// assert.NotNil(err)
	// assert.Nil(jwt)

	assert.Nil(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestJWTListMethods(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	// consumer for the JWT
	consumer1 := &Consumer{
		Username: String("foo"),
	}

	consumer1, err = client.Consumers.Create(defaultCtx, consumer1)
	assert.Nil(err)
	assert.NotNil(consumer1)

	consumer2 := &Consumer{
		Username: String("bar"),
	}

	consumer2, err = client.Consumers.Create(defaultCtx, consumer2)
	assert.Nil(err)
	assert.NotNil(consumer2)

	// fixtures
	jwts := []*JWTAuth{
		{
			Key:      String("B3c7hBrvkjp7GFsgrAVCTAzc93nLgpXG"),
			Consumer: consumer1,
		},
		{
			Key:      String("NpTDhujWV48RxkSP9qWE8xxewxy79PNF"),
			Consumer: consumer1,
		},
		{
			Key:      String("agvmRNhJjM5YtksJVKx5bHSgb4s8jXAE"),
			Consumer: consumer2,
		},
		{
			Key:      String("LB3m9ygX9zzpmtdw2GaaqtTfHrugDene"),
			Consumer: consumer2,
		},
	}

	// create fixturs
	for i := 0; i < len(jwts); i++ {
		jwt, err := client.JWTAuths.Create(defaultCtx,
			jwts[i].Consumer.ID, jwts[i])
		assert.Nil(err)
		assert.NotNil(jwt)
		jwts[i] = jwt
	}

	jwtsFromKong, next, err := client.JWTAuths.List(defaultCtx, nil)
	assert.Nil(err)
	assert.Nil(next)
	assert.NotNil(jwtsFromKong)
	assert.Equal(4, len(jwtsFromKong))

	// first page
	// page1, next, err := client.JWTAuths.List(defaultCtx, &ListOpt{Size: 1})
	// assert.Nil(err)
	// assert.NotNil(next)
	// assert.NotNil(page1)
	// assert.Equal(1, len(page1))

	// // last page
	// next.Size = 3
	// page2, next, err := client.JWTAuths.List(defaultCtx, next)
	// assert.Nil(err)
	// assert.Nil(next)
	// assert.NotNil(page2)
	// assert.Equal(3, len(page2))

	// jwtsForConsumer, next, err := client.JWTAuths.ListForConsumer(defaultCtx,
	// 	consumer1.ID, nil)
	// assert.Nil(err)
	// assert.Nil(next)
	// assert.NotNil(jwtsForConsumer)
	// assert.Equal(2, len(jwtsForConsumer))

	jwts, err = client.JWTAuths.ListAll(defaultCtx)
	assert.Nil(err)
	assert.NotNil(jwts)
	assert.Equal(4, len(jwts))

	assert.Nil(client.Consumers.Delete(defaultCtx, consumer1.ID))
	assert.Nil(client.Consumers.Delete(defaultCtx, consumer2.ID))
}
