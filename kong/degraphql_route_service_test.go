package kong

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

type ByID []*DegraphqlRoute

func (a ByID) Len() int           { return len(a) }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByID) Less(i, j int) bool { return *a[i].ID < *a[j].ID }

func TestDegraphqlRouteService(t *testing.T) {
	RunWhenEnterprise(t, ">=2.1.0", RequiredFeatures{})

	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	t.Run("can't create empty route", func(t *testing.T) {
		dgr := &DegraphqlRoute{}
		notCreated, err := client.DegraphqlRoutes.Create(defaultCtx, dgr)
		require.Error(t, err)
		require.Nil(t, notCreated)
	})

	t.Run("plain create/get/update/delete cycle", func(t *testing.T) {
		// service for the route
		service := &Service{
			Name: String("foo2"),
			Host: String("upstream"),
			Port: Int(42),
			Path: String("/path"),
		}

		service, err = client.Services.Create(defaultCtx, service)
		require.NoError(t, err)
		require.NotNil(t, service)
		require.NotEmpty(t, *service.ID)
		require.Equal(t, "foo2", *service.Name)

		// create the route
		route := &DegraphqlRoute{
			Service: service,
			URI:     String("/front"),
			Query:   String("query { viewer { login } }"),
		}
		createdRoute, err := client.DegraphqlRoutes.Create(defaultCtx, route)
		require.NoError(t, err)
		require.NotNil(t, createdRoute)
		require.NotEmpty(t, *createdRoute.ID)
		require.Equal(t, "/front", *createdRoute.URI)
		require.Equal(t, "query { viewer { login } }", *createdRoute.Query)
		require.Equal(t, *service.ID, *createdRoute.Service.ID)
		require.NotEmpty(t, createdRoute.Service.Name)
		require.Equal(t, "foo2", *createdRoute.Service.Name)

		// fetch route via serviceName + routeID
		fetchedRoute, err := client.DegraphqlRoutes.Get(defaultCtx, service.Name, createdRoute.ID)
		require.NoError(t, err)
		require.NotNil(t, fetchedRoute)
		require.Equal(t, *createdRoute.ID, *fetchedRoute.ID)
		require.Equal(t, "/front", *fetchedRoute.URI)
		require.Equal(t, "query { viewer { login } }", *fetchedRoute.Query)
		require.Equal(t, *service.ID, *fetchedRoute.Service.ID)
		require.NotEmpty(t, fetchedRoute.Service.Name)
		require.Equal(t, "foo2", *fetchedRoute.Service.Name)

		// fetch route again, via serviceID + routeID
		fetchedRoute2, err := client.DegraphqlRoutes.Get(defaultCtx, service.ID, createdRoute.ID)
		require.NoError(t, err)
		require.NotNil(t, fetchedRoute2)
		require.Equal(t, *createdRoute.ID, *fetchedRoute2.ID)
		require.Equal(t, "/front", *fetchedRoute2.URI)
		require.Equal(t, "query { viewer { login } }", *fetchedRoute2.Query)
		require.Equal(t, *service.ID, *fetchedRoute2.Service.ID)
		require.NotEmpty(t, fetchedRoute2.Service.Name)
		require.Equal(t, "foo2", *fetchedRoute2.Service.Name)

		// modify route
		fetchedRoute.URI = String("/v2/front")
		fetchedRoute.Methods = StringSlice("GET", "POST")
		updatedRoute, err := client.DegraphqlRoutes.Update(defaultCtx, fetchedRoute)
		require.NoError(t, err)
		require.NotNil(t, updatedRoute)
		require.Equal(t, *createdRoute.ID, *updatedRoute.ID)
		require.Equal(t, "/v2/front", *updatedRoute.URI)
		require.Equal(t, "query { viewer { login } }", *updatedRoute.Query)
		require.Equal(t, *service.ID, *updatedRoute.Service.ID)
		require.NotEmpty(t, updatedRoute.Service.Name)
		require.Equal(t, "foo2", *updatedRoute.Service.Name)

		// delete it
		err = client.DegraphqlRoutes.Delete(defaultCtx, service.Name, createdRoute.ID)
		require.NoError(t, err)

		// should not be there
		fetchedRoute, err = client.DegraphqlRoutes.Get(defaultCtx, service.Name, createdRoute.ID)
		require.Error(t, err)
		require.Nil(t, fetchedRoute)

		// delete the service
		err = client.Services.Delete(defaultCtx, service.ID)
		require.NoError(t, err)
	})
}

func TestDegraphqlRouteList(t *testing.T) {
	RunWhenEnterprise(t, ">=2.1.0", RequiredFeatures{})

	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	// serviceA for the route
	serviceA := &Service{
		Name: String("fooone"),
		Host: String("upstream"),
		Port: Int(42),
		Path: String("/path"),
	}

	serviceA, err = client.Services.Create(defaultCtx, serviceA)
	require.NoError(t, err)
	require.NotNil(t, serviceA)
	require.NotEmpty(t, *serviceA.ID)
	require.Equal(t, "fooone", *serviceA.Name)

	// serviceB for the route
	serviceB := &Service{
		Name: String("foo2"),
		Host: String("upstream"),
		Port: Int(42),
		Path: String("/path"),
	}

	serviceB, err = client.Services.Create(defaultCtx, serviceB)
	require.NoError(t, err)
	require.NotNil(t, serviceB)
	require.NotEmpty(t, *serviceB.ID)
	require.Equal(t, "foo2", *serviceB.Name)

	t.Run("add routes to a service, list them", func(t *testing.T) {
		routeTemplates := []*DegraphqlRoute{
			{
				Service: serviceA,
				URI:     String("/front"),
				Query:   String("query { viewer { login } }"),
			},
			{
				Service: serviceA,
				URI:     String("/frontBis"),
				Query:   String("query { viewer { login } }"),
			},
			{
				Service: serviceA,
				URI:     String("/frontThird"),
				Query:   String("query { viewer { login } }"),
			},
			{
				Service: serviceA,
				URI:     String("/front4"),
				Query:   String("query { viewer { login } }"),
			},
			{
				Service: serviceA,
				URI:     String("/frontLast"),
				Query:   String("query { viewer { login } }"),
			},
		}

		for _, r := range routeTemplates {
			createdRoute, err := client.DegraphqlRoutes.Create(defaultCtx, r)
			require.NoError(t, err)
			require.NotNil(t, createdRoute)
			require.NotEmpty(t, *createdRoute.ID)
			require.NotEmpty(t, createdRoute.Service.Name)

			*r = *createdRoute
		}

		sort.Sort(ByID(routeTemplates))

		routes, next, err := client.DegraphqlRoutes.List(defaultCtx, serviceA.Name, &ListOpt{})
		require.NoError(t, err)
		require.Nil(t, next)

		sort.Sort(ByID(routes))
		for i, r := range routeTemplates {
			require.Equal(t, *r.ID, *routes[i].ID)
			require.Equal(t, *r.Service.ID, *routes[i].Service.ID)
			require.Equal(t, *r.URI, *routes[i].URI)
			require.Equal(t, *r.Query, *routes[i].Query)
		}

		for _, r := range routes {
			require.NoError(t, client.DegraphqlRoutes.Delete(defaultCtx, r.Service.ID, r.ID))
		}
	})

	t.Run("add routes to different service, list each", func(t *testing.T) {
		routeTemplates := []*DegraphqlRoute{
			{
				Service: serviceA,
				URI:     String("/front"),
				Query:   String("query { viewer { login } }"),
			},
			{
				Service: serviceB,
				URI:     String("/frontBis"),
				Query:   String("query { viewer { login } }"),
			},
			{
				Service: serviceB,
				URI:     String("/frontThird"),
				Query:   String("query { viewer { login } }"),
			},
			{
				Service: serviceA,
				URI:     String("/front4"),
				Query:   String("query { viewer { login } }"),
			},
			{
				Service: serviceA,
				URI:     String("/frontLast"),
				Query:   String("query { viewer { login } }"),
			},
		}

		for _, r := range routeTemplates {
			createdRoute, err := client.DegraphqlRoutes.Create(defaultCtx, r)
			require.NoError(t, err)
			require.NotNil(t, createdRoute)
			require.NotEmpty(t, *createdRoute.ID)
			require.NotEmpty(t, createdRoute.Service.Name)

			*r = *createdRoute
		}

		routeTemplatesA := []*DegraphqlRoute{routeTemplates[0], routeTemplates[3], routeTemplates[4]}
		sort.Sort(ByID(routeTemplatesA))

		routeTemplatesB := []*DegraphqlRoute{routeTemplates[1], routeTemplates[2]}
		sort.Sort(ByID(routeTemplatesB))

		{
			routes, next, err := client.DegraphqlRoutes.List(defaultCtx, serviceA.Name, &ListOpt{})
			require.NoError(t, err)
			require.Nil(t, next)

			sort.Sort(ByID(routes))
			for i, r := range routeTemplatesA {
				require.Equal(t, *r.ID, *routes[i].ID)
				require.Equal(t, *r.Service.ID, *routes[i].Service.ID)
				require.Equal(t, *r.URI, *routes[i].URI)
				require.Equal(t, *r.Query, *routes[i].Query)
			}
		}

		{
			routes, next, err := client.DegraphqlRoutes.List(defaultCtx, serviceB.Name, &ListOpt{})
			require.NoError(t, err)
			require.Nil(t, next)

			sort.Sort(ByID(routes))
			for i, r := range routeTemplatesB {
				require.Equal(t, *r.ID, *routes[i].ID)
				require.Equal(t, *r.Service.ID, *routes[i].Service.ID)
				require.Equal(t, *r.URI, *routes[i].URI)
				require.Equal(t, *r.Query, *routes[i].Query)
			}
		}

		for _, r := range routeTemplates {
			require.NoError(t, client.DegraphqlRoutes.Delete(defaultCtx, r.Service.Name, r.ID))
		}
	})

	require.NoError(t, client.Services.Delete(defaultCtx, serviceA.ID))
	require.NoError(t, client.Services.Delete(defaultCtx, serviceB.ID))
}
