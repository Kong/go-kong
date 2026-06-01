package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

// AbstractCustomPluginService handles CustomPluginDefintions in Kong.
type AbstractCustomPluginService interface {
	// Create creates a CustomPluginDefintion in Kong.
	Create(ctx context.Context, plugin *CustomPluginDefinition) (*CustomPluginDefinition, error)
	// Get fetches a CustomPluginDefintion in Kong.
	Get(ctx context.Context, nameOrID *string) (*CustomPluginDefinition, error)
	// Update updates a CustomPluginDefintion in Kong.
	Update(ctx context.Context, plugin *CustomPluginDefinition) (*CustomPluginDefinition, error)
	// Delete deletes a CustomPluginDefintion in Kong.
	Delete(ctx context.Context, nameOrID *string) error
	// List fetches a list of CustomPluginDefintions in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*CustomPluginDefinition, *ListOpt, error)
	// ListAll fetches all CustomPluginDefintions in Kong.
	ListAll(ctx context.Context) ([]*CustomPluginDefinition, error)
}

// CustomPluginService handles CustomPluginDefintions in Kong.
type CustomPluginService service

// Create creates a CustomPluginDefintion in Kong.
// If an ID is specified, it will be used to
// create the custom plugin definition in Kong, otherwise an ID
// is auto-generated.
func (s *CustomPluginService) Create(ctx context.Context,
	plugin *CustomPluginDefinition,
) (*CustomPluginDefinition, error) {
	queryPath := "/custom-plugins"
	method := "POST"
	if plugin.ID != nil {
		queryPath = queryPath + "/" + *plugin.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, queryPath, nil, plugin)
	if err != nil {
		return nil, err
	}

	var createdPlugin CustomPluginDefinition
	_, err = s.client.Do(ctx, req, &createdPlugin)
	if err != nil {
		return nil, err
	}
	return &createdPlugin, nil
}

// Get fetches a CustomPluginDefintion in Kong.
func (s *CustomPluginService) Get(ctx context.Context,
	nameOrID *string,
) (*CustomPluginDefinition, error) {
	if isEmptyString(nameOrID) {
		return nil, fmt.Errorf("nameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/custom-plugins/%v", *nameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var plugin CustomPluginDefinition
	_, err = s.client.Do(ctx, req, &plugin)
	if err != nil {
		return nil, err
	}
	return &plugin, nil
}

// Update updates a CustomPluginDefintion in Kong.
func (s *CustomPluginService) Update(ctx context.Context,
	plugin *CustomPluginDefinition,
) (*CustomPluginDefinition, error) {
	if isEmptyString(plugin.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	method := "PATCH"
	if s.client.isKonnect {
		method = "PUT"
	}

	endpoint := fmt.Sprintf("/custom-plugins/%v", *plugin.ID)
	req, err := s.client.NewRequest(method, endpoint, nil, plugin)
	if err != nil {
		return nil, err
	}

	var updatedPlugin CustomPluginDefinition
	_, err = s.client.Do(ctx, req, &updatedPlugin)
	if err != nil {
		return nil, err
	}
	return &updatedPlugin, nil
}

// Delete deletes a CustomPluginDefintion in Kong.
func (s *CustomPluginService) Delete(ctx context.Context,
	nameOrID *string,
) error {
	if isEmptyString(nameOrID) {
		return fmt.Errorf("nameOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/custom-plugins/%v", *nameOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of CustomPluginDefintions in Kong.
// opt can be used to control pagination.
func (s *CustomPluginService) List(ctx context.Context,
	opt *ListOpt,
) ([]*CustomPluginDefinition, *ListOpt, error) {
	data, next, err := s.client.list(ctx, "/custom-plugins", opt)
	if err != nil {
		return nil, nil, err
	}
	var plugins []*CustomPluginDefinition

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var plugin CustomPluginDefinition
		err = json.Unmarshal(b, &plugin)
		if err != nil {
			return nil, nil, err
		}
		plugins = append(plugins, &plugin)
	}

	return plugins, next, nil
}

// ListAll fetches all CustomPluginDefintions in Kong.
// This method can take a while if there
// are a lot of CustomPluginDefintions present.
func (s *CustomPluginService) ListAll(ctx context.Context) ([]*CustomPluginDefinition, error) {
	var plugins, data []*CustomPluginDefinition
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, data...)
	}
	return plugins, nil
}
