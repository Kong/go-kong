package kong_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kong/go-kong/kong"
)

type fillEntityIDTestCase struct {
	name         string
	entity       kong.IDFillable
	workspace    string
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
		{
			name: "consumer with username and workspace",
			entity: &kong.Consumer{
				Username: kong.String("some.username"),
			},
			workspace: "ws.1",
			assertEntity: func(t *testing.T, e kong.IDFillable) {
				consumer := e.(*kong.Consumer)
				require.NotNil(t, consumer.ID)

				const expectedID = "4701c235-f7f0-51ee-8782-0d8e145d6771"
				require.Equal(t, expectedID, *consumer.ID, "ID should be deterministic")
			},
		},
		// Consumer Group
		{
			name:      "consumer group nil pointer",
			entity:    (*kong.ConsumerGroup)(nil),
			expectErr: true,
		},
		{
			name:      "consumer group with nil name",
			entity:    &kong.ConsumerGroup{},
			expectErr: true,
		},
		{
			name:      "consumer group with empty name",
			entity:    &kong.ConsumerGroup{Name: kong.String("")},
			expectErr: true,
		},
		{
			name: "consumer group with name",
			entity: &kong.ConsumerGroup{
				Name: kong.String("some.consumer.group"),
			},
			assertEntity: func(t *testing.T, e kong.IDFillable) {
				cg := e.(*kong.ConsumerGroup)
				require.NotNil(t, cg.ID)

				const expectedID = "e5643801-37c6-5d04-9d3f-c1c84c747e90"
				require.Equal(t, expectedID, *cg.ID, "ID should be deterministic")
			},
		},
		// Vault
		{
			name:      "vault with nil pointer",
			entity:    (*kong.Vault)(nil),
			expectErr: true,
		},
		{
			name:      "vault with nil prefix",
			entity:    &kong.Vault{},
			expectErr: true,
		},
		{
			name: "vault with empty prefix",
			entity: &kong.Vault{
				Name: kong.String(""),
			},
			expectErr: true,
		},
		{
			name: "vault with prefix",
			entity: &kong.Vault{
				Name:   kong.String("env"),
				Prefix: kong.String("test-env"),
			},
			assertEntity: func(t *testing.T, e kong.IDFillable) {
				v := e.(*kong.Vault)
				require.NotNil(t, v.ID)

				const expectedID = "837665c3-856f-5ca2-9db4-52a1cf8a32be"
				require.Equal(t, expectedID, *v.ID, "ID should be deterministic")
			},
		},
		// Plugin
		{
			name:      "plugin with nil pointer",
			entity:    (*kong.Plugin)(nil),
			expectErr: true,
		},
		{
			name: "plugin with empty name",
			entity: &kong.Plugin{
				Name: kong.String(""),
			},
			expectErr: true,
		},
		{
			name: "plugin with id should not be modified",
			entity: &kong.Plugin{
				Name: kong.String("rate-limiting"),
				ID:   kong.String("abcd1234-5678-abcd-0123-abcdeffedcba"),
			},
			assertEntity: func(t *testing.T, e kong.IDFillable) {
				p := e.(*kong.Plugin)
				require.NotNil(t, p.ID)

				const expectedID = "abcd1234-5678-abcd-0123-abcdeffedcba"
				require.Equal(t, expectedID, *p.ID, "ID should not be changed")
			},
		},
		{
			name: "plugin with name",
			entity: &kong.Plugin{
				Name: kong.String("rate-limiting"),
			},
			assertEntity: func(t *testing.T, e kong.IDFillable) {
				p := e.(*kong.Plugin)
				require.NotNil(t, p.ID)

				const expectedID = "f0f8012b-f709-5685-9812-0130e5d83c5a"
				require.Equal(t, expectedID, *p.ID, "ID should be deterministic")
			},
		},
		{
			name: "plugin with name and service",
			entity: &kong.Plugin{
				Name: kong.String("rate-limiting"),
				Service: &kong.Service{
					Name: kong.String("service-1"),
				},
			},
			assertEntity: func(t *testing.T, e kong.IDFillable) {
				p := e.(*kong.Plugin)
				require.NotNil(t, p.ID)

				const expectedID = "b0df9683-8b2f-5557-8b66-204b1529ed7f"
				require.Equal(t, expectedID, *p.ID, "ID should be deterministic")
			},
		},
		{
			name: "plugin with with name, route and consumer",
			entity: &kong.Plugin{
				Name: kong.String("rate-limiting"),
				Route: &kong.Route{
					Name: kong.String("route-1"),
				},
				Consumer: &kong.Consumer{
					Username: kong.String("consumer-1"),
				},
			},
			assertEntity: func(t *testing.T, e kong.IDFillable) {
				p := e.(*kong.Plugin)
				require.NotNil(t, p.ID)

				const expectedID = "d6a8d8fd-f156-5809-a2d0-c75b23db07a5"
				require.Equal(t, expectedID, *p.ID, "ID should be deterministic")
			},
		},
		{
			name: "plugin with name, instance name and consumer group",
			entity: &kong.Plugin{
				Name:         kong.String("rate-limiting"),
				InstanceName: kong.String("rl-1"),
				ConsumerGroup: &kong.ConsumerGroup{
					Name: kong.String("group-1"),
				},
			},
			assertEntity: func(t *testing.T, e kong.IDFillable) {
				p := e.(*kong.Plugin)
				require.NotNil(t, p.ID)

				const expectedID = "01af8dbc-e3e8-5ccd-b20d-d55227c15cbf"
				require.Equal(t, expectedID, *p.ID, "ID should be deterministic")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.entity.FillID(tc.workspace)
			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			tc.assertEntity(t, tc.entity)
		})
	}
}
