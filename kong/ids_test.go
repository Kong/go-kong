package kong_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/go-kong/kong"
)

type fillEntityIDTestCase struct {
	name   string
	entity kong.IDFillable

	assertEntity func(t *testing.T, entity kong.IDFillable)
	expectErr    bool
}

func TestFillEntityID(t *testing.T) {
	testCases := []fillEntityIDTestCase{
		// Service
		{
			name:      "service nil pointer",
			entity:    (*kong.Service)(nil),
			expectErr: true,
		},
		{
			name:      "service with nil name",
			entity:    &kong.Service{},
			expectErr: true,
		},
		{
			name:      "service with empty name",
			entity:    &kong.Service{Name: kong.String("")},
			expectErr: true,
		},
		{
			name: "service with name",
			entity: &kong.Service{
				Name: kong.String("some.service.name"),
			},
			assertEntity: func(t *testing.T, e kong.IDFillable) {
				svc := e.(*kong.Service)
				require.NotNil(t, svc.ID)
				const expectedID = "d9bee1f8-db6e-5a37-9281-fd4aca16dc00"
				require.Equal(t, expectedID, *svc.ID, "ID should be deterministic")
			},
		},
		// Route
		{
			name:      "route nil pointer",
			entity:    (*kong.Route)(nil),
			expectErr: true,
		},
		{
			name:      "route with nil name",
			entity:    &kong.Route{},
			expectErr: true,
		},
		{
			name:      "route with empty name",
			entity:    &kong.Route{Name: kong.String("")},
			expectErr: true,
		},
		{
			name: "route with name",
			entity: &kong.Route{
				Name: kong.String("some.service.name"),
			},
			assertEntity: func(t *testing.T, e kong.IDFillable) {
				route := e.(*kong.Route)
				require.NotNil(t, route.ID)

				const expectedID = "7f0753cd-bff2-5f74-85b0-e8e2057d9500"
				require.Equal(t, expectedID, *route.ID, "ID should be deterministic")
			},
		},
		// Consumer
		{
			name:      "consumer nil pointer",
			entity:    (*kong.Consumer)(nil),
			expectErr: true,
		},
		{
			name:      "consumer with nil username",
			entity:    &kong.Consumer{},
			expectErr: true,
		},
		{
			name:      "consumer with empty username",
			entity:    &kong.Consumer{Username: kong.String("")},
			expectErr: true,
		},
		{
			name: "consumer with username",
			entity: &kong.Consumer{
				Username: kong.String("some.username"),
			},
			assertEntity: func(t *testing.T, e kong.IDFillable) {
				consumer := e.(*kong.Consumer)
				require.NotNil(t, consumer.ID)

				const expectedID = "782780c4-d345-5682-bdf6-722dccaf93e0"
				require.Equal(t, expectedID, *consumer.ID, "ID should be deterministic")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.entity.FillID()
			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			tc.assertEntity(t, tc.entity)
		})
	}
}
