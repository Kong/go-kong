package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// AbstractDeveloperRoleService handles Developer Roles in Kong.
type AbstractDeveloperRoleService interface {
	// Create creates a Developer Role in Kong.
	Create(ctx context.Context, role *DeveloperRole) (*DeveloperRole, error)
	// Get fetches a Developer Role in Kong.
	Get(ctx context.Context, nameOrID *string) (*DeveloperRole, error)
	// Update updates a Developer Role in Kong.
	Update(ctx context.Context, role *DeveloperRole) (*DeveloperRole, error)
	// Delete deletes a Developer Role in Kong
	Delete(ctx context.Context, RoleOrID *string) error
	// List fetches a list of Developer Roles in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*DeveloperRole, *ListOpt, error)
	// List fetches a list of all Developer Roles in Kong.
	ListAll(ctx context.Context) ([]*DeveloperRole, error)
}

// DeveloperRoleService handles Developer Roles in Kong.
type DeveloperRoleService service

// Create creates a Developer Role in Kong.
func (s *DeveloperRoleService) Create(ctx context.Context,
	role *DeveloperRole) (*DeveloperRole, error) {

	if role == nil {
		return nil, errors.New("cannot create a nil role")
	}

	endpoint := "/developers/roles"
	method := "POST"
	req, err := s.client.NewRequest(method, endpoint, nil, role)

	if err != nil {
		return nil, err
	}

	var createdRole DeveloperRole
	_, err = s.client.Do(ctx, req, &createdRole)
	if err != nil {
		return nil, err
	}
	return &createdRole, nil
}

// Get fetches a Developer Role in Kong.
func (s *DeveloperRoleService) Get(ctx context.Context,
	nameOrID *string) (*DeveloperRole, error) {

	if isEmptyString(nameOrID) {
		return nil, errors.New("nameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/developers/roles/%v", *nameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var Role DeveloperRole
	_, err = s.client.Do(ctx, req, &Role)
	if err != nil {
		return nil, err
	}
	return &Role, nil
}

// Update updates a Developer Role in Kong.
func (s *DeveloperRoleService) Update(ctx context.Context,
	role *DeveloperRole) (*DeveloperRole, error) {

	if role == nil {
		return nil, errors.New("cannot update a nil Role")
	}

	if isEmptyString(role.ID) {
		return nil, errors.New("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/developers/roles/%v", *role.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, role)
	if err != nil {
		return nil, err
	}

	var updatedRole DeveloperRole
	_, err = s.client.Do(ctx, req, &updatedRole)
	if err != nil {
		return nil, err
	}
	return &updatedRole, nil
}

// Delete deletes a Developer Role in Kong
func (s *DeveloperRoleService) Delete(ctx context.Context,
	RoleOrID *string) error {

	if isEmptyString(RoleOrID) {
		return errors.New("RoleOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/developers/roles/%v", *RoleOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of all Developer Roles in Kong.
// opt can be used to control pagination.
func (s *DeveloperRoleService) List(ctx context.Context,
	opt *ListOpt) ([]*DeveloperRole, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/developers/roles/", opt)
}

// ListAll fetches all Developer Roles in Kong.
// This method can take a while if there
// a lot of Developer Roles present.
func (s *DeveloperRoleService) ListAll(ctx context.Context) ([]*DeveloperRole, error) {
	return s.listAllByEndpointAndOpt(ctx, "/developers/roles/", newOpt(nil))
}

func (s *DeveloperRoleService) listByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*DeveloperRole, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	roles, err := asDeveloperRole(data)
	if err != nil {
		return nil, nil, err
	}
	return roles, next, nil
}

func (s *DeveloperRoleService) listAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*DeveloperRole, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, false)
	if err != nil {
		return nil, err
	}
	return asDeveloperRole(data)
}

func asDeveloperRole(data [][]byte) ([]*DeveloperRole, error) {
	var roles []*DeveloperRole
	for _, object := range data {
		var role DeveloperRole
		err := json.Unmarshal(object, &role)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}
	return roles, nil
}
