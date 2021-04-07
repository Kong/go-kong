package kong

import (
	"context"
	"encoding/json"
)

// AbstractBasicAuthService handles basic-auth credentials in Kong.
type AbstractBasicAuthService interface {
	// Create creates a basic-auth credential in Kong
	// is auto-generated.
	Create(ctx context.Context, consumerUsernameOrID *string, basicAuth *BasicAuth) (*BasicAuth, error)
	// Get fetches a basic-auth credential from Kong.
	Get(ctx context.Context, consumerUsernameOrID, usernameOrID *string) (*BasicAuth, error)
	// Update updates a basic-auth credential in Kong
	Update(ctx context.Context, consumerUsernameOrID *string, basicAuth *BasicAuth) (*BasicAuth, error)
	// Delete deletes a basic-auth credential in Kong
	Delete(ctx context.Context, consumerUsernameOrID, usernameOrID *string) error
	// List fetches a list of basic-auth credentials in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*BasicAuth, *ListOpt, error)
	// ListAll fetches all basic-auth credentials in Kong.
	ListAll(ctx context.Context) ([]*BasicAuth, error)
	// ListAllByOpt fetches all basic-auth credentials filtered by opt in Kong.
	ListAllByOpt(ctx context.Context, opt *ListOpt) ([]*BasicAuth, error)
	// ListAllByTags fetches all basic-auth credentials filtered by tags in Kong.
	ListAllByTags(ctx context.Context, tags []string) ([]*BasicAuth, error)
	// ListForConsumer fetches a list of basic-auth credentials
	// in Kong associated with a specific consumer.
	ListForConsumer(ctx context.Context, consumerUsernameOrID *string, opt *ListOpt) ([]*BasicAuth, *ListOpt, error)
}

// BasicAuthService handles basic-auth credentials in Kong.
type BasicAuthService service

// Create creates a basic-auth credential in Kong
// If an ID is specified, it will be used to
// create a basic-auth in Kong, otherwise an ID
// is auto-generated.
func (s *BasicAuthService) Create(ctx context.Context,
	consumerUsernameOrID *string, basicAuth *BasicAuth) (*BasicAuth, error) {

	cred, err := s.client.credentials.Create(ctx, "basic-auth",
		consumerUsernameOrID, basicAuth)
	if err != nil {
		return nil, err
	}

	var createdBasicAuth BasicAuth
	err = json.Unmarshal(cred, &createdBasicAuth)
	if err != nil {
		return nil, err
	}

	return &createdBasicAuth, nil
}

// Get fetches a basic-auth credential from Kong.
func (s *BasicAuthService) Get(ctx context.Context,
	consumerUsernameOrID, usernameOrID *string) (*BasicAuth, error) {

	cred, err := s.client.credentials.Get(ctx, "basic-auth",
		consumerUsernameOrID, usernameOrID)
	if err != nil {
		return nil, err
	}

	var basicAuth BasicAuth
	err = json.Unmarshal(cred, &basicAuth)
	if err != nil {
		return nil, err
	}

	return &basicAuth, nil
}

// Update updates a basic-auth credential in Kong
func (s *BasicAuthService) Update(ctx context.Context,
	consumerUsernameOrID *string, basicAuth *BasicAuth) (*BasicAuth, error) {

	cred, err := s.client.credentials.Update(ctx, "basic-auth",
		consumerUsernameOrID, basicAuth)
	if err != nil {
		return nil, err
	}

	var updatedBasicAuth BasicAuth
	err = json.Unmarshal(cred, &updatedBasicAuth)
	if err != nil {
		return nil, err
	}

	return &updatedBasicAuth, nil
}

// Delete deletes a basic-auth credential in Kong
func (s *BasicAuthService) Delete(ctx context.Context,
	consumerUsernameOrID, usernameOrID *string) error {
	return s.client.credentials.Delete(ctx, "basic-auth",
		consumerUsernameOrID, usernameOrID)
}

// List fetches a list of basic-auth credentials in Kong.
// opt can be used to control pagination.
func (s *BasicAuthService) List(ctx context.Context,
	opt *ListOpt) ([]*BasicAuth, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/basic-auths", opt)
}

// ListAll fetches all basic-auth credentials in Kong.
// This method can take a while if there
// a lot of basic-auth credentials present.
func (s *BasicAuthService) ListAll(ctx context.Context) ([]*BasicAuth, error) {
	return s.ListAllByTags(ctx, nil)
}

// ListAllByTags fetches all basic-auth credentials filtered by tags in Kong.
// This method can take a while if there
// a lot of basic-auth credentials present.
func (s *BasicAuthService) ListAllByTags(ctx context.Context, tags []string) ([]*BasicAuth, error) {
	return s.ListAllByOpt(ctx, newOpt(tags))
}

// ListAllByTags fetches all basic-auth credentials filtered by opt in Kong.
// This method can take a while if there
// a lot of basic-auth credentials present.
// opt can be used to control pagination and tags.
func (s *BasicAuthService) ListAllByOpt(ctx context.Context, opt *ListOpt) ([]*BasicAuth, error) {
	return s.listAllByEndpointAndOpt(ctx, "/basic-auths", opt)
}

// ListForConsumer fetches a list of basic-auth credentials
// in Kong associated with a specific consumer.
// opt can be used to control pagination.
func (s *BasicAuthService) ListForConsumer(ctx context.Context,
	consumerUsernameOrID *string, opt *ListOpt) ([]*BasicAuth, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/consumers/"+*consumerUsernameOrID+"/basic-auth", opt)
}

func (s *BasicAuthService) listAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*BasicAuth, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, true)
	if err != nil {
		return nil, err
	}
	basicAuths, err := asBasicAuth(data)
	if err != nil {
		return nil, err
	}
	return basicAuths, nil
}

func (s *BasicAuthService) listByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*BasicAuth, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	basicAuths, err := asBasicAuth(data)
	if err != nil {
		return nil, nil, err
	}
	return basicAuths, next, nil
}

func asBasicAuth(data [][]byte) ([]*BasicAuth, error) {
	var basicAuths []*BasicAuth
	for _, object := range data {
		var basicAuth BasicAuth
		err := json.Unmarshal(object, &basicAuth)
		if err != nil {
			return nil, err
		}
		basicAuths = append(basicAuths, &basicAuth)
	}
	return basicAuths, nil
}
