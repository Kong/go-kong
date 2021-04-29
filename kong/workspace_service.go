package kong

import (
	"net/http"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// AbstractWorkspaceService handles Workspaces in Kong.
type AbstractWorkspaceService interface {
	// Exists checks the exitence of a Workspace in Kong.
	Exists(ctx context.Context, nameOrID *string) (*bool, error)
	// Create creates a Workspace in Kong.
	Create(ctx context.Context, workspace *Workspace) (*Workspace, error)
	// Get fetches a Workspace in Kong.
	Get(ctx context.Context, nameOrID *string) (*Workspace, error)
	// Update updates a Workspace in Kong.
	Update(ctx context.Context, workspace *Workspace) (*Workspace, error)
	// Delete deletes a Workspace in Kong
	Delete(ctx context.Context, WorkspaceOrID *string) error
	// List fetches a list of all Workspaces in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*Workspace, *ListOpt, error)
	// ListAll fetches all workspaces in Kong.
	ListAll(ctx context.Context) ([]*Workspace, error)
	// AddEntities adds entity ids given as a a comma delimited string
	// to a given workspace in Kong. The response is a representation
	// of the entity that was added to the workspace.
	//
	// Deprecated: Kong 2.x removed this endpoint.
	AddEntities(ctx context.Context, workspaceNameOrID *string, entityIds *string) (*[]map[string]interface{}, error)
	// DeleteEntities deletes entity ids given as a a comma delimited string
	// to a given workspace in Kong.
	//
	// Deprecated: Kong 2.x removed this endpoint.
	DeleteEntities(ctx context.Context, workspaceNameOrID *string, entityIds *string) error
	// ListEntities fetches a list of all workspace entities in Kong.
	//
	// Deprecated: Kong 2.x removed this endpoint.
	ListEntities(ctx context.Context, workspaceNameOrID *string) ([]*WorkspaceEntity, error)
}

// WorkspaceService handles Workspaces in Kong.
type WorkspaceService service


// Exists checks the exitence of the Workspace in Kong.
func (s *WorkspaceService) Exists(ctx context.Context,
	nameOrID *string) (*bool, error) {
	if isEmptyString(nameOrID) {
		return nil, errors.New("nameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/workspaces/%v", *nameOrID)
	req, err := s.client.NewRequest("HEAD", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}
	var status = false
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		if IsNotFoundErr(err) {
			return &status, nil
		}
		return nil, err
	}
	status = resp.StatusCode == http.StatusOK
	return &status, nil
}

// Create creates a Workspace in Kong.
func (s *WorkspaceService) Create(ctx context.Context,
	workspace *Workspace) (*Workspace, error) {

	if workspace == nil {
		return nil, errors.New("cannot create a nil workspace")
	}

	endpoint := "/workspaces"
	method := "POST"
	if workspace.ID != nil {
		endpoint = endpoint + "/" + *workspace.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, endpoint, nil, workspace)

	if err != nil {
		return nil, err
	}

	var createdWorkspace Workspace
	_, err = s.client.Do(ctx, req, &createdWorkspace)
	if err != nil {
		return nil, err
	}
	return &createdWorkspace, nil
}

// Get fetches a Workspace in Kong.
func (s *WorkspaceService) Get(ctx context.Context,
	nameOrID *string) (*Workspace, error) {

	if isEmptyString(nameOrID) {
		return nil, errors.New("nameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/workspaces/%v", *nameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var Workspace Workspace
	_, err = s.client.Do(ctx, req, &Workspace)
	if err != nil {
		return nil, err
	}
	return &Workspace, nil
}

// Update updates a Workspace in Kong. Only updates to the
// `comment` field are supported. To rename a workspace use Create.
func (s *WorkspaceService) Update(ctx context.Context,
	workspace *Workspace) (*Workspace, error) {

	if workspace == nil {
		return nil, errors.New("cannot update a nil Workspace")
	}

	if isEmptyString(workspace.ID) {
		return nil, errors.New("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/workspaces/%v", *workspace.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, workspace)
	if err != nil {
		return nil, err
	}

	var updatedWorkspace Workspace
	_, err = s.client.Do(ctx, req, &updatedWorkspace)
	if err != nil {
		return nil, err
	}
	return &updatedWorkspace, nil
}

// Delete deletes a Workspace in Kong
func (s *WorkspaceService) Delete(ctx context.Context,
	WorkspaceOrID *string) error {

	if isEmptyString(WorkspaceOrID) {
		return errors.New("WorkspaceOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/workspaces/%v", *WorkspaceOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of all Workspaces in Kong.
func (s *WorkspaceService) List(ctx context.Context,
	opt *ListOpt) ([]*Workspace, *ListOpt, error) {

	data, next, err := s.client.list(ctx, "/workspaces/", opt)
	if err != nil {
		return nil, nil, err
	}
	var workspaces []*Workspace
	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var workspace Workspace
		err = json.Unmarshal(b, &workspace)
		if err != nil {
			return nil, nil, err
		}
		workspaces = append(workspaces, &workspace)
	}

	return workspaces, next, nil
}

// ListAll fetches all workspaces in Kong.
func (s *WorkspaceService) ListAll(ctx context.Context) ([]*Workspace, error) {

	var workspaces, data []*Workspace
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		workspaces = append(workspaces, data...)
	}

	return workspaces, nil
}

// AddEntities adds entity ids given as a a comma delimited string
// to a given workspace in Kong. The response is a representation
// of the entity that was added to the workspace.
//
// Deprecated: Kong 2.x removed this endpoint.
func (s *WorkspaceService) AddEntities(ctx context.Context,
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

// DeleteEntities deletes entity ids given as a a comma delimited string
// to a given workspace in Kong.
//
// Deprecated: Kong 2.x removed this endpoint.
func (s *WorkspaceService) DeleteEntities(ctx context.Context,
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

// ListEntities fetches a list of all workspace entities in Kong.
//
// Deprecated: Kong 2.x removed this endpoint.
func (s *WorkspaceService) ListEntities(ctx context.Context,
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
