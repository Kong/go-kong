package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeveloperRoleService(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{Portal: true})
	// NOTE: Developer Portal is not available in Kong < 3.5.0. Requires special config/license to enable.
	RunWhenEnterprise(T, "<3.5.0", RequiredFeatures{Portal: true})
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	testWs, err := NewTestWorkspace(client, "default")
	require.NoError(err)
	T.Cleanup(func() {
		require.NoError(testWs.Reset())
	})

	require.NoError(testWs.UpdateConfig(map[string]interface{}{"portal": true}))
	role := &DeveloperRole{
		Name: String("roleA"),
	}

	createdRole, err := client.DeveloperRoles.Create(defaultCtx, role)
	require.NoError(err)
	assert.NotNil(createdRole)

	role, err = client.DeveloperRoles.Get(defaultCtx, createdRole.ID)
	require.NoError(err)
	require.NotNil(role)

	role.Comment = String("new comment")
	role, err = client.DeveloperRoles.Update(defaultCtx, role)
	require.NoError(err)
	assert.NotNil(role)
	assert.Equal("roleA", *role.Name)

	err = client.DeveloperRoles.Delete(defaultCtx, createdRole.ID)
	require.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	role = &DeveloperRole{
		Name: String("teamB"),
		ID:   String(id),
	}

	createdRole, err = client.DeveloperRoles.Create(defaultCtx, role)
	require.NoError(err)
	assert.NotNil(createdRole)
	assert.Equal(id, *createdRole.ID)

	err = client.DeveloperRoles.Delete(defaultCtx, createdRole.ID)
	require.NoError(err)
}

func TestDeveloperRoleServiceList(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{Portal: true})
	// NOTE: Developer Portal is not available in Kong < 3.5.0. Requires special config/license to enable.
	RunWhenEnterprise(T, "<3.5.0", RequiredFeatures{Portal: true})
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	testWs, err := NewTestWorkspace(client, "default")
	require.NoError(err)
	T.Cleanup(func() {
		require.NoError(testWs.Reset())
	})

	require.NoError(testWs.UpdateConfig(map[string]interface{}{"portal": true}))

	roleA := &DeveloperRole{
		Name: String("roleA"),
	}
	roleB := &DeveloperRole{
		Name: String("roleB"),
	}

	createdRoleA, err := client.DeveloperRoles.Create(defaultCtx, roleA)
	require.NoError(err)
	T.Cleanup(func() {
		require.NoError(client.DeveloperRoles.Delete(defaultCtx, createdRoleA.ID))
	})

	createdRoleB, err := client.DeveloperRoles.Create(defaultCtx, roleB)
	require.NoError(err)
	T.Cleanup(func() {
		require.NoError(client.DeveloperRoles.Delete(defaultCtx, createdRoleB.ID))
	})

	roles, next, err := client.DeveloperRoles.List(defaultCtx, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(roles)
	assert.Len(roles, 2)
}

func TestDeveloperRoleListEndpoint(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{Portal: true})
	// NOTE: Developer Portal is not available in Kong < 3.5.0. Requires special config/license to enable.
	RunWhenEnterprise(T, "<3.5.0", RequiredFeatures{Portal: true})
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	testWs, err := NewTestWorkspace(client, "default")
	require.NoError(err)
	T.Cleanup(func() {
		require.NoError(testWs.Reset())
	})
	require.NoError(testWs.UpdateConfig(map[string]interface{}{"portal": true}))

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
		require.NoError(err)
		assert.NotNil(role)
		T.Cleanup(func() {
			id := *role.ID
			require.NoError(client.DeveloperRoles.Delete(defaultCtx, &id))
		})
		roles[i] = role
	}

	rolesFromKong, next, err := client.DeveloperRoles.List(defaultCtx, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(rolesFromKong)
	assert.Len(rolesFromKong, 3)

	// check if we see all developerRoles
	assert.True(compareDeveloperRoles(roles, rolesFromKong))

	// Test pagination
	rolesFromKong = []*DeveloperRole{}

	// first page
	page1, next, err := client.DeveloperRoles.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)
	rolesFromKong = append(rolesFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.DeveloperRoles.List(defaultCtx, next)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 2)
	rolesFromKong = append(rolesFromKong, page2...)

	assert.True(compareDeveloperRoles(roles, rolesFromKong))

	roles, err = client.DeveloperRoles.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(roles)
	assert.Len(roles, 3)
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
