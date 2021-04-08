package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// AbstractTargetService handles Targets in Kong.
type AbstractTargetService interface {
	// Create creates a Target in Kong under upstreamID.
	Create(ctx context.Context, upstreamNameOrID *string, target *Target) (*Target, error)
	// Delete deletes a Target in Kong
	Delete(ctx context.Context, upstreamNameOrID *string, targetOrID *string) error
	// List fetches a list of Targets in Kong.
	List(ctx context.Context, upstreamNameOrID *string, opt *ListOpt) ([]*Target, *ListOpt, error)
	// ListAll fetches all Targets in Kong for an upstream.
	ListAll(ctx context.Context, upstreamNameOrID *string) ([]*Target, error)
	// MarkHealthy marks target belonging to upstreamNameOrID as healthy in
	// Kong's load balancer.
	MarkHealthy(ctx context.Context, upstreamNameOrID *string, target *Target) error
	// MarkUnhealthy marks target belonging to upstreamNameOrID as unhealthy in
	// Kong's load balancer.
	MarkUnhealthy(ctx context.Context, upstreamNameOrID *string, target *Target) error
}

// TargetService handles Targets in Kong.
type TargetService service

// TODO foreign key can be read directly from the embedded key itself
// upstreamNameOrID need not be an explicit parameter.

// Create creates a Target in Kong under upstreamID.
// If an ID is specified, it will be used to
// create a target in Kong, otherwise an ID
// is auto-generated.
func (s *TargetService) Create(ctx context.Context,
	upstreamNameOrID *string, target *Target) (*Target, error) {
	if isEmptyString(upstreamNameOrID) {
		return nil, errors.New("upstreamNameOrID can not be nil")
	}
	queryPath := "/upstreams/" + *upstreamNameOrID + "/targets"
	method := "POST"
	// if target.ID != nil {
	// 	queryPath = queryPath + "/" + *target.ID
	// 	method = "PUT"
	// }
	req, err := s.client.NewRequest(method, queryPath, nil, target)

	if err != nil {
		return nil, err
	}

	var createdTarget Target
	_, err = s.client.Do(ctx, req, &createdTarget)
	if err != nil {
		return nil, err
	}
	return &createdTarget, nil
}

// Delete deletes a Target in Kong
func (s *TargetService) Delete(ctx context.Context,
	upstreamNameOrID *string, targetOrID *string) error {
	if isEmptyString(upstreamNameOrID) {
		return errors.New("upstreamNameOrID cannot be nil for Get operation")
	}
	if isEmptyString(targetOrID) {
		return errors.New("targetOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/upstreams/%v/targets/%v",
		*upstreamNameOrID, *targetOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of Targets in Kong.
// opt can be used to control pagination.
func (s *TargetService) List(ctx context.Context,
	upstreamNameOrID *string, opt *ListOpt) ([]*Target, *ListOpt, error) {
	if isEmptyString(upstreamNameOrID) {
		return nil, nil, errors.New("upstreamNameOrID cannot be nil for Get operation")
	}
	return s.listByEndpointAndOpt(ctx, "/upstreams/"+*upstreamNameOrID+"/targets", opt)
}

// ListAll fetches all Targets in Kong for an upstream.
func (s *TargetService) ListAll(ctx context.Context,
	upstreamNameOrID *string) ([]*Target, error) {
	return s.listAllByTags(ctx, upstreamNameOrID, nil)
}

// ListAll fetches all Targets filtered by tags in Kong for an upstream.
func (s *TargetService) listAllByTags(ctx context.Context,
	upstreamNameOrID *string, tags []string) ([]*Target, error) {
	return s.listAllByOpt(ctx, upstreamNameOrID, newOpt(tags))
}

// ListAll fetches all Targets filtered by opt in Kong for an upstream.
func (s *TargetService) listAllByOpt(ctx context.Context,
	upstreamNameOrID *string, opt *ListOpt) ([]*Target, error) {
	if isEmptyString(upstreamNameOrID) {
		return nil, errors.New("upstreamNameOrID cannot be nil for Get operation")
	}
	return s.listAllByEndpointAndOpt(ctx, "/upstreams/"+*upstreamNameOrID+"/targets", opt)
}

// MarkHealthy marks target belonging to upstreamNameOrID as healthy in
// Kong's load balancer.
func (s *TargetService) MarkHealthy(ctx context.Context,
	upstreamNameOrID *string, target *Target) error {
	if target == nil {
		return errors.New("cannot set health status for a nil target")
	}
	if isEmptyString(target.ID) && isEmptyString(target.Target) {
		return errors.New("need at least one of target or ID to" +
			" set health status")
	}
	if isEmptyString(upstreamNameOrID) {
		return errors.New("upstreamNameOrID cannot be nil " +
			"for updating health check")
	}

	tid := target.ID
	if target.ID == nil {
		tid = target.Target
	}

	endpoint := fmt.Sprintf("/upstreams/%v/targets/%v/healthy",
		*upstreamNameOrID, *tid)
	req, err := s.client.NewRequest("POST", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// MarkUnhealthy marks target belonging to upstreamNameOrID as unhealthy in
// Kong's load balancer.
func (s *TargetService) MarkUnhealthy(ctx context.Context,
	upstreamNameOrID *string, target *Target) error {
	if target == nil {
		return errors.New("cannot set health status for a nil target")
	}
	if isEmptyString(target.ID) && isEmptyString(target.Target) {
		return errors.New("need at least one of target or ID to" +
			" set health status")
	}
	if isEmptyString(upstreamNameOrID) {
		return errors.New("upstreamNameOrID cannot be nil " +
			"for updating health check")
	}

	tid := target.ID
	if target.ID == nil {
		tid = target.Target
	}

	endpoint := fmt.Sprintf("/upstreams/%v/targets/%v/unhealthy",
		*upstreamNameOrID, *tid)
	req, err := s.client.NewRequest("POST", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *TargetService) listByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*Target, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	targets, err := asTarget(data)
	if err != nil {
		return nil, nil, err
	}
	return targets, next, nil
}

func (s *TargetService) listAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*Target, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, false)
	if err != nil {
		return nil, err
	}
	return asTarget(data)
}

func asTarget(data [][]byte) ([]*Target, error) {
	var targets []*Target
	for _, object := range data {
		var target Target
		err := json.Unmarshal(object, &target)
		if err != nil {
			return nil, err
		}
		targets = append(targets, &target)
	}
	return targets, nil
}
