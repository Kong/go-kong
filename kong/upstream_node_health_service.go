package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// AbstractUpstreamNodeHealthService handles Upstream Node Healths in Kong.
type AbstractUpstreamNodeHealthService interface {
	// List fetches a list of Upstream Node Healths in Kong.
	List(ctx context.Context, upstreamNameOrID *string, opt *ListOpt) ([]*UpstreamNodeHealth, *ListOpt, error)
	// ListAll fetches all Upstream Node Healths in Kong.
	ListAll(ctx context.Context, upstreamNameOrID *string) ([]*UpstreamNodeHealth, error)
	// ListAllByOpt fetches all Upstream Node Healths filtered by opt in Kong.
	ListAllByOpt(ctx context.Context, upstreamNameOrID *string, opt *ListOpt) ([]*UpstreamNodeHealth, error)
	// ListAllByTags fetches all Upstream Node Healths filtered by opt in Kong.
	ListAllByTags(ctx context.Context, upstreamNameOrID *string, tags []string) ([]*UpstreamNodeHealth, error)
}

// UpstreamNodeHealthService handles Upstream Node Healths in Kong.
type UpstreamNodeHealthService service

// List fetches a list of Upstream Node Healths in Kong.
// opt can be used to control pagination.
func (s *UpstreamNodeHealthService) List(ctx context.Context,
	upstreamNameOrID *string, opt *ListOpt) ([]*UpstreamNodeHealth, *ListOpt, error) {
	if isEmptyString(upstreamNameOrID) {
		return nil, nil, errors.New("upstreamNameOrID cannot be nil for Get operation")
	}
	endpoint := fmt.Sprintf("/upstreams/%v/health", *upstreamNameOrID)

	return s.listByEndpointAndOpt(ctx, endpoint, opt)
}

// ListAll fetches all Upstream Node Healths in Kong.
// This method can take a while if there are
// a lot of Upstream Node Healths present.
func (s *UpstreamNodeHealthService) ListAll(ctx context.Context,
	upstreamNameOrID *string) ([]*UpstreamNodeHealth, error) {
	return s.ListAllByTags(ctx, upstreamNameOrID, nil)
}

// ListAll fetches all Upstream Node Healths filtered by tags in Kong.
// This method can take a while if there are
// a lot of Upstream Node Healths present.
func (s *UpstreamNodeHealthService) ListAllByTags(ctx context.Context,
	upstreamNameOrID *string, tags []string) ([]*UpstreamNodeHealth, error) {
	return s.ListAllByOpt(ctx, upstreamNameOrID, newOpt(tags))
}

// ListAll fetches all Upstream Node Healths filtered by opt in Kong.
// This method can take a while if there are
// a lot of Upstream Node Healths present.
func (s *UpstreamNodeHealthService) ListAllByOpt(ctx context.Context,
	upstreamNameOrID *string, opt *ListOpt) ([]*UpstreamNodeHealth, error) {
	if isEmptyString(upstreamNameOrID) {
		return nil, errors.New("upstreamNameOrID cannot be nil for Get operation")
	}
	endpoint := fmt.Sprintf("/upstreams/%v/health", *upstreamNameOrID)
	return s.listAllByEndpointAndOpt(ctx, endpoint, opt)
}

func (s *UpstreamNodeHealthService) listByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*UpstreamNodeHealth, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	upstreamNodeHealths, err := asUpstreamNodeHealth(data)
	if err != nil {
		return nil, nil, err
	}
	return upstreamNodeHealths, next, nil
}

func (s *UpstreamNodeHealthService) listAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*UpstreamNodeHealth, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, false)
	if err != nil {
		return nil, err
	}
	upstreamNodeHealths, err := asUpstreamNodeHealth(data)
	if err != nil {
		return nil, err
	}
	return upstreamNodeHealths, nil
}

func asUpstreamNodeHealth(data [][]byte) ([]*UpstreamNodeHealth, error) {
	var upstreamNodeHealths []*UpstreamNodeHealth
	for _, object := range data {
		var upstreamNodeHealth UpstreamNodeHealth
		err := json.Unmarshal(object, &upstreamNodeHealth)
		if err != nil {
			return nil, err
		}
		upstreamNodeHealths = append(upstreamNodeHealths, &upstreamNodeHealth)
	}
	return upstreamNodeHealths, nil
}
