package kong

import (
	"context"
	"encoding/json"
)

// AbstractMTLSAuthService handles MTLS credentials in Kong.
type AbstractMTLSAuthService interface {
	// Create creates an MTLS credential in Kong
	Create(ctx context.Context, consumerUsernameOrID *string, mtlsAuth *MTLSAuth) (*MTLSAuth, error)
	// Get fetches an MTLS credential from Kong.
	Get(ctx context.Context, consumerUsernameOrID, keyOrID *string) (*MTLSAuth, error)
	// Update updates an MTLS credential in Kong
	Update(ctx context.Context, consumerUsernameOrID *string, mtlsAuth *MTLSAuth) (*MTLSAuth, error)
	// Delete deletes an MTLS credential in Kong
	Delete(ctx context.Context, consumerUsernameOrID, keyOrID *string) error
	// List fetches a list of MTLS credentials in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*MTLSAuth, *ListOpt, error)
	// ListAll fetches all MTLS credentials in Kong.
	ListAll(ctx context.Context) ([]*MTLSAuth, error)
	// ListForConsumer fetches a list of mtls credentials
	// in Kong associated with a specific consumer.
	ListForConsumer(ctx context.Context, consumerUsernameOrID *string, opt *ListOpt) ([]*MTLSAuth, *ListOpt, error)
}

// MTLSAuthService handles MTLS credentials in Kong.
type MTLSAuthService service

// Create creates an MTLS credential in Kong
// If an ID is specified, it will be used to
// create a MTLS in Kong, otherwise an ID
// is auto-generated.
func (s *MTLSAuthService) Create(ctx context.Context,
	consumerUsernameOrID *string, mtlsAuth *MTLSAuth) (*MTLSAuth, error) {

	cred, err := s.client.credentials.Create(ctx, "mtls-auth",
		consumerUsernameOrID, mtlsAuth)
	if err != nil {
		return nil, err
	}

	var createdMTLS MTLSAuth
	err = json.Unmarshal(cred, &createdMTLS)
	if err != nil {
		return nil, err
	}

	return &createdMTLS, nil
}

// Get fetches an MTLS credential from Kong.
func (s *MTLSAuthService) Get(ctx context.Context,
	consumerUsernameOrID, keyOrID *string) (*MTLSAuth, error) {

	cred, err := s.client.credentials.Get(ctx, "mtls-auth",
		consumerUsernameOrID, keyOrID)
	if err != nil {
		return nil, err
	}

	var mtlsAuth MTLSAuth
	err = json.Unmarshal(cred, &mtlsAuth)
	if err != nil {
		return nil, err
	}

	return &mtlsAuth, nil
}

// Update updates an MTLS credential in Kong
func (s *MTLSAuthService) Update(ctx context.Context,
	consumerUsernameOrID *string, mtlsAuth *MTLSAuth) (*MTLSAuth, error) {

	cred, err := s.client.credentials.Update(ctx, "mtls-auth",
		consumerUsernameOrID, mtlsAuth)
	if err != nil {
		return nil, err
	}

	var updatedMTLS MTLSAuth
	err = json.Unmarshal(cred, &updatedMTLS)
	if err != nil {
		return nil, err
	}

	return &updatedMTLS, nil
}

// Delete deletes an MTLS credential in Kong
func (s *MTLSAuthService) Delete(ctx context.Context,
	consumerUsernameOrID, keyOrID *string) error {
	return s.client.credentials.Delete(ctx, "mtls-auth",
		consumerUsernameOrID, keyOrID)
}

// List fetches a list of MTLS credentials in Kong.
// opt can be used to control pagination.
func (s *MTLSAuthService) List(ctx context.Context,
	opt *ListOpt) ([]*MTLSAuth, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/mtls-auths", opt)
}

// ListAll fetches all MTLS credentials in Kong.
// This method can take a while if there
// a lot of MTLS credentials present.
func (s *MTLSAuthService) ListAll(ctx context.Context) ([]*MTLSAuth, error) {
	return s.listAllByEndpointAndOpt(ctx, "/mtls-auths", NewOpt(nil))
}

// ListForConsumer fetches a list of mtls credentials
// in Kong associated with a specific consumer.
// opt can be used to control pagination.
func (s *MTLSAuthService) ListForConsumer(ctx context.Context,
	consumerUsernameOrID *string, opt *ListOpt) ([]*MTLSAuth, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/consumers/"+*consumerUsernameOrID+"/mtls-auth", opt)
}

func (s *MTLSAuthService) listByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*MTLSAuth, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	mtlss, err := asMTLSAuth(data)
	if err != nil {
		return nil, nil, err
	}
	return mtlss, next, nil
}

func (s *MTLSAuthService) listAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*MTLSAuth, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, true)
	if err != nil {
		return nil, err
	}
	return asMTLSAuth(data)
}

func asMTLSAuth(data [][]byte) ([]*MTLSAuth, error) {
	var mtlss []*MTLSAuth
	for _, object := range data {
		var mtlsAuth MTLSAuth
		err := json.Unmarshal(object, &mtlsAuth)
		if err != nil {
			return nil, err
		}
		mtlss = append(mtlss, &mtlsAuth)
	}
	return mtlss, nil
}
