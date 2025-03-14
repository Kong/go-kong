//go:build enterprise

package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMTLSCreate(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	mtls, err := client.MTLSAuths.Create(defaultCtx, String("foo"), nil)
	require.Error(err)
	assert.Nil(mtls)

	mtls = &MTLSAuth{}
	mtls, err = client.MTLSAuths.Create(defaultCtx, String(""), mtls)
	require.Error(err)
	assert.Nil(mtls)

	mtls, err = client.MTLSAuths.Create(defaultCtx,
		String("does-not-exist"), mtls)
	require.Error(err)
	assert.Nil(mtls)

	// consumer for the MTLS
	consumer := &Consumer{
		Username: String("foo"),
	}

	// without a CA certificate attached
	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	mtls = &MTLSAuth{
		SubjectName: String("test@example.com"),
	}
	mtls, err = client.MTLSAuths.Create(defaultCtx, consumer.ID, mtls)
	require.NoError(err)
	assert.NotNil(mtls)
	assert.Equal("test@example.com", *mtls.SubjectName)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))

	// with a CA certificate attached
	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	certificate := &CACertificate{
		Cert: String(caCert1),
	}
	createdCertificate, err := client.CACertificates.Create(defaultCtx,
		certificate)
	require.NoError(err)

	assert.NotNil(createdCertificate)
	mtls = &MTLSAuth{
		SubjectName:   String("test@example.com"),
		CACertificate: createdCertificate,
	}
	mtls, err = client.MTLSAuths.Create(defaultCtx, consumer.ID, mtls)
	require.NoError(err)
	assert.NotNil(mtls)
	assert.Equal("test@example.com", *mtls.SubjectName)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestMTLSCreateWithID(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	mtls := &MTLSAuth{
		ID:          String(uuid),
		SubjectName: String("test@example.com"),
	}

	// consumer for the mtls
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdMTLS, err := client.MTLSAuths.Create(defaultCtx, consumer.ID,
		mtls)
	require.NoError(err)
	assert.NotNil(createdMTLS)

	assert.Equal(uuid, *createdMTLS.ID)
	assert.Equal("test@example.com", *mtls.SubjectName)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestMTLSGet(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	mtls := &MTLSAuth{
		ID:          String(uuid),
		SubjectName: String("test@example.com"),
	}

	// consumer for the mtls
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdMTLS, err := client.MTLSAuths.Create(defaultCtx, consumer.ID, mtls)
	require.NoError(err)
	assert.NotNil(createdMTLS)

	mtls, err = client.MTLSAuths.Get(defaultCtx, consumer.ID, mtls.ID)
	require.NoError(err)
	assert.Equal("test@example.com", *mtls.SubjectName)

	mtls, err = client.MTLSAuths.Get(defaultCtx, consumer.ID,
		String("does-not-exists"))
	assert.Nil(mtls)
	require.Error(err)

	mtls, err = client.MTLSAuths.Get(defaultCtx, consumer.ID, String(""))
	assert.Nil(mtls)
	require.Error(err)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestMTLSUpdate(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	mtls := &MTLSAuth{
		ID:          String(uuid),
		SubjectName: String("test@example.com"),
	}

	// consumer for the mtls
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdMTLS, err := client.MTLSAuths.Create(defaultCtx, consumer.ID, mtls)
	require.NoError(err)
	assert.NotNil(createdMTLS)

	mtls, err = client.MTLSAuths.Get(defaultCtx, consumer.ID, mtls.ID)
	require.NoError(err)
	assert.Equal("test@example.com", *mtls.SubjectName)

	mtls.SubjectName = String("different@example.com")
	updatedMTLS, err := client.MTLSAuths.Update(defaultCtx, consumer.ID, mtls)
	require.NoError(err)
	assert.NotNil(updatedMTLS)
	assert.Equal("different@example.com", *updatedMTLS.SubjectName)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestMTLSDelete(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	mtls := &MTLSAuth{
		ID:          String(uuid),
		SubjectName: String("test@example.com"),
	}

	// consumer for the mtls
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdMTLS, err := client.MTLSAuths.Create(defaultCtx, consumer.ID, mtls)
	require.NoError(err)
	assert.NotNil(createdMTLS)

	err = client.MTLSAuths.Delete(defaultCtx, consumer.ID, mtls.ID)
	require.NoError(err)

	mtls, err = client.MTLSAuths.Get(defaultCtx, consumer.ID, mtls.ID)
	require.Error(err)
	assert.Nil(mtls)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestMTLSListMethods(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	// consumer for the MTLS
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
	mtlss := []*MTLSAuth{
		{
			SubjectName: String("username11@example.com"),
			Consumer:    consumer1,
		},
		{
			SubjectName: String("username12@example.com"),
			Consumer:    consumer1,
		},
		{
			SubjectName: String("username21@example.com"),
			Consumer:    consumer2,
		},
		{
			SubjectName: String("username22@example.com"),
			Consumer:    consumer2,
		},
	}

	// create fixturs
	for i := 0; i < len(mtlss); i++ {
		mtls, err := client.MTLSAuths.Create(defaultCtx,
			mtlss[i].Consumer.ID, mtlss[i])
		require.NoError(err)
		assert.NotNil(mtls)
		mtlss[i] = mtls
	}

	mtlssFromKong, next, err := client.MTLSAuths.List(defaultCtx, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(mtlssFromKong)
	assert.Len(mtlssFromKong, 4)

	// first page
	page1, next, err := client.MTLSAuths.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)

	// last page
	next.Size = 3
	page2, next, err := client.MTLSAuths.List(defaultCtx, next)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 3)

	mtlssForConsumer, next, err := client.MTLSAuths.ListForConsumer(defaultCtx,
		consumer1.ID, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(mtlssForConsumer)
	assert.Len(mtlssForConsumer, 2)

	mtlss, err = client.MTLSAuths.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(mtlss)
	assert.Len(mtlss, 4)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer1.ID))
	require.NoError(client.Consumers.Delete(defaultCtx, consumer2.ID))
}
