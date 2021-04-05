package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// AbstractRBACRoleService handles Roles in Kong.
type AbstractRBACRoleService interface {
	// Create creates a Role in Kong.
	Create(ctx context.Context, role *RBACRole) (*RBACRole, error)
	// Get fetches a Role in Kong.
	Get(ctx context.Context, nameOrID *string) (*RBACRole, error)
	// Update updates a Role in Kong.
	Update(ctx context.Context, role *RBACRole) (*RBACRole, error)
	// Delete deletes a Role in Kong
	Delete(ctx context.Context, RoleOrID *string) error
	// List fetches a list of Roles in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*RBACRole, *ListOpt, error)
	// List fetches a list of all Roles in Kong.
	ListAll(ctx context.Context) ([]*RBACRole, error)
}

// RBACRoleService handles Roles in Kong.
type RBACRoleService service

// Create creates a Role in Kong.
func (s *RBACRoleService) Create(ctx context.Context,
	role *RBACRole) (*RBACRole, error) {

	if role == nil {
		return nil, errors.New("cannot create a nil role")
	}

	endpoint := "/rbac/roles"
	method := "POST"
	if role.ID != nil {
		endpoint = endpoint + "/" + *role.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, endpoint, nil, role)

	if err != nil {
		return nil, err
	}

	var createdRole RBACRole
	_, err = s.client.Do(ctx, req, &createdRole)
	if err != nil {
		return nil, err
	}
	return &createdRole, nil
}

// Get fetches a Role in Kong.
func (s *RBACRoleService) Get(ctx context.Context,
	nameOrID *string) (*RBACRole, error) {

	if isEmptyString(nameOrID) {
		return nil, errors.New("nameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/rbac/roles/%v", *nameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var Role RBACRole
	_, err = s.client.Do(ctx, req, &Role)
	if err != nil {
		return nil, err
	}
	return &Role, nil
}

// Update updates a Role in Kong.
func (s *RBACRoleService) Update(ctx context.Context,
	role *RBACRole) (*RBACRole, error) {

	if role == nil {
		return nil, errors.New("cannot update a nil Role")
	}

	if isEmptyString(role.ID) {
		return nil, errors.New("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/rbac/roles/%v", *role.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, role)
	if err != nil {
		return nil, err
	}

	var updatedRole RBACRole
	_, err = s.client.Do(ctx, req, &updatedRole)
	if err != nil {
		return nil, err
	}
	return &updatedRole, nil
}

// Delete deletes a Role in Kong
func (s *RBACRoleService) Delete(ctx context.Context,
	RoleOrID *string) error {

	if isEmptyString(RoleOrID) {
		return errors.New("RoleOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/rbac/roles/%v", *RoleOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of all Roles in Kong.
func (s *RBACRoleService) List(ctx context.Context,
	opt *ListOpt) ([]*RBACRole, *ListOpt, error) {

	data, next, err := s.client.list(ctx, "/rbac/roles/", opt)
	if err != nil {
		return nil, nil, err
	}
	var roles []*RBACRole
	for _, object := range data {
		var role RBACRole
		err = json.Unmarshal(object, &role)
		if err != nil {
			return nil, nil, err
		}
		roles = append(roles, &role)
	}

	return roles, next, nil
}

// ListAll fetches all  Roles in Kong.
// This method can take a while if there
// a lot of Roles present.
func (s *RBACRoleService) ListAll(ctx context.Context) ([]*RBACRole, error) {
	var roles, data []*RBACRole
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, data...)
	}
	return roles, nil
}
