package kong

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestServicesService(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	service := &Service{
		Name: String("foo"),
		Host: String("upstream"),
		Port: Int(42),
		Path: String("/path"),
	}

	createdService, err := client.Services.Create(defaultCtx, service)
	assert.Nil(err)
	assert.NotNil(createdService)

	service, err = client.Services.Get(defaultCtx, createdService.ID)
	assert.Nil(err)
	assert.NotNil(service)

	service.Name = String("bar")
	service.Host = String("newUpstream")
	service, err = client.Services.Update(defaultCtx, service)
	assert.Nil(err)
	assert.NotNil(service)
	assert.Equal("bar", *service.Name)
	assert.Equal("newUpstream", *service.Host)
	assert.Equal(42, *service.Port)

	route, err := client.Routes.CreateInService(defaultCtx, service.ID, &Route{
		Paths: StringSlice("/route"),
	})
	assert.Nil(err)
	assert.NotNil(route)

	serviceForRoute, err := client.Services.GetForRoute(defaultCtx, route.ID)
	assert.Nil(err)
	assert.NotNil(serviceForRoute)

	assert.Equal(*service.ID, *serviceForRoute.ID)

	err = client.Routes.Delete(defaultCtx, route.ID)
	assert.Nil(err)

	err = client.Services.Delete(defaultCtx, service.ID)
	assert.Nil(err)

	// ID can be specified
	id := uuid.NewString()
	service = &Service{
		Name: String("fizz"),
		ID:   String(id),
		Host: String("buzz"),
	}

	createdService, err = client.Services.Create(defaultCtx, service)
	assert.Nil(err)
	assert.NotNil(createdService)
	assert.Equal(id, *createdService.ID)
	assert.Equal("buzz", *createdService.Host)

	err = client.Services.Delete(defaultCtx, createdService.ID)
	assert.Nil(err)

	_, err = client.Services.Create(defaultCtx, nil)
	assert.NotNil(err)

	_, err = client.Services.Update(defaultCtx, nil)
	assert.NotNil(err)
}

func TestServiceWithTags(T *testing.T) {
	runWhenKong(T, ">=1.1.0")
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	service := &Service{
		Name: String("key-auth"),
		Host: String("example.com"),
		Tags: StringSlice("tag1", "tag2"),
	}

	createdService, err := client.Services.Create(defaultCtx, service)
	assert.Nil(err)
	assert.NotNil(createdService)
	assert.Equal(StringSlice("tag1", "tag2"), createdService.Tags)

	err = client.Services.Delete(defaultCtx, createdService.ID)
	assert.Nil(err)
}

func TestServiceListEndpoint(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
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
		assert.Nil(err)
		assert.NotNil(service)
		services[i] = service
	}

	servicesFromKong, next, err := client.Services.List(defaultCtx, nil)
	assert.Nil(err)
	assert.Nil(next)
	assert.NotNil(servicesFromKong)
	assert.Equal(3, len(servicesFromKong))

	// check if we see all services
	assert.True(compareServices(services, servicesFromKong))

	// Test pagination
	servicesFromKong = []*Service{}

	// first page
	page1, next, err := client.Services.List(defaultCtx, &ListOpt{Size: 1})
	assert.Nil(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	servicesFromKong = append(servicesFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.Services.List(defaultCtx, next)
	assert.Nil(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Equal(2, len(page2))
	servicesFromKong = append(servicesFromKong, page2...)

	assert.True(compareServices(services, servicesFromKong))

	services, err = client.Services.ListAll(defaultCtx)
	assert.Nil(err)
	assert.NotNil(services)
	assert.Equal(3, len(services))

	for i := 0; i < len(services); i++ {
		assert.Nil(client.Services.Delete(defaultCtx, services[i].ID))
	}
}

func compareServices(expected, actual []*Service) bool {
	var expectedUsernames, actualUsernames []string
	for _, service := range expected {
		expectedUsernames = append(expectedUsernames, *service.Name)
	}

	for _, service := range actual {
		actualUsernames = append(actualUsernames, *service.Name)
	}

	return (compareSlices(expectedUsernames, actualUsernames))
}

func TestServiceWithClientCert(T *testing.T) {
	runWhenKong(T, ">=1.3.0")
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	certificate := &Certificate{
		Key:  String(key1),
		Cert: String(cert1),
	}
	createdCertificate, err := client.Certificates.Create(defaultCtx, certificate)
	assert.Nil(err)
	assert.NotNil(createdCertificate)

	service := &Service{
		Name:              String("foo"),
		Host:              String("example.com"),
		Protocol:          String("https"),
		ClientCertificate: createdCertificate,
	}

	createdService, err := client.Services.Create(defaultCtx, service)
	assert.Nil(err)
	assert.NotNil(createdService)
	assert.Equal(*createdCertificate.ID, *createdService.ClientCertificate.ID)

	err = client.Services.Delete(defaultCtx, createdService.ID)
	assert.Nil(err)

	err = client.Certificates.Delete(defaultCtx, createdCertificate.ID)
	assert.Nil(err)
}

func TestServiceGetFullSchema(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	schema, err := client.Services.GetFullSchema(defaultCtx)
	_, ok := schema["fields"]
	assert.True(ok)
	assert.Nil(err)
}

func TestFillServiceDefaults(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	tests := []struct {
		name     string
		service  *Service
		expected *Service
	}{
		{
			name: "name and host only",
			service: &Service{
				Name: String("svc1"),
				Host: String("mockbin.org"),
			},
			expected: &Service{
				Name:           String("svc1"),
				Host:           String("mockbin.org"),
				Port:           Int(80),
				Protocol:       String("http"),
				ConnectTimeout: Int(60000),
				ReadTimeout:    Int(60000),
				Retries:        Int(5),
				WriteTimeout:   Int(60000),
			},
		},
		{
			name: "name host and port",
			service: &Service{
				Name: String("svc1"),
				Host: String("mockbin.org"),
				Port: Int(8080),
			},
			expected: &Service{
				Name:           String("svc1"),
				Host:           String("mockbin.org"),
				Port:           Int(8080),
				Protocol:       String("http"),
				ConnectTimeout: Int(60000),
				ReadTimeout:    Int(60000),
				Retries:        Int(5),
				WriteTimeout:   Int(60000),
			},
		},
		{
			name: "name host port and tags",
			service: &Service{
				Name: String("svc1"),
				Host: String("mockbin.org"),
				Port: Int(8080),
				Tags: []*string{String("tag1"), String("tag2")},
			},
			expected: &Service{
				Name:           String("svc1"),
				Host:           String("mockbin.org"),
				Port:           Int(8080),
				Protocol:       String("http"),
				ConnectTimeout: Int(60000),
				ReadTimeout:    Int(60000),
				Retries:        Int(5),
				WriteTimeout:   Int(60000),
				Tags:           []*string{String("tag1"), String("tag2")},
			},
		},
	}

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			s := tc.service
			fullSchema, err := client.Services.GetFullSchema(defaultCtx)
			assert.Nil(err)
			assert.NotNil(fullSchema)
			if err := FillServicesDefaults(s, fullSchema); err != nil {
				t.Errorf(err.Error())
			}
			if diff := cmp.Diff(s, tc.expected); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
