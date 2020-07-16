package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// UserService handles Users in Kong.
type UserService service

// Create creates a User in Kong.
func (s *UserService) Create(ctx context.Context,
	user *User) (*User, error) {

	if user == nil {
		return nil, errors.New("cannot create a nil user")
	}

	endpoint := "/rbac/users"
	method := "POST"
	if user.ID != nil {
		endpoint = endpoint + "/" + *user.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, endpoint, nil, user)

	if err != nil {
		return nil, err
	}

	var createdUser User
	_, err = s.client.Do(ctx, req, &createdUser)
	if err != nil {
		return nil, err
	}
	return &createdUser, nil
}

// Get fetches a User in Kong.
func (s *UserService) Get(ctx context.Context,
	nameOrID *string) (*User, error) {

	if isEmptyString(nameOrID) {
		return nil, errors.New("nameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/rbac/users/%v", *nameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var User User
	_, err = s.client.Do(ctx, req, &User)
	if err != nil {
		return nil, err
	}
	return &User, nil
}

// Update updates a User in Kong.
func (s *UserService) Update(ctx context.Context,
	user *User) (*User, error) {

	if user == nil {
		return nil, errors.New("cannot update a nil User")
	}

	if isEmptyString(user.ID) && isEmptyString(user.Name) {
		return nil, errors.New("ID and Name cannot both be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/rbac/users/%v", *user.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, user)
	if err != nil {
		return nil, err
	}

	var updatedUser User
	_, err = s.client.Do(ctx, req, &updatedUser)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

// Delete deletes a User in Kong
func (s *UserService) Delete(ctx context.Context,
	userOrID *string) error {

	if isEmptyString(userOrID) {
		return errors.New("UserOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/rbac/users/%v", *userOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of Users in Kong.
// opt can be used to control pagination.
func (s *UserService) List(ctx context.Context,
	opt *ListOpt) ([]*User, *ListOpt, error) {

	data, next, err := s.client.list(ctx, "/rbac/users/", opt)
	if err != nil {
		return nil, nil, err
	}
	var users []*User
	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var user User
		err = json.Unmarshal(b, &user)
		if err != nil {
			return nil, nil, err
		}
		users = append(users, &user)
	}

	return users, next, nil
}

// ListAll fetches all users in Kong.
func (s *UserService) ListAll(ctx context.Context) ([]*User, error) {

	var users, data []*User
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		users = append(users, data...)
	}

	return users, nil
}

// TODO: After implementing the roles service add:
// * AddRoles
// * DeleteRoles
// * ListRoles

// TODO: After implementing the permissions service add:
// * ListPermissions
