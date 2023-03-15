package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

type AbstractDegraphqlRouteService interface {
	// Creates a DeGraphQL route in kong.
	Create(ctx context.Context, route *DegraphqlRoute) (*DegraphqlRoute, error)
	// Fetches a DeGraphQL route from kong.
	Get(ctx context.Context, serviceNameOrID *string, ID *string) (*DegraphqlRoute, error)
	// Updates a DeGraphQL route in kong.
	Update(ctx context.Context, route *DegraphqlRoute) (*DegraphqlRoute, error)
	// Deletes a DeGraphQL route in kong.
	Delete(ctx context.Context, serviceNameOrID *string, id *string) error
	// Retrieves a page of DeGraphQL routes in kong.
	List(ctx context.Context, serviceNameOrID *string, listopt *ListOpt) ([]*DegraphqlRoute, *ListOpt, error)
	// Retrieves all DeGraphQL routes in kong.
	ListAll(ctx context.Context, serviceNameOrID *string) ([]*DegraphqlRoute, error)
}

type DegraphqlRouteService service

func getServiceNameOrID(service *Service) *string {
	if service == nil {
		return nil
	}
	if service.Name != nil {
		return service.Name
	}
	return service.ID
}

func (s *DegraphqlRouteService) fetchService(ctx context.Context, service *Service) (*Service, error) {
	if service == nil || service.ID == nil || *service.ID == "" {
		return nil, fmt.Errorf("invalid service in fetched DeGraphQL route")
	}

	serviceService := Svcservice{client: s.client}

	return serviceService.Get(ctx, service.ID)
}

// Create creates a DeGraphQL route in Kong.
// Note the Service must be specified.
func (s *DegraphqlRouteService) Create(ctx context.Context, route *DegraphqlRoute) (*DegraphqlRoute, error) {
	if route == nil {
		return nil, fmt.Errorf("cannot create a nil route")
	}
	serviceNameOrID := getServiceNameOrID(route.Service)
	if serviceNameOrID == nil {
		return nil, fmt.Errorf("cannot create a DeGraphQL route for a service without name or ID")
	}

	endpoint := fmt.Sprintf("/services/%s/degraphql/routes", *serviceNameOrID)
	req, err := s.client.NewRequest("POST", endpoint, nil, route)
	if err != nil {
		return nil, err
	}

	var createdRoute DegraphqlRoute
	_, err = s.client.Do(ctx, req, &createdRoute)
	if err != nil {
		return nil, err
	}

	createdRoute.Service, err = s.fetchService(ctx, createdRoute.Service)
	if err != nil {
		return nil, err
	}

	return &createdRoute, nil
}

// Get returns an existing DeGraphQL route from Kong, given a Service and Route ID.
func (s *DegraphqlRouteService) Get(
	ctx context.Context,
	serviceNameOrID *string,
	ID *string,
) (*DegraphqlRoute, error) {
	if ID == nil {
		return nil, fmt.Errorf("ID cannot be nil for a Get operation")
	}
	if serviceNameOrID == nil {
		return nil, fmt.Errorf("serviceID cannot be nil for fetching a DeGraphQL route")
	}

	endpoint := fmt.Sprintf("/services/%s/degraphql/routes/%s", *serviceNameOrID, *ID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var route DegraphqlRoute
	_, err = s.client.Do(ctx, req, &route)
	if err != nil {
		return nil, err
	}

	route.Service, err = s.fetchService(ctx, route.Service)
	if err != nil {
		return nil, err
	}

	return &route, nil
}

// Update modifies an existing Degraphql route in Kong.
// Note that the Service must be specified.
func (s *DegraphqlRouteService) Update(ctx context.Context, route *DegraphqlRoute) (*DegraphqlRoute, error) {
	if route == nil {
		return nil, fmt.Errorf("cannot update a nil route")
	}
	if isEmptyString(route.ID) {
		return nil, fmt.Errorf("ID cannot be nil or empty for Update operation")
	}
	serviceNameOrID := getServiceNameOrID(route.Service)
	if serviceNameOrID == nil {
		return nil, fmt.Errorf("cannot update a DeGraphQL route without a valid service")
	}

	endpoint := fmt.Sprintf("/services/%s/degraphql/routes/%s", *serviceNameOrID, *route.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, route)
	if err != nil {
		return nil, err
	}

	var updatedRoute DegraphqlRoute
	_, err = s.client.Do(ctx, req, &updatedRoute)
	if err != nil {
		return nil, err
	}

	updatedRoute.Service, err = s.fetchService(ctx, updatedRoute.Service)
	if err != nil {
		return nil, err
	}

	return &updatedRoute, nil
}

// Delete removes an existing DeGraphQL route from Kong given a Service and Route ID.
func (s *DegraphqlRouteService) Delete(
	ctx context.Context,
	serviceNameOrID *string,
	ID *string,
) error {
	if isEmptyString(ID) {
		return fmt.Errorf("id cannot be nil for Delete operation")
	}
	if isEmptyString(serviceNameOrID) {
		return fmt.Errorf("serviceID cannot be nil for fetching a DeGraphQL route")
	}

	endpoint := fmt.Sprintf("/services/%s/degraphql/routes/%s", *serviceNameOrID, *ID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List returns a page of DeGraphQL routes from Kong,
// all associated to the specified Service.
func (s *DegraphqlRouteService) List(
	ctx context.Context,
	serviceNameOrID *string,
	opt *ListOpt,
) ([]*DegraphqlRoute, *ListOpt, error) {
	if isEmptyString(serviceNameOrID) {
		return nil, nil, fmt.Errorf("serviceID cannot be nil for fetching a DeGraphQL route")
	}

	endpoint := fmt.Sprintf("/services/%s/degraphql/routes", *serviceNameOrID)
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}

	routes := make([]*DegraphqlRoute, 0, len(data))
	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var route DegraphqlRoute
		err = json.Unmarshal(b, &route)
		if err != nil {
			return nil, nil, err
		}
		routes = append(routes, &route)
	}

	return routes, next, nil
}

// ListAll fetches all DeGraphQL routes associated with the given Service present in Kong.
// This method can take a while to pull all pages of content
// if there are many items present
func (s *DegraphqlRouteService) ListAll(
	ctx context.Context,
	serviceNameOrID *string,
) ([]*DegraphqlRoute, error) {
	var routes, data []*DegraphqlRoute
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, serviceNameOrID, opt)
		if err != nil {
			return nil, err
		}
		routes = append(routes, data...)
	}
	return routes, nil
}
