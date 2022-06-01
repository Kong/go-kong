package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSNIsCertificate(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	sni := &SNI{
		Name: String("host1.com"),
	}

	// cert is required
	badSNI, err := client.SNIs.Create(defaultCtx, sni)
	assert.NotNil(err)
	assert.Nil(badSNI)

	// create a cert
	fixtureCertificate, err := client.Certificates.Create(defaultCtx,
		&Certificate{
			Key:  String(key1),
			Cert: String(cert1),
		})
	assert.NoError(err)
	require.NotNil(fixtureCertificate)
	assert.NotNil(fixtureCertificate.ID)

	createdSNI, err := client.SNIs.Create(defaultCtx, &SNI{
		Name:        String("host1.com"),
		Certificate: fixtureCertificate,
	})
	assert.NoError(err)
	assert.NotNil(createdSNI)

	sni, err = client.SNIs.Get(defaultCtx, createdSNI.ID)
	assert.NoError(err)
	assert.NotNil(sni)

	sni.Name = String("host2.com")
	sni, err = client.SNIs.Update(defaultCtx, sni)
	assert.NoError(err)
	assert.NotNil(sni)
	assert.Equal("host2.com", *sni.Name)

	err = client.SNIs.Delete(defaultCtx, createdSNI.ID)
	assert.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	sni = &SNI{
		Name:        String("host3.com"),
		ID:          String(id),
		Certificate: fixtureCertificate,
	}

	createdSNI, err = client.SNIs.Create(defaultCtx, sni)
	assert.NoError(err)
	assert.NotNil(createdSNI)
	assert.Equal(id, *createdSNI.ID)

	err = client.Certificates.Delete(defaultCtx, fixtureCertificate.ID)
	assert.NoError(err)
}

func TestSNIWithTags(T *testing.T) {
	RunWhenKong(T, ">=1.1.0")
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	fixtureCertificate, err := client.Certificates.Create(defaultCtx,
		&Certificate{
			Key:  String(key1),
			Cert: String(cert1),
		})
	assert.NoError(err)

	createdSNI, err := client.SNIs.Create(defaultCtx, &SNI{
		Name:        String("host1.com"),
		Certificate: fixtureCertificate,
		Tags:        StringSlice("tag1", "tag2"),
	})
	assert.NoError(err)
	assert.NotNil(createdSNI)
	assert.Equal(StringSlice("tag1", "tag2"), createdSNI.Tags)

	err = client.Certificates.Delete(defaultCtx, fixtureCertificate.ID)
	assert.NoError(err)
}

func TestSNIListEndpoint(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	certificate := &Certificate{
		Cert: String(cert2),
		Key:  String(key2),
	}

	createdCertificate, err := client.Certificates.Create(defaultCtx,
		certificate)
	assert.NoError(err)
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
		assert.NoError(err)
		assert.NotNil(sni)
		snis[i] = sni
	}

	snisFromKong, next, err := client.SNIs.List(defaultCtx, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(snisFromKong)
	assert.Equal(3, len(snisFromKong))

	// check if we see all snis
	assert.True(compareSNIs(T, snis, snisFromKong))

	// Test pagination
	snisFromKong = []*SNI{}

	// first page
	page1, next, err := client.SNIs.List(defaultCtx, &ListOpt{Size: 1})
	assert.NoError(err)
	require.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	snisFromKong = append(snisFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.SNIs.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Equal(2, len(page2))
	snisFromKong = append(snisFromKong, page2...)

	assert.True(compareSNIs(T, snis, snisFromKong))

	snisForCert, next, err := client.SNIs.ListForCertificate(defaultCtx,
		createdCertificate.ID, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(snisForCert)

	assert.True(compareSNIs(T, snis, snisForCert))

	snis, err = client.SNIs.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(snis)
	assert.Equal(3, len(snis))

	for i := 0; i < len(snis); i++ {
		assert.NoError(client.SNIs.Delete(defaultCtx, snis[i].ID))
	}

	assert.NoError(client.Certificates.Delete(defaultCtx, createdCertificate.ID))
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
