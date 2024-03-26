package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeveloperRoleService(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{Portal: true})
	RunWhenEnterprise(T, "<3.7.0", RequiredFeatures{Portal: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	testWs, err := NewTestWorkspace(client, "default")
	assert.NoError(err)
	assert.NoError(testWs.UpdateConfig(map[string]interface{}{"portal": true}))

	role := &DeveloperRole{
		Name: String("roleA"),
	}

	createdRole, err := client.DeveloperRoles.Create(defaultCtx, role)
	assert.NoError(err)
	assert.NotNil(createdRole)

	role, err = client.DeveloperRoles.Get(defaultCtx, createdRole.ID)
	assert.NoError(err)
	assert.NotNil(role)

	role.Comment = String("new comment")
	role, err = client.DeveloperRoles.Update(defaultCtx, role)
	assert.NoError(err)
	assert.NotNil(role)
	assert.Equal("roleA", *role.Name)

	err = client.DeveloperRoles.Delete(defaultCtx, createdRole.ID)
	assert.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	role = &DeveloperRole{
		Name: String("teamB"),
		ID:   String(id),
	}

	createdRole, err = client.DeveloperRoles.Create(defaultCtx, role)
	assert.NoError(err)
	assert.NotNil(createdRole)
	assert.Equal(id, *createdRole.ID)

	err = client.DeveloperRoles.Delete(defaultCtx, createdRole.ID)
	assert.NoError(err)

	assert.NoError(testWs.Reset())
}

func TestDeveloperRoleServiceList(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{Portal: true})
	RunWhenEnterprise(T, "<3.7.0", RequiredFeatures{Portal: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	testWs, err := NewTestWorkspace(client, "default")
	assert.NoError(err)
	assert.NoError(testWs.UpdateConfig(map[string]interface{}{"portal": true}))

	roleA := &DeveloperRole{
		Name: String("roleA"),
	}
	roleB := &DeveloperRole{
		Name: String("roleB"),
	}

	createdRoleA, err := client.DeveloperRoles.Create(defaultCtx, roleA)
	assert.NoError(err)
	createdRoleB, err := client.DeveloperRoles.Create(defaultCtx, roleB)
	assert.NoError(err)

	roles, next, err := client.DeveloperRoles.List(defaultCtx, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(roles)
	assert.Equal(2, len(roles))

	err = client.DeveloperRoles.Delete(defaultCtx, createdRoleA.ID)
	assert.NoError(err)
	err = client.DeveloperRoles.Delete(defaultCtx, createdRoleB.ID)
	assert.NoError(err)

	assert.NoError(testWs.Reset())
}

func TestDeveloperRoleListEndpoint(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{Portal: true})
	RunWhenEnterprise(T, "<3.7.0", RequiredFeatures{Portal: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	testWs, err := NewTestWorkspace(client, "default")
	assert.NoError(err)
	assert.NoError(testWs.UpdateConfig(map[string]interface{}{"portal": true}))

	// fixtures
	roles := []*DeveloperRole{
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
		role, err := client.DeveloperRoles.Create(defaultCtx, roles[i])
		assert.NoError(err)
		assert.NotNil(role)
		roles[i] = role
	}

	rolesFromKong, next, err := client.DeveloperRoles.List(defaultCtx, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(rolesFromKong)
	assert.Equal(3, len(rolesFromKong))

	// check if we see all developerRoles
	assert.True(compareDeveloperRoles(roles, rolesFromKong))

	// Test pagination
	rolesFromKong = []*DeveloperRole{}

	// first page
	page1, next, err := client.DeveloperRoles.List(defaultCtx, &ListOpt{Size: 1})
	assert.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	rolesFromKong = append(rolesFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.DeveloperRoles.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Equal(2, len(page2))
	rolesFromKong = append(rolesFromKong, page2...)

	assert.True(compareDeveloperRoles(roles, rolesFromKong))

	roles, err = client.DeveloperRoles.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(roles)
	assert.Equal(3, len(roles))

	for i := 0; i < len(roles); i++ {
		assert.NoError(client.DeveloperRoles.Delete(defaultCtx, roles[i].ID))
	}

	assert.NoError(testWs.Reset())
}

func compareDeveloperRoles(expected, actual []*DeveloperRole) bool {
	var expectedNames, actualNames []string
	for _, role := range expected {
		expectedNames = append(expectedNames, *role.Name)
	}

	for _, role := range actual {
		actualNames = append(actualNames, *role.Name)
	}

	return (compareSlices(expectedNames, actualNames))
}
