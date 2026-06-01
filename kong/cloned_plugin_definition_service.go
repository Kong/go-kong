package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

// AbstractClonedPluginService handles ClonedPluginDefintions in Kong.
type AbstractClonedPluginService interface {
	// Create creates a ClonedPluginDefintion in Kong.
	Create(ctx context.Context, plugin *ClonedPluginDefinition) (*ClonedPluginDefinition, error)
	// Get fetches a ClonedPluginDefintion in Kong.
	Get(ctx context.Context, nameOrID *string) (*ClonedPluginDefinition, error)
	// Update updates a ClonedPluginDefintion in Kong.
	Update(ctx context.Context, plugin *ClonedPluginDefinition) (*ClonedPluginDefinition, error)
	// Delete deletes a ClonedPluginDefintion in Kong.
	Delete(ctx context.Context, nameOrID *string) error
	// List fetches a list of ClonedPluginDefintions in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*ClonedPluginDefinition, *ListOpt, error)
	// ListAll fetches all ClonedPluginDefintions in Kong.
	ListAll(ctx context.Context) ([]*ClonedPluginDefinition, error)
}

// ClonedPluginService handles ClonedPluginDefintions in Kong.
type ClonedPluginService service

// Create creates a ClonedPluginDefintion in Kong.
// If an ID is specified, it will be used to
// create the cloned plugin definition in Kong, otherwise an ID
// is auto-generated.
func (s *ClonedPluginService) Create(ctx context.Context,
	plugin *ClonedPluginDefinition,
) (*ClonedPluginDefinition, error) {
	queryPath := "/cloned-plugins"
	method := "POST"
	if plugin.ID != nil {
		queryPath = queryPath + "/" + *plugin.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, queryPath, nil, plugin)
	if err != nil {
		return nil, err
	}

	var createdPlugin ClonedPluginDefinition
	_, err = s.client.Do(ctx, req, &createdPlugin)
	if err != nil {
		return nil, err
	}
	return &createdPlugin, nil
}

// Get fetches a ClonedPluginDefintion in Kong.
func (s *ClonedPluginService) Get(ctx context.Context,
	nameOrID *string,
) (*ClonedPluginDefinition, error) {
	if isEmptyString(nameOrID) {
		return nil, fmt.Errorf("nameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/cloned-plugins/%v", *nameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var plugin ClonedPluginDefinition
	_, err = s.client.Do(ctx, req, &plugin)
	if err != nil {
		return nil, err
	}
	return &plugin, nil
}

// Update updates a ClonedPluginDefintion in Kong.
func (s *ClonedPluginService) Update(ctx context.Context,
	plugin *ClonedPluginDefinition,
) (*ClonedPluginDefinition, error) {
	if isEmptyString(plugin.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	method := "PATCH"
	if s.client.isKonnect {
		method = "PUT"
	}

	endpoint := fmt.Sprintf("/cloned-plugins/%v", *plugin.ID)
	req, err := s.client.NewRequest(method, endpoint, nil, plugin)
	if err != nil {
		return nil, err
	}

	var updatedPlugin ClonedPluginDefinition
	_, err = s.client.Do(ctx, req, &updatedPlugin)
	if err != nil {
		return nil, err
	}
	return &updatedPlugin, nil
}

// Delete deletes a ClonedPluginDefintion in Kong.
func (s *ClonedPluginService) Delete(ctx context.Context,
	nameOrID *string,
) error {
	if isEmptyString(nameOrID) {
		return fmt.Errorf("nameOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/cloned-plugins/%v", *nameOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of ClonedPluginDefintions in Kong.
// opt can be used to control pagination.
func (s *ClonedPluginService) List(ctx context.Context,
	opt *ListOpt,
) ([]*ClonedPluginDefinition, *ListOpt, error) {
	data, next, err := s.client.list(ctx, "/cloned-plugins", opt)
	if err != nil {
		return nil, nil, err
	}
	var plugins []*ClonedPluginDefinition

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var plugin ClonedPluginDefinition
		err = json.Unmarshal(b, &plugin)
		if err != nil {
			return nil, nil, err
		}
		plugins = append(plugins, &plugin)
	}

	return plugins, next, nil
}

// ListAll fetches all ClonedPluginDefintions in Kong.
// This method can take a while if there
// are a lot of ClonedPluginDefintions present.
func (s *ClonedPluginService) ListAll(ctx context.Context) ([]*ClonedPluginDefinition, error) {
	var plugins, data []*ClonedPluginDefinition
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
