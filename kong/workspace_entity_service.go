package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// WorkspaceEntityService handles WorkspaceEntitys in Kong.
type WorkspaceEntityService service

// Add adds entity ids given as a a comma delimited string
// to a given workspace in Kong. The response is a representation
// of the entity that was added to the workspace.
func (s *WorkspaceEntityService) Add(ctx context.Context,
	workspaceNameOrID *string, entityIds *string) (*[]map[string]interface{}, error) {

	if entityIds == nil {
		return nil, errors.New("entityIds cannot be nil")
	}

	endpoint := fmt.Sprintf("/workspaces/%v/entities", *workspaceNameOrID)
	var entities struct {
		Entities *string `json:"entities,omitempty"`
	}
	entities.Entities = entityIds

	req, err := s.client.NewRequest("POST", endpoint, nil, entities)

	if err != nil {
		return nil, err
	}

	var createdWorkspaceEntities []map[string]interface{}

	_, err = s.client.Do(ctx, req, &createdWorkspaceEntities)
	if err != nil {
		return nil, err
	}
	return &createdWorkspaceEntities, nil
}

// Delete deletes entity ids given as a a comma delimited string
// to a given workspace in Kong.
func (s *WorkspaceEntityService) Delete(ctx context.Context,
	workspaceNameOrID *string, entityIds *string) error {

	if entityIds == nil {
		return errors.New("entityIds cannot be nil")
	}

	endpoint := fmt.Sprintf("/workspaces/%v/entities", *workspaceNameOrID)
	var entities struct {
		Entities *string `json:"entities,omitempty"`
	}
	entities.Entities = entityIds

	req, err := s.client.NewRequest("DELETE", endpoint, nil, entities)

	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}
	return nil
}

// List fetches a list of all workspace entities in Kong.
func (s *WorkspaceEntityService) List(ctx context.Context,
	workspaceNameOrID *string) ([]*WorkspaceEntity, error) {

	endpoint := fmt.Sprintf("/workspaces/%v/entities", *workspaceNameOrID)

	data, _, err := s.client.list(ctx, endpoint, nil)
	if err != nil {
		return nil, err
	}
	var workspaceEntities []*WorkspaceEntity
	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, err
		}
		var workspaceEntity WorkspaceEntity
		err = json.Unmarshal(b, &workspaceEntity)
		if err != nil {
			return nil, err
		}
		workspaceEntities = append(workspaceEntities, &workspaceEntity)
	}

	return workspaceEntities, nil
}
