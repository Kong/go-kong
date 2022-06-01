package kong

import (
	"net/url"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRBACUserService(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	user := &RBACUser{
		Name:      String("newUser"),
		Enabled:   Bool(true),
		Comment:   String("testing"),
		UserToken: String("foo"),
	}

	createdUser, err := client.RBACUsers.Create(defaultCtx, user)
	assert.NoError(err)
	assert.NotNil(createdUser)

	user, err = client.RBACUsers.Get(defaultCtx, createdUser.ID)
	assert.NoError(err)
	assert.NotNil(user)

	user.Comment = String("new comment")
	user, err = client.RBACUsers.Update(defaultCtx, user)
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("new comment", *user.Comment)

	err = client.RBACUsers.Delete(defaultCtx, createdUser.ID)
	assert.NoError(err)
}

func TestRBACUserServiceWorkspace(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	workspace := Workspace{
		Name: String("test-workspace"),
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, &workspace)
	assert.NoError(err)
	assert.NotNil(createdWorkspace)
	// Setup Workspace aware client
	url, err := url.Parse(defaultBaseURL)
	assert.NoError(err)
	url.Path = path.Join(url.Path, *createdWorkspace.Name)
	workspaceClient, err := NewTestClient(String(url.String()), nil)
	assert.NoError(err)
	assert.NotNil(workspaceClient)

	user := &RBACUser{
		Name:      String("newUser"),
		Enabled:   Bool(true),
		Comment:   String("testing"),
		UserToken: String("foo"),
	}

	createdUser, err := workspaceClient.RBACUsers.Create(defaultCtx, user)
	assert.NoError(err)
	assert.NotNil(createdUser)

	user, err = workspaceClient.RBACUsers.Get(defaultCtx, createdUser.ID)
	assert.NoError(err)
	assert.NotNil(user)

	user.Comment = String("new comment")
	user, err = workspaceClient.RBACUsers.Update(defaultCtx, user)
	assert.NoError(err)
	assert.NotNil(user)
	assert.Equal("new comment", *user.Comment)

	err = workspaceClient.RBACUsers.Delete(defaultCtx, createdUser.ID)
	assert.NoError(err)

	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.Name)
	assert.NoError(err)
}

func TestUserRoles(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)
	client, err := NewTestClient(nil, nil)

	assert.NoError(err)
	assert.NotNil(client)

	roleA := &RBACRole{
		Name: String("roleA"),
	}
	roleB := &RBACRole{
		Name: String("roleB"),
	}

	createdRoleA, err := client.RBACRoles.Create(defaultCtx, roleA)
	assert.NoError(err)
	createdRoleB, err := client.RBACRoles.Create(defaultCtx, roleB)
	assert.NoError(err)

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
	assert.NoError(err)
	assert.NotNil(createdEndpointPermission)

	user := &RBACUser{
		Name:      String("newUser"),
		Enabled:   Bool(true),
		Comment:   String("testing"),
		UserToken: String("foo"),
	}

	createdUser, err := client.RBACUsers.Create(defaultCtx, user)
	assert.NoError(err)
	assert.NotNil(createdUser)

	roles := []*RBACRole{
		createdRoleA,
		createdRoleB,
	}

	updatedUser, err := client.RBACUsers.AddRoles(defaultCtx, createdUser.ID, roles)
	assert.NoError(err)
	assert.NotNil(updatedUser)

	roleList, err := client.RBACUsers.ListRoles(defaultCtx, createdUser.ID)
	assert.NoError(err)
	assert.NotNil(roleList)
	assert.Equal(2, len(roleList))

	permissionsList, err := client.RBACUsers.ListPermissions(defaultCtx, createdUser.ID)
	assert.NoError(err)
	assert.NotNil(permissionsList)
	assert.Equal(1, len(permissionsList.Endpoints))

	err = client.RBACEndpointPermissions.Delete(
		defaultCtx, createdRoleA.ID, String("default"), createdEndpointPermission.Endpoint)
	assert.NoError(err)
	err = client.RBACUsers.Delete(defaultCtx, createdUser.ID)
	assert.NoError(err)
	err = client.RBACRoles.Delete(defaultCtx, createdRoleA.ID)
	assert.NoError(err)
	err = client.RBACRoles.Delete(defaultCtx, createdRoleB.ID)
	assert.NoError(err)
}
