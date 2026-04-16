package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGraphqlRateLimitingCostDecorationService(t *testing.T) {
	RunWhenEnterprise(t, ">=2.1.0", RequiredFeatures{})
	RunWhenDBMode(t, "postgres")

	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	t.Run("plain create/get/update/delete cycle", func(t *testing.T) {
		deco := &GraphqlRateLimitingCostDecoration{
			TypePath:     String("Vehicle.name"),
			AddConstant:  Float64(8),
			MulArguments: []*string{String("first")},
		}

		createdDeco, err := client.GraphqlRateLimitingCostDecorations.Create(defaultCtx, deco)
		require.NoError(t, err)
		require.NotNil(t, createdDeco)

		deco, err = client.GraphqlRateLimitingCostDecorations.Get(defaultCtx, createdDeco.ID)
		require.NoError(t, err)
		require.NotNil(t, deco)

		deco.TypePath = String("car.designation")
		deco, err = client.GraphqlRateLimitingCostDecorations.Update(defaultCtx, deco)
		require.NoError(t, err)
		require.NotNil(t, deco)
		require.Equal(t, "car.designation", *deco.TypePath)

		err = client.GraphqlRateLimitingCostDecorations.Delete(defaultCtx, createdDeco.ID)
		require.NoError(t, err)
	})

	t.Run("Can't create Cost Decoration with given ID", func(t *testing.T) {
		id := uuid.NewString()
		deco := &GraphqlRateLimitingCostDecoration{
			ID:           String(id),
			TypePath:     String("car.name"),
			AddConstant:  Float64(8),
			MulArguments: []*string{String("first")},
		}

		createdDeco, err := client.GraphqlRateLimitingCostDecorations.Create(defaultCtx, deco)
		require.Error(t, err, "can't specify an ID for creating new Cost Decoration")
		require.Nil(t, createdDeco)
	})

	t.Run("CreateWithID create/get/update/delete cycle", func(t *testing.T) {
		id := uuid.NewString()
		deco := &GraphqlRateLimitingCostDecoration{
			ID:           String(id),
			TypePath:     String("Vehicle.name"),
			AddConstant:  Float64(8),
			MulArguments: []*string{String("first")},
		}

		createdDeco, err := client.GraphqlRateLimitingCostDecorations.CreateWithID(defaultCtx, deco)
		require.NoError(t, err)
		require.NotNil(t, createdDeco)
		require.Equal(t, id, *createdDeco.ID)

		deco, err = client.GraphqlRateLimitingCostDecorations.Get(defaultCtx, createdDeco.ID)
		require.NoError(t, err)
		require.NotNil(t, deco)
		require.Equal(t, id, *deco.ID)

		deco.TypePath = String("car.designation")
		deco, err = client.GraphqlRateLimitingCostDecorations.Update(defaultCtx, deco)
		require.NoError(t, err)
		require.NotNil(t, deco)
		require.Equal(t, id, *deco.ID)
		require.Equal(t, "car.designation", *deco.TypePath)

		err = client.GraphqlRateLimitingCostDecorations.Delete(defaultCtx, createdDeco.ID)
		require.NoError(t, err)
	})

	t.Run("Can't CreateWithID Cost Decoration without ID", func(t *testing.T) {
		deco := &GraphqlRateLimitingCostDecoration{
			TypePath:     String("car.name"),
			AddConstant:  Float64(8),
			MulArguments: []*string{String("first")},
		}

		createdDeco, err := client.GraphqlRateLimitingCostDecorations.CreateWithID(defaultCtx, deco)
		require.Error(t, err, "ID cannot be nil for CreateWithID operation")
		require.Nil(t, createdDeco)
	})
}

func TestGraphqlRateLimitingCostDecorationServiceForService(t *testing.T) {
	RunWhenEnterprise(t, "=3.4.3.25 || =3.10.0.10 || =3.11.0.9 || =3.12.0.5 || =3.13.0.3 || >=3.14.0.2",
		RequiredFeatures{})
	RunWhenDBMode(t, "postgres")

	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Create a service to associate cost decorations with.
	service := &Service{
		Name: String("test-graphql-service2"),
		Host: String("mockbin.org"),
	}
	service, err = client.Services.Create(defaultCtx, service)
	require.NoError(t, err)
	require.NotNil(t, service)
	defer func() {
		_ = client.Services.Delete(defaultCtx, service.ID)
	}()

	// Enable the graphql-rate-limiting-advanced plugin on the service.
	plugin := &Plugin{
		Name: String("graphql-rate-limiting-advanced"),
		Service: &Service{
			ID: service.ID,
		},
		Config: Configuration{
			"limit":       []interface{}{float64(10)},
			"window_size": []interface{}{float64(60)},
			"strategy":    "cluster",
			"sync_rate":   float64(-1),
		},
	}
	plugin, err = client.Plugins.Create(defaultCtx, plugin)
	require.NoError(t, err)
	require.NotNil(t, plugin)
	defer func() {
		_ = client.Plugins.Delete(defaultCtx, plugin.ID)
	}()

	t.Run("CreateForService without ID (POST) create/get/update/delete cycle", func(t *testing.T) {
		deco := &GraphqlRateLimitingCostDecoration{
			Service:      service,
			TypePath:     String("Vehicle.name"),
			AddConstant:  Float64(8),
			MulArguments: []*string{String("first")},
		}

		createdDeco, err := client.GraphqlRateLimitingCostDecorations.CreateForService(defaultCtx, deco)
		require.NoError(t, err)
		require.NotNil(t, createdDeco)

		deco, err = client.GraphqlRateLimitingCostDecorations.Get(defaultCtx, createdDeco.ID)
		require.NoError(t, err)
		require.NotNil(t, deco)

		deco.Service = service
		deco.TypePath = String("car.designation")
		deco, err = client.GraphqlRateLimitingCostDecorations.UpdateForService(defaultCtx, deco)
		require.NoError(t, err)
		require.NotNil(t, deco)
		require.Equal(t, "car.designation", *deco.TypePath)

		err = client.GraphqlRateLimitingCostDecorations.DeleteForService(defaultCtx, createdDeco)
		require.NoError(t, err)
	})

	t.Run("CreateForService with ID (PUT) create/get/update/delete cycle", func(t *testing.T) {
		id := uuid.NewString()
		deco := &GraphqlRateLimitingCostDecoration{
			ID:           String(id),
			Service:      service,
			TypePath:     String("Vehicle.name"),
			AddConstant:  Float64(8),
			MulArguments: []*string{String("first")},
		}

		createdDeco, err := client.GraphqlRateLimitingCostDecorations.CreateForService(defaultCtx, deco)
		require.NoError(t, err)
		require.NotNil(t, createdDeco)
		require.Equal(t, id, *createdDeco.ID)

		deco, err = client.GraphqlRateLimitingCostDecorations.Get(defaultCtx, createdDeco.ID)
		require.NoError(t, err)
		require.NotNil(t, deco)
		require.Equal(t, id, *deco.ID)

		deco.Service = service
		deco.TypePath = String("car.designation")
		deco, err = client.GraphqlRateLimitingCostDecorations.UpdateForService(defaultCtx, deco)
		require.NoError(t, err)
		require.NotNil(t, deco)
		require.Equal(t, id, *deco.ID)
		require.Equal(t, "car.designation", *deco.TypePath)

		err = client.GraphqlRateLimitingCostDecorations.DeleteForService(defaultCtx, createdDeco)
		require.NoError(t, err)
	})

	t.Run("Can't CreateForService Cost Decoration without service", func(t *testing.T) {
		deco := &GraphqlRateLimitingCostDecoration{
			TypePath:     String("car.name"),
			AddConstant:  Float64(8),
			MulArguments: []*string{String("first")},
		}

		createdDeco, err := client.GraphqlRateLimitingCostDecorations.CreateForService(defaultCtx, deco)
		require.Error(t, err, "cannot create a cost decoration for a service without name or ID")
		require.Nil(t, createdDeco)
	})
}
