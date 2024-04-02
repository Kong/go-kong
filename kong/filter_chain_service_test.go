package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterChainsService(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=3.4.0")
	SkipWhenKongRouterFlavor(T, Expressions)

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	service := &Service{
		Name: String("fooWithFilterChain1"),
		Host: String("example.com"),
		Port: Int(42),
		Path: String("/"),
	}
	err = client.Services.Delete(defaultCtx, service.Name)
	assert.NoError(err)

	_, err = client.Services.Create(defaultCtx, service)
	assert.NoError(err)

	filterChain := &FilterChain{
		Filters: []*Filter{
			{
				Name:   String("example-filter"),
				Config: JSONRawMessage(`"{ \"my_greeting\": \"Howdy\" }"`),
			},
		},
		Service: service,
	}
	assert.NotNil(filterChain)

	createdFilterChain, err := client.FilterChains.Create(defaultCtx, filterChain)
	assert.NoError(err)
	require.NotNil(createdFilterChain)
	require.Nil(createdFilterChain.Name)

	filterChain, err = client.FilterChains.Get(defaultCtx, createdFilterChain.ID)
	assert.NoError(err)
	assert.NotNil(filterChain)

	filterChain.Name = String("my-chain")
	filterChain, err = client.FilterChains.Update(defaultCtx, filterChain)
	assert.NoError(err)
	assert.NotNil(filterChain)
	assert.Equal(String("my-chain"), filterChain.Name)

	err = client.FilterChains.Delete(defaultCtx, createdFilterChain.ID)
	assert.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	filterChain = &FilterChain{
		Filters: []*Filter{
			{
				Name:   String("example-filter"),
				Config: JSONRawMessage(`"{ \"my_greeting\": \"Howdy\" }"`),
			},
		},
		Service: service,
		ID:      String(id),
	}

	createdFilterChain, err = client.FilterChains.Create(defaultCtx, filterChain)
	assert.NoError(err)
	assert.NotNil(createdFilterChain)
	assert.Equal(id, *createdFilterChain.ID)

	err = client.FilterChains.Delete(defaultCtx, createdFilterChain.ID)
	assert.NoError(err)

	service = &Service{
		Name: String("fooWithFilterChain2"),
		Host: String("upstream"),
		Port: Int(42),
		Path: String("/path"),
	}
	// Clean Data
	err = client.Services.Delete(defaultCtx, service.Name)
	assert.NoError(err)
	// Test to create filter chain from service endpoint
	createdService, err := client.Services.Create(defaultCtx, service)
	assert.NoError(err)

	id = uuid.NewString()
	FilterChainForService := &FilterChain{
		Filters: []*Filter{
			{
				Name:    String("example-filter"),
				Config:  JSONRawMessage(`"{ \"my_greeting\": \"Howdy\" }"`),
				Enabled: Bool(true),
			},
		},
		ID: String(id),
	}

	createdFilterChain, err = client.FilterChains.CreateForService(defaultCtx, createdService.Name, FilterChainForService)
	assert.NoError(err)
	assert.NotNil(createdFilterChain)
	assert.Equal(id, *createdFilterChain.ID)
	assert.Equal(Bool(true), createdFilterChain.Filters[0].Enabled)

	createdFilterChain.Filters[0].Enabled = Bool(false)
	updatedFilterChain, err := client.FilterChains.UpdateForService(defaultCtx, createdService.Name, createdFilterChain)
	assert.NoError(err)
	assert.NotNil(updatedFilterChain)
	assert.Equal(id, *updatedFilterChain.ID)
	assert.Equal(Bool(false), createdFilterChain.Filters[0].Enabled)

	err = client.FilterChains.DeleteForService(defaultCtx, createdService.Name, updatedFilterChain.ID)
	assert.NoError(err)

	// Create filter chain without ID
	createdFilterChain, err = client.FilterChains.CreateForService(defaultCtx, createdService.Name, &FilterChain{
		Filters: []*Filter{
			{
				Name:    String("example-filter"),
				Config:  JSONRawMessage(`"{ \"my_greeting\": \"Howdy\" }"`),
				Enabled: Bool(true),
			},
		},
	})
	assert.NoError(err)
	assert.NotNil(createdFilterChain)
	assert.NotNil(createdFilterChain.ID)

	assert.NoError(client.Services.Delete(defaultCtx, createdService.ID))

	route := &Route{
		Name:  String("route_filter_chain"),
		Paths: []*string{String("/route_filter_chain")},
	}
	// Clean Data
	err = client.Routes.Delete(defaultCtx, route.Name)
	assert.NoError(err)
	// Test to create filter chain from route endpoint
	createdRoute, err := client.Routes.Create(defaultCtx, route)
	assert.NoError(err)

	id = uuid.NewString()
	FilterChainForRoute := &FilterChain{
		Filters: []*Filter{
			{
				Name:    String("example-filter"),
				Config:  JSONRawMessage(`"{ \"my_greeting\": \"Howdy\" }"`),
				Enabled: Bool(true),
			},
		},
		ID: String(id),
	}

	createdFilterChain, err = client.FilterChains.CreateForRoute(defaultCtx, createdRoute.Name, FilterChainForRoute)
	assert.NoError(err)
	assert.NotNil(createdFilterChain)
	assert.Equal(id, *createdFilterChain.ID)
	assert.Equal(Bool(true), createdFilterChain.Filters[0].Enabled)

	createdFilterChain.Filters[0].Enabled = Bool(false)
	updatedFilterChain, err = client.FilterChains.UpdateForRoute(defaultCtx, createdRoute.Name, createdFilterChain)
	assert.NoError(err)
	assert.NotNil(updatedFilterChain)
	assert.Equal(id, *updatedFilterChain.ID)
	assert.Equal(Bool(false), createdFilterChain.Filters[0].Enabled)

	err = client.FilterChains.DeleteForRoute(defaultCtx, createdRoute.Name, updatedFilterChain.ID)
	assert.NoError(err)

	// Create filter chain without ID
	createdFilterChain, err = client.FilterChains.CreateForRoute(defaultCtx, createdRoute.Name, &FilterChain{
		Filters: []*Filter{
			{
				Name:    String("example-filter"),
				Config:  JSONRawMessage(`"{ \"my_greeting\": \"Howdy\" }"`),
				Enabled: Bool(true),
			},
		},
	})
	assert.NoError(err)
	assert.NotNil(createdFilterChain)
	assert.NotNil(createdFilterChain.ID)

	assert.NoError(client.Routes.Delete(defaultCtx, createdRoute.ID))
}

func TestFilterChainWithTags(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=3.4.0")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	service := &Service{
		Name: String("fooWithFilterChain1"),
		Host: String("example.com"),
		Port: Int(42),
		Path: String("/"),
	}
	err = client.Services.Delete(defaultCtx, service.Name)
	assert.NoError(err)

	createdService, err := client.Services.Create(defaultCtx, service)
	assert.NoError(err)

	filterChain := &FilterChain{
		Filters: []*Filter{
			{
				Name:   String("example-filter"),
				Config: JSONRawMessage(`"{ \"my_greeting\": \"Howdy\" }"`),
			},
		},
		Service: createdService,
		Tags:    StringSlice("tag1", "tag2"),
	}

	createdFilterChain, err := client.FilterChains.Create(defaultCtx, filterChain)
	assert.NoError(err)
	require.NotNil(createdFilterChain)
	require.Equal(StringSlice("tag1", "tag2"), createdFilterChain.Tags)

	err = client.FilterChains.Delete(defaultCtx, createdFilterChain.ID)
	require.NoError(err)

	err = client.Services.Delete(defaultCtx, createdService.ID)
	require.NoError(err)
}

func TestUnknownFilterChain(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=3.4.0")

	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	service := &Service{
		Name: String("fooWithFilterChain1"),
		Host: String("example.com"),
		Port: Int(42),
		Path: String("/"),
	}
	err = client.Services.Delete(defaultCtx, service.Name)
	assert.NoError(err)

	createdService, err := client.Services.Create(defaultCtx, service)
	assert.NoError(err)

	filterChain := &FilterChain{
		Filters: []*Filter{
			{
				Name:   String("filter-chain-not-present"),
				Config: JSONRawMessage(`"{ \"option\": true }"`),
			},
		},
		Service: createdService,
		Tags:    StringSlice("tag1", "tag2"),
	}

	createdFilterChain, err := client.FilterChains.Create(defaultCtx, filterChain)
	require.Error(T, err)
	require.Nil(T, createdFilterChain)

	err = client.Services.Delete(defaultCtx, createdService.ID)
	assert.NoError(err)
}

func TestFilterChainListEndpoint(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=3.4.0")

	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	// fixtures
	filterChains := []*FilterChain{
		{
			Name: String("chain-1"),
			Filters: []*Filter{
				{
					Name:   String("example-filter"),
					Config: JSONRawMessage(`"{ \"my_greeting\": \"Hi\" }"`),
				},
			},
		},
		{
			Name: String("chain-2"),
			Filters: []*Filter{
				{
					Name:   String("example-filter"),
					Config: JSONRawMessage(`"{ \"my_greeting\": \"Hey\" }"`),
				},
			},
		},
		{
			Name: String("chain-3"),
			Filters: []*Filter{
				{
					Name:   String("example-filter"),
					Config: JSONRawMessage(`"{ \"my_greeting\": \"Howdy\" }"`),
				},
			},
		},
	}

	// create fixtures
	for i := 0; i < len(filterChains); i++ {
		service, err := client.Services.Create(defaultCtx, &Service{
			Name: String("service-for-" + *filterChains[i].Name),
			Host: String("example.com"),
			Port: Int(42),
			Path: String("/"),
		})

		assert.NoError(err)
		assert.NotNil(service)
		filterChain, err := client.FilterChains.CreateForService(defaultCtx, service.Name, filterChains[i])
		assert.NoError(err)
		assert.NotNil(filterChain)
		filterChains[i] = filterChain
	}

	filterChainsFromKong, next, err := client.FilterChains.List(defaultCtx, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(filterChainsFromKong)
	assert.Equal(3, len(filterChainsFromKong))

	// check if we see all filterChains
	assert.True(compareFilterChains(T, filterChains, filterChainsFromKong))

	// Test pagination
	filterChainsFromKong = []*FilterChain{}

	// first page
	page1, next, err := client.FilterChains.List(defaultCtx, &ListOpt{Size: 1})
	assert.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	filterChainsFromKong = append(filterChainsFromKong, page1...)

	// second page
	page2, next, err := client.FilterChains.List(defaultCtx, next)
	assert.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page2)
	assert.Equal(1, len(page2))
	filterChainsFromKong = append(filterChainsFromKong, page2...)

	// last page
	page3, next, err := client.FilterChains.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(page3)
	assert.Equal(1, len(page3))
	filterChainsFromKong = append(filterChainsFromKong, page3...)

	assert.True(compareFilterChains(T, filterChains, filterChainsFromKong))

	filterChains, err = client.FilterChains.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(filterChains)
	assert.Equal(3, len(filterChains))

	for i := 0; i < len(filterChains); i++ {
		assert.NoError(client.Services.Delete(defaultCtx, filterChains[i].Service.ID))
	}
}

func TestFilterChainListAllForEntityEndpoint(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=3.4.0")
	SkipWhenKongRouterFlavor(T, Expressions)

	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	// fixtures

	createdService, err := client.Services.Create(defaultCtx, &Service{
		Name: String("foo"),
		Host: String("upstream"),
		Port: Int(42),
		Path: String("/path"),
	})
	assert.NoError(err)
	assert.NotNil(createdService)

	createdRoute, err := client.Routes.Create(defaultCtx, &Route{
		Hosts:   StringSlice("example.com", "example.test"),
		Service: createdService,
	})
	assert.NoError(err)
	assert.NotNil(createdRoute)

	filterChains := []*FilterChain{
		// specific to route
		{
			Name: String("route-chain"),
			Filters: []*Filter{
				{
					Name:   String("example-filter"),
					Config: JSONRawMessage(`"{ \"my_greeting\": \"Hello, route\" }"`),
				},
				{
					Name:   String("example-filter"),
					Config: JSONRawMessage(`"{ \"option\": false }"`),
				},
			},
			Route: createdRoute,
		},
		// specific to service
		{
			Name: String("service-chain"),
			Filters: []*Filter{
				{
					Name:   String("example-filter"),
					Config: JSONRawMessage(`"{ \"option\": false }"`),
				},
				{
					Name:   String("example-filter"),
					Config: JSONRawMessage(`"{ \"my_greeting\": \"Hello, service\" }"`),
				},
			},
			Service: createdService,
		},
	}

	// create fixtures
	for i := 0; i < len(filterChains); i++ {
		filterChain, err := client.FilterChains.Create(defaultCtx, filterChains[i])
		assert.NoError(err)
		assert.NotNil(filterChain)
		filterChains[i] = filterChain
	}

	filterChainsFromKong, err := client.FilterChains.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(filterChainsFromKong)
	assert.Equal(len(filterChains), len(filterChainsFromKong))

	// check if we see all filterChains
	assert.True(compareFilterChains(T, filterChains, filterChainsFromKong))

	filterChainsFromKong, err = client.FilterChains.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(filterChainsFromKong)
	assert.Equal(2, len(filterChainsFromKong))

	filterChainsFromKong, err = client.FilterChains.ListAllForService(defaultCtx,
		createdService.ID)
	assert.NoError(err)
	assert.NotNil(filterChainsFromKong)
	assert.Equal(1, len(filterChainsFromKong))

	filterChainsFromKong, err = client.FilterChains.ListAllForRoute(defaultCtx,
		createdRoute.ID)
	assert.NoError(err)
	assert.NotNil(filterChainsFromKong)
	assert.Equal(1, len(filterChainsFromKong))

	for i := 0; i < len(filterChains); i++ {
		assert.NoError(client.FilterChains.Delete(defaultCtx, filterChains[i].ID))
	}

	assert.NoError(client.Routes.Delete(defaultCtx, createdRoute.ID))
	assert.NoError(client.Services.Delete(defaultCtx, createdService.ID))
}

func compareFilterChains(T *testing.T, expected, actual []*FilterChain) bool {
	var expectedNames, actualNames []string
	for _, filterChain := range expected {
		if !assert.NotNil(T, filterChain) {
			continue
		}
		expectedNames = append(expectedNames, *filterChain.Name)
	}

	for _, filterChain := range actual {
		actualNames = append(actualNames, *filterChain.Name)
	}

	return (compareSlices(expectedNames, actualNames))
}
