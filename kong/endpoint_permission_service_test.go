package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpointPermissionService(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", true)
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	// Create Workspace
	workspace := &Workspace{
		Name: String("endpoint-test-workspace"),
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, workspace)
	assert.Nil(err)
	assert.NotNil(createdWorkspace)

	// Use new client in workspace context.
	role := &RBACRole{
		Name: String("test-role-endpoint-perm"),
	}

	createdRole, err := client.Roles.Create(defaultCtx, role)
	assert.Nil(err)
	assert.NotNil(createdRole)

	// Add Endpoint Permission to Role
	ep := &EndpointPermission{
		Role: &PermissionRoleID{
			ID: createdRole.ID,
		},
		Endpoint: String("/rbac"),
		Actions: []*string{
			String("create"),
			String("read"),
		},
	}

	createdEndpointPermission, err := client.EndpointPermissions.Create(defaultCtx, ep)
	assert.Nil(err)
	assert.NotNil(createdEndpointPermission)

	ep, err = client.EndpointPermissions.Get(
		defaultCtx, createdRole.ID, createdWorkspace.ID, createdEndpointPermission.Endpoint)
	assert.Nil(err)
	assert.NotNil(ep)

	ep.Comment = String("new comment")
	ep, err = client.EndpointPermissions.Update(defaultCtx, ep)
	assert.Nil(err)
	assert.NotNil(ep)
	assert.Equal("new comment", *ep.Comment)

	err = client.EndpointPermissions.Delete(
		defaultCtx, createdRole.ID, createdWorkspace.ID, createdEndpointPermission.Endpoint)
	assert.Nil(err)
	err = client.Roles.Delete(defaultCtx, createdRole.ID)
	assert.Nil(err)
	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	assert.Nil(err)

}
