package kong

import (
	"context"
	"encoding/json"
)

// AbstractHMACAuthService handles hmac-auth credentials in Kong.
type AbstractHMACAuthService interface {
	// Create creates a hmac-auth credential in Kong
	Create(ctx context.Context, consumerUsernameOrID *string, hmacAuth *HMACAuth) (*HMACAuth, error)
	// Get fetches a hmac-auth credential from Kong.
	Get(ctx context.Context, consumerUsernameOrID, usernameOrID *string) (*HMACAuth, error)
	// Update updates a hmac-auth credential in Kong
	Update(ctx context.Context, consumerUsernameOrID *string, hmacAuth *HMACAuth) (*HMACAuth, error)
	// Delete deletes a hmac-auth credential in Kong
	Delete(ctx context.Context, consumerUsernameOrID, usernameOrID *string) error
	// List fetches a list of hmac-auth credentials in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*HMACAuth, *ListOpt, error)
	// ListAll fetches all hmac-auth credentials in Kong.
	ListAll(ctx context.Context) ([]*HMACAuth, error)
	// ListAllByOpt fetches all hmac-auth credentials filtered by opt in Kong.
	ListAllByOpt(ctx context.Context, opt *ListOpt) ([]*HMACAuth, error)
	// ListAllByTags fetches all hmac-auth credentials filtered by tags in Kong.
	ListAllByTags(ctx context.Context, tags []string) ([]*HMACAuth, error)
	// ListForConsumer fetches a list of hmac-auth credentials
	// in Kong associated with a specific consumer.
	ListForConsumer(ctx context.Context, consumerUsernameOrID *string, opt *ListOpt) ([]*HMACAuth, *ListOpt, error)
}

// HMACAuthService handles hmac-auth credentials in Kong.
type HMACAuthService service

// Create creates a hmac-auth credential in Kong
// If an ID is specified, it will be used to
// create a hmac-auth in Kong, otherwise an ID
// is auto-generated.
func (s *HMACAuthService) Create(ctx context.Context,
	consumerUsernameOrID *string, hmacAuth *HMACAuth) (*HMACAuth, error) {

	cred, err := s.client.credentials.Create(ctx, "hmac-auth",
		consumerUsernameOrID, hmacAuth)
	if err != nil {
		return nil, err
	}

	var createdHMACAuth HMACAuth
	err = json.Unmarshal(cred, &createdHMACAuth)
	if err != nil {
		return nil, err
	}

	return &createdHMACAuth, nil
}

// Get fetches a hmac-auth credential from Kong.
func (s *HMACAuthService) Get(ctx context.Context,
	consumerUsernameOrID, usernameOrID *string) (*HMACAuth, error) {

	cred, err := s.client.credentials.Get(ctx, "hmac-auth",
		consumerUsernameOrID, usernameOrID)
	if err != nil {
		return nil, err
	}

	var hmacAuth HMACAuth
	err = json.Unmarshal(cred, &hmacAuth)
	if err != nil {
		return nil, err
	}

	return &hmacAuth, nil
}

// Update updates a hmac-auth credential in Kong
func (s *HMACAuthService) Update(ctx context.Context,
	consumerUsernameOrID *string, hmacAuth *HMACAuth) (*HMACAuth, error) {

	cred, err := s.client.credentials.Update(ctx, "hmac-auth",
		consumerUsernameOrID, hmacAuth)
	if err != nil {
		return nil, err
	}

	var updatedHMACAuth HMACAuth
	err = json.Unmarshal(cred, &updatedHMACAuth)
	if err != nil {
		return nil, err
	}

	return &updatedHMACAuth, nil
}

// Delete deletes a hmac-auth credential in Kong
func (s *HMACAuthService) Delete(ctx context.Context,
	consumerUsernameOrID, usernameOrID *string) error {
	return s.client.credentials.Delete(ctx, "hmac-auth",
		consumerUsernameOrID, usernameOrID)
}

// List fetches a list of hmac-auth credentials in Kong.
// opt can be used to control pagination.
func (s *HMACAuthService) List(ctx context.Context,
	opt *ListOpt) ([]*HMACAuth, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/hmac-auths", opt)
}

// ListAll fetches all hmac-auth credentials in Kong.
// This method can take a while if there
// a lot of hmac-auth credentials present.
func (s *HMACAuthService) ListAll(ctx context.Context) ([]*HMACAuth, error) {
	return s.ListAllByTags(ctx, nil)
}

// ListAllByTags fetches all hmac-auth credentials filtered by tags in Kong.
// This method can take a while if there
// a lot of hmac-auth credentials present.
func (s *HMACAuthService) ListAllByTags(ctx context.Context, tags []string) ([]*HMACAuth, error) {
	return s.ListAllByOpt(ctx, newOpt(tags))
}

// ListAllByOpt fetches all hmac-auth credentials in Kong.
// This method can take a while if there
// a lot of hmac-auth credentials present.
func (s *HMACAuthService) ListAllByOpt(ctx context.Context, opt *ListOpt) ([]*HMACAuth, error) {
	return s.listAllByEndpointAndOpt(ctx, "/hmac-auths", opt)
}

// ListForConsumer fetches a list of hmac-auth credentials
// in Kong associated with a specific consumer.
// opt can be used to control pagination.
func (s *HMACAuthService) ListForConsumer(ctx context.Context,
	consumerUsernameOrID *string, opt *ListOpt) ([]*HMACAuth, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/consumers/"+*consumerUsernameOrID+"/hmac-auth", opt)
}

func (s *HMACAuthService) listByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*HMACAuth, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	var hmacAuths []*HMACAuth
	for _, object := range data {
		var hmacAuth HMACAuth
		err = json.Unmarshal(object, &hmacAuth)
		if err != nil {
			return nil, nil, err
		}
		hmacAuths = append(hmacAuths, &hmacAuth)
	}

	return hmacAuths, next, nil
}

func (s *HMACAuthService) listAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*HMACAuth, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, true)
	if err != nil {
		return nil, err
	}
	var hmacAuths []*HMACAuth
	for _, object := range data {
		var hmacAuth HMACAuth
		err = json.Unmarshal(object, &hmacAuth)
		if err != nil {
			return nil, err
		}
		hmacAuths = append(hmacAuths, &hmacAuth)
	}

	return hmacAuths, nil
}
