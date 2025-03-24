package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOauth2CredentialCreate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	oauth2Cred, err := client.Oauth2Credentials.Create(defaultCtx,
		String("foo"), nil)
	require.Error(err)
	assert.Nil(oauth2Cred)

	oauth2Cred = &Oauth2Credential{}
	oauth2Cred, err = client.Oauth2Credentials.Create(defaultCtx, String(""),
		oauth2Cred)
	require.Error(err)
	assert.Nil(oauth2Cred)

	// consumer for the oauth2 cred
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	oauth2Cred = &Oauth2Credential{
		ClientID:     String("foo"),
		Name:         String("name-foo"),
		RedirectURIs: StringSlice("http://foo.com", "http://bar.com"),
	}
	oauth2Cred, err = client.Oauth2Credentials.Create(defaultCtx,
		consumer.ID, oauth2Cred)
	require.NoError(err)
	assert.NotNil(oauth2Cred)
	assert.NotNil(oauth2Cred.ClientSecret)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestOauth2CredentialCreateWithID(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	oauth2Cred := &Oauth2Credential{
		ID:           String(uuid),
		Name:         String("name"),
		ClientSecret: String("my-client-secret"),
		RedirectURIs: StringSlice("http://foo.com", "http://bar.com"),
		ClientID:     String("my-clientid"),
	}

	// consumer for the oauth2 cred
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdOauth2Credential, err := client.Oauth2Credentials.Create(
		defaultCtx, consumer.ID, oauth2Cred)
	require.NoError(err)
	assert.NotNil(createdOauth2Credential)

	assert.Equal(uuid, *createdOauth2Credential.ID)
	assert.Equal("my-clientid", *createdOauth2Credential.ClientID)
	assert.Equal("my-client-secret", *createdOauth2Credential.ClientSecret)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestOauth2CredentialCreatePublicClientType(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	uuid := uuid.NewString()
	oauth2Cred := &Oauth2Credential{
		ID:           String(uuid),
		Name:         String("name"),
		ClientID:     String("public-client"),
		ClientType:   String("public"),
		RedirectURIs: StringSlice("http://foo.com", "http://bar.com"),
	}

	// consumer for the oauth2 cred
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdOauth2Credential, err := client.Oauth2Credentials.Create(
		defaultCtx, consumer.ID, oauth2Cred)
	require.NoError(err)
	require.NotNil(createdOauth2Credential)

	assert.Equal(uuid, *createdOauth2Credential.ID)
	assert.Equal("public", *createdOauth2Credential.ClientType)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestOauth2CredentialGet(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	oauth2Cred := &Oauth2Credential{
		ID:           String(uuid),
		Name:         String("name-foo"),
		ClientID:     String("foo-clientid"),
		RedirectURIs: StringSlice("http://foo.com", "http://bar.com"),
	}

	// consumer for the oauth2 cred
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdOauth2Credential, err := client.Oauth2Credentials.Create(defaultCtx,
		consumer.ID, oauth2Cred)
	require.NoError(err)
	assert.NotNil(createdOauth2Credential)

	oauth2Cred, err = client.Oauth2Credentials.Get(defaultCtx,
		consumer.ID, oauth2Cred.ID)
	require.NoError(err)
	assert.Equal("foo-clientid", *oauth2Cred.ClientID)

	oauth2Cred, err = client.Oauth2Credentials.Get(defaultCtx, consumer.ID,
		String("foo-clientid"))
	require.NoError(err)
	assert.NotNil(oauth2Cred)

	oauth2Cred, err = client.Oauth2Credentials.Get(defaultCtx, consumer.ID,
		String("does-not-exists"))
	assert.Nil(oauth2Cred)
	require.Error(err)

	oauth2Cred, err = client.Oauth2Credentials.Get(defaultCtx,
		consumer.ID, String(""))
	assert.Nil(oauth2Cred)
	require.Error(err)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestOauth2CredentialGetById(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	oauth2Cred := &Oauth2Credential{
		ID:           String(uuid),
		Name:         String("name-foo"),
		ClientID:     String("foo-clientid"),
		RedirectURIs: StringSlice("http://foo.com", "http://bar.com"),
	}

	// consumer for the oauth2 cred
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdOauth2Credential, err := client.Oauth2Credentials.Create(defaultCtx,
		consumer.ID, oauth2Cred)
	require.NoError(err)
	assert.NotNil(createdOauth2Credential)

	oauth2Cred, err = client.Oauth2Credentials.GetById(defaultCtx, oauth2Cred.ID)
	require.NoError(err)
	assert.Equal("foo-clientid", *oauth2Cred.ClientID)

	oauth2Cred, err = client.Oauth2Credentials.GetById(defaultCtx, String("does-not-exist"))
	assert.Nil(oauth2Cred)
	require.Error(err)

	oauth2Cred, err = client.Oauth2Credentials.GetById(defaultCtx, String(""))
	assert.Nil(oauth2Cred)
	require.Error(err)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestOauth2CredentialUpdate(T *testing.T) {
	RunWhenKong(T, "<=2.0.5")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	oauth2Cred := &Oauth2Credential{
		ID:           String(uuid),
		ClientID:     String("client-id"),
		Name:         String("foo-name"),
		RedirectURIs: StringSlice("http://foo.com", "http://bar.com"),
	}

	// consumer for the oauth2 cred
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdOauth2Credential, err := client.Oauth2Credentials.Create(
		defaultCtx, consumer.ID, oauth2Cred)
	require.NoError(err)
	assert.NotNil(createdOauth2Credential)

	oauth2Cred, err = client.Oauth2Credentials.Get(defaultCtx,
		consumer.ID, oauth2Cred.ID)
	require.NoError(err)
	assert.Equal("foo-name", *oauth2Cred.Name)

	oauth2Cred.Name = String("new-foo-name")
	oauth2Cred.ClientSecret = String("my-new-secret")
	updatedOauth2Credential, err := client.Oauth2Credentials.Update(defaultCtx, consumer.ID, oauth2Cred)
	require.NoError(err)
	assert.NotNil(updatedOauth2Credential)
	assert.Equal("new-foo-name", *updatedOauth2Credential.Name)
	assert.Equal("my-new-secret", *updatedOauth2Credential.ClientSecret)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestOauth2CredentialDelete(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	oauth2Cred := &Oauth2Credential{
		ID:           String(uuid),
		ClientID:     String("my-client-id"),
		Name:         String("my-name"),
		RedirectURIs: StringSlice("http://foo.com", "http://bar.com"),
	}

	// consumer for the oauth2 cred
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdOauth2Credential, err := client.Oauth2Credentials.Create(defaultCtx, consumer.ID, oauth2Cred)
	require.NoError(err)
	assert.NotNil(createdOauth2Credential)

	err = client.Oauth2Credentials.Delete(defaultCtx,
		consumer.ID, oauth2Cred.ClientID)
	require.NoError(err)

	oauth2Cred, err = client.Oauth2Credentials.Get(defaultCtx,
		consumer.ID, oauth2Cred.ClientID)
	require.Error(err)
	assert.Nil(oauth2Cred)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestOauth2CredentialListMethods(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	// consumer for the oauth2 cred
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
	oauth2Creds := []*Oauth2Credential{
		{
			ClientID:     String("clientid11"),
			Name:         String("name11"),
			RedirectURIs: StringSlice("http://foo.com", "http://bar.com"),
			Consumer:     consumer1,
		},
		{
			ClientID:     String("clientid12"),
			Name:         String("name12"),
			RedirectURIs: StringSlice("http://foo.com", "http://bar.com"),
			Consumer:     consumer1,
		},
		{
			ClientID:     String("clientid21"),
			Name:         String("name21"),
			RedirectURIs: StringSlice("http://foo.com", "http://bar.com"),
			Consumer:     consumer2,
		},
		{
			ClientID:     String("clientid22"),
			Name:         String("name22"),
			RedirectURIs: StringSlice("http://foo.com", "http://bar.com"),
			Consumer:     consumer2,
		},
	}

	// create fixturs
	for i := 0; i < len(oauth2Creds); i++ {
		oauth2Cred, err := client.Oauth2Credentials.Create(defaultCtx,
			oauth2Creds[i].Consumer.ID, oauth2Creds[i])
		require.NoError(err)
		assert.NotNil(oauth2Cred)
		oauth2Creds[i] = oauth2Cred
	}

	oauth2CredsFromKong, next, err := client.Oauth2Credentials.List(defaultCtx, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(oauth2CredsFromKong)
	assert.Len(oauth2CredsFromKong, 4)

	// first page
	page1, next, err := client.Oauth2Credentials.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)

	// last page
	next.Size = 3
	page2, next, err := client.Oauth2Credentials.List(defaultCtx, next)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 3)

	oauth2CredsForConsumer, next, err := client.Oauth2Credentials.ListForConsumer(defaultCtx, consumer1.ID, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(oauth2CredsForConsumer)
	assert.Len(oauth2CredsForConsumer, 2)

	oauth2Creds, err = client.Oauth2Credentials.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(oauth2Creds)
	assert.Len(oauth2Creds, 4)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer1.ID))
	require.NoError(client.Consumers.Delete(defaultCtx, consumer2.ID))
}
