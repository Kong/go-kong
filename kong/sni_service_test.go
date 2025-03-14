package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSNIsCertificate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	sni := &SNI{
		Name: String("host1.com"),
	}

	// cert is required
	badSNI, err := client.SNIs.Create(defaultCtx, sni)
	require.Error(err)
	assert.Nil(badSNI)

	// create a cert
	fixtureCertificate, err := client.Certificates.Create(defaultCtx,
		&Certificate{
			Key:  String(key1),
			Cert: String(cert1),
		})
	require.NoError(err)
	require.NotNil(fixtureCertificate)
	assert.NotNil(fixtureCertificate.ID)

	createdSNI, err := client.SNIs.Create(defaultCtx, &SNI{
		Name:        String("host1.com"),
		Certificate: fixtureCertificate,
	})
	require.NoError(err)
	assert.NotNil(createdSNI)

	sni, err = client.SNIs.Get(defaultCtx, createdSNI.ID)
	require.NoError(err)
	assert.NotNil(sni)

	sni.Name = String("host2.com")
	sni, err = client.SNIs.Update(defaultCtx, sni)
	require.NoError(err)
	assert.NotNil(sni)
	assert.Equal("host2.com", *sni.Name)

	err = client.SNIs.Delete(defaultCtx, createdSNI.ID)
	require.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	sni = &SNI{
		Name:        String("host3.com"),
		ID:          String(id),
		Certificate: fixtureCertificate,
	}

	createdSNI, err = client.SNIs.Create(defaultCtx, sni)
	require.NoError(err)
	assert.NotNil(createdSNI)
	assert.Equal(id, *createdSNI.ID)

	err = client.Certificates.Delete(defaultCtx, fixtureCertificate.ID)
	require.NoError(err)
}

func TestSNIWithTags(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=1.1.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	fixtureCertificate, err := client.Certificates.Create(defaultCtx,
		&Certificate{
			Key:  String(key1),
			Cert: String(cert1),
		})
	require.NoError(err)

	createdSNI, err := client.SNIs.Create(defaultCtx, &SNI{
		Name:        String("host1.com"),
		Certificate: fixtureCertificate,
		Tags:        StringSlice("tag1", "tag2"),
	})
	require.NoError(err)
	require.NotNil(createdSNI)
	assert.Equal(StringSlice("tag1", "tag2"), createdSNI.Tags)

	err = client.Certificates.Delete(defaultCtx, fixtureCertificate.ID)
	require.NoError(err)
}

func TestSNIListEndpoint(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	certificate := &Certificate{
		Cert: String(cert2),
		Key:  String(key2),
	}

	createdCertificate, err := client.Certificates.Create(defaultCtx,
		certificate)
	require.NoError(err)
	assert.NotNil(createdCertificate)

	// fixtures
	snis := []*SNI{
		{
			Name:        String("sni1"),
			Certificate: createdCertificate,
		},
		{
			Name:        String("sni2"),
			Certificate: createdCertificate,
		},
		{
			Name:        String("sni3"),
			Certificate: createdCertificate,
		},
	}

	// create fixturs
	for i := 0; i < len(snis); i++ {
		sni, err := client.SNIs.Create(defaultCtx, snis[i])
		require.NoError(err)
		assert.NotNil(sni)
		snis[i] = sni
	}

	snisFromKong, next, err := client.SNIs.List(defaultCtx, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(snisFromKong)
	assert.Len(snisFromKong, 3)

	// check if we see all snis
	assert.True(compareSNIs(T, snis, snisFromKong))

	// Test pagination
	snisFromKong = []*SNI{}

	// first page
	page1, next, err := client.SNIs.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(err)
	require.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)
	snisFromKong = append(snisFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.SNIs.List(defaultCtx, next)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 2)
	snisFromKong = append(snisFromKong, page2...)

	assert.True(compareSNIs(T, snis, snisFromKong))

	snisForCert, next, err := client.SNIs.ListForCertificate(defaultCtx,
		createdCertificate.ID, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(snisForCert)

	assert.True(compareSNIs(T, snis, snisForCert))

	snis, err = client.SNIs.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(snis)
	assert.Len(snis, 3)

	for i := 0; i < len(snis); i++ {
		require.NoError(client.SNIs.Delete(defaultCtx, snis[i].ID))
	}

	require.NoError(client.Certificates.Delete(defaultCtx, createdCertificate.ID))
}

func compareSNIs(T *testing.T, expected, actual []*SNI) bool {
	var expectedUsernames, actualUsernames []string
	for _, sni := range expected {
		if !assert.NotNil(T, sni) {
			continue
		}
		expectedUsernames = append(expectedUsernames, *sni.Name)
	}

	for _, sni := range actual {
		actualUsernames = append(actualUsernames, *sni.Name)
	}

	return (compareSlices(expectedUsernames, actualUsernames))
}
