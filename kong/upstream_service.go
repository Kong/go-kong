package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// AbstractUpstreamService handles Upstreams in Kong.
type AbstractUpstreamService interface {
	// Create creates a Upstream in Kong.
	Create(ctx context.Context, upstream *Upstream) (*Upstream, error)
	// Get fetches a Upstream in Kong.
	Get(ctx context.Context, upstreamNameOrID *string) (*Upstream, error)
	// Update updates a Upstream in Kong
	Update(ctx context.Context, upstream *Upstream) (*Upstream, error)
	// Delete deletes a Upstream in Kong
	Delete(ctx context.Context, upstreamNameOrID *string) error
	// List fetches a list of Upstreams in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*Upstream, *ListOpt, error)
	// ListAll fetches all Upstreams in Kong.
	ListAll(ctx context.Context) ([]*Upstream, error)
	// ListAll fetches all Upstreams filtered by tags in Kong.
	ListAllByTags(ctx context.Context, tags []string) ([]*Upstream, error)
}

// UpstreamService handles Upstreams in Kong.
type UpstreamService service

// Create creates a Upstream in Kong.
// If an ID is specified, it will be used to
// create a upstream in Kong, otherwise an ID
// is auto-generated.
func (s *UpstreamService) Create(ctx context.Context,
	upstream *Upstream) (*Upstream, error) {

	queryPath := "/upstreams"
	method := "POST"
	if upstream.ID != nil {
		queryPath = queryPath + "/" + *upstream.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, queryPath, nil, upstream)

	if err != nil {
		return nil, err
	}

	var createdUpstream Upstream
	_, err = s.client.Do(ctx, req, &createdUpstream)
	if err != nil {
		return nil, err
	}
	return &createdUpstream, nil
}

// Get fetches a Upstream in Kong.
func (s *UpstreamService) Get(ctx context.Context,
	upstreamNameOrID *string) (*Upstream, error) {

	if isEmptyString(upstreamNameOrID) {
		return nil, errors.New("upstreamNameOrID cannot" +
			" be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/upstreams/%v", *upstreamNameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var upstream Upstream
	_, err = s.client.Do(ctx, req, &upstream)
	if err != nil {
		return nil, err
	}
	return &upstream, nil
}

// Update updates a Upstream in Kong
func (s *UpstreamService) Update(ctx context.Context,
	upstream *Upstream) (*Upstream, error) {

	if isEmptyString(upstream.ID) {
		return nil, errors.New("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/upstreams/%v", *upstream.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, upstream)
	if err != nil {
		return nil, err
	}

	var updatedUpstream Upstream
	_, err = s.client.Do(ctx, req, &updatedUpstream)
	if err != nil {
		return nil, err
	}
	return &updatedUpstream, nil
}

// Delete deletes a Upstream in Kong
func (s *UpstreamService) Delete(ctx context.Context,
	upstreamNameOrID *string) error {

	if isEmptyString(upstreamNameOrID) {
		return errors.New("upstreamNameOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/upstreams/%v", *upstreamNameOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of Upstreams in Kong.
// opt can be used to control pagination.
func (s *UpstreamService) List(ctx context.Context,
	opt *ListOpt) ([]*Upstream, *ListOpt, error) {
	return s.ListByEndpointAndOpt(ctx, "/upstreams", opt)
}

// ListAll fetches all Upstreams in Kong.
// This method can take a while if there
// a lot of Upstreams present.
func (s *UpstreamService) ListAll(ctx context.Context) ([]*Upstream, error) {
	return s.ListAllByTags(ctx, nil)
}

func (s *UpstreamService) ListAllByTags(ctx context.Context, tags []string) ([]*Upstream, error) {
	return s.ListAllByEndpointAndOpt(ctx, "/upstreams", newOpt(tags))
}

func (s *UpstreamService) ListByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*Upstream, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	var upstreams []*Upstream

	for _, object := range data {
		var upstream Upstream
		err = json.Unmarshal(object, &upstream)
		if err != nil {
			return nil, nil, err
		}
		upstreams = append(upstreams, &upstream)
	}

	return upstreams, next, nil
}

func (s *UpstreamService) ListAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*Upstream, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, false)
	if err != nil {
		return nil, err
	}
	var upstreams []*Upstream

	for _, object := range data {
		var upstream Upstream
		err = json.Unmarshal(object, &upstream)
		if err != nil {
			return nil, err
		}
		upstreams = append(upstreams, &upstream)
	}

	return upstreams, nil
}
