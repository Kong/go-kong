package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServicesService(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	SkipWhenKongRouterFlavor(T, Expressions)

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	service := &Service{
		Name: String("foo"),
		Host: String("upstream"),
		Port: Int(42),
		Path: String("/path"),
	}

	createdService, err := client.Services.Create(defaultCtx, service)
	require.NoError(err)
	require.NotNil(createdService)

	service, err = client.Services.Get(defaultCtx, createdService.ID)
	require.NoError(err)
	assert.NotNil(service)

	service.Name = String("bar")
	service.Host = String("newUpstream")
	service, err = client.Services.Update(defaultCtx, service)
	require.NoError(err)
	assert.NotNil(service)
	assert.Equal("bar", *service.Name)
	assert.Equal("newUpstream", *service.Host)
	assert.Equal(42, *service.Port)

	route, err := client.Routes.CreateInService(defaultCtx, service.ID, &Route{
		Paths: StringSlice("/route"),
	})
	require.NoError(err)
	assert.NotNil(route)

	serviceForRoute, err := client.Services.GetForRoute(defaultCtx, route.ID)
	require.NoError(err)
	assert.NotNil(serviceForRoute)

	assert.Equal(*service.ID, *serviceForRoute.ID)

	err = client.Routes.Delete(defaultCtx, route.ID)
	require.NoError(err)

	err = client.Services.Delete(defaultCtx, service.ID)
	require.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	service = &Service{
		Name: String("fizz"),
		ID:   String(id),
		Host: String("buzz"),
	}

	createdService, err = client.Services.Create(defaultCtx, service)
	require.NoError(err)
	assert.NotNil(createdService)
	assert.Equal(id, *createdService.ID)
	assert.Equal("buzz", *createdService.Host)

	err = client.Services.Delete(defaultCtx, createdService.ID)
	require.NoError(err)

	_, err = client.Services.Create(defaultCtx, nil)
	require.Error(err)

	_, err = client.Services.Update(defaultCtx, nil)
	require.Error(err)
}

func TestServiceWithTags(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	RunWhenKong(T, ">=1.1.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	service := &Service{
		Name: String("key-auth"),
		Host: String("example.com"),
		Tags: StringSlice("tag1", "tag2"),
	}

	createdService, err := client.Services.Create(defaultCtx, service)
	require.NoError(err)
	require.NotNil(createdService)
	assert.Equal(StringSlice("tag1", "tag2"), createdService.Tags)

	err = client.Services.Delete(defaultCtx, createdService.ID)
	require.NoError(err)
}

func TestServiceListEndpoint(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	// fixtures
	services := []*Service{
		{
			Name: String("foo1"),
			Host: String("upstream1.com"),
		},
		{
			Name: String("foo2"),
			Host: String("upstream2.com"),
		},
		{
			Name: String("foo3"),
			Host: String("upstream3.com"),
		},
	}

	// create fixturs
	for i := 0; i < len(services); i++ {
		service, err := client.Services.Create(defaultCtx, services[i])
		require.NoError(err)
		require.NotNil(service)
		services[i] = service
	}

	servicesFromKong, next, err := client.Services.List(defaultCtx, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(servicesFromKong)
	assert.Len(servicesFromKong, 3)

	// check if we see all services
	assert.True(compareServices(T, services, servicesFromKong))

	// Test pagination
	servicesFromKong = []*Service{}

	// first page
	page1, next, err := client.Services.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(err)
	require.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)
	servicesFromKong = append(servicesFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.Services.List(defaultCtx, next)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 2)
	servicesFromKong = append(servicesFromKong, page2...)

	assert.True(compareServices(T, services, servicesFromKong))

	services, err = client.Services.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(services)
	assert.Len(services, 3)

	for i := 0; i < len(services); i++ {
		require.NoError(client.Services.Delete(defaultCtx, services[i].ID))
	}
}

func compareServices(T *testing.T, expected, actual []*Service) bool {
	var expectedUsernames, actualUsernames []string
	for _, service := range expected {
		if !assert.NotNil(T, service) {
			continue
		}
		expectedUsernames = append(expectedUsernames, *service.Name)
	}

	for _, service := range actual {
		actualUsernames = append(actualUsernames, *service.Name)
	}

	return (compareSlices(expectedUsernames, actualUsernames))
}

func TestServiceWithClientCert(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=1.3.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	certificate := &Certificate{
		Key:  String(key1),
		Cert: String(cert1),
	}
	createdCertificate, err := client.Certificates.Create(defaultCtx, certificate)
	require.NoError(err)
	require.NotNil(createdCertificate)

	service := &Service{
		Name:              String("foo"),
		Host:              String("example.com"),
		Protocol:          String("https"),
		ClientCertificate: createdCertificate,
	}

	createdService, err := client.Services.Create(defaultCtx, service)
	require.NoError(err)
	require.NotNil(createdService)
	assert.Equal(*createdCertificate.ID, *createdService.ClientCertificate.ID)

	err = client.Services.Delete(defaultCtx, createdService.ID)
	require.NoError(err)

	err = client.Certificates.Delete(defaultCtx, createdCertificate.ID)
	require.NoError(err)
}

func TestServiceWithTLSSANs(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=3.10.0.6")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	// Create a service with tls_sans property
	service := &Service{
		Name:     String("service-with-tls-sans"),
		Host:     String("example.com"),
		Protocol: String("https"),
		TLSSANs: &SANs{
			DNSNames: StringSlice("example.com", "test.example.com"),
			Uris:     StringSlice("https://example.com"),
		},
	}

	createdService, err := client.Services.Create(defaultCtx, service)
	require.NoError(err)
	require.NotNil(createdService)
	require.NotNil(createdService.TLSSANs)
	assert.Equal(StringSlice("example.com", "test.example.com"), createdService.TLSSANs.DNSNames)
	assert.Equal(StringSlice("https://example.com"), createdService.TLSSANs.Uris)

	// Update tls_sans property
	createdService.TLSSANs = &SANs{
		DNSNames: StringSlice("new-example.com"),
		Uris:     StringSlice("https://new-example.com", "https://api.new-example.com"),
	}

	updatedService, err := client.Services.Update(defaultCtx, createdService)
	require.NoError(err)
	require.NotNil(updatedService)
	require.Equal(updatedService.ID, createdService.ID)
	require.NotNil(updatedService.TLSSANs)
	assert.Equal(StringSlice("new-example.com"), updatedService.TLSSANs.DNSNames)
	assert.Equal(StringSlice("https://new-example.com", "https://api.new-example.com"), updatedService.TLSSANs.Uris)

	// Delete service
	err = client.Services.Delete(defaultCtx, updatedService.ID)
	require.NoError(err)

	// Verify deletion
	_, err = client.Services.Get(defaultCtx, updatedService.ID)
	require.Error(err)
}
