package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// AbstractDeveloperService handles Developers in Kong.
type AbstractDeveloperService interface {
	// Create creates a Developer in Kong.
	Create(ctx context.Context, developer *Developer) (*Developer, error)
	// Get fetches a Developer in Kong.
	Get(ctx context.Context, emailOrID *string) (*Developer, error)
	// GetByCustomID fetches a Developer in Kong.
	GetByCustomID(ctx context.Context, customID *string) (*Developer, error)
	// Update updates a Developer in Kong
	Update(ctx context.Context, developer *Developer) (*Developer, error)
	// Delete deletes a Developer in Kong
	Delete(ctx context.Context, emailOrID *string) error
	// List fetches a list of Developers in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*Developer, *ListOpt, error)
	// ListAll fetches all Developers in Kong.
	ListAll(ctx context.Context) ([]*Developer, error)
}

// DeveloperService handles Developers in Kong.
type DeveloperService service

// Create creates a Developer in Kong.
// If an ID is specified, it will be used to
// create a developer in Kong, otherwise an ID
// is auto-generated.
// This call does _not_ use a PUT when provided an ID.
// Although /developers accepts PUTs, PUTs do not accept passwords and do not create
// the hidden consumer that backs the developer. Subsequent attempts to use such developers
// result in fatal errors.
func (s *DeveloperService) Create(ctx context.Context,
	developer *Developer) (*Developer, error) {

	queryPath := "/developers"
	method := "POST"
	req, err := s.client.NewRequest(method, queryPath, nil, developer)

	if err != nil {
		return nil, err
	}

	var createdDeveloper Developer
	_, err = s.client.Do(ctx, req, &createdDeveloper)
	if err != nil {
		return nil, err
	}
	return &createdDeveloper, nil
}

// Get fetches a Developer in Kong.
func (s *DeveloperService) Get(ctx context.Context,
	emailOrID *string) (*Developer, error) {

	if isEmptyString(emailOrID) {
		return nil, errors.New("emailOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/developers/%v", *emailOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var developer Developer
	_, err = s.client.Do(ctx, req, &developer)
	if err != nil {
		return nil, err
	}
	return &developer, nil
}

// GetByCustomID fetches a Developer in Kong.
func (s *DeveloperService) GetByCustomID(ctx context.Context,
	customID *string) (*Developer, error) {

	if isEmptyString(customID) {
		return nil, errors.New("customID cannot be nil for Get operation")
	}

	type QS struct {
		CustomID string `url:"custom_id,omitempty"`
	}

	req, err := s.client.NewRequest("GET", "/developers",
		&QS{CustomID: *customID}, nil)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data []Developer
	}
	var resp Response
	_, err = s.client.Do(ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, NewAPIError(http.StatusNotFound, "Not found")
	}

	return &resp.Data[0], nil
}

// Update updates a Developer in Kong
func (s *DeveloperService) Update(ctx context.Context,
	developer *Developer) (*Developer, error) {

	if isEmptyString(developer.ID) {
		return nil, errors.New("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/developers/%v", *developer.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, developer)
	if err != nil {
		return nil, err
	}
	type Response struct {
		Developer Developer
	}
	var resp Response
	_, err = s.client.Do(ctx, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp.Developer, nil
}

// Delete deletes a Developer in Kong
func (s *DeveloperService) Delete(ctx context.Context,
	emailOrID *string) error {

	if isEmptyString(emailOrID) {
		return errors.New("emailOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/developers/%v", *emailOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of Developers in Kong.
// opt can be used to control pagination.
func (s *DeveloperService) List(ctx context.Context,
	opt *ListOpt) ([]*Developer, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/developers", opt)
}

// ListAll fetches all Developers in Kong.
// This method can take a while if there
// a lot of Developers present.
func (s *DeveloperService) ListAll(ctx context.Context) ([]*Developer, error) {
	return s.listAllByEndpointAndOpt(ctx, "/developers", newOpt(nil))
}

func (s *DeveloperService) listByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*Developer, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	developers, err := asDeveloper(data)
	if err != nil {
		return nil, nil, err
	}
	return developers, next, nil
}

func (s *DeveloperService) listAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*Developer, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, false)
	if err != nil {
		return nil, err
	}
	return asDeveloper(data)
}

func asDeveloper(data [][]byte) ([]*Developer, error) {
	var developers []*Developer
	for _, object := range data {
		var developer Developer
		err := json.Unmarshal(object, &developer)
		if err != nil {
			return nil, err
		}
		developers = append(developers, &developer)
	}
	return developers, nil
}
