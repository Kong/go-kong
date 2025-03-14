//go:build enterprise

package kong

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminService(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	admin := &Admin{
		Email:            String("admin@test.com"),
		Username:         String("newAdmin"),
		CustomID:         String("admin123"),
		RBACTokenEnabled: Bool(true),
	}

	createdAdmin, err := client.Admins.Create(defaultCtx, admin)
	require.NoError(T, err)
	assert.NotNil(createdAdmin)

	admin, err = client.Admins.Get(defaultCtx, createdAdmin.ID)
	require.NoError(T, err)
	assert.NotNil(admin)

	admin.CustomID = String("admin321")
	admin, err = client.Admins.Update(defaultCtx, admin)
	require.NoError(T, err)
	assert.NotNil(admin)
	assert.Equal("admin321", *admin.CustomID)

	err = client.Admins.Delete(defaultCtx, createdAdmin.ID)
	require.NoError(T, err)
}

func TestAdminServiceWorkspace(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{})
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

	workspaceClient, err := NewTestClient(String(path.Join(defaultBaseURL, *createdWorkspace.Name)), nil)
	require.NoError(T, err)
	assert.NotNil(workspaceClient)

	admin := &Admin{
		Email:            String("admin@test.com"),
		Username:         String("newAdmin"),
		CustomID:         String("admin123"),
		RBACTokenEnabled: Bool(true),
	}

	createdAdmin, err := client.Admins.Create(defaultCtx, admin)
	require.NoError(T, err)
	assert.NotNil(createdAdmin)

	admin, err = client.Admins.Get(defaultCtx, createdAdmin.ID)
	require.NoError(T, err)
	assert.NotNil(admin)

	admin.CustomID = String("admin321")
	admin, err = client.Admins.Update(defaultCtx, admin)
	require.NoError(T, err)
	assert.NotNil(admin)
	assert.Equal("admin321", *admin.CustomID)

	err = client.Admins.Delete(defaultCtx, createdAdmin.ID)
	require.NoError(T, err)

	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.Name)
	require.NoError(T, err)
}

func TestAdminServiceList(T *testing.T) {
	assert := assert.New(T)
	client, err := NewTestClient(nil, nil)
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{})

	require.NoError(T, err)
	assert.NotNil(client)

	admin1 := &Admin{
		Email:            String("admin1@test.com"),
		Username:         String("newAdmin1"),
		CustomID:         String("admin1"),
		RBACTokenEnabled: Bool(true),
	}
	admin2 := &Admin{
		Email:            String("admin2@test.com"),
		Username:         String("newAdmin2"),
		CustomID:         String("admin2"),
		RBACTokenEnabled: Bool(true),
	}

	createdAdmin1, err := client.Admins.Create(defaultCtx, admin1)
	require.NoError(T, err)
	assert.NotNil(createdAdmin1)
	createdAdmin2, err := client.Admins.Create(defaultCtx, admin2)
	require.NoError(T, err)
	assert.NotNil(createdAdmin2)

	admins, _, err := client.Admins.List(defaultCtx, nil)
	require.NoError(T, err)
	assert.NotNil(admins)

	// Check if RBAC is enabled
	res, err := client.Root(defaultCtx)
	require.NoError(T, err)
	rbac := res["configuration"].(map[string]interface{})["rbac"].(string)
	expectedAdmins := 3
	if rbac == "off" {
		expectedAdmins = 2
	}
	require.Len(T, admins, expectedAdmins)

	err = client.Admins.Delete(defaultCtx, createdAdmin1.ID)
	require.NoError(T, err)
	err = client.Admins.Delete(defaultCtx, createdAdmin2.ID)
	require.NoError(T, err)
}

// XXX:
// This test requires RBAC to be enabled.
func TestAdminServiceRegisterCredentials(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	admin := &Admin{
		Email:            String("admin1@test.com"),
		Username:         String("newAdmin1"),
		CustomID:         String("admin1"),
		RBACTokenEnabled: Bool(true),
	}

	admin, err = client.Admins.Invite(defaultCtx, admin)
	require.NoError(T, err)
	assert.NotNil(admin)

	// Generate a new registration URL for the Admin
	admin, err = client.Admins.GenerateRegisterURL(defaultCtx, admin.ID)
	require.NoError(T, err)
	assert.NotNil(admin)

	admin.Password = String("bar")

	err = client.Admins.RegisterCredentials(defaultCtx, admin)
	require.NoError(T, err)

	admin, err = client.Admins.Get(defaultCtx, admin.ID)
	require.NoError(T, err)
	assert.NotNil(admin)

	err = client.Admins.Delete(defaultCtx, admin.ID)
	require.NoError(T, err)
}
