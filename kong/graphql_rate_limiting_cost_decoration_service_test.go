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
}
