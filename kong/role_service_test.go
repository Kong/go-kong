package kong

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestRoleService(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", true)
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	role := &Role{
		Name: String("roleA"),
	}

	createdRole, err := client.Roles.Create(defaultCtx, role)
	assert.Nil(err)
	assert.NotNil(createdRole)

	role, err = client.Roles.Get(defaultCtx, createdRole.ID)
	assert.Nil(err)
	assert.NotNil(role)

	role.Comment = String("new comment")
	role, err = client.Roles.Update(defaultCtx, role)
	assert.Nil(err)
	assert.NotNil(role)
	assert.Equal("roleA", *role.Name)

	err = client.Roles.Delete(defaultCtx, createdRole.ID)
	assert.Nil(err)

	// ID can be specified
	id := uuid.NewV4().String()
	role = &Role{
		Name: String("teamB"),
		ID:   String(id),
	}

	createdRole, err = client.Roles.Create(defaultCtx, role)
	assert.Nil(err)
	assert.NotNil(createdRole)
	assert.Equal(id, *createdRole.ID)

	err = client.Roles.Delete(defaultCtx, createdRole.ID)
	assert.Nil(err)
}
func TestRoleServiceList(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", true)
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	roleA := &Role{
		Name: String("roleA"),
	}
	roleB := &Role{
		Name: String("roleB"),
	}

	createdRoleA, err := client.Roles.Create(defaultCtx, roleA)
	assert.Nil(err)
	createdRoleB, err := client.Roles.Create(defaultCtx, roleB)
	assert.Nil(err)

	roles, err := client.Roles.List(defaultCtx)
	assert.Nil(err)
	assert.NotNil(roles)
	// Counts default roles (super-admin, admin, read-only)
	assert.Equal(5, len(roles))

	err = client.Roles.Delete(defaultCtx, createdRoleA.ID)
	assert.Nil(err)
	err = client.Roles.Delete(defaultCtx, createdRoleB.ID)
	assert.Nil(err)
}
