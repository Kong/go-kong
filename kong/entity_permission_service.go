package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// EntityPermissionService handles EntityPermissions in Kong.
type EntityPermissionService service

// Create creates an EntityPermission in Kong.
func (s *EntityPermissionService) Create(ctx context.Context,
	ep *EntityPermission) (*EntityPermission, error) {

	if ep == nil {
		return nil, errors.New("cannot create a nil entitypermission")
	}
	if ep.Role.ID == nil {
		return nil, errors.New("cannot create entity permission with role id undefined")
	}

	method := "POST"
	entity := fmt.Sprintf("/rbac/roles/%v/entities", *ep.Role.ID)
	req, err := s.client.NewRequest(method, entity, nil, ep)

	if err != nil {
		return nil, err
	}

	var createdEntityPermission EntityPermission

	_, err = s.client.Do(ctx, req, &createdEntityPermission)
	if err != nil {
		return nil, err
	}
	return &createdEntityPermission, nil
}

// Get fetches an EntityPermission in Kong.
func (s *EntityPermissionService) Get(ctx context.Context,
	roleNameOrID *string, entityName *string) (*EntityPermission, error) {

	if isEmptyString(entityName) {
		return nil, errors.New("entityName cannot be nil for Get operation")
	}

	entity := fmt.Sprintf("/rbac/roles/%v/entities/%v", *roleNameOrID, *entityName)
	req, err := s.client.NewRequest("GET", entity, nil, nil)
	if err != nil {
		return nil, err
	}

	var EntityPermission EntityPermission
	_, err = s.client.Do(ctx, req, &EntityPermission)
	if err != nil {
		return nil, err
	}
	return &EntityPermission, nil
}

// Update updates an EntityPermission in Kong.
func (s *EntityPermissionService) Update(ctx context.Context,
	ep *EntityPermission) (*EntityPermission, error) {

	if ep == nil {
		return nil, errors.New("cannot update a nil EntityPermission")
	}

	if ep.Role.ID == nil {
		return nil, errors.New("cannot update an EntityPermission with role ID as nil")
	}

	if isEmptyString(ep.EntityID) {
		return nil, errors.New("ID cannot be nil for Update operation")
	}

	entity := fmt.Sprintf("/rbac/roles/%v/entities/%v",
		*ep.Role.ID, *ep.EntityID)
	req, err := s.client.NewRequest("PATCH", entity, nil, ep)
	if err != nil {
		return nil, err
	}

	var updatedEntityPermission EntityPermission
	_, err = s.client.Do(ctx, req, &updatedEntityPermission)
	if err != nil {
		return nil, err
	}
	return &updatedEntityPermission, nil
}

// Delete deletes an EntityPermission in Kong
func (s *EntityPermissionService) Delete(ctx context.Context,
	roleNameOrID *string, entityID *string) error {

	if roleNameOrID == nil {
		return errors.New("cannot update an EntityPermission with role as nil")
	}
	if entityID == nil {
		return errors.New("cannot update an EntityPermission with entity ID as nil")
	}

	endpoint := fmt.Sprintf("/rbac/roles/%v/entities/%v",
		*roleNameOrID, *entityID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// ListAllForRole fetches a list of all EntityPermissions in Kong for a given role.
func (s *EntityPermissionService) ListAllForRole(ctx context.Context,
	roleNameOrID *string) ([]*EntityPermission, error) {

	endpoint := fmt.Sprintf("/rbac/roles/%v/entities", *roleNameOrID)
	data, _, err := s.client.list(ctx, endpoint, nil)
	if err != nil {
		return nil, err
	}
	var eps []*EntityPermission
	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, err
		}
		var ep EntityPermission
		err = json.Unmarshal(b, &ep)
		if err != nil {
			return nil, err
		}
		eps = append(eps, &ep)
	}

	return eps, nil
}
