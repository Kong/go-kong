package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRBACRoleService(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	workspace := &Workspace{
		Name: String("rbac-role-test-workspace"),
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, workspace)

	require.NoError(T, err)
	assert.NotNil(createdWorkspace)

	workspaced, err := NewTestClient(String(defaultBaseURL+"/rbac-role-test-workspace"), nil)
	require.NoError(T, err)

	role := &RBACRole{
		Name: String("roleA"),
	}

	createdRole, err := workspaced.RBACRoles.Create(defaultCtx, role)
	require.NoError(T, err)
	assert.NotNil(createdRole)

	role, err = workspaced.RBACRoles.Get(defaultCtx, createdRole.ID)
	require.NoError(T, err)
	assert.NotNil(role)

	role.Comment = String("new comment")
	role, err = workspaced.RBACRoles.Update(defaultCtx, role)
	require.NoError(T, err)
	assert.NotNil(role)
	assert.Equal("roleA", *role.Name)

	err = workspaced.RBACRoles.Delete(defaultCtx, createdRole.ID)
	require.NoError(T, err)

	// ID can be specified
	id := uuid.NewString()
	role = &RBACRole{
		Name: String("teamB"),
		ID:   String(id),
	}

	createdRole, err = workspaced.RBACRoles.Create(defaultCtx, role)
	require.NoError(T, err)
	assert.NotNil(createdRole)
	assert.Equal(id, *createdRole.ID)

	err = workspaced.RBACRoles.Delete(defaultCtx, createdRole.ID)
	require.NoError(T, err)

	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	require.NoError(T, err)
}

func TestRBACRoleServiceList(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	workspace := &Workspace{
		Name: String("rbac-role-list-test-workspace"),
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, workspace)

	require.NoError(T, err)
	assert.NotNil(createdWorkspace)

	workspaced, err := NewTestClient(String(defaultBaseURL+"/rbac-role-list-test-workspace"), nil)
	require.NoError(T, err)

	roleA := &RBACRole{
		Name: String("roleA"),
	}
	roleB := &RBACRole{
		Name: String("roleB"),
	}

	createdRoleA, err := workspaced.RBACRoles.Create(defaultCtx, roleA)
	require.NoError(T, err)
	createdRoleB, err := workspaced.RBACRoles.Create(defaultCtx, roleB)
	require.NoError(T, err)

	roles, next, err := workspaced.RBACRoles.List(defaultCtx, nil)
	require.NoError(T, err)
	assert.Nil(next)
	assert.NotNil(roles)
	assert.Len(roles, 2)

	err = workspaced.RBACRoles.Delete(defaultCtx, createdRoleA.ID)
	require.NoError(T, err)
	err = workspaced.RBACRoles.Delete(defaultCtx, createdRoleB.ID)
	require.NoError(T, err)

	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	require.NoError(T, err)
}

func TestRBACRoleListEndpoint(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	workspace := &Workspace{
		Name: String("rbac-role-list-endpoint-test-workspace"),
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, workspace)

	require.NoError(T, err)
	assert.NotNil(createdWorkspace)

	workspaced, err := NewTestClient(String(defaultBaseURL+"/rbac-role-list-endpoint-test-workspace"), nil)
	require.NoError(T, err)

	// fixtures
	roles := []*RBACRole{
		{
			Name: String("roleA"),
		},
		{
			Name: String("roleB"),
		},
		{
			Name: String("roleC"),
		},
	}

	// create fixturs
	for i := 0; i < len(roles); i++ {
		role, err := workspaced.RBACRoles.Create(defaultCtx, roles[i])
		require.NoError(T, err)
		assert.NotNil(role)
		roles[i] = role
	}

	rolesFromKong, next, err := workspaced.RBACRoles.List(defaultCtx, nil)
	require.NoError(T, err)
	assert.Nil(next)
	assert.NotNil(rolesFromKong)
	assert.Len(rolesFromKong, 3)

	// check if we see all roles
	assert.True(compareRBACRoles(roles, rolesFromKong))

	// Test pagination
	rolesFromKong = []*RBACRole{}

	// first page
	page1, next, err := workspaced.RBACRoles.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(T, err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)
	rolesFromKong = append(rolesFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := workspaced.RBACRoles.List(defaultCtx, next)
	require.NoError(T, err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 2)
	rolesFromKong = append(rolesFromKong, page2...)

	assert.True(compareRBACRoles(roles, rolesFromKong))

	roles, err = workspaced.RBACRoles.ListAll(defaultCtx)
	require.NoError(T, err)
	assert.NotNil(roles)
	assert.Len(roles, 3)

	for i := 0; i < len(roles); i++ {
		require.NoError(T, workspaced.RBACRoles.Delete(defaultCtx, roles[i].ID))
	}

	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	require.NoError(T, err)
}

func compareRBACRoles(expected, actual []*RBACRole) bool {
	var expectedNames, actualNames []string
	for _, role := range expected {
		expectedNames = append(expectedNames, *role.Name)
	}

	for _, role := range actual {
		actualNames = append(actualNames, *role.Name)
	}

	return (compareSlices(expectedNames, actualNames))
}
