package kong

import (
	"context"
	"fmt"
)

// AbstractEventHooks handles event hooks in Kong
type AbstractEventHooksService interface {
	AddWebhook(ctx context.Context, eventHooks *EventHooks) (*EventHooks, error)

	ListAllEventHooks(ctx context.Context) (interface{}, error)
	ListAllSources(ctx context.Context) (interface{}, error)
	ListAllEventsForSouce(ctx context.Context, sourceid *string) (interface{}, error)
}

type EventHooksService service

// AddWebhook make json post request to a required URL with event data as a payload
func (s *EventHooksService) AddWebhook(ctx context.Context, eventHooks *EventHooks) (*EventHooks, error) {
	endpoint := "/event-hooks/"
	req, err := s.client.NewRequest("POST", endpoint, nil, eventHooks)
	if err != nil {
		panic(err)
	}
	fmt.Printf(" to be created eventhook %v", req.Body)

	var createdEventHooks EventHooks
	_, err = s.client.Do(ctx, req, &createdEventHooks)
	if err != nil {
		panic(err)
	}
	return &createdEventHooks, nil
}

func (s *EventHooksService) ListAllEventHooks(ctx context.Context) (interface{}, error) {
	endpoint := "/event-hooks/"
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *EventHooksService) ListAllSources(ctx context.Context) (interface{}, error) {
	endpoint := "/event-hooks/sources"
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *EventHooksService) ListAllEventsForSouce(ctx context.Context, sourceid *string) (interface{}, error) {
	if sourceid == nil {
		return nil, fmt.Errorf("source id cannot be nil")
	}

	endpoint := fmt.Sprintf("/event-hooks/sources/%s", *sourceid)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
