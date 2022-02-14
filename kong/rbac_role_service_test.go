package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRBACRoleService(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", requiredFeatures{rbac: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	role := &RBACRole{
		Name: String("roleA"),
	}

	createdRole, err := client.RBACRoles.Create(defaultCtx, role)
	assert.Nil(err)
	assert.NotNil(createdRole)

	role, err = client.RBACRoles.Get(defaultCtx, createdRole.ID)
	assert.Nil(err)
	assert.NotNil(role)

	role.Comment = String("new comment")
	role, err = client.RBACRoles.Update(defaultCtx, role)
	assert.Nil(err)
	assert.NotNil(role)
	assert.Equal("roleA", *role.Name)

	err = client.RBACRoles.Delete(defaultCtx, createdRole.ID)
	assert.Nil(err)

	// ID can be specified
	id := uuid.NewString()
	role = &RBACRole{
		Name: String("teamB"),
		ID:   String(id),
	}

	createdRole, err = client.RBACRoles.Create(defaultCtx, role)
	assert.Nil(err)
	assert.NotNil(createdRole)
	assert.Equal(id, *createdRole.ID)

	err = client.RBACRoles.Delete(defaultCtx, createdRole.ID)
	assert.Nil(err)
}

func TestRBACRoleServiceList(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", requiredFeatures{rbac: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	roleA := &RBACRole{
		Name: String("roleA"),
	}
	roleB := &RBACRole{
		Name: String("roleB"),
	}

	createdRoleA, err := client.RBACRoles.Create(defaultCtx, roleA)
	assert.Nil(err)
	createdRoleB, err := client.RBACRoles.Create(defaultCtx, roleB)
	assert.Nil(err)

	roles, next, err := client.RBACRoles.List(defaultCtx, nil)
	assert.Nil(err)
	assert.Nil(next)
	assert.NotNil(roles)
	// Counts default roles (super-admin, admin, read-only)
	assert.Equal(5, len(roles))

	err = client.RBACRoles.Delete(defaultCtx, createdRoleA.ID)
	assert.Nil(err)
	err = client.RBACRoles.Delete(defaultCtx, createdRoleB.ID)
	assert.Nil(err)
}

func TestRBACRoleListEndpoint(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", requiredFeatures{rbac: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

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
		role, err := client.RBACRoles.Create(defaultCtx, roles[i])
		assert.Nil(err)
		assert.NotNil(role)
		roles[i] = role
	}

	rolesFromKong, next, err := client.RBACRoles.List(defaultCtx, nil)
	assert.Nil(err)
	assert.Nil(next)
	assert.NotNil(rolesFromKong)
	assert.Equal(3, len(rolesFromKong))

	// check if we see all roles
	assert.True(compareRBACRoles(roles, rolesFromKong))

	// Test pagination
	rolesFromKong = []*RBACRole{}

	// first page
	page1, next, err := client.RBACRoles.List(defaultCtx, &ListOpt{Size: 1})
	assert.Nil(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	rolesFromKong = append(rolesFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.RBACRoles.List(defaultCtx, next)
	assert.Nil(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Equal(2, len(page2))
	rolesFromKong = append(rolesFromKong, page2...)

	assert.True(compareRBACRoles(roles, rolesFromKong))

	roles, err = client.RBACRoles.ListAll(defaultCtx)
	assert.Nil(err)
	assert.NotNil(roles)
	assert.Equal(3, len(roles))

	for i := 0; i < len(roles); i++ {
		assert.Nil(client.RBACRoles.Delete(defaultCtx, roles[i].ID))
	}
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
