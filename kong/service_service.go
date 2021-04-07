package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// AbstractSvcService handles services in Kong.
type AbstractSvcService interface {
	// Create creates an Service in Kong
	Create(ctx context.Context, service *Service) (*Service, error)
	// Get fetches an Service in Kong.
	Get(ctx context.Context, nameOrID *string) (*Service, error)
	// GetForRoute fetches a Service associated with routeID in Kong.
	GetForRoute(ctx context.Context, routeID *string) (*Service, error)
	// Update updates an Service in Kong
	Update(ctx context.Context, service *Service) (*Service, error)
	// Delete deletes an Service in Kong
	Delete(ctx context.Context, nameOrID *string) error
	// List fetches a list of Services in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*Service, *ListOpt, error)
	// ListAll fetches all Services in Kong.
	ListAll(ctx context.Context) ([]*Service, error)
	// ListAllByTags fetches all Services filtered by opt in Kong.
	ListAllByOpt(ctx context.Context, opt *ListOpt) ([]*Service, error)
	// ListAllByTags fetches all Services filtered by tags in Kong.
	ListAllByTags(ctx context.Context, tags []string) ([]*Service, error)
}

// Svcservice handles services in Kong.
type Svcservice service

// Create creates an Service in Kong
// If an ID is specified, it will be used to
// create a service in Kong, otherwise an ID
// is auto-generated.
func (s *Svcservice) Create(ctx context.Context,
	service *Service) (*Service, error) {

	if service == nil {
		return nil, errors.New("cannot create a nil service")
	}

	endpoint := "/services"
	method := "POST"
	if service.ID != nil {
		endpoint = endpoint + "/" + *service.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, endpoint, nil, service)
	if err != nil {
		return nil, err
	}

	var createdService Service
	_, err = s.client.Do(ctx, req, &createdService)
	if err != nil {
		return nil, err
	}
	return &createdService, nil
}

// Get fetches an Service in Kong.
func (s *Svcservice) Get(ctx context.Context,
	nameOrID *string) (*Service, error) {

	if isEmptyString(nameOrID) {
		return nil, errors.New("nameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/services/%v", *nameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var Service Service
	_, err = s.client.Do(ctx, req, &Service)
	if err != nil {
		return nil, err
	}
	return &Service, nil
}

// GetForRoute fetches a Service associated with routeID in Kong.
func (s *Svcservice) GetForRoute(ctx context.Context,
	routeID *string) (*Service, error) {

	if isEmptyString(routeID) {
		return nil, errors.New("routeID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/routes/%v/service", *routeID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var Service Service
	_, err = s.client.Do(ctx, req, &Service)
	if err != nil {
		return nil, err
	}
	return &Service, nil
}

// Update updates an Service in Kong
func (s *Svcservice) Update(ctx context.Context,
	service *Service) (*Service, error) {

	if service == nil {
		return nil, errors.New("cannot update a nil service")
	}

	if isEmptyString(service.ID) {
		return nil, errors.New("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/services/%v", *service.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, service)
	if err != nil {
		return nil, err
	}

	var updatedService Service
	_, err = s.client.Do(ctx, req, &updatedService)
	if err != nil {
		return nil, err
	}
	return &updatedService, nil
}

// Delete deletes an Service in Kong
func (s *Svcservice) Delete(ctx context.Context, nameOrID *string) error {

	if isEmptyString(nameOrID) {
		return errors.New("nameOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/services/%v", *nameOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of Services in Kong.
// opt can be used to control pagination.
func (s *Svcservice) List(ctx context.Context,
	opt *ListOpt) ([]*Service, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/services", opt)
}

// ListAll fetches all Services in Kong.
// This method can take a while if there
// a lot of Services present.
func (s *Svcservice) ListAll(ctx context.Context) ([]*Service, error) {
	return s.ListAllByTags(ctx, nil)
}

// ListAllByTags fetches all Services filtered by tags in Kong.
// This method can take a while if there
// a lot of Services present.
func (s *Svcservice) ListAllByTags(ctx context.Context, tags []string) ([]*Service, error) {
	return s.listAllByEndpointAndOpt(ctx, "/services", newOpt(tags))
}

// ListAllByOpt fetches all Services filtered by opt in Kong.
// This method can take a while if there
// a lot of Services present.
func (s *Svcservice) ListAllByOpt(ctx context.Context, opt *ListOpt) ([]*Service, error) {
	return s.listAllByEndpointAndOpt(ctx, "/services", opt)
}

func (s *Svcservice) listByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*Service, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	services, err := asService(data)
	if err != nil {
		return nil, nil, err
	}
	return services, next, nil
}

func (s *Svcservice) listAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*Service, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, false)
	if err != nil {
		return nil, err
	}
	return asService(data)
}

func asService(data [][]byte) ([]*Service, error) {
	var services []*Service
	for _, object := range data {
		var service Service
		err := json.Unmarshal(object, &service)
		if err != nil {
			return nil, err
		}
		services = append(services, &service)
	}
	return services, nil
}
