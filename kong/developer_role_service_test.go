package kong

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeveloperRoleService(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", true)
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	role := &DeveloperRole{
		Name: String("roleA"),
	}

	createdRole, err := client.DeveloperRoles.Create(defaultCtx, role)
	assert.Nil(err)
	assert.NotNil(createdRole)

	role, err = client.DeveloperRoles.Get(defaultCtx, createdRole.ID)
	assert.Nil(err)
	assert.NotNil(role)

	role.Comment = String("new comment")
	role, err = client.DeveloperRoles.Update(defaultCtx, role)
	assert.Nil(err)
	assert.NotNil(role)
	assert.Equal("roleA", *role.Name)

	err = client.DeveloperRoles.Delete(defaultCtx, createdRole.ID)
	assert.Nil(err)

	// ID can be specified
	id := uuid.NewV4().String()
	role = &DeveloperRole{
		Name: String("teamB"),
		ID:   String(id),
	}

	createdRole, err = client.DeveloperRoles.Create(defaultCtx, role)
	assert.Nil(err)
	assert.NotNil(createdRole)
	assert.Equal(id, *createdRole.ID)

	err = client.DeveloperRoles.Delete(defaultCtx, createdRole.ID)
	assert.Nil(err)
}
func TestDeveloperRoleServiceList(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", true)
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	roleA := &DeveloperRole{
		Name: String("roleA"),
	}
	roleB := &DeveloperRole{
		Name: String("roleB"),
	}

	createdRoleA, err := client.DeveloperRoles.Create(defaultCtx, roleA)
	assert.Nil(err)
	createdRoleB, err := client.DeveloperRoles.Create(defaultCtx, roleB)
	assert.Nil(err)

	roles, err := client.DeveloperRoles.List(defaultCtx)
	assert.Nil(err)
	assert.NotNil(roles)
	// Counts default roles (super-admin, admin, read-only)
	assert.Equal(5, len(roles))

	err = client.DeveloperRoles.Delete(defaultCtx, createdRoleA.ID)
	assert.Nil(err)
	err = client.DeveloperRoles.Delete(defaultCtx, createdRoleB.ID)
	assert.Nil(err)
}
