package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

// AbstractRBACEndpointPermissionService handles RBACEndpointPermissions in Kong.
type AbstractRBACEndpointPermissionService interface {
	// Create creates a RBACEndpointPermission in Kong.
	Create(ctx context.Context, ep *RBACEndpointPermission) (*RBACEndpointPermission, error)
	// Get fetches a RBACEndpointPermission in Kong.
	Get(ctx context.Context, roleNameOrID *string, workspaceNameOrID *string,
		endpointName *string) (*RBACEndpointPermission, error)
	// Update updates a RBACEndpointPermission in Kong.
	Update(ctx context.Context, ep *RBACEndpointPermission) (*RBACEndpointPermission, error)
	// Delete deletes a EndpointPermission in Kong
	Delete(ctx context.Context, roleNameOrID *string, workspaceNameOrID *string, endpoint *string) error
	// ListAllForRole fetches a list of all RBACEndpointPermissions in Kong for a given role.
	ListAllForRole(ctx context.Context, roleNameOrID *string) ([]*RBACEndpointPermission, error)
}

// RBACEndpointPermissionService handles RBACEndpointPermissions in Kong.
type RBACEndpointPermissionService service

// Create creates a RBACEndpointPermission in Kong.
func (s *RBACEndpointPermissionService) Create(ctx context.Context,
	ep *RBACEndpointPermission,
) (*RBACEndpointPermission, error) {
	if ep == nil {
		return nil, fmt.Errorf("cannot create a nil endpointpermission")
	}
	if ep.Role == nil || ep.Role.ID == nil {
		return nil, fmt.Errorf("cannot create endpoint permission with role or role id undefined")
	}

	method := "POST"
	endpoint := fmt.Sprintf("/rbac/roles/%v/endpoints", *ep.Role.ID)
	req, err := s.client.NewRequest(method, endpoint, nil, ep)
	if err != nil {
		return nil, err
	}

	var createdEndpointPermission RBACEndpointPermission

	_, err = s.client.Do(ctx, req, &createdEndpointPermission)
	if err != nil {
		return nil, err
	}
	return &createdEndpointPermission, nil
}

// Get fetches a RBACEndpointPermission in Kong.
func (s *RBACEndpointPermissionService) Get(ctx context.Context,
	roleNameOrID *string, workspaceNameOrID *string, endpointName *string,
) (*RBACEndpointPermission, error) {
	if isEmptyString(endpointName) {
		return nil, fmt.Errorf("endpointName cannot be nil for Get operation")
	}
	if *endpointName == "*" {
		endpointName = String("/" + *endpointName)
	}
	endpoint := fmt.Sprintf("/rbac/roles/%v/endpoints/%v%v", *roleNameOrID, *workspaceNameOrID, *endpointName)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var EndpointPermission RBACEndpointPermission
	_, err = s.client.Do(ctx, req, &EndpointPermission)
	if err != nil {
		return nil, err
	}
	return &EndpointPermission, nil
}

// Update updates a RBACEndpointPermission in Kong.
func (s *RBACEndpointPermissionService) Update(ctx context.Context,
	ep *RBACEndpointPermission,
) (*RBACEndpointPermission, error) {
	if ep == nil {
		return nil, fmt.Errorf("cannot update a nil EndpointPermission")
	}
	if ep.Workspace == nil {
		return nil, fmt.Errorf("cannot update an EndpointPermission with workspace as nil")
	}
	if ep.Role == nil || ep.Role.ID == nil {
		return nil, fmt.Errorf("cannot create endpoint permission with role or role id undefined")
	}

	if isEmptyString(ep.Endpoint) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	endpointName := ep.Endpoint
	if *endpointName == "*" {
		endpointName = String("/" + *endpointName)
	}
	endpoint := fmt.Sprintf("/rbac/roles/%v/endpoints/%v%v",
		*ep.Role.ID, *ep.Workspace, *endpointName)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, ep)
	if err != nil {
		return nil, err
	}

	var updatedEndpointPermission RBACEndpointPermission
	_, err = s.client.Do(ctx, req, &updatedEndpointPermission)
	if err != nil {
		return nil, err
	}
	return &updatedEndpointPermission, nil
}

// Delete deletes a EndpointPermission in Kong
func (s *RBACEndpointPermissionService) Delete(ctx context.Context,
	roleNameOrID *string, workspaceNameOrID *string, endpointName *string,
) error {
	if endpointName == nil {
		return fmt.Errorf("cannot update a nil EndpointPermission")
	}
	if workspaceNameOrID == nil {
		return fmt.Errorf("cannot update an EndpointPermission with workspace as nil")
	}
	if roleNameOrID == nil {
		return fmt.Errorf("cannot update an EndpointPermission with role as nil")
	}

	if *endpointName == "*" {
		endpointName = String("/" + *endpointName)
	}
	reqEndpoint := fmt.Sprintf("/rbac/roles/%v/endpoints/%v/%v",
		*roleNameOrID, *workspaceNameOrID, *endpointName)
	req, err := s.client.NewRequest("DELETE", reqEndpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// ListAllForRole fetches a list of all RBACEndpointPermissions in Kong for a given role.
func (s *RBACEndpointPermissionService) ListAllForRole(ctx context.Context,
	roleNameOrID *string,
) ([]*RBACEndpointPermission, error) {
	var data, page []json.RawMessage
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		page, opt, err = s.client.list(ctx, fmt.Sprintf("/rbac/roles/%v/endpoints", *roleNameOrID), opt)
		if err != nil {
			return nil, err
		}
		data = append(data, page...)
	}

	var eps []*RBACEndpointPermission
	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, err
		}
		var ep RBACEndpointPermission
		err = json.Unmarshal(b, &ep)
		if err != nil {
			return nil, err
		}
		eps = append(eps, &ep)
	}

	return eps, nil
}
