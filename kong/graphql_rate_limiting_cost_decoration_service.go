package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

type AbstractGraphqlRateLimitingCostDecorationService interface {
	// Creates a cost decoration for the GraphQL rate-limiting plugin in Kong.
	Create(ctx context.Context, costDeco *GraphqlRateLimitingCostDecoration) (*GraphqlRateLimitingCostDecoration, error)
	// Creates a cost decoration with a specified ID for the GraphQL rate-limiting plugin in Kong.
	CreateWithID(ctx context.Context, costDeco *GraphqlRateLimitingCostDecoration) (
		*GraphqlRateLimitingCostDecoration, error)
	// Fetches a cost decoration for the GraphQL rate-limiting plugin from Kong.
	Get(ctx context.Context, ID *string) (*GraphqlRateLimitingCostDecoration, error)
	// Updates a cost decoration for the GraphQL rate-limiting plugin in Kong.
	Update(ctx context.Context, costDeco *GraphqlRateLimitingCostDecoration) (*GraphqlRateLimitingCostDecoration, error)
	// Deletes a cost decoration for the GraphQL rate-limiting plugin in Kong.
	Delete(ctx context.Context, ID *string) error
	// Retrieves a page of cost decorations for the GraphQL rate-limiting plugin in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*GraphqlRateLimitingCostDecoration, *ListOpt, error)
	// Retrieves all decorations for the GraphQL rate-limiting plugin in Kong.
	ListAll(ctx context.Context) ([]*GraphqlRateLimitingCostDecoration, error)
	// Creates a cost decoration scoped to a Service for the GraphQL rate-limiting plugin in Kong.
	// Uses PUT if an ID is provided, POST otherwise.
	CreateForService(ctx context.Context,
		costDeco *GraphqlRateLimitingCostDecoration) (*GraphqlRateLimitingCostDecoration, error)
	// Updates a cost decoration scoped to a Service for the GraphQL rate-limiting plugin in Kong.
	UpdateForService(ctx context.Context,
		costDeco *GraphqlRateLimitingCostDecoration) (*GraphqlRateLimitingCostDecoration, error)
	// Deletes a cost decoration scoped to a Service for the GraphQL rate-limiting plugin in Kong.
	DeleteForService(ctx context.Context,
		costDeco *GraphqlRateLimitingCostDecoration) error
	// Retrieves a page of cost decorations scoped to a Service for the GraphQL rate-limiting plugin in Kong.
	ListForService(ctx context.Context, serviceNameOrID *string,
		opt *ListOpt) ([]*GraphqlRateLimitingCostDecoration, *ListOpt, error)
	// Retrieves all decorations scoped to a Service for the GraphQL rate-limiting plugin in Kong.
	ListAllForService(ctx context.Context, serviceNameOrID *string) ([]*GraphqlRateLimitingCostDecoration, error)
}

type GraphqlRateLimitingCostDecorationService service

// Create creates a CostDecoration item in Kong for the GraphQL rate limiting advanced plugin.
func (s *GraphqlRateLimitingCostDecorationService) Create(
	ctx context.Context,
	costDeco *GraphqlRateLimitingCostDecoration,
) (*GraphqlRateLimitingCostDecoration, error) {
	const queryPath = "/graphql-rate-limiting-advanced/costs"

	if costDeco.ID != nil {
		return nil, fmt.Errorf("can't specify an ID for creating new Cost Decoration")
	}
	req, err := s.client.NewRequest("POST", queryPath, nil, costDeco)
	if err != nil {
		return nil, err
	}

	var createdCostDeco GraphqlRateLimitingCostDecoration
	err = ErrorOrResponseError(s.client.Do(ctx, req, &createdCostDeco))
	if err != nil {
		return nil, err
	}

	return &createdCostDeco, nil
}

// CreateWithID creates a CostDecoration item in Kong for the GraphQL rate limiting advanced plugin
// with a specified ID.
func (s *GraphqlRateLimitingCostDecorationService) CreateWithID(
	ctx context.Context,
	costDeco *GraphqlRateLimitingCostDecoration,
) (*GraphqlRateLimitingCostDecoration, error) {
	if isEmptyString(costDeco.ID) {
		return nil, fmt.Errorf("ID cannot be nil for CreateWithID operation")
	}

	const queryPath = "/graphql-rate-limiting-advanced/costs"

	req, err := s.client.NewRequest("POST", queryPath, nil, costDeco)
	if err != nil {
		return nil, err
	}

	var createdCostDeco GraphqlRateLimitingCostDecoration
	err = ErrorOrResponseError(s.client.Do(ctx, req, &createdCostDeco))
	if err != nil {
		return nil, err
	}

	return &createdCostDeco, nil
}

// Get fetches a CostDecoration item from Kong given an ID.
func (s *GraphqlRateLimitingCostDecorationService) Get(
	ctx context.Context,
	id *string,
) (*GraphqlRateLimitingCostDecoration, error) {
	if isEmptyString(id) {
		return nil, fmt.Errorf("id cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/graphql-rate-limiting-advanced/costs/%v", *id)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var costDeco GraphqlRateLimitingCostDecoration
	err = ErrorOrResponseError(s.client.Do(ctx, req, &costDeco))
	if err != nil {
		return nil, err
	}

	return &costDeco, nil
}

// Update updates a CostDecoration item in Kong.
// The given data must include the ID of an existing item.
func (s *GraphqlRateLimitingCostDecorationService) Update(
	ctx context.Context,
	costDeco *GraphqlRateLimitingCostDecoration,
) (*GraphqlRateLimitingCostDecoration, error) {
	if isEmptyString(costDeco.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/graphql-rate-limiting-advanced/costs/%v", *costDeco.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, costDeco)
	if err != nil {
		return nil, err
	}

	var updatedAPI GraphqlRateLimitingCostDecoration
	err = ErrorOrResponseError(s.client.Do(ctx, req, &updatedAPI))
	if err != nil {
		return nil, err
	}

	return &updatedAPI, nil
}

// Delete deletes a CostDecoration item from Kong given an ID.
func (s *GraphqlRateLimitingCostDecorationService) Delete(
	ctx context.Context,
	id *string,
) error {
	if isEmptyString(id) {
		return fmt.Errorf("ID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/graphql-rate-limiting-advanced/costs/%v", *id)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	err = ErrorOrResponseError(s.client.Do(ctx, req, nil))
	return err
}

// List fetches a list of CostDecoration items from Kong.
// opt can be used to control pagination.
func (s *GraphqlRateLimitingCostDecorationService) List(
	ctx context.Context,
	opt *ListOpt,
) ([]*GraphqlRateLimitingCostDecoration, *ListOpt, error) {
	data, next, err := s.client.list(ctx, "/graphql-rate-limiting-advanced/costs", opt)
	if err != nil {
		return nil, nil, err
	}
	costDecos := make([]*GraphqlRateLimitingCostDecoration, 0, len(data))

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var deco GraphqlRateLimitingCostDecoration
		err = json.Unmarshal(b, &deco)
		if err != nil {
			return nil, nil, err
		}
		costDecos = append(costDecos, &deco)
	}

	return costDecos, next, nil
}

// ListAll fetches all CostDecoration items present in Kong.
// This method can take a while to pull all pages of content
// if there are many items present.
func (s *GraphqlRateLimitingCostDecorationService) ListAll(
	ctx context.Context,
) ([]*GraphqlRateLimitingCostDecoration, error) {
	var decos, data []*GraphqlRateLimitingCostDecoration
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		decos = append(decos, data...)
	}
	return decos, nil
}

// CreateForService creates a CostDecoration item in Kong for the GraphQL rate limiting
// advanced plugin, scoped to a specific Service.
// The Service must be specified in the cost decoration.
// If no ID is provided, create operation is performed otherwise update operation is performed.
func (s *GraphqlRateLimitingCostDecorationService) CreateForService(
	ctx context.Context,
	costDeco *GraphqlRateLimitingCostDecoration,
) (*GraphqlRateLimitingCostDecoration, error) {
	if costDeco == nil {
		return nil, fmt.Errorf("cannot create a nil cost decoration")
	}
	serviceNameOrID := getServiceNameOrID(costDeco.Service)
	if serviceNameOrID == nil {
		return nil, fmt.Errorf("cannot create a cost decoration for a service without name or ID")
	}
	endpoint := fmt.Sprintf("/services/%s/graphql-rate-limiting-advanced/costs", *serviceNameOrID)
	method := "POST"
	if costDeco.ID != nil {
		method = "PUT"
		endpoint = fmt.Sprintf("/services/%s/graphql-rate-limiting-advanced/costs/%s", *serviceNameOrID, *costDeco.ID)
	}
	req, err := s.client.NewRequest(method, endpoint, nil, costDeco)
	if err != nil {
		return nil, err
	}

	var createdCostDeco GraphqlRateLimitingCostDecoration
	err = ErrorOrResponseError(s.client.Do(ctx, req, &createdCostDeco))
	if err != nil {
		return nil, err
	}

	return &createdCostDeco, nil
}

// UpdateForService updates a CostDecoration item in Kong, scoped to a specific Service.
// The given data must include the ID and Service of an existing item.
func (s *GraphqlRateLimitingCostDecorationService) UpdateForService(
	ctx context.Context,
	costDeco *GraphqlRateLimitingCostDecoration,
) (*GraphqlRateLimitingCostDecoration, error) {
	if isEmptyString(costDeco.ID) {
		return nil, fmt.Errorf("ID cannot be nil for UpdateForService operation")
	}
	serviceNameOrID := getServiceNameOrID(costDeco.Service)
	if serviceNameOrID == nil {
		return nil, fmt.Errorf("cannot update a cost decoration without a valid service")
	}

	endpoint := fmt.Sprintf("/services/%s/graphql-rate-limiting-advanced/costs/%s",
		*serviceNameOrID, *costDeco.ID)
	req, err := s.client.NewRequest("PUT", endpoint, nil, costDeco)
	if err != nil {
		return nil, err
	}

	var updatedCostDeco GraphqlRateLimitingCostDecoration
	err = ErrorOrResponseError(s.client.Do(ctx, req, &updatedCostDeco))
	if err != nil {
		return nil, err
	}

	return &updatedCostDeco, nil
}

// DeleteForService deletes a CostDecoration item in Kong, scoped to a specific Service.
// The given data must include the ID and Service of an existing item.
func (s *GraphqlRateLimitingCostDecorationService) DeleteForService(
	ctx context.Context,
	costDeco *GraphqlRateLimitingCostDecoration,
) error {
	if isEmptyString(costDeco.ID) {
		return fmt.Errorf("ID cannot be nil for DeleteForService operation")
	}
	serviceNameOrID := getServiceNameOrID(costDeco.Service)
	if serviceNameOrID == nil {
		return fmt.Errorf("cannot delete a cost decoration without a valid service")
	}

	endpoint := fmt.Sprintf("/services/%s/graphql-rate-limiting-advanced/costs/%s",
		*serviceNameOrID, *costDeco.ID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, costDeco)
	if err != nil {
		return err
	}
	err = ErrorOrResponseError(s.client.Do(ctx, req, nil))
	return err
}

// ListForService fetches a list of CostDecoration items from Kong,
// all associated to the specified Service.
// opt can be used to control pagination.
func (s *GraphqlRateLimitingCostDecorationService) ListForService(
	ctx context.Context,
	serviceNameOrID *string,
	opt *ListOpt,
) ([]*GraphqlRateLimitingCostDecoration, *ListOpt, error) {
	if isEmptyString(serviceNameOrID) {
		return nil, nil, fmt.Errorf("serviceNameOrID cannot be nil for listing cost decorations")
	}

	endpoint := fmt.Sprintf("/services/%s/graphql-rate-limiting-advanced/costs", *serviceNameOrID)
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	costDecos := make([]*GraphqlRateLimitingCostDecoration, 0, len(data))

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var deco GraphqlRateLimitingCostDecoration
		err = json.Unmarshal(b, &deco)
		if err != nil {
			return nil, nil, err
		}
		costDecos = append(costDecos, &deco)
	}

	return costDecos, next, nil
}

// ListAllForService fetches all CostDecoration items associated with the given Service present in Kong.
// This method can take a while to pull all pages of content
// if there are many items present.
func (s *GraphqlRateLimitingCostDecorationService) ListAllForService(
	ctx context.Context,
	serviceNameOrID *string,
) ([]*GraphqlRateLimitingCostDecoration, error) {
	var decos, data []*GraphqlRateLimitingCostDecoration
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.ListForService(ctx, serviceNameOrID, opt)
		if err != nil {
			return nil, err
		}
		decos = append(decos, data...)
	}
	return decos, nil
}
