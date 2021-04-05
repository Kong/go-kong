package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// AbstractSNIService handles SNIs in Kong.
type AbstractSNIService interface {
	// Create creates a SNI in Kong.
	Create(ctx context.Context, sni *SNI) (*SNI, error)
	// Get fetches a SNI in Kong.
	Get(ctx context.Context, usernameOrID *string) (*SNI, error)
	// Update updates a SNI in Kong
	Update(ctx context.Context, sni *SNI) (*SNI, error)
	// Delete deletes a SNI in Kong
	Delete(ctx context.Context, usernameOrID *string) error
	// List fetches a list of SNIs in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*SNI, *ListOpt, error)
	// ListForCertificate fetches a list of SNIs
	ListForCertificate(ctx context.Context, certificateID *string, opt *ListOpt) ([]*SNI, *ListOpt, error)
	// ListAll fetches all SNIs in Kong.
	ListAll(ctx context.Context) ([]*SNI, error)
	// ListAllByTags fetches all SNIs filtered by tags in Kong.
	ListAllByTags(ctx context.Context, tags []string) ([]*SNI, error)
}

// SNIService handles SNIs in Kong.
type SNIService service

// Create creates a SNI in Kong.
// If an ID is specified, it will be used to
// create a sni in Kong, otherwise an ID
// is auto-generated.
func (s *SNIService) Create(ctx context.Context, sni *SNI) (*SNI, error) {

	queryPath := "/snis"
	method := "POST"
	if sni.ID != nil {
		queryPath = queryPath + "/" + *sni.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, queryPath, nil, sni)

	if err != nil {
		return nil, err
	}

	var createdSNI SNI
	_, err = s.client.Do(ctx, req, &createdSNI)
	if err != nil {
		return nil, err
	}
	return &createdSNI, nil
}

// Get fetches a SNI in Kong.
func (s *SNIService) Get(ctx context.Context,
	usernameOrID *string) (*SNI, error) {

	if isEmptyString(usernameOrID) {
		return nil, errors.New(
			"usernameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/snis/%v", *usernameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var sni SNI
	_, err = s.client.Do(ctx, req, &sni)
	if err != nil {
		return nil, err
	}
	return &sni, nil
}

// Update updates a SNI in Kong
func (s *SNIService) Update(ctx context.Context, sni *SNI) (*SNI, error) {

	if isEmptyString(sni.ID) {
		return nil, errors.New("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/snis/%v", *sni.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, sni)
	if err != nil {
		return nil, err
	}

	var updatedAPI SNI
	_, err = s.client.Do(ctx, req, &updatedAPI)
	if err != nil {
		return nil, err
	}
	return &updatedAPI, nil
}

// Delete deletes a SNI in Kong
func (s *SNIService) Delete(ctx context.Context, usernameOrID *string) error {

	if isEmptyString(usernameOrID) {
		return errors.New("usernameOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/snis/%v", *usernameOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of SNIs in Kong.
// opt can be used to control pagination.
func (s *SNIService) List(ctx context.Context,
	opt *ListOpt) ([]*SNI, *ListOpt, error) {
	return s.ListByEndpointAndOpt(ctx, "/snis", opt)
}

// ListForCertificate fetches a list of SNIs
// in Kong associated with certificateID.
// opt can be used to control pagination.
func (s *SNIService) ListForCertificate(ctx context.Context,
	certificateID *string, opt *ListOpt) ([]*SNI, *ListOpt, error) {
	return s.ListByEndpointAndOpt(ctx, "/certificates/"+*certificateID+"/snis", opt)
}

// ListAll fetches all SNIs in Kong.
// This method can take a while if there
// a lot of SNIs present.
func (s *SNIService) ListAll(ctx context.Context) ([]*SNI, error) {
	return s.ListAllByTags(ctx, nil)
}

func (s *SNIService) ListAllByTags(ctx context.Context, tags []string) ([]*SNI, error) {
	return s.ListAllByEndpointAndOpt(ctx, "/snis", newOpt(tags))
}

func (s *SNIService) ListByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*SNI, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	var snis []*SNI
	for _, object := range data {
		var sni SNI
		err = json.Unmarshal(object, &sni)
		if err != nil {
			return nil, nil, err
		}
		snis = append(snis, &sni)
	}

	return snis, next, nil
}

func (s *SNIService) ListAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*SNI, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, false)
	if err != nil {
		return nil, err
	}
	var snis []*SNI
	for _, object := range data {
		var sni SNI
		err = json.Unmarshal(object, &sni)
		if err != nil {
			return nil, err
		}
		snis = append(snis, &sni)
	}

	return snis, nil
}
