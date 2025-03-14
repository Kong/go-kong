package kong

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRBACEndpointPermissionservice(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	// Create Workspace
	workspace := &Workspace{
		Name: String("endpoint-test-workspace"),
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, workspace)
	require.NoError(T, err)
	assert.NotNil(createdWorkspace)

	// Use new client in workspace context.
	workspaced, err := NewTestClient(String(defaultBaseURL+"/endpoint-test-workspace"), nil)
	require.NoError(T, err)
	role := &RBACRole{
		Name: String("test-role-endpoint-perm"),
	}

	createdRole, err := workspaced.RBACRoles.Create(defaultCtx, role)
	require.NoError(T, err)
	assert.NotNil(createdRole)

	// Add Endpoint Permission to Role
	origEp := &RBACEndpointPermission{
		Role: &RBACRole{
			ID: createdRole.ID,
		},
		Endpoint: String("/rbac"),
		Actions: []*string{
			String("create"),
			String("read"),
		},
	}

	createdEndpointPermission, err := workspaced.RBACEndpointPermissions.Create(defaultCtx, origEp)
	require.NoError(T, err)
	assert.NotNil(createdEndpointPermission)

	ep, err := workspaced.RBACEndpointPermissions.Get(
		defaultCtx, createdRole.ID, createdWorkspace.Name, createdEndpointPermission.Endpoint)
	require.NoError(T, err)
	assert.NotNil(ep)
	// we test this equality specifically because the Kong API handles this field oddly
	// see https://github.com/Kong/go-kong/pull/148
	var origActions []string
	for _, action := range origEp.Actions {
		origActions = append(origActions, *action)
	}
	var actions []string
	for _, action := range ep.Actions {
		actions = append(actions, *action)
	}
	sort.Strings(origActions)
	sort.Strings(actions)
	assert.Equal(origActions, actions)

	negative := true
	ep.Comment = String("new comment")
	ep.Negative = &negative
	ep, err = workspaced.RBACEndpointPermissions.Update(defaultCtx, ep)
	require.NoError(T, err)
	assert.NotNil(ep)
	assert.Equal("new comment", *ep.Comment)
	assert.Equal(negative, *ep.Negative)

	err = workspaced.RBACEndpointPermissions.Delete(
		defaultCtx, createdRole.ID, createdWorkspace.ID, createdEndpointPermission.Endpoint)
	require.NoError(T, err)
	err = workspaced.RBACRoles.Delete(defaultCtx, createdRole.ID)
	require.NoError(T, err)
	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	require.NoError(T, err)
}
