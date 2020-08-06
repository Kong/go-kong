package kong

import (
	"net/url"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserService(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", true)
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	user := &User{
		Name:      String("newUser"),
		Enabled:   Bool(true),
		Comment:   String("testing"),
		UserToken: String("foo"),
	}

	createdUser, err := client.Users.Create(defaultCtx, user)
	assert.Nil(err)
	assert.NotNil(createdUser)

	user, err = client.Users.Get(defaultCtx, createdUser.ID)
	assert.Nil(err)
	assert.NotNil(user)

	user.Comment = String("new comment")
	user, err = client.Users.Update(defaultCtx, user)
	assert.Nil(err)
	assert.NotNil(user)
	assert.Equal("new comment", *user.Comment)

	err = client.Users.Delete(defaultCtx, createdUser.ID)
	assert.Nil(err)
}

func TestUserServiceWorkspace(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", true)
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	workspace := Workspace{
		Name: String("test-workspace"),
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, &workspace)
	assert.Nil(err)
	assert.NotNil(createdWorkspace)
	// Setup Workspace aware client
	url, err := url.Parse(defaultBaseURL)
	assert.Nil(err)
	url.Path = path.Join(url.Path, *createdWorkspace.Name)
	workspaceClient, err := NewTestClient(String(url.String()), nil)
	assert.Nil(err)
	assert.NotNil(workspaceClient)

	user := &User{
		Name:      String("newUser"),
		Enabled:   Bool(true),
		Comment:   String("testing"),
		UserToken: String("foo"),
	}

	createdUser, err := workspaceClient.Users.Create(defaultCtx, user)
	assert.Nil(err)
	assert.NotNil(createdUser)

	user, err = workspaceClient.Users.Get(defaultCtx, createdUser.ID)
	assert.Nil(err)
	assert.NotNil(user)

	user.Comment = String("new comment")
	user, err = workspaceClient.Users.Update(defaultCtx, user)
	assert.Nil(err)
	assert.NotNil(user)
	assert.Equal("new comment", *user.Comment)

	err = workspaceClient.Users.Delete(defaultCtx, createdUser.ID)
	assert.Nil(err)

	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.Name)
	assert.Nil(err)
}

// TODO: After implementing roles service, test interaction with users
