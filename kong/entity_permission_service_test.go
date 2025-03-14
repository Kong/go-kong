package kong

import (
	"net/url"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRBACEntityPermissionservice(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	// Create Workspace
	workspace := &Workspace{
		Name: String("entity-test-workspace"),
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, workspace)
	require.NoError(T, err)
	assert.NotNil(createdWorkspace)
	// Create new workspace client
	url, err := url.Parse(defaultBaseURL)
	require.NoError(T, err)
	url.Path = path.Join(url.Path, *createdWorkspace.Name)
	workspaceClient, err := NewTestClient(String(url.String()), nil)
	require.NoError(T, err)
	assert.NotNil(workspaceClient)
	// Use new client in workspace context.
	role := &RBACRole{
		Name: String("test-role-entity-perm"),
	}

	createdRole, err := workspaceClient.RBACRoles.Create(defaultCtx, role)
	require.NoError(T, err)
	assert.NotNil(createdRole)

	// Add Entity Permission to Role
	ep := &RBACEntityPermission{
		Role: &RBACRole{
			ID: createdRole.ID,
		},
		EntityID: String("*"),
		Actions: []*string{
			String("create"),
			String("read"),
		},
	}

	createdEntityPermission, err := workspaceClient.RBACEntityPermissions.Create(defaultCtx, ep)
	require.NoError(T, err)
	assert.NotNil(createdEntityPermission)

	ep, err = workspaceClient.RBACEntityPermissions.Get(defaultCtx, createdRole.ID, createdEntityPermission.EntityID)
	require.NoError(T, err)
	assert.NotNil(ep)

	negative := true
	ep.Comment = String("new comment")
	ep.Negative = &negative
	ep, err = workspaceClient.RBACEntityPermissions.Update(defaultCtx, ep)
	require.NoError(T, err)
	assert.NotNil(ep)
	assert.Equal("new comment", *ep.Comment)
	assert.Equal(negative, *ep.Negative)

	err = workspaceClient.RBACEntityPermissions.Delete(defaultCtx, createdRole.ID, createdEntityPermission.EntityID)
	require.NoError(T, err)
	err = workspaceClient.RBACRoles.Delete(defaultCtx, createdRole.ID)
	require.NoError(T, err)
	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	require.NoError(T, err)
}

func TestRBACEntityPermissionserviceList(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	// Create a workspace
	workspace := &Workspace{
		Name: String("entity-test-list-workspace"),
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, workspace)
	require.NoError(T, err)
	assert.NotNil(createdWorkspace)

	// Create a role
	role := &RBACRole{
		Name: String("test-role-entity-perm"),
	}
	createdRole, err := client.RBACRoles.Create(defaultCtx, role)
	require.NoError(T, err)
	assert.NotNil(createdRole)

	ep1 := &RBACEntityPermission{
		Role: &RBACRole{
			ID: createdRole.ID,
		},
		EntityID: String("*"),
		Actions: []*string{
			String("create"),
			String("read"),
			String("update"),
			String("delete"),
		},
	}
	ep2 := &RBACEntityPermission{
		Role: &RBACRole{
			ID: createdRole.ID,
		},
		EntityID:   createdWorkspace.ID,
		EntityType: String("workspaces"),
		Actions: []*string{
			String("update"),
			String("delete"),
		},
	}

	createdEntityPermissionA, err := client.RBACEntityPermissions.Create(defaultCtx, ep1)
	require.NoError(T, err)
	createdEntityPermissionB, err := client.RBACEntityPermissions.Create(defaultCtx, ep2)
	require.NoError(T, err)

	eps, err := client.RBACEntityPermissions.ListAllForRole(defaultCtx, createdRole.ID)
	require.NoError(T, err)
	assert.NotNil(eps)
	// Counts default ep
	assert.Len(eps, 2)

	err = client.RBACEntityPermissions.Delete(defaultCtx, createdRole.ID, createdEntityPermissionA.EntityID)
	require.NoError(T, err)
	err = client.RBACEntityPermissions.Delete(defaultCtx, createdRole.ID, createdEntityPermissionB.EntityID)
	require.NoError(T, err)
	err = client.RBACRoles.Delete(defaultCtx, createdRole.ID)
	require.NoError(T, err)
	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	require.NoError(T, err)
}
