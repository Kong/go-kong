package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// AbstractConsumerService handles Consumers in Kong.
type AbstractConsumerService interface {
	// Create creates a Consumer in Kong.
	Create(ctx context.Context, consumer *Consumer) (*Consumer, error)
	// Get fetches a Consumer in Kong.
	Get(ctx context.Context, usernameOrID *string) (*Consumer, error)
	// GetByCustomID fetches a Consumer in Kong.
	GetByCustomID(ctx context.Context, customID *string) (*Consumer, error)
	// Update updates a Consumer in Kong
	Update(ctx context.Context, consumer *Consumer) (*Consumer, error)
	// Delete deletes a Consumer in Kong
	Delete(ctx context.Context, usernameOrID *string) error
	// List fetches a list of Consumers in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*Consumer, *ListOpt, error)
	// ListAll fetches all Consumers in Kong.
	ListAll(ctx context.Context) ([]*Consumer, error)
	// ListAllByTags fetches all Consumers filtered by tags in Kong.
	ListAllByTags(ctx context.Context, tags []string) ([]*Consumer, error)
}

// ConsumerService handles Consumers in Kong.
type ConsumerService service

// Create creates a Consumer in Kong.
// If an ID is specified, it will be used to
// create a consumer in Kong, otherwise an ID
// is auto-generated.
func (s *ConsumerService) Create(ctx context.Context,
	consumer *Consumer) (*Consumer, error) {

	queryPath := "/consumers"
	method := "POST"
	if consumer.ID != nil {
		queryPath = queryPath + "/" + *consumer.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, queryPath, nil, consumer)

	if err != nil {
		return nil, err
	}

	var createdConsumer Consumer
	_, err = s.client.Do(ctx, req, &createdConsumer)
	if err != nil {
		return nil, err
	}
	return &createdConsumer, nil
}

// Get fetches a Consumer in Kong.
func (s *ConsumerService) Get(ctx context.Context,
	usernameOrID *string) (*Consumer, error) {

	if isEmptyString(usernameOrID) {
		return nil, errors.New("usernameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/consumers/%v", *usernameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var consumer Consumer
	_, err = s.client.Do(ctx, req, &consumer)
	if err != nil {
		return nil, err
	}
	return &consumer, nil
}

// GetByCustomID fetches a Consumer in Kong.
func (s *ConsumerService) GetByCustomID(ctx context.Context,
	customID *string) (*Consumer, error) {

	if isEmptyString(customID) {
		return nil, errors.New("customID cannot be nil for Get operation")
	}

	type QS struct {
		CustomID string `url:"custom_id,omitempty"`
	}

	req, err := s.client.NewRequest("GET", "/consumers",
		&QS{CustomID: *customID}, nil)
	if err != nil {
		return nil, err
	}

	type Response struct {
		Data []Consumer
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

// Update updates a Consumer in Kong
func (s *ConsumerService) Update(ctx context.Context,
	consumer *Consumer) (*Consumer, error) {

	if isEmptyString(consumer.ID) {
		return nil, errors.New("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/consumers/%v", *consumer.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, consumer)
	if err != nil {
		return nil, err
	}

	var updatedAPI Consumer
	_, err = s.client.Do(ctx, req, &updatedAPI)
	if err != nil {
		return nil, err
	}
	return &updatedAPI, nil
}

// Delete deletes a Consumer in Kong
func (s *ConsumerService) Delete(ctx context.Context,
	usernameOrID *string) error {

	if isEmptyString(usernameOrID) {
		return errors.New("usernameOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/consumers/%v", *usernameOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of Consumers in Kong.
// opt can be used to control pagination.
func (s *ConsumerService) List(ctx context.Context,
	opt *ListOpt) ([]*Consumer, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/consumers", opt)
}

// ListAll fetches all Consumers in Kong.
// This method can take a while if there
// a lot of Consumers present.
func (s *ConsumerService) ListAll(ctx context.Context) ([]*Consumer, error) {
	return s.ListAllByTags(ctx, nil)
}

// ListAll fetches all Consumers filtered by tags in Kong.
// This method can take a while if there
// a lot of Consumers present.
func (s *ConsumerService) ListAllByTags(ctx context.Context, tags []string) ([]*Consumer, error) {
	return s.listAllByEndpointAndOpt(ctx, "/consumers", newOpt(tags))
}

func (s *ConsumerService) listByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*Consumer, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	var consumers []*Consumer

	for _, object := range data {
		var consumer Consumer
		err = json.Unmarshal(object, &consumer)
		if err != nil {
			return nil, nil, err
		}
		consumers = append(consumers, &consumer)
	}

	return consumers, next, nil
}

func (s *ConsumerService) listAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*Consumer, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, false)
	if err != nil {
		return nil, err
	}
	var consumers []*Consumer

	for _, object := range data {
		var consumer Consumer
		err = json.Unmarshal(object, &consumer)
		if err != nil {
			return nil, err
		}
		consumers = append(consumers, &consumer)
	}

	return consumers, nil
}
