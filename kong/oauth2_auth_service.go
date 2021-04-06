package kong

import (
	"context"
	"encoding/json"
)

// AbstractOauth2Service handles oauth2 credentials in Kong.
type AbstractOauth2Service interface {
	// Create creates an oauth2 credential in Kong
	Create(ctx context.Context, consumerUsernameOrID *string, oauth2Cred *Oauth2Credential) (*Oauth2Credential, error)
	// Get fetches an oauth2 credential from Kong.
	Get(ctx context.Context, consumerUsernameOrID, clientIDorID *string) (*Oauth2Credential, error)
	// Update updates an oauth2 credential in Kong.
	Update(ctx context.Context, consumerUsernameOrID *string, oauth2Cred *Oauth2Credential) (*Oauth2Credential, error)
	// Delete deletes an oauth2 credential in Kong.
	Delete(ctx context.Context, consumerUsernameOrID, clientIDorID *string) error
	// List fetches a list of oauth2 credentials in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*Oauth2Credential, *ListOpt, error)
	// ListAll fetches all oauth2 credentials in Kong.
	ListAll(ctx context.Context) ([]*Oauth2Credential, error)
	// ListAllByTags fetches all oauth2 credentials filtered by tags in Kong.
	ListAllByTags(ctx context.Context, tags []string) ([]*Oauth2Credential, error)
	// ListForConsumer fetches a list of oauth2 credentials
	// in Kong associated with a specific consumer.
	ListForConsumer(ctx context.Context, consumerUsernameOrID *string, opt *ListOpt) ([]*Oauth2Credential, *ListOpt, error)
}

// Oauth2Service handles oauth2 credentials in Kong.
type Oauth2Service service

// Create creates an oauth2 credential in Kong
// If an ID is specified, it will be used to
// create a oauth2 credential in Kong, otherwise an ID
// is auto-generated.
func (s *Oauth2Service) Create(ctx context.Context,
	consumerUsernameOrID *string,
	oauth2Cred *Oauth2Credential) (*Oauth2Credential, error) {

	cred, err := s.client.credentials.Create(ctx, "oauth2",
		consumerUsernameOrID, oauth2Cred)
	if err != nil {
		return nil, err
	}

	var createdOauth2Cred Oauth2Credential
	err = json.Unmarshal(cred, &createdOauth2Cred)
	if err != nil {
		return nil, err
	}

	return &createdOauth2Cred, nil
}

// Get fetches an oauth2 credential from Kong.
func (s *Oauth2Service) Get(ctx context.Context,
	consumerUsernameOrID, clientIDorID *string) (*Oauth2Credential, error) {

	cred, err := s.client.credentials.Get(ctx, "oauth2",
		consumerUsernameOrID, clientIDorID)
	if err != nil {
		return nil, err
	}

	var oauth2Cred Oauth2Credential
	err = json.Unmarshal(cred, &oauth2Cred)
	if err != nil {
		return nil, err
	}

	return &oauth2Cred, nil
}

// Update updates an oauth2 credential in Kong.
func (s *Oauth2Service) Update(ctx context.Context,
	consumerUsernameOrID *string,
	oauth2Cred *Oauth2Credential) (*Oauth2Credential, error) {

	cred, err := s.client.credentials.Update(ctx, "oauth2",
		consumerUsernameOrID, oauth2Cred)
	if err != nil {
		return nil, err
	}

	var updatedHMACAuth Oauth2Credential
	err = json.Unmarshal(cred, &updatedHMACAuth)
	if err != nil {
		return nil, err
	}

	return &updatedHMACAuth, nil
}

// Delete deletes an oauth2 credential in Kong.
func (s *Oauth2Service) Delete(ctx context.Context,
	consumerUsernameOrID, clientIDorID *string) error {
	return s.client.credentials.Delete(ctx, "oauth2",
		consumerUsernameOrID, clientIDorID)
}

// List fetches a list of oauth2 credentials in Kong.
// opt can be used to control pagination.
func (s *Oauth2Service) List(ctx context.Context,
	opt *ListOpt) ([]*Oauth2Credential, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/oauth2", opt)
}

// ListAll fetches all oauth2 credentials in Kong.
// This method can take a while if there
// a lot of oauth2 credentials present.
func (s *Oauth2Service) ListAll(ctx context.Context) ([]*Oauth2Credential, error) {
	return s.ListAllByTags(ctx, nil)
}

// ListAllByTags fetches all oauth2 credentials filtered by tags in Kong.
// This method can take a while if there
// a lot of oauth2 credentials present.
func (s *Oauth2Service) ListAllByTags(ctx context.Context, tags []string) ([]*Oauth2Credential, error) {
	return s.listAllByEndpointAndOpt(ctx, "/oauth2", newOpt(tags))
}

// ListForConsumer fetches a list of oauth2 credentials
// in Kong associated with a specific consumer.
// opt can be used to control pagination.
func (s *Oauth2Service) ListForConsumer(ctx context.Context,
	consumerUsernameOrID *string, opt *ListOpt) ([]*Oauth2Credential,
	*ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/consumers/"+*consumerUsernameOrID+"/oauth2", opt)
}

func (s *Oauth2Service) listByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*Oauth2Credential, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	var oauth2Creds []*Oauth2Credential
	for _, object := range data {
		var oauth2Cred Oauth2Credential
		err = json.Unmarshal(object, &oauth2Cred)
		if err != nil {
			return nil, nil, err
		}
		oauth2Creds = append(oauth2Creds, &oauth2Cred)
	}

	return oauth2Creds, next, nil
}

func (s *Oauth2Service) listAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*Oauth2Credential, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, true)
	if err != nil {
		return nil, err
	}
	var oauth2Creds []*Oauth2Credential
	for _, object := range data {
		var oauth2Cred Oauth2Credential
		err = json.Unmarshal(object, &oauth2Cred)
		if err != nil {
			return nil, err
		}
		oauth2Creds = append(oauth2Creds, &oauth2Cred)
	}

	return oauth2Creds, nil
}
