package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRBACEndpointPermissionservice(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
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
	workspaced, err := NewTestClient(String(defaultBaseURL+"/endpoint-test-workspace"), nil)
	role := &RBACRole{
		Name: String("test-role-endpoint-perm"),
	}

	createdRole, err := workspaced.RBACRoles.Create(defaultCtx, role)
	assert.Nil(err)
	assert.NotNil(createdRole)

	// Add Endpoint Permission to Role
	ep := &RBACEndpointPermission{
		Role: &RBACRole{
			ID: createdRole.ID,
		},
		Endpoint: String("/rbac"),
		Actions: []*string{
			String("create"),
			String("read"),
		},
	}

	createdEndpointPermission, err := workspaced.RBACEndpointPermissions.Create(defaultCtx, ep)
	assert.Nil(err)
	assert.NotNil(createdEndpointPermission)

	ep, err = workspaced.RBACEndpointPermissions.Get(
		defaultCtx, createdRole.ID, createdWorkspace.Name, createdEndpointPermission.Endpoint)
	assert.Nil(err)
	assert.NotNil(ep)

	negative := true
	ep.Comment = String("new comment")
	ep.Negative = &negative
	ep, err = workspaced.RBACEndpointPermissions.Update(defaultCtx, ep)
	assert.Nil(err)
	assert.NotNil(ep)
	assert.Equal("new comment", *ep.Comment)
	assert.Equal(negative, *ep.Negative)

	err = workspaced.RBACEndpointPermissions.Delete(
		defaultCtx, createdRole.ID, createdWorkspace.ID, createdEndpointPermission.Endpoint)
	assert.Nil(err)
	err = workspaced.RBACRoles.Delete(defaultCtx, createdRole.ID)
	assert.Nil(err)
	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	assert.Nil(err)
}
