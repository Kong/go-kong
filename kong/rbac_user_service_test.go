package kong

import (
	"net/url"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRBACUserService(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	user := &RBACUser{
		Name:      String("newUser"),
		Enabled:   Bool(true),
		Comment:   String("testing"),
		UserToken: String("foo"),
	}

	createdUser, err := client.RBACUsers.Create(defaultCtx, user)
	require.NoError(T, err)
	assert.NotNil(createdUser)

	user, err = client.RBACUsers.Get(defaultCtx, createdUser.ID)
	require.NoError(T, err)
	assert.NotNil(user)

	user.Comment = String("new comment")
	user.UserTokenIdent = nil
	user, err = client.RBACUsers.Update(defaultCtx, user)
	require.NoError(T, err)
	assert.NotNil(user)
	assert.Equal("new comment", *user.Comment)

	err = client.RBACUsers.Delete(defaultCtx, createdUser.ID)
	require.NoError(T, err)
}

func TestRBACUserServiceWorkspace(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	workspace := Workspace{
		Name: String("test-workspace"),
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, &workspace)
	require.NoError(T, err)
	assert.NotNil(createdWorkspace)
	// Setup Workspace aware client
	url, err := url.Parse(defaultBaseURL)
	require.NoError(T, err)
	url.Path = path.Join(url.Path, *createdWorkspace.Name)
	workspaceClient, err := NewTestClient(String(url.String()), nil)
	require.NoError(T, err)
	assert.NotNil(workspaceClient)

	user := &RBACUser{
		Name:      String("newUser"),
		Enabled:   Bool(true),
		Comment:   String("testing"),
		UserToken: String("foo"),
	}

	createdUser, err := workspaceClient.RBACUsers.Create(defaultCtx, user)
	require.NoError(T, err)
	assert.NotNil(createdUser)

	user, err = workspaceClient.RBACUsers.Get(defaultCtx, createdUser.ID)
	require.NoError(T, err)
	assert.NotNil(user)

	user.Comment = String("new comment")
	user.UserTokenIdent = nil
	user, err = workspaceClient.RBACUsers.Update(defaultCtx, user)
	require.NoError(T, err)
	assert.NotNil(user)
	assert.Equal("new comment", *user.Comment)

	err = workspaceClient.RBACUsers.Delete(defaultCtx, createdUser.ID)
	require.NoError(T, err)

	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.Name)
	require.NoError(T, err)
}

func TestUserRoles(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)
	client, err := NewTestClient(nil, nil)

	require.NoError(T, err)
	assert.NotNil(client)

	roleA := &RBACRole{
		Name: String("roleA"),
	}
	roleB := &RBACRole{
		Name: String("roleB"),
	}

	createdRoleA, err := client.RBACRoles.Create(defaultCtx, roleA)
	require.NoError(T, err)
	createdRoleB, err := client.RBACRoles.Create(defaultCtx, roleB)
	require.NoError(T, err)

	ep := &RBACEndpointPermission{
		Role: &RBACRole{
			ID: createdRoleA.ID,
		},
		Endpoint: String("/rbac"),
		Actions: []*string{
			String("create"),
			String("read"),
		},
	}

	createdEndpointPermission, err := client.RBACEndpointPermissions.Create(defaultCtx, ep)
	require.NoError(T, err)
	assert.NotNil(createdEndpointPermission)

	user := &RBACUser{
		Name:      String("newUser"),
		Enabled:   Bool(true),
		Comment:   String("testing"),
		UserToken: String("foo"),
	}

	createdUser, err := client.RBACUsers.Create(defaultCtx, user)
	require.NoError(T, err)
	assert.NotNil(createdUser)

	roles := []*RBACRole{
		createdRoleA,
		createdRoleB,
	}

	updatedUser, err := client.RBACUsers.AddRoles(defaultCtx, createdUser.ID, roles)
	require.NoError(T, err)
	assert.NotNil(updatedUser)

	roleList, err := client.RBACUsers.ListRoles(defaultCtx, createdUser.ID)
	require.NoError(T, err)
	assert.NotNil(roleList)
	assert.Len(roleList, 2)

	permissionsList, err := client.RBACUsers.ListPermissions(defaultCtx, createdUser.ID)
	require.NoError(T, err)
	assert.NotNil(permissionsList)
	assert.Len(permissionsList.Endpoints, 1)

	err = client.RBACEndpointPermissions.Delete(
		defaultCtx, createdRoleA.ID, String("default"), createdEndpointPermission.Endpoint)
	require.NoError(T, err)
	err = client.RBACUsers.Delete(defaultCtx, createdUser.ID)
	require.NoError(T, err)
	err = client.RBACRoles.Delete(defaultCtx, createdRoleA.ID)
	require.NoError(T, err)
	err = client.RBACRoles.Delete(defaultCtx, createdRoleB.ID)
	require.NoError(T, err)
}
