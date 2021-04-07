package kong

import (
	"context"
	"encoding/json"
)

// AbstractJWTAuthService handles JWT credentials in Kong.
type AbstractJWTAuthService interface {
	// Create creates a JWT credential in Kong
	Create(ctx context.Context, consumerUsernameOrID *string, jwtAuth *JWTAuth) (*JWTAuth, error)
	// Get fetches a JWT credential from Kong.
	Get(ctx context.Context, consumerUsernameOrID, keyOrID *string) (*JWTAuth, error)
	// Update updates a JWT credential in Kong
	Update(ctx context.Context, consumerUsernameOrID *string, jwtAuth *JWTAuth) (*JWTAuth, error)
	// Delete deletes a JWT credential in Kong
	Delete(ctx context.Context, consumerUsernameOrID, keyOrID *string) error
	// List fetches a list of JWT credentials in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*JWTAuth, *ListOpt, error)
	// ListAll fetches all JWT credentials in Kong.
	ListAll(ctx context.Context) ([]*JWTAuth, error)
	// ListAllByTags fetches all JWT credentials filtered by opt in Kong.
	ListAllByOpt(ctx context.Context, opt *ListOpt) ([]*JWTAuth, error)
	// ListAllByTags fetches all JWT credentials filtered by tags in Kong.
	ListAllByTags(ctx context.Context, tags []string) ([]*JWTAuth, error)
	// ListForConsumer fetches a list of jwt credentials
	// in Kong associated with a specific consumer.
	ListForConsumer(ctx context.Context, consumerUsernameOrID *string, opt *ListOpt) ([]*JWTAuth, *ListOpt, error)
}

// JWTAuthService handles JWT credentials in Kong.
type JWTAuthService service

// Create creates a JWT credential in Kong
// If an ID is specified, it will be used to
// create a JWT in Kong, otherwise an ID
// is auto-generated.
func (s *JWTAuthService) Create(ctx context.Context,
	consumerUsernameOrID *string, jwtAuth *JWTAuth) (*JWTAuth, error) {

	cred, err := s.client.credentials.Create(ctx, "jwt-auth",
		consumerUsernameOrID, jwtAuth)
	if err != nil {
		return nil, err
	}

	var createdJWT JWTAuth
	err = json.Unmarshal(cred, &createdJWT)
	if err != nil {
		return nil, err
	}

	return &createdJWT, nil
}

// Get fetches a JWT credential from Kong.
func (s *JWTAuthService) Get(ctx context.Context,
	consumerUsernameOrID, keyOrID *string) (*JWTAuth, error) {

	cred, err := s.client.credentials.Get(ctx, "jwt-auth",
		consumerUsernameOrID, keyOrID)
	if err != nil {
		return nil, err
	}

	var jwtAuth JWTAuth
	err = json.Unmarshal(cred, &jwtAuth)
	if err != nil {
		return nil, err
	}

	return &jwtAuth, nil
}

// Update updates a JWT credential in Kong
func (s *JWTAuthService) Update(ctx context.Context,
	consumerUsernameOrID *string, jwtAuth *JWTAuth) (*JWTAuth, error) {

	cred, err := s.client.credentials.Update(ctx, "jwt-auth",
		consumerUsernameOrID, jwtAuth)
	if err != nil {
		return nil, err
	}

	var updatedJWT JWTAuth
	err = json.Unmarshal(cred, &updatedJWT)
	if err != nil {
		return nil, err
	}

	return &updatedJWT, nil
}

// Delete deletes a JWT credential in Kong
func (s *JWTAuthService) Delete(ctx context.Context,
	consumerUsernameOrID, keyOrID *string) error {
	return s.client.credentials.Delete(ctx, "jwt-auth",
		consumerUsernameOrID, keyOrID)
}

// List fetches a list of JWT credentials in Kong.
// opt can be used to control pagination.
func (s *JWTAuthService) List(ctx context.Context,
	opt *ListOpt) ([]*JWTAuth, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/jwts", opt)
}

// ListAll fetches all JWT credentials in Kong.
// This method can take a while if there
// a lot of JWT credentials present.
func (s *JWTAuthService) ListAll(ctx context.Context) ([]*JWTAuth, error) {
	return s.ListAllByTags(ctx, nil)
}

// ListAllByTags fetches all JWT credentials filtered by tags in Kong.
// This method can take a while if there
// a lot of JWT credentials present.
func (s *JWTAuthService) ListAllByTags(ctx context.Context, tags []string) ([]*JWTAuth, error) {
	return s.ListAllByOpt(ctx, newOpt(tags))
}

// ListAllByOpt fetches all JWT credentials filtered by opt in Kong.
// This method can take a while if there
// a lot of JWT credentials present.
func (s *JWTAuthService) ListAllByOpt(ctx context.Context, opt *ListOpt) ([]*JWTAuth, error) {
	return s.listAllByEndpointAndOpt(ctx, "/jwts", opt)
}

// ListForConsumer fetches a list of jwt credentials
// in Kong associated with a specific consumer.
// opt can be used to control pagination.
func (s *JWTAuthService) ListForConsumer(ctx context.Context,
	consumerUsernameOrID *string, opt *ListOpt) ([]*JWTAuth, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/consumers/"+*consumerUsernameOrID+"/jwt", opt)
}

func (s *JWTAuthService) listByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*JWTAuth, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	jwts, err := asJWTAuth(data)
	if err != nil {
		return nil, nil, err
	}
	return jwts, next, nil
}

func (s *JWTAuthService) listAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*JWTAuth, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, true)
	if err != nil {
		return nil, err
	}
	return asJWTAuth(data)
}

func asJWTAuth(data [][]byte) ([]*JWTAuth, error) {
	var jwts []*JWTAuth
	for _, object := range data {
		var jwtAuth JWTAuth
		err := json.Unmarshal(object, &jwtAuth)
		if err != nil {
			return nil, err
		}
		jwts = append(jwts, &jwtAuth)
	}
	return jwts, nil
}
