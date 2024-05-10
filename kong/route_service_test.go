package kong

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoutesRoute(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	SkipWhenKongRouterFlavor(T, Expressions)

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	route := &Route{}

	routeNotCreated, err := client.Routes.Create(defaultCtx, route)
	assert.NotNil(err)
	assert.Nil(routeNotCreated)

	// service for the route
	service := &Service{
		Name: String("foo2"),
		Host: String("upstream"),
		Port: Int(42),
		Path: String("/path"),
	}

	service, err = client.Services.Create(defaultCtx, service)
	assert.NoError(err)
	assert.NotNil(service)

	route = &Route{
		Hosts:   StringSlice("host1.com", "host2.com"),
		Service: service,
	}
	createdRoute, err := client.Routes.Create(defaultCtx, route)
	assert.NoError(err)
	require.NotNil(createdRoute)

	route, err = client.Routes.Get(defaultCtx, createdRoute.ID)
	assert.NoError(err)
	assert.NotNil(route)
	assert.Empty(route.Methods)
	assert.Empty(route.Paths)

	route.Hosts = StringSlice("newHost.com")
	route.Methods = StringSlice("GET", "POST")
	route, err = client.Routes.Update(defaultCtx, route)
	assert.NoError(err)
	assert.NotNil(route)
	assert.Equal(1, len(route.Hosts))
	assert.Equal("newHost.com", *route.Hosts[0])

	err = client.Routes.Delete(defaultCtx, createdRoute.ID)
	assert.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	route = &Route{
		ID:        String(id),
		Name:      String("new-route"),
		SNIs:      StringSlice("snihost1.com", "snihost2.com"),
		Protocols: StringSlice("tcp", "tls"),
		Destinations: []*CIDRPort{
			{
				IP:   String("10.0.0.0/8"),
				Port: Int(80),
			},
		},
		Service: service,
	}

	createdRoute, err = client.Routes.Create(defaultCtx, route)
	assert.NoError(err)
	assert.NotNil(createdRoute)
	assert.Equal(id, *createdRoute.ID)
	assert.Equal(2, len(createdRoute.SNIs))
	assert.Equal("snihost1.com", *createdRoute.SNIs[0])
	assert.Equal("snihost2.com", *createdRoute.SNIs[1])
	assert.Equal("10.0.0.0/8", *createdRoute.Destinations[0].IP)
	assert.Equal(80, *createdRoute.Destinations[0].Port)

	err = client.Routes.Delete(defaultCtx, createdRoute.ID)
	assert.NoError(err)

	err = client.Services.Delete(defaultCtx, service.ID)
	assert.NoError(err)

	_, err = client.Routes.Create(defaultCtx, nil)
	assert.NotNil(err)

	_, err = client.Routes.Update(defaultCtx, nil)
	assert.NotNil(err)
}

func TestRouteWithTags(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=1.1.0")
	SkipWhenKongRouterFlavor(T, Expressions)

	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	route := &Route{
		Name:  String("key-auth"),
		Paths: StringSlice("/"),
		Tags:  StringSlice("tag1", "tag2"),
	}

	createdRoute, err := client.Routes.Create(defaultCtx, route)
	require.NoError(err)
	require.NotNil(createdRoute)
	require.Equal(StringSlice("tag1", "tag2"), createdRoute.Tags)

	err = client.Routes.Delete(defaultCtx, createdRoute.ID)
	require.NoError(err)
}

func TestCreateExpressionRoutes(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKongRouterFlavor(T, Expressions)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	require.NotNil(T, client)

	for _, tc := range []struct {
		name   string
		route  *Route
		valid  bool
		assert func(t *testing.T, route *Route)
	}{
		{
			name: "route with expression and high priority",
			route: &Route{
				Expression: String(`http.path == "/"`),
				Priority:   Uint64(33820977671045),
			},
			valid: true,
			assert: func(t *testing.T, route *Route) {
				assert.Equal(t, `http.path == "/"`, *route.Expression)
				assert.Equal(t, uint64(33820977671045), *route.Priority)
			},
		},
		{
			name: "route with expression and priority of 1",
			route: &Route{
				Expression: String(`http.path == "/"`),
				Priority:   Uint64(1),
			},
			valid: true,
			assert: func(t *testing.T, route *Route) {
				assert.Equal(t, `http.path == "/"`, *route.Expression)
				assert.Equal(t, uint64(1), *route.Priority)
			},
		},
		// TODO: this fails now because Gateway returns priority in scientific notation:
		// failed decoding response body: json: cannot unmarshal number 3.3820977671045e+15 into Go struct field Route.priority of type int64
		// Ref: https://konghq.atlassian.net/browse/FTI-5515
		// {
		// 	route: &Route{
		// 		Expression: String(`lower(http.path) ^= "/"`),
		// 		Priority:   Int64(3382097767104500),
		// 	},
		// 	valid: true,
		// },
	} {
		T.Run(tc.name, func(T *testing.T) {
			createdRoute, err := client.Routes.Create(defaultCtx, tc.route)
			if tc.valid {
				assert.NoError(T, err)
				require.NotNil(T, createdRoute)
				T.Cleanup(func() {
					assert.NoError(T, client.Routes.Delete(defaultCtx, createdRoute.ID))
				})
				tc.assert(T, createdRoute)
			} else {
				assert.Error(T, err)
			}
		})
	}
}

func TestCreateInRoute(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	SkipWhenKongRouterFlavor(T, Expressions)

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	service := &Service{
		Name: String("foo"),
		Host: String("upstream"),
		Port: Int(42),
		Path: String("/path"),
	}

	createdService, err := client.Services.Create(defaultCtx, service)
	assert.NoError(err)
	require.NotNil(createdService)

	route := &Route{
		Hosts: StringSlice("host1.com", "host2.com"),
	}

	// specifying name won't work
	routeNotCreated, err := client.Routes.CreateInService(defaultCtx,
		createdService.Name, route)
	assert.Nil(routeNotCreated)
	assert.NotNil(err)

	createdRoute, err := client.Routes.CreateInService(defaultCtx,
		createdService.ID, route)
	assert.NoError(err)
	assert.NotNil(createdRoute)

	assert.NoError(client.Routes.Delete(defaultCtx, createdRoute.ID))
	assert.NoError(client.Services.Delete(defaultCtx, createdService.ID))
}

func TestRouteListEndpoint(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	SkipWhenKongRouterFlavor(T, Expressions)

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	service := &Service{
		Name: String("foo"),
		Host: String("upstream"),
		Port: Int(42),
		Path: String("/path"),
	}

	createdService, err := client.Services.Create(defaultCtx, service)
	require.NoError(err)
	require.NotNil(createdService)

	// fixtures
	routes := []*Route{
		{
			Paths:   StringSlice("/foo1"),
			Service: createdService,
		},
		{
			Paths:   StringSlice("/foo2"),
			Service: createdService,
		},
		{
			Paths:   StringSlice("/foo3"),
			Service: createdService,
		},
	}

	// create fixturs
	for i := 0; i < len(routes); i++ {
		route, err := client.Routes.Create(defaultCtx, routes[i])
		assert.NoError(err)
		assert.NotNil(route)
		routes[i] = route
	}

	routesFromKong, next, err := client.Routes.List(defaultCtx, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(routesFromKong)
	assert.Equal(3, len(routesFromKong))

	// check if we see all routes
	assert.True(compareRoutes(T, routes, routesFromKong))

	// Test pagination
	routesFromKong = []*Route{}

	// first page
	page1, next, err := client.Routes.List(defaultCtx, &ListOpt{Size: 1})
	assert.NoError(err)
	require.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	routesFromKong = append(routesFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.Routes.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Equal(2, len(page2))
	routesFromKong = append(routesFromKong, page2...)

	assert.True(compareRoutes(T, routes, routesFromKong))

	routesForService, next, err := client.Routes.ListForService(defaultCtx,
		createdService.ID, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(routesForService)
	assert.True(compareRoutes(T, routes, routesForService))

	routes, err = client.Routes.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(routes)
	assert.Equal(3, len(routes))

	for i := 0; i < len(routes); i++ {
		assert.NoError(client.Routes.Delete(defaultCtx, routes[i].ID))
	}

	assert.NoError(client.Services.Delete(defaultCtx, createdService.ID))
}

func compareRoutes(T *testing.T, expected, actual []*Route) bool {
	var expectedUsernames, actualUsernames []string
	for _, route := range expected {
		if !assert.NotNil(T, route) {
			continue
		}
		expectedUsernames = append(expectedUsernames, *route.Paths[0])
	}

	for _, route := range actual {
		actualUsernames = append(actualUsernames, *route.Paths[0])
	}

	return (compareSlices(expectedUsernames, actualUsernames))
}

func TestRouteWithHeaders(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=1.3.0")
	SkipWhenKongRouterFlavor(T, Expressions)

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	route := &Route{
		Name: String("route-by-header"),
		Headers: map[string][]string{
			"foo": {"bar"},
		},
		Tags: StringSlice("tag1", "tag2"),
	}

	createdRoute, err := client.Routes.Create(defaultCtx, route)
	require.NoError(err)
	require.NotNil(createdRoute)
	assert.Equal(StringSlice("tag1", "tag2"), createdRoute.Tags)
	assert.Equal(map[string][]string{"foo": {"bar"}}, createdRoute.Headers)

	err = client.Routes.Delete(defaultCtx, createdRoute.ID)
	assert.NoError(err)
}

func TestRoutesValidationExpressions(T *testing.T) {
	RunWhenKong(T, ">=3.0.0")
	SkipWhenKongRouterFlavor(T, Traditional, TraditionalCompatible)

	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	const errMsgStart = "schema violation (Router Expression failed validation:"
	for _, tC := range []struct {
		name         string
		route        *Route
		valid        bool
		msgStartWith string
	}{
		{
			name: "invalid expression - nonexisting LHS field",
			route: &Route{
				Expression: String("net.foo == 3000"),
			},
			msgStartWith: errMsgStart,
		},
		{
			name: "invalid expression - invalid regex",
			route: &Route{
				Expression: String(`lower(http.path) ~ "pref~[[[[[ix"`),
			},
			msgStartWith: errMsgStart,
		},
		{
			name: "valid expression",
			route: &Route{
				Expression: String(`lower(http.path) ^= "/prefix/"`),
			},
			valid: true,
		},
		{
			name: "valid expression with priority",
			route: &Route{
				Expression: String(`lower(http.path) ^= "/prefix/"`),
				Priority:   Uint64(MaxRoutePriority),
			},
			valid: true,
		},
	} {
		T.Run(tC.name, func(_ *testing.T) {
			ok, msg, err := client.Routes.Validate(defaultCtx, tC.route)
			require.NoError(err)
			require.Equal(tC.valid, ok)
			if !ok {
				require.NotEmpty(tC.msgStartWith)
				require.True(strings.HasPrefix(msg, tC.msgStartWith))
			}
		})
	}
}

func TestRoutesValidationTraditionalCompatible(T *testing.T) {
	RunWhenKong(T, ">=3.0.0")
	SkipWhenKongRouterFlavor(T, Traditional, Expressions)

	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	var (
		validPath    = String("/prefix/")
		validRegex   = String("~/payment/(docs|health)$")
		invalidRegex = String("~/payment/(docs|health))")
	)
	for _, tC := range []struct {
		name         string
		route        *Route
		valid        bool
		msgStartWith string
	}{
		{
			name: "valid path - prefix",
			route: &Route{
				Paths: []*string{validPath},
			},
			valid: true,
		},
		{
			name: "valid path - regex",
			route: &Route{
				Paths: []*string{validRegex},
			},
			valid: true,
		},
		{
			name: "multiple valid paths - prefix and regex",
			route: &Route{
				Paths: []*string{validPath, validRegex},
			},
			valid: true,
		},
		{
			name: "invalid path - invalid regex (unmatched parentheses)",
			route: &Route{
				Paths: []*string{invalidRegex},
			},
			msgStartWith: "schema violation (paths.1: invalid regex:",
		},
		{
			name: "multiple paths - one path with invalid regex",
			route: &Route{
				Paths: []*string{validPath, invalidRegex, String("/foo")},
			},
			msgStartWith: "schema violation (paths.2: invalid regex:",
		},
	} {
		T.Run(tC.name, func(_ *testing.T) {
			ok, msg, err := client.Routes.Validate(defaultCtx, tC.route)
			require.NoError(err)
			require.Equal(tC.valid, ok)
			if !ok {
				require.NotEmpty(tC.msgStartWith)
				require.True(strings.HasPrefix(msg, tC.msgStartWith))
			}
		})
	}
}
