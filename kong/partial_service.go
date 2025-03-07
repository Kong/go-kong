package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

// AbstractPartialService handles Partials in Kong.
type AbstractPartialService interface {
	// Create creates a Partial in Kong.
	Create(ctx context.Context, partial *Partial) (*Partial, error)
	// Get fetches a Partial in Kong.
	Get(ctx context.Context, partialID *string) (*Partial, error)
	// Update updates a Partial in Kong
	Update(ctx context.Context, partial *Partial) (*Partial, error)
	// Delete deletes a Partial in Kong
	Delete(ctx context.Context, partialID *string) error
	// List fetches a list of Partials in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*Partial, *ListOpt, error)
	// ListAll fetches all Partials in Kong.
	ListAll(ctx context.Context) ([]*Partial, error)
	// GetFullSchema retrieves the full schema of a partial.
	// This makes the use of `/schemas` endpoint in Kong.
	GetFullSchema(ctx context.Context, partialName *string) (Schema, error)
	// GetLinkedPlugins fetches a list of Plugins in Kong,
	// linked with the Partial. ListOpt can be used to control pagination.
	GetLinkedPlugins(ctx context.Context, partialID *string, opt *ListOpt) ([]*Plugin, *ListOpt, error)
}

// PartialService handles Partials in Kong.
type PartialService service

// Create creates a Partial in Kong.
// If an ID is specified, it will be used to
// create a Partial in Kong, otherwise an ID
// is auto-generated.
func (s *PartialService) Create(ctx context.Context,
	partial *Partial,
) (*Partial, error) {
	if partial == nil {
		return nil, fmt.Errorf("cannot create a nil partial")
	}

	if partial.Type == nil {
		return nil, fmt.Errorf("partial type cannot be nil")
	}

	queryPath := "/partials"
	method := "POST"
	if partial.ID != nil {
		queryPath = queryPath + "/" + *partial.ID
		method = "PUT"
	}

	req, err := s.client.NewRequest(method, queryPath, nil, partial)
	if err != nil {
		return nil, err
	}

	var createdPartial Partial
	_, err = s.client.Do(ctx, req, &createdPartial)
	if err != nil {
		return nil, err
	}
	return &createdPartial, nil
}

// Get fetches a Partial in Kong.
func (s *PartialService) Get(ctx context.Context,
	partialID *string,
) (*Partial, error) {
	if isEmptyString(partialID) {
		return nil, fmt.Errorf("partialID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/partials/%v", *partialID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var partial Partial
	_, err = s.client.Do(ctx, req, &partial)
	if err != nil {
		return nil, err
	}
	return &partial, nil
}

// Update updates a Partial in Kong
func (s *PartialService) Update(ctx context.Context,
	partial *Partial,
) (*Partial, error) {
	if partial == nil {
		return nil, fmt.Errorf("cannot update a nil partial")
	}

	if isEmptyString(partial.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/partials/%v", *partial.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, partial)
	if err != nil {
		return nil, err
	}

	var updatedPartial Partial
	_, err = s.client.Do(ctx, req, &updatedPartial)
	if err != nil {
		return nil, err
	}
	return &updatedPartial, nil
}

// Delete deletes a Partial in Kong
func (s *PartialService) Delete(ctx context.Context,
	partialID *string,
) error {
	if isEmptyString(partialID) {
		return fmt.Errorf("partialID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/partials/%v", *partialID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of Partials in Kong.
// opt can be used to control pagination.
func (s *PartialService) List(ctx context.Context,
	opt *ListOpt,
) ([]*Partial, *ListOpt, error) {
	data, next, err := s.client.list(ctx, "/partials", opt)
	if err != nil {
		return nil, nil, err
	}
	var partials []*Partial

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var p Partial
		err = json.Unmarshal(b, &p)
		if err != nil {
			return nil, nil, err
		}
		partials = append(partials, &p)
	}

	return partials, next, nil
}

// ListAll fetches all Partials in Kong.
func (s *PartialService) ListAll(ctx context.Context) ([]*Partial, error) {
	var partials, data []*Partial
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		partials = append(partials, data...)
	}
	return partials, nil
}

// GetLinkedPlugins fetches a list of Plugins in Kong,
// linked with the Partial.
// opt can be used to control pagination.
func (s *PartialService) GetLinkedPlugins(ctx context.Context,
	partialID *string, opt *ListOpt,
) ([]*Plugin, *ListOpt, error) {
	if isEmptyString(partialID) {
		return nil, nil, fmt.Errorf("partialID cannot be nil for GetLinkedPlugins operation")
	}

	endpoint := fmt.Sprintf("/partials/%v/links", *partialID)
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}

	var plugins []*Plugin

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var p Plugin
		err = json.Unmarshal(b, &p)
		if err != nil {
			return nil, nil, err
		}
		plugins = append(plugins, &p)
	}

	return plugins, next, nil
}

// GetFullSchema retrieves the full schema of a partial.
func (s *PartialService) GetFullSchema(ctx context.Context,
	partialName *string,
) (Schema, error) {
	if isEmptyString(partialName) {
		return nil, fmt.Errorf("partialName cannot be nil for GetFullSchema operation")
	}

	endpoint := fmt.Sprintf("/schemas/partials/%v", *partialName)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var schema Schema
	_, err = s.client.Do(ctx, req, &schema)
	if err != nil {
		return nil, err
	}
	return schema, nil
}
