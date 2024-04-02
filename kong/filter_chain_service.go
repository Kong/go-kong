package kong

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AbstractFilterChainService handles FilterChains in Kong.
type AbstractFilterChainService interface {
	// Create creates a FilterChain in Kong.
	Create(ctx context.Context, filterChain *FilterChain) (*FilterChain, error)
	// CreateForService creates a FilterChain in Kong.
	CreateForService(ctx context.Context, serviceIDorName *string, filterChain *FilterChain) (*FilterChain, error)
	// CreateForRoute creates a FilterChain in Kong.
	CreateForRoute(ctx context.Context, routeIDorName *string, filterChain *FilterChain) (*FilterChain, error)
	// Get fetches a FilterChain in Kong.
	Get(ctx context.Context, nameOrID *string) (*FilterChain, error)
	// Update updates a FilterChain in Kong
	Update(ctx context.Context, filterChain *FilterChain) (*FilterChain, error)
	// UpdateForService updates a FilterChain in Kong for a service
	UpdateForService(ctx context.Context, serviceIDorName *string, filterChain *FilterChain) (*FilterChain, error)
	// UpdateForRoute updates a FilterChain in Kong for a service
	UpdateForRoute(ctx context.Context, routeIDorName *string, filterChain *FilterChain) (*FilterChain, error)
	// Delete deletes a FilterChain in Kong
	Delete(ctx context.Context, nameOrID *string) error
	// DeleteForService deletes a FilterChain in Kong
	DeleteForService(ctx context.Context, serviceIDorName *string, filterChainID *string) error
	// DeleteForRoute deletes a FilterChain in Kong
	DeleteForRoute(ctx context.Context, routeIDorName *string, filterChainID *string) error
	// List fetches a list of FilterChains in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*FilterChain, *ListOpt, error)
	// ListAll fetches all FilterChains in Kong.
	ListAll(ctx context.Context) ([]*FilterChain, error)
	// ListAllForService fetches all FilterChains in Kong enabled for a service.
	ListAllForService(ctx context.Context, serviceIDorName *string) ([]*FilterChain, error)
	// ListAllForRoute fetches all FilterChains in Kong enabled for a service.
	ListAllForRoute(ctx context.Context, routeID *string) ([]*FilterChain, error)
}

// FilterChainService handles FilterChains in Kong.
type FilterChainService service

// Create creates a FilterChain in Kong.
// If an ID is specified, it will be used to
// create a filter chain in Kong, otherwise an ID
// is auto-generated.
func (s *FilterChainService) Create(ctx context.Context,
	filterChain *FilterChain,
) (*FilterChain, error) {
	queryPath := "/filter-chains"
	method := "POST"
	if filterChain.ID != nil {
		queryPath = queryPath + "/" + *filterChain.ID
		method = "PUT"
	}
	return s.sendRequest(ctx, filterChain, queryPath, method)
}

// CreateForService creates a FilterChain in Kong at Service level.
// If an ID is specified, it will be used to
// create a filter chain in Kong, otherwise an ID
// is auto-generated.
func (s *FilterChainService) CreateForService(ctx context.Context,
	serviceIDorName *string, filterChain *FilterChain,
) (*FilterChain, error) {
	queryPath := "/filter-chains"
	method := "POST"
	if filterChain.ID != nil {
		queryPath = queryPath + "/" + *filterChain.ID
		method = "PUT"
	}
	if isEmptyString(serviceIDorName) {
		return nil, fmt.Errorf("serviceIDorName cannot be nil")
	}

	return s.sendRequest(ctx, filterChain, fmt.Sprintf("/services/%v"+queryPath, *serviceIDorName), method)
}

// CreateForRoute creates a FilterChain in Kong at Route level.
// If an ID is specified, it will be used to
// create a filter chain in Kong, otherwise an ID
// is auto-generated.
func (s *FilterChainService) CreateForRoute(ctx context.Context,
	routeIDorName *string, filterChain *FilterChain,
) (*FilterChain, error) {
	queryPath := "/filter-chains"
	method := "POST"

	if filterChain.ID != nil {
		queryPath = queryPath + "/" + *filterChain.ID
		method = "PUT"
	}
	if isEmptyString(routeIDorName) {
		return nil, fmt.Errorf("routeIDorName cannot be nil")
	}

	return s.sendRequest(ctx, filterChain, fmt.Sprintf("/routes/%v"+queryPath, *routeIDorName), method)
}

// Get fetches a FilterChain in Kong.
func (s *FilterChainService) Get(ctx context.Context,
	nameOrID *string,
) (*FilterChain, error) {
	if isEmptyString(nameOrID) {
		return nil, fmt.Errorf("nameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/filter-chains/%v", *nameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var filterChain FilterChain
	_, err = s.client.Do(ctx, req, &filterChain)
	if err != nil {
		return nil, err
	}
	return &filterChain, nil
}

// Update updates a FilterChain in Kong
func (s *FilterChainService) Update(ctx context.Context,
	filterChain *FilterChain,
) (*FilterChain, error) {
	if isEmptyString(filterChain.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/filter-chains/%v", *filterChain.ID)
	return s.sendRequest(ctx, filterChain, endpoint, "PATCH")
}

// UpdateForService updates a FilterChain in Kong at Service level.
func (s *FilterChainService) UpdateForService(ctx context.Context,
	serviceIDorName *string, filterChain *FilterChain,
) (*FilterChain, error) {
	if isEmptyString(filterChain.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}
	if isEmptyString(serviceIDorName) {
		return nil, fmt.Errorf("serviceIDorName cannot be nil")
	}

	endpoint := fmt.Sprintf("/services/%v/filter-chains/%v", *serviceIDorName, *filterChain.ID)
	return s.sendRequest(ctx, filterChain, endpoint, "PATCH")
}

// UpdateForRoute updates a FilterChain in Kong at Route level.
func (s *FilterChainService) UpdateForRoute(ctx context.Context,
	routeIDorName *string, filterChain *FilterChain,
) (*FilterChain, error) {
	if isEmptyString(filterChain.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}
	if isEmptyString(routeIDorName) {
		return nil, fmt.Errorf("routeIDorName cannot be nil")
	}

	endpoint := fmt.Sprintf("/routes/%v/filter-chains/%v", *routeIDorName, *filterChain.ID)
	return s.sendRequest(ctx, filterChain, endpoint, "PATCH")
}

// Delete deletes a FilterChain in Kong
func (s *FilterChainService) Delete(ctx context.Context,
	filterChainID *string,
) error {
	if isEmptyString(filterChainID) {
		return fmt.Errorf("filterChainID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/filter-chains/%v", *filterChainID)
	_, err := s.sendRequest(ctx, nil, endpoint, "DELETE")
	if err != nil {
		return err
	}
	return err
}

// DeleteForService deletes a FilterChain in Kong at Service level.
func (s *FilterChainService) DeleteForService(ctx context.Context,
	serviceIDorName *string, filterChainID *string,
) error {
	if isEmptyString(filterChainID) {
		return fmt.Errorf("filterChain ID cannot be nil for Delete operation")
	}
	if isEmptyString(serviceIDorName) {
		return fmt.Errorf("serviceIDorName cannot be nil")
	}

	endpoint := fmt.Sprintf("/services/%v/filter-chains/%v", *serviceIDorName, *filterChainID)
	_, err := s.sendRequest(ctx, nil, endpoint, "DELETE")
	if err != nil {
		return err
	}
	return err
}

// DeleteForRoute deletes a FilterChain in Kong at Route level.
func (s *FilterChainService) DeleteForRoute(ctx context.Context,
	routeIDorName *string, filterChainID *string,
) error {
	if isEmptyString(filterChainID) {
		return fmt.Errorf("filterChain ID cannot be nil for Delete operation")
	}
	if isEmptyString(routeIDorName) {
		return fmt.Errorf("routeIDorName cannot be nil")
	}

	endpoint := fmt.Sprintf("/routes/%v/filter-chains/%v", *routeIDorName, *filterChainID)
	_, err := s.sendRequest(ctx, nil, endpoint, "DELETE")
	if err != nil {
		return err
	}
	return nil
}

// listByPath fetches a list of FilterChains in Kong
// on a specific path.
// This is a helper method for listing all filter chains
// or filter chains for specific entities.
func (s *FilterChainService) listByPath(ctx context.Context,
	path string, opt *ListOpt,
) ([]*FilterChain, *ListOpt, error) {
	data, next, err := s.client.list(ctx, path, opt)
	if err != nil {
		return nil, nil, err
	}
	var filterChains []*FilterChain

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var filterChain FilterChain
		err = json.Unmarshal(b, &filterChain)
		if err != nil {
			return nil, nil, err
		}
		filterChains = append(filterChains, &filterChain)
	}

	return filterChains, next, nil
}

// ListAll fetches all FilterChains in Kong.
// This method can take a while if there
// a lot of FilterChains present.
func (s *FilterChainService) listAllByPath(ctx context.Context,
	path string,
) ([]*FilterChain, error) {
	var filterChains, data []*FilterChain
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.listByPath(ctx, path, opt)
		if err != nil {
			return nil, err
		}
		filterChains = append(filterChains, data...)
	}
	return filterChains, nil
}

// List fetches a list of FilterChains in Kong.
// opt can be used to control pagination.
func (s *FilterChainService) List(ctx context.Context,
	opt *ListOpt,
) ([]*FilterChain, *ListOpt, error) {
	return s.listByPath(ctx, "/filter-chains", opt)
}

// ListAll fetches all FilterChains in Kong.
// This method can take a while if there
// a lot of FilterChains present.
func (s *FilterChainService) ListAll(ctx context.Context) ([]*FilterChain, error) {
	return s.listAllByPath(ctx, "/filter-chains")
}

// ListAllForService fetches all FilterChains in Kong enabled for a service.
func (s *FilterChainService) ListAllForService(ctx context.Context,
	serviceIDorName *string,
) ([]*FilterChain, error) {
	if isEmptyString(serviceIDorName) {
		return nil, fmt.Errorf("serviceIDorName cannot be nil")
	}
	return s.listAllByPath(ctx, "/services/"+*serviceIDorName+"/filter-chains")
}

// ListAllForRoute fetches all FilterChains in Kong enabled for a service.
func (s *FilterChainService) ListAllForRoute(ctx context.Context,
	routeID *string,
) ([]*FilterChain, error) {
	if isEmptyString(routeID) {
		return nil, fmt.Errorf("routeID cannot be nil")
	}
	return s.listAllByPath(ctx, "/routes/"+*routeID+"/filter-chains")
}

func (s *FilterChainService) sendRequest(ctx context.Context,
	filterChain *FilterChain, endpoint, method string,
) (*FilterChain, error) {
	var req *http.Request
	var err error
	if method == "DELETE" {
		req, err = s.client.NewRequest(method, endpoint, nil, nil)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = s.client.NewRequest(method, endpoint, nil, filterChain)
		if err != nil {
			return nil, err
		}
	}
	var createdFilterChain FilterChain
	if method == "DELETE" {
		_, err = s.client.Do(ctx, req, nil)
		if err != nil {
			return nil, err
		}
	} else {
		_, err = s.client.Do(ctx, req, &createdFilterChain)
		if err != nil {
			return nil, err
		}
	}
	return &createdFilterChain, nil
}
