package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestACLGroupCreate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	acl, err := client.ACLs.Create(defaultCtx,
		String("foo"), nil)
	require.Error(err)
	assert.Nil(acl)

	acl = &ACLGroup{}
	acl, err = client.ACLs.Create(defaultCtx, String(""),
		acl)
	require.Error(err)
	assert.Nil(acl)

	acl, err = client.ACLs.Create(defaultCtx,
		String("does-not-exist"), acl)
	require.Error(err)
	assert.Nil(acl)

	// consumer for the ACL group
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	acl = &ACLGroup{
		Group: String("my-group"),
	}
	createdACL, err := client.ACLs.Create(defaultCtx, consumer.ID, acl)
	require.NoError(err)
	assert.NotNil(createdACL)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestACLGroupCreateWithID(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	acl := &ACLGroup{
		ID:    String(uuid),
		Group: String("my-group"),
	}

	// consumer for the ACLGroup
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdACL, err := client.ACLs.Create(defaultCtx, consumer.ID, acl)
	require.NoError(err)
	require.NotNil(createdACL)

	assert.Equal(uuid, *createdACL.ID)
	assert.Equal("my-group", *createdACL.Group)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestACLGroupGet(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	acl := &ACLGroup{
		ID:    String(uuid),
		Group: String("my-group"),
	}

	// consumer for the ACLGroup
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdACL, err := client.ACLs.Create(defaultCtx, consumer.ID, acl)
	require.NoError(err)
	assert.NotNil(createdACL)

	aclGroup, err := client.ACLs.Get(defaultCtx, consumer.ID, acl.ID)
	require.NoError(err)
	assert.Equal("my-group", *aclGroup.Group)

	aclGroup, err = client.ACLs.Get(defaultCtx, consumer.ID, acl.Group)
	require.NoError(err)
	assert.Equal("my-group", *aclGroup.Group)

	aclGroup, err = client.ACLs.Get(defaultCtx, consumer.ID,
		String("does-not-exists"))
	assert.Nil(aclGroup)
	require.Error(err)

	aclGroup, err = client.ACLs.Get(defaultCtx, consumer.ID, String(""))
	assert.Nil(aclGroup)
	require.Error(err)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestACLGroupGetByID(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	acl := &ACLGroup{
		ID:    String(uuid),
		Group: String("my-group"),
	}

	// consumer for the ACLGroup
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdACL, err := client.ACLs.Create(defaultCtx, consumer.ID, acl)
	require.NoError(err)
	assert.NotNil(createdACL)

	aclGroup, err := client.ACLs.GetByID(defaultCtx, acl.ID)
	require.NoError(err)
	assert.Equal("my-group", *aclGroup.Group)

	aclGroup, err = client.ACLs.GetByID(defaultCtx, String("does-not-exist"))
	assert.Nil(aclGroup)
	require.Error(err)

	aclGroup, err = client.ACLs.GetByID(defaultCtx, String(""))
	assert.Nil(aclGroup)
	require.Error(err)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestACLGroupUpdate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	acl := &ACLGroup{
		ID:    String(uuid),
		Group: String("my-group"),
	}

	// consumer for the ACLGroup
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdACL, err := client.ACLs.Create(defaultCtx, consumer.ID, acl)
	require.NoError(err)
	assert.NotNil(createdACL)

	aclGroup, err := client.ACLs.Get(defaultCtx, consumer.ID, acl.ID)
	require.NoError(err)
	assert.Equal("my-group", *aclGroup.Group)

	acl.Group = String("my-new-group")
	updatedACLGroup, err := client.ACLs.Update(defaultCtx, consumer.ID, acl)
	require.NoError(err)
	assert.NotNil(updatedACLGroup)
	assert.Equal("my-new-group", *updatedACLGroup.Group)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestACLGroupDelete(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	uuid := uuid.NewString()
	acl := &ACLGroup{
		ID:    String(uuid),
		Group: String("my-group"),
	}

	// consumer for the ACLGroup
	consumer := &Consumer{
		Username: String("foo"),
	}

	consumer, err = client.Consumers.Create(defaultCtx, consumer)
	require.NoError(err)
	require.NotNil(consumer)

	createdACL, err := client.ACLs.Create(defaultCtx, consumer.ID, acl)
	require.NoError(err)
	assert.NotNil(createdACL)

	err = client.ACLs.Delete(defaultCtx, consumer.ID, acl.Group)
	require.NoError(err)

	aclGroup, err := client.ACLs.Get(defaultCtx, consumer.ID, acl.ID)
	require.Error(err)
	assert.Nil(aclGroup)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}

func TestACLGroupListMethods(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	// consumer for the ACLGroup
	consumer1 := &Consumer{
		Username: String("foo"),
	}

	consumer1, err = client.Consumers.Create(defaultCtx, consumer1)
	require.NoError(err)
	require.NotNil(consumer1)

	consumer2 := &Consumer{
		Username: String("bar"),
	}

	consumer2, err = client.Consumers.Create(defaultCtx, consumer2)
	require.NoError(err)
	require.NotNil(consumer2)

	// fixtures
	aclGroups := []*ACLGroup{
		{
			Group:    String("acl11"),
			Consumer: consumer1,
		},
		{
			Group:    String("acl12"),
			Consumer: consumer1,
		},
		{
			Group:    String("acl21"),
			Consumer: consumer2,
		},
		{
			Group:    String("acl22"),
			Consumer: consumer2,
		},
	}

	// create fixturs
	for i := 0; i < len(aclGroups); i++ {
		acl, err := client.ACLs.Create(defaultCtx,
			aclGroups[i].Consumer.ID, aclGroups[i])
		require.NoError(err)
		assert.NotNil(acl)
		aclGroups[i] = acl
	}

	aclGroupsFromKong, next, err := client.ACLs.List(defaultCtx, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(aclGroupsFromKong)
	assert.Len(aclGroupsFromKong, 4)

	// first page
	page1, next, err := client.ACLs.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)

	// last page
	next.Size = 3
	page2, next, err := client.ACLs.List(defaultCtx, next)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 3)

	aclGroupsForConsumer, next, err := client.ACLs.ListForConsumer(defaultCtx,
		consumer1.ID, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(aclGroupsForConsumer)
	assert.Len(aclGroupsForConsumer, 2)

	aclGroups, err = client.ACLs.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(aclGroups)
	assert.Len(aclGroups, 4)

	require.NoError(client.Consumers.Delete(defaultCtx, consumer1.ID))
	require.NoError(client.Consumers.Delete(defaultCtx, consumer2.ID))
}
