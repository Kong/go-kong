package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// AbstractCACertificateService handles Certificates in Kong.
type AbstractCACertificateService interface {
	// Create creates a CACertificate in Kong.
	Create(ctx context.Context, certificate *CACertificate) (*CACertificate, error)
	// Get fetches a CACertificate in Kong.
	Get(ctx context.Context, ID *string) (*CACertificate, error)
	// Update updates a CACertificate in Kong
	Update(ctx context.Context, certificate *CACertificate) (*CACertificate, error)
	// Delete deletes a CACertificate in Kong
	Delete(ctx context.Context, ID *string) error
	// List fetches a list of certificate in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*CACertificate, *ListOpt, error)
	// ListAll fetches all Certificates in Kong.
	ListAll(ctx context.Context) ([]*CACertificate, error)
}

// CACertificateService handles Certificates in Kong.
type CACertificateService service

// Create creates a CACertificate in Kong.
// If an ID is specified, it will be used to
// create a certificate in Kong, otherwise an ID
// is auto-generated.
func (s *CACertificateService) Create(ctx context.Context,
	certificate *CACertificate) (*CACertificate, error) {

	queryPath := "/ca_certificates"
	method := "POST"
	if certificate.ID != nil {
		queryPath = queryPath + "/" + *certificate.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, queryPath, nil, certificate)

	if err != nil {
		return nil, err
	}

	var createdCACertificate CACertificate
	_, err = s.client.Do(ctx, req, &createdCACertificate)
	if err != nil {
		return nil, err
	}
	return &createdCACertificate, nil
}

// Get fetches a CACertificate in Kong.
func (s *CACertificateService) Get(ctx context.Context,
	ID *string) (*CACertificate, error) {

	if isEmptyString(ID) {
		return nil, errors.New("ID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/ca_certificates/%v", *ID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var certificate CACertificate
	_, err = s.client.Do(ctx, req, &certificate)
	if err != nil {
		return nil, err
	}
	return &certificate, nil
}

// Update updates a CACertificate in Kong
func (s *CACertificateService) Update(ctx context.Context,
	certificate *CACertificate) (*CACertificate, error) {

	if isEmptyString(certificate.ID) {
		return nil, errors.New("ID cannot be nil for Update op           eration")
	}

	endpoint := fmt.Sprintf("/ca_certificates/%v", *certificate.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, certificate)
	if err != nil {
		return nil, err
	}

	var updatedAPI CACertificate
	_, err = s.client.Do(ctx, req, &updatedAPI)
	if err != nil {
		return nil, err
	}
	return &updatedAPI, nil
}

// Delete deletes a CACertificate in Kong
func (s *CACertificateService) Delete(ctx context.Context,
	ID *string) error {

	if isEmptyString(ID) {
		return errors.New("ID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/ca_certificates/%v", *ID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of certificate in Kong.
// opt can be used to control pagination.
func (s *CACertificateService) List(ctx context.Context,
	opt *ListOpt) ([]*CACertificate, *ListOpt, error) {
	return s.ListByEndpointAndOpt(ctx, "/ca_certificates", opt)
}

// ListAll fetches all Certificates in Kong.
// This method can take a while if there
// a lot of Certificates present.
func (s *CACertificateService) ListAll(ctx context.Context) ([]*CACertificate,
	error) {
	return s.ListAllByTags(ctx, nil)
}

// ListAll fetches all Certificates in Kong, filtered by tags
// This method can take a while if there
// a lot of Certificates present.
func (s *CACertificateService) ListAllByTags(ctx context.Context, tags []string) ([]*CACertificate,
	error) {
	return s.ListAllByEndpointAndOpt(ctx, "/ca_certificates", newOpt(tags))
}

func (s *CACertificateService) ListByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*CACertificate, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	var certificates []*CACertificate
	for _, object := range data {
		var certificate CACertificate
		err = json.Unmarshal(object, &certificate)
		if err != nil {
			return nil, nil, err
		}
		certificates = append(certificates, &certificate)
	}

	return certificates, next, nil
}

func (s *CACertificateService) ListAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*CACertificate, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, true)
	if err != nil {
		return nil, err
	}
	var certificates []*CACertificate
	for _, object := range data {
		var certificate CACertificate
		err = json.Unmarshal(object, &certificate)
		if err != nil {
			return nil, err
		}
		certificates = append(certificates, &certificate)
	}

	return certificates, nil
}
