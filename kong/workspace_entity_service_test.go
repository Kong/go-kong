// +build enterprise

package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkspaceEntitiesService(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", false)
	assert := assert.New(T)

	client, err := NewClient(nil, nil)
	assert.Nil(err)
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
	assert.Nil(err)
	assert.NotNil(createdWorkspace)

	service := &Service{
		Name: String("foo"),
		Host: String("upstream"),
		Port: Int(42),
		Path: String("/path"),
	}

	// Create a service
	createdService, err := client.Services.Create(defaultCtx, service)
	assert.Nil(err)
	assert.NotNil(createdService)

	// Add the service to the workspace
	entities, err := client.WorkspaceEntities.Add(
		defaultCtx, createdWorkspace.ID, createdService.ID)
	assert.Nil(err)
	assert.NotNil(entities)

	// List Entities attached to the workspace
	entitiesAdded, err := client.WorkspaceEntities.List(defaultCtx, createdWorkspace.ID)
	assert.Nil(err)
	assert.NotNil(entitiesAdded)
	// The two entities are records capturing the service name and id
	assert.Equal(2, len(entitiesAdded))

	// Delete the service from the workspace
	err = client.WorkspaceEntities.Delete(defaultCtx, createdWorkspace.ID, createdService.ID)
	assert.Nil(err)

	// Delete the service
	err = client.Services.Delete(defaultCtx, createdService.ID)
	assert.Nil(err)

	// Delete the workspace
	err = client.Workspaces.Delete(defaultCtx, createdWorkspace.ID)
	assert.Nil(err)

}
