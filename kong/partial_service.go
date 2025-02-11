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
	Get(ctx context.Context, nameOrID *string) (*Partial, error)
	// Update updates a Partial in Kong
	Update(ctx context.Context, partial *Partial) (*Partial, error)
	// Delete deletes a Partial in Kong
	Delete(ctx context.Context, nameOrID *string) error
	// List fetches a list of Partials in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*Partial, *ListOpt, error)
	// ListAll fetches all Partials in Kong.
	ListAll(ctx context.Context) ([]*Partial, error)
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
	nameOrID *string,
) (*Partial, error) {
	if isEmptyString(nameOrID) {
		return nil, fmt.Errorf("nameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/partials/%v", *nameOrID)
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
	nameOrID *string,
) error {
	if isEmptyString(nameOrID) {
		return fmt.Errorf("nameOrID cannot be nil for Delete operation")
	}

	// TODO: Uncomment when /partials/:id/links endpoint is enabled
	// plugins, _, err := s.GetLinkedPlugins(ctx, nameOrID, nil)
	// if err != nil {
	// 	return err
	// }

	// if len(plugins) != 0 {
	// 	return fmt.Errorf("cannot delete partial %v, it is still linked to plugins: %v",
	// 		*nameOrID, plugins)
	// }

	endpoint := fmt.Sprintf("/partials/%v", *nameOrID)
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

// TODO: Uncomment when /partials/:id/links endpoint is enabled
// GetLinkedPlugins fetches a list of Plugins in Kong,
// linked with the Partial.
// opt can be used to control pagination.
// func (s *PartialService) GetLinkedPlugins(ctx context.Context,
// 	nameOrID *string, opt *ListOpt,
// ) ([]*Plugin, *ListOpt, error) {
// 	if isEmptyString(nameOrID) {
// 		return nil, nil, fmt.Errorf("nameOrID cannot be nil for GetLinkedPlugins operation")
// 	}

// 	endpoint := fmt.Sprintf("/partials/%v/links", *nameOrID)
// 	data, next, err := s.client.list(ctx, endpoint, opt)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	var plugins []*Plugin

// 	for _, object := range data {
// 		b, err := object.MarshalJSON()
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		var p Plugin
// 		err = json.Unmarshal(b, &p)
// 		if err != nil {
// 			return nil, nil, err
// 		}
// 		plugins = append(plugins, &p)
// 	}

// 	return plugins, next, nil
// }
