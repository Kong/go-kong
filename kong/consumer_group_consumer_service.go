package kong

import (
	"context"
	"fmt"
)

// AbstractConsumerGroupConsumerService handles ConsumerGroups' Consumers in Kong.
type AbstractConsumerGroupConsumerService interface {
	// Create creates a ConsumerGroupConsumer in Kong.
	Create(ctx context.Context, consumerGroupNameOrID *string, consumerNameOrID *string) (*ConsumerGroupObject, error)
	// Delete deletes a ConsumerGroupConsumer in Kong
	Delete(ctx context.Context, consumerGroupNameOrID *string, consumerNameOrID *string) error
	// ListAll fetches all ConsumerGroup's Consumers in Kong.
	ListAll(ctx context.Context, consumerGroupNameOrID *string) ([]*Consumer, error)
}

// ConsumerGroupService handles ConsumerGroup in Kong.
type ConsumerGroupConsumerService service

// Create creates a ConsumerGroupConsumer in Kong.
// If an ID is specified, it will be used to
// create a consumer for a consumer group in Kong, otherwise an ID
// is auto-generated.
func (s *ConsumerGroupConsumerService) Create(ctx context.Context,
	consumerGroupNameOrID *string, consumerNameOrID *string,
) (*ConsumerGroupObject, error) {
	if isEmptyString(consumerGroupNameOrID) {
		return nil, fmt.Errorf("consumerGroupNameOrID cannot be nil for Create operation")
	}
	if isEmptyString(consumerNameOrID) {
		return nil, fmt.Errorf("consumerNameOrID cannot be nil for Create operation")
	}

	queryPath := "/consumer_groups/" + *consumerGroupNameOrID + "/consumers"
	consumerGroupConsumerReq := &consumerGroupConsumerRequestBody{Consumer: consumerNameOrID}

	req, err := s.client.NewRequest("POST", queryPath, nil, consumerGroupConsumerReq)
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

// Delete deletes a ConsumerGroupConsumer in Kong
func (s *ConsumerGroupConsumerService) Delete(ctx context.Context,
	consumerGroupNameOrID *string, consumerNameOrID *string,
) error {
	if isEmptyString(consumerGroupNameOrID) {
		return fmt.Errorf("consumerGroupNameOrID cannot be nil for Delete operation")
	}
	if isEmptyString(consumerNameOrID) {
		return fmt.Errorf("consumerNameOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf(
		"/consumer_groups/%v/consumers/%v", *consumerGroupNameOrID, *consumerNameOrID,
	)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list all of ConsumerGroup's consumers in Kong.
func (s *ConsumerGroupConsumerService) ListAll(
	ctx context.Context, consumerGroupNameOrID *string,
) ([]*Consumer, error) {
	if isEmptyString(consumerGroupNameOrID) {
		return nil, fmt.Errorf("consumerGroupNameOrID cannot be nil for ListAll operation")
	}
	endpoint := fmt.Sprintf(
		"/consumer_groups/%v/consumers", *consumerGroupNameOrID,
	)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var cg consumerGroupConsumers
	_, err = s.client.Do(ctx, req, &cg)
	if err != nil {
		return nil, err
	}

	var consumers []*Consumer

	// See https://github.com/Kong/kong-ee/pull/6421
	if cg.Consumers != nil {
		// Kong < 3.5
		consumers = cg.Consumers
	} else if cg.Data != nil {
		// Kong >= 3.5
		consumers = cg.Data
	} else {
		consumers = make([]*Consumer, 0)
	}

	return consumers, nil
}
