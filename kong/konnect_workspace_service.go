package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

type AbstractKonnectWorkspaceService interface {
	ExistsByName(ctx context.Context, name *string) (bool, error)

	Create(ctx context.Context, workspace *Workspace) (*Workspace, error)

	Get(ctx context.Context, nameOrID *string) (*Workspace, error)

	ListAll(ctx context.Context) ([]*Workspace, error)
}

// KonnectWorkspaceService handles Workspaces in Konnect.
type KonnectWorkspaceService service

var _ AbstractKonnectWorkspaceService = &KonnectWorkspaceService{}

// ExistsByName checks for a Workspace by name in the given control plane.
func (s *KonnectWorkspaceService) ExistsByName(ctx context.Context, name *string) (bool, error) {
	if isEmptyString(name) {
		return false, fmt.Errorf("name cannot be nil")
	}
	endpoint := fmt.Sprintf("/workspaces/%s", *name)
	return s.client.existsForKonnect(ctx, endpoint)
}

// Create creates a new Workspace in the given control plane in Konnect.
func (s *KonnectWorkspaceService) Create(ctx context.Context, workspace *Workspace) (*Workspace, error) {
	if workspace == nil {
		return nil, fmt.Errorf("cannot create a nil workspace")
	}
	endpoint := "/workspaces"
	req, err := s.client.NewKonnectWorkspaceRequest("POST", endpoint, nil, workspace)
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

// Get fetches a Workspace by name or ID from the given control plane.
func (s *KonnectWorkspaceService) Get(ctx context.Context, nameOrID *string) (*Workspace, error) {
	if isEmptyString(nameOrID) {
		return nil, fmt.Errorf("nameOrID cannot be nil")
	}
	endpoint := fmt.Sprintf("/workspaces/%s", *nameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}
	var ws Workspace
	_, err = s.client.Do(ctx, req, &ws)
	if err != nil {
		return nil, err
	}
	return &ws, nil
}

// ListAll fetches all Workspaces from the given control plane.
func (s *KonnectWorkspaceService) ListAll(ctx context.Context) ([]*Workspace, error) {
	var workspaces []*Workspace
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		endpoint := "/workspaces"
		data, next, err := s.client.list(ctx, endpoint, opt)
		if err != nil {
			return nil, err
		}
		for _, object := range data {
			b, err := object.MarshalJSON()
			if err != nil {
				return nil, err
			}
			var ws Workspace
			err = json.Unmarshal(b, &ws)
			if err != nil {
				return nil, err
			}
			workspaces = append(workspaces, &ws)
		}
		opt = next
	}
	return workspaces, nil
}
