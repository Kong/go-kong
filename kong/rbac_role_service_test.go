package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRBACRoleService(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	workspace := &Workspace{
		Name: String("rbac-role-test-workspace"),
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, workspace)

	assert.NoError(err)
	assert.NotNil(createdWorkspace)

	workspaced, err := NewTestClient(String(defaultBaseURL+"/rbac-role-test-workspace"), nil)
	assert.NoError(err)

	role := &RBACRole{
		Name: String("roleA"),
	}

	createdRole, err := workspaced.RBACRoles.Create(defaultCtx, role)
	assert.NoError(err)
	assert.NotNil(createdRole)

	role, err = workspaced.RBACRoles.Get(defaultCtx, createdRole.ID)
	assert.NoError(err)
	assert.NotNil(role)

	role.Comment = String("new comment")
	role, err = workspaced.RBACRoles.Update(defaultCtx, role)
	assert.NoError(err)
	assert.NotNil(role)
	assert.Equal("roleA", *role.Name)

	err = workspaced.RBACRoles.Delete(defaultCtx, createdRole.ID)
	assert.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	role = &RBACRole{
		Name: String("teamB"),
		ID:   String(id),
	}

	createdRole, err = workspaced.RBACRoles.Create(defaultCtx, role)
	assert.NoError(err)
	assert.NotNil(createdRole)
	assert.Equal(id, *createdRole.ID)

	err = workspaced.RBACRoles.Delete(defaultCtx, createdRole.ID)
	assert.NoError(err)

	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	assert.NoError(err)
}

func TestRBACRoleServiceList(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	workspace := &Workspace{
		Name: String("rbac-role-list-test-workspace"),
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, workspace)

	assert.NoError(err)
	assert.NotNil(createdWorkspace)

	workspaced, err := NewTestClient(String(defaultBaseURL+"/rbac-role-list-test-workspace"), nil)
	assert.NoError(err)

	roleA := &RBACRole{
		Name: String("roleA"),
	}
	roleB := &RBACRole{
		Name: String("roleB"),
	}

	createdRoleA, err := workspaced.RBACRoles.Create(defaultCtx, roleA)
	assert.NoError(err)
	createdRoleB, err := workspaced.RBACRoles.Create(defaultCtx, roleB)
	assert.NoError(err)

	roles, next, err := workspaced.RBACRoles.List(defaultCtx, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(roles)
	assert.Equal(2, len(roles))

	err = workspaced.RBACRoles.Delete(defaultCtx, createdRoleA.ID)
	assert.NoError(err)
	err = workspaced.RBACRoles.Delete(defaultCtx, createdRoleB.ID)
	assert.NoError(err)

	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	assert.NoError(err)
}

func TestRBACRoleListEndpoint(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{RBAC: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	workspace := &Workspace{
		Name: String("rbac-role-list-endpoint-test-workspace"),
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, workspace)

	assert.NoError(err)
	assert.NotNil(createdWorkspace)

	workspaced, err := NewTestClient(String(defaultBaseURL+"/rbac-role-list-endpoint-test-workspace"), nil)
	assert.NoError(err)

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
		assert.NoError(err)
		assert.NotNil(role)
		roles[i] = role
	}

	rolesFromKong, next, err := workspaced.RBACRoles.List(defaultCtx, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(rolesFromKong)
	assert.Equal(3, len(rolesFromKong))

	// check if we see all roles
	assert.True(compareRBACRoles(roles, rolesFromKong))

	// Test pagination
	rolesFromKong = []*RBACRole{}

	// first page
	page1, next, err := workspaced.RBACRoles.List(defaultCtx, &ListOpt{Size: 1})
	assert.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	rolesFromKong = append(rolesFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := workspaced.RBACRoles.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Equal(2, len(page2))
	rolesFromKong = append(rolesFromKong, page2...)

	assert.True(compareRBACRoles(roles, rolesFromKong))

	roles, err = workspaced.RBACRoles.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(roles)
	assert.Equal(3, len(roles))

	for i := 0; i < len(roles); i++ {
		assert.Nil(workspaced.RBACRoles.Delete(defaultCtx, roles[i].ID))
	}

	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	assert.NoError(err)
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
