package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

// AbstractConsumerGroupService handles ConsumerGroups in Kong.
type AbstractConsumerGroupService interface {
	// Create creates a ConsumerGroup in Kong.
	Create(ctx context.Context, consumerGroup *ConsumerGroup) (*ConsumerGroup, error)
	// Get fetches a ConsumerGroup from Kong.
	Get(ctx context.Context, nameOrID *string) (*ConsumerGroupObject, error)
	// Update updates a ConsumerGroup in Kong
	Update(ctx context.Context, consumerGroup *ConsumerGroup) (*ConsumerGroup, error)
	// Delete deletes a ConsumerGroup in Kong
	Delete(ctx context.Context, usernameOrID *string) error
	// List fetches a list of ConsumerGroups in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*ConsumerGroup, *ListOpt, error)
	// ListAll fetches all ConsumerGroups in Kong.
	ListAll(ctx context.Context) ([]*ConsumerGroup, error)

	// UpdateRateLimitingAdvancedPlugin upsert a RLA plugin for ConsumerGroups in Kong.
	UpdateRateLimitingAdvancedPlugin(
		ctx context.Context, nameOrID *string, config map[string]Configuration,
	) (*ConsumerGroupRLA, error)
}

// ConsumerGroupService handles ConsumerGroup in Kong.
type ConsumerGroupService service

// Create creates a ConsumerGroup in Kong.
// If an ID is specified, it will be used to create a consumer group in Kong,
// otherwise an ID is auto-generated.
func (s *ConsumerGroupService) Create(ctx context.Context,
	consumerGroup *ConsumerGroup,
) (*ConsumerGroup, error) {
	queryPath := "/consumer_groups"
	method := "POST"
	if consumerGroup.ID != nil {
		queryPath = queryPath + "/" + *consumerGroup.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, queryPath, nil, consumerGroup)
	if err != nil {
		return nil, err
	}

	var cg ConsumerGroup
	_, err = s.client.Do(ctx, req, &cg)
	if err != nil {
		return nil, err
	}
	return &cg, nil
}

// Get fetches a ConsumerGroup from Kong.
func (s *ConsumerGroupService) Get(ctx context.Context,
	nameOrID *string,
) (*ConsumerGroupObject, error) {
	if isEmptyString(nameOrID) {
		return nil, fmt.Errorf("nameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/consumer_groups/%v", *nameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var cg ConsumerGroupObject
	_, err = s.client.Do(ctx, req, &cg)
	if err != nil {
		return nil, err
	}
	return &cg, nil
}

// Update updates a ConsumerGroup in Kong
func (s *ConsumerGroupService) Update(ctx context.Context,
	consumerGroup *ConsumerGroup,
) (*ConsumerGroup, error) {
	if isEmptyString(consumerGroup.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/consumer_groups/%v", *consumerGroup.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, consumerGroup)
	if err != nil {
		return nil, err
	}

	var updatedAPI ConsumerGroup
	_, err = s.client.Do(ctx, req, &updatedAPI)
	if err != nil {
		return nil, err
	}
	return &updatedAPI, nil
}

// Delete deletes a ConsumerGroup in Kong
func (s *ConsumerGroupService) Delete(ctx context.Context,
	usernameOrID *string,
) error {
	if isEmptyString(usernameOrID) {
		return fmt.Errorf("usernameOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/consumer_groups/%v", *usernameOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of ConsumerGroup in Kong.
// opt can be used to control pagination.
func (s *ConsumerGroupService) List(ctx context.Context,
	opt *ListOpt,
) ([]*ConsumerGroup, *ListOpt, error) {
	data, next, err := s.client.list(ctx, "/consumer_groups", opt)
	if err != nil {
		return nil, nil, err
	}
	var consumers []*ConsumerGroup

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var consumer ConsumerGroup
		err = json.Unmarshal(b, &consumer)
		if err != nil {
			return nil, nil, err
		}
		consumers = append(consumers, &consumer)
	}

	return consumers, next, nil
}

// ListAll fetches all ConsumerGroup in Kong.
func (s *ConsumerGroupService) ListAll(ctx context.Context) ([]*ConsumerGroup, error) {
	var consumerGroups, data []*ConsumerGroup
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		consumerGroups = append(consumerGroups, data...)
	}
	return consumerGroups, nil
}

// UpdateRateLimitingAdvancedPlugin upsert a RLA plugin for ConsumerGroups in Kong.
func (s *ConsumerGroupService) UpdateRateLimitingAdvancedPlugin(
	ctx context.Context, nameOrID *string, config map[string]Configuration,
) (*ConsumerGroupRLA, error) {
	if isEmptyString(nameOrID) {
		return nil, fmt.Errorf("ID cannot be nil")
	}

	endpoint := fmt.Sprintf(
		"/consumer_groups/%v/overrides/plugins/rate-limiting-advanced", *nameOrID,
	)
	req, err := s.client.NewRequest("PUT", endpoint, nil, config)
	if err != nil {
		return nil, err
	}

	var rla ConsumerGroupRLA
	_, err = s.client.Do(ctx, req, &rla)
	if err != nil {
		return nil, err
	}
	return &rla, nil
}
