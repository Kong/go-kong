package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

type AbstractGraphqlRateLimitingCostDecorationService interface {
	// Creates a cost decoration for the GraphQL rate-limiting plugin in Kong.
	Create(ctx context.Context, costDeco *GraphqlRateLimitingCostDecoration) (*GraphqlRateLimitingCostDecoration, error)
	// Fetches a cost decoration for the GraphQL rate-limiting plugin from Kong.
	Get(ctx context.Context, ID *string) (*GraphqlRateLimitingCostDecoration, error)
	// UPdates a cost decoration for the GraphQL rate-limiting plugin in Kong.
	Update(ctx context.Context, costDeco *GraphqlRateLimitingCostDecoration) (*GraphqlRateLimitingCostDecoration, error)
	// Deletes a cost decoration for the GraphQL rate-limiting plugin in Kong.
	Delete(ctx context.Context, ID *string) error
	// Retrieves a page of cost decorations for the GraphQL rate-limiting plugin in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*GraphqlRateLimitingCostDecoration, *ListOpt, error)
	// Retrieves all decorations for the GraphQL rate-limiting plugin in Kong.
	ListAll(ctx context.Context) ([]*GraphqlRateLimitingCostDecoration, error)
}

type GraphqlRateLimitingCostDecorationService service

// Create creates a CostDecoration item in Kong for the GraphQL rate limiting advanced plugin.
func (s *GraphqlRateLimitingCostDecorationService) Create(
	ctx context.Context,
	costDeco *GraphqlRateLimitingCostDecoration,
) (*GraphqlRateLimitingCostDecoration, error) {
	queryPath := "/graphql-rate-limiting-advanced/costs"
	method := "POST"

	if costDeco.ID != nil {
		return nil, fmt.Errorf("can't specify an ID for creating new Cost Decoration")
	}
	req, err := s.client.NewRequest(method, queryPath, nil, costDeco)
	if err != nil {
		return nil, err
	}

	var createdCostDeco GraphqlRateLimitingCostDecoration
	_, err = s.client.Do(ctx, req, &createdCostDeco)
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
	_, err = s.client.Do(ctx, req, &costDeco)
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
	_, err = s.client.Do(ctx, req, &updatedAPI)
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

	_, err = s.client.Do(ctx, req, nil)
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
// if there are many items present
func (s *GraphqlRateLimitingCostDecorationService) ListAll(
	ctx context.Context,
) (
	[]*GraphqlRateLimitingCostDecoration, error,
) {
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
