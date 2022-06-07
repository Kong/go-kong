//go:build enterprise
// +build enterprise

package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceService(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	workspace := &Workspace{
		Name: String("teamA"),
		Meta: map[string]interface{}{
			"color":     "#814CA6",
			"thumbnail": nil,
		},
	}

	createdWorkspace, err := client.Workspaces.Create(defaultCtx, workspace)
	assert.NoError(err)
	require.NotNil(createdWorkspace)

	workspace, err = client.Workspaces.Get(defaultCtx, createdWorkspace.ID)
	assert.NoError(err)
	assert.NotNil(workspace)

	exists, err := client.Workspaces.Exists(defaultCtx, createdWorkspace.ID)
	assert.NoError(err)
	assert.True(exists)

	exists, err = client.Workspaces.ExistsByName(defaultCtx, createdWorkspace.Name)
	assert.NoError(err)
	assert.True(exists)

	fakeID := *createdWorkspace.ID + "garbage"
	exists, err = client.Workspaces.Exists(defaultCtx, &fakeID)
	assert.NoError(err)
	assert.False(exists)

	fakeName := *createdWorkspace.Name + "garbage"
	exists, err = client.Workspaces.ExistsByName(defaultCtx, &fakeName)
	assert.NoError(err)
	assert.False(exists)

	workspace.Comment = String("new comment")
	workspace, err = client.Workspaces.Update(defaultCtx, workspace)
	assert.NoError(err)
	assert.NotNil(workspace)
	assert.NotNil(workspace.Config)
	assert.Equal("teamA", *workspace.Name)
	assert.Equal("new comment", *workspace.Comment)
	assert.Equal("#814CA6", workspace.Meta["color"])

	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	assert.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	workspace = &Workspace{
		Name: String("teamB"),
		ID:   String(id),
	}

	createdWorkspace, err = client.Workspaces.Create(defaultCtx, workspace)
	assert.NoError(err)
	assert.NotNil(createdWorkspace)
	assert.Equal(id, *createdWorkspace.ID)

	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	assert.NoError(err)
}

func TestWorkspaceServiceList(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	require.NotNil(client)

	workspaceA := &Workspace{
		Name: String("teamA"),
	}
	workspaceB := &Workspace{
		Name: String("teamB"),
	}

	createdWorkspaceA, err := client.Workspaces.Create(defaultCtx, workspaceA)
	require.NoError(err)
	createdWorkspaceB, err := client.Workspaces.Create(defaultCtx, workspaceB)
	require.NoError(err)
	// paged List
	page1, next, err := client.Workspaces.List(defaultCtx, &ListOpt{Size: 1})
	assert.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	// nil ListOpt List
	workspaces, next, err := client.Workspaces.List(defaultCtx, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(workspaces)
	// Counts default workspace
	assert.Equal(3, len(workspaces))

	err = client.Workspaces.Delete(defaultCtx, createdWorkspaceA.ID)
	assert.NoError(err)
	err = client.Workspaces.Delete(defaultCtx, createdWorkspaceB.ID)
	assert.NoError(err)
}

func TestWorkspaceServiceListAll(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	workspaceA := &Workspace{
		Name: String("teamA"),
	}
	workspaceB := &Workspace{
		Name: String("teamB"),
	}

	createdWorkspaceA, err := client.Workspaces.Create(defaultCtx, workspaceA)
	require.NoError(err)
	createdWorkspaceB, err := client.Workspaces.Create(defaultCtx, workspaceB)
	require.NoError(err)

	workspaces, err := client.Workspaces.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(workspaces)
	// Counts default workspace
	assert.Equal(3, len(workspaces))

	err = client.Workspaces.Delete(defaultCtx, createdWorkspaceA.ID)
	require.NoError(err)
	err = client.Workspaces.Delete(defaultCtx, createdWorkspaceB.ID)
	require.NoError(err)
}

// Workspace entities

func TestWorkspaceService_Entities(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0 <=2.0.5", RequiredFeatures{})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	workspace := &Workspace{
		Name: String("teamA"),
		Meta: map[string]interface{}{
			"color":     "#814CA6",
			"thumbnail": nil,
		},
	}

	// Create a workspace
	createdWorkspace, err := client.Workspaces.Create(defaultCtx, workspace)
	assert.NoError(err)
	assert.NotNil(createdWorkspace)

	service := &Service{
		Name: String("foo"),
		Host: String("upstream"),
		Port: Int(42),
		Path: String("/path"),
	}

	// Create a service
	createdService, err := client.Services.Create(defaultCtx, service)
	assert.NoError(err)
	assert.NotNil(createdService)

	// Add the service to the workspace
	entities, err := client.Workspaces.AddEntities(
		defaultCtx, createdWorkspace.ID, createdService.ID)
	assert.NoError(err)
	assert.NotNil(entities)

	// List Entities attached to the workspace
	entitiesAdded, err := client.Workspaces.ListEntities(defaultCtx, createdWorkspace.ID)
	assert.NoError(err)
	assert.NotNil(entitiesAdded)
	// The two entities are records capturing the service name and id
	assert.Equal(2, len(entitiesAdded))

	// Delete the service from the workspace
	err = client.Workspaces.DeleteEntities(defaultCtx, createdWorkspace.ID, createdService.ID)
	assert.NoError(err)

	// Delete the service
	err = client.Services.Delete(defaultCtx, createdService.ID)
	assert.NoError(err)

	// Delete the workspace
	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	assert.NoError(err)
}
