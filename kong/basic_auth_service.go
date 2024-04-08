package kong

import (
	"context"
	"encoding/json"
)

// AbstractBasicAuthService handles basic-auth credentials in Kong.
type AbstractBasicAuthService interface {
	// Create creates a basic-auth credential in Kong
	// is auto-generated.
	Create(ctx context.Context, consumerUsernameOrID *string, basicAuth *BasicAuthSkipHash) (*BasicAuth, error)
	// Get fetches a basic-auth credential from Kong.
	Get(ctx context.Context, consumerUsernameOrID, usernameOrID *string) (*BasicAuth, error)
	// Update updates a basic-auth credential in Kong
	Update(ctx context.Context, consumerUsernameOrID *string, basicAuth *BasicAuthSkipHash) (*BasicAuth, error)
	// Delete deletes a basic-auth credential in Kong
	Delete(ctx context.Context, consumerUsernameOrID, usernameOrID *string) error
	// List fetches a list of basic-auth credentials in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*BasicAuth, *ListOpt, error)
	// ListAll fetches all basic-auth credentials in Kong.
	ListAll(ctx context.Context) ([]*BasicAuth, error)
	// ListForConsumer fetches a list of basic-auth credentials
	// in Kong associated with a specific consumer.
	ListForConsumer(ctx context.Context, consumerUsernameOrID *string, opt *ListOpt) ([]*BasicAuth, *ListOpt, error)
}

// BasicAuthService handles basic-auth credentials in Kong.
type BasicAuthService service

type BasicAuthSkipHash struct {
	Consumer  *Consumer `json:"consumer,omitempty" yaml:"consumer,omitempty"`
	CreatedAt *int      `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	Username  *string   `json:"username,omitempty" yaml:"username,omitempty"`
	Password  *string   `json:"password,omitempty" yaml:"password,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
	SkipHash  *bool     `json:"_skip_hash,omitempty" yaml:"_skip_hash,omitempty"`
}

// Create creates a basic-auth credential in Kong
// If an ID is specified, it will be used to
// create a basic-auth in Kong, otherwise an ID
// is auto-generated.
func (s *BasicAuthService) Create(ctx context.Context,
	consumerUsernameOrID *string, basicAuth *BasicAuthSkipHash,
) (*BasicAuth, error) {
	var skipHash bool
	if basicAuth.SkipHash != nil {
		skipHash = *basicAuth.SkipHash
	}
	newBasicAuth := &BasicAuth{
		Consumer: basicAuth.Consumer,
		ID:       basicAuth.ID,
		Password: basicAuth.Password,
		Username: basicAuth.Username,
		Tags:     basicAuth.Tags,
	}
	cred, err := s.client.credentials.Create(ctx, "basic-auth",
		consumerUsernameOrID, newBasicAuth, skipHash)
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
	consumerUsernameOrID, usernameOrID *string,
) (*BasicAuth, error) {
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
	consumerUsernameOrID *string, basicAuth *BasicAuthSkipHash,
) (*BasicAuth, error) {
	var skipHash bool
	if basicAuth.SkipHash != nil {
		skipHash = *basicAuth.SkipHash
		basicAuth.SkipHash = nil
	}

	newBasicAuth := &BasicAuth{
		Consumer: basicAuth.Consumer,
		ID:       basicAuth.ID,
		Password: basicAuth.Password,
		Username: basicAuth.Username,
		Tags:     basicAuth.Tags,
	}
	cred, err := s.client.credentials.Update(ctx, "basic-auth",
		consumerUsernameOrID, newBasicAuth, skipHash)
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
	consumerUsernameOrID, usernameOrID *string,
) error {
	return s.client.credentials.Delete(ctx, "basic-auth",
		consumerUsernameOrID, usernameOrID)
}

// List fetches a list of basic-auth credentials in Kong.
// opt can be used to control pagination.
func (s *BasicAuthService) List(ctx context.Context,
	opt *ListOpt,
) ([]*BasicAuth, *ListOpt, error) {
	data, next, err := s.client.list(ctx, "/basic-auths", opt)
	if err != nil {
		return nil, nil, err
	}
	var basicAuths []*BasicAuth
	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var basicAuth BasicAuth
		err = json.Unmarshal(b, &basicAuth)
		if err != nil {
			return nil, nil, err
		}
		basicAuths = append(basicAuths, &basicAuth)
	}

	return basicAuths, next, nil
}

// ListAll fetches all basic-auth credentials in Kong.
// This method can take a while if there
// a lot of basic-auth credentials present.
func (s *BasicAuthService) ListAll(ctx context.Context) ([]*BasicAuth, error) {
	var basicAuths, data []*BasicAuth
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		basicAuths = append(basicAuths, data...)
	}
	return basicAuths, nil
}

// ListForConsumer fetches a list of basic-auth credentials
// in Kong associated with a specific consumer.
// opt can be used to control pagination.
func (s *BasicAuthService) ListForConsumer(ctx context.Context,
	consumerUsernameOrID *string, opt *ListOpt,
) ([]*BasicAuth, *ListOpt, error) {
	data, next, err := s.client.list(ctx,
		"/consumers/"+*consumerUsernameOrID+"/basic-auth", opt)
	if err != nil {
		return nil, nil, err
	}
	var basicAuths []*BasicAuth
	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var basicAuth BasicAuth
		err = json.Unmarshal(b, &basicAuth)
		if err != nil {
			return nil, nil, err
		}
		basicAuths = append(basicAuths, &basicAuth)
	}

	return basicAuths, next, nil
}
