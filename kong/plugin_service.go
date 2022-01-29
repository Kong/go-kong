package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// AbstractPluginService handles Plugins in Kong.
type AbstractPluginService interface {
	// Create creates a Plugin in Kong.
	Create(ctx context.Context, plugin *Plugin) (*Plugin, error)
	// Get fetches a Plugin in Kong.
	Get(ctx context.Context, usernameOrID *string) (*Plugin, error)
	// Update updates a Plugin in Kong
	Update(ctx context.Context, plugin *Plugin) (*Plugin, error)
	// Delete deletes a Plugin in Kong
	Delete(ctx context.Context, usernameOrID *string) error
	// List fetches a list of Plugins in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*Plugin, *ListOpt, error)
	// ListAll fetches all Plugins in Kong.
	ListAll(ctx context.Context) ([]*Plugin, error)
	// ListAllForConsumer fetches all Plugins in Kong enabled for a consumer.
	ListAllForConsumer(ctx context.Context, consumerIDorName *string) ([]*Plugin, error)
	// ListAllForService fetches all Plugins in Kong enabled for a service.
	ListAllForService(ctx context.Context, serviceIDorName *string) ([]*Plugin, error)
	// ListAllForRoute fetches all Plugins in Kong enabled for a service.
	ListAllForRoute(ctx context.Context, routeID *string) ([]*Plugin, error)
	// Validate validates a Plugin against its schema
	Validate(ctx context.Context, plugin *Plugin) (bool, string, error)
	// GetSchema retrieves the config schema of a plugin.
	//
	// Deprecated: Use GetFullSchema instead.
	GetSchema(ctx context.Context, pluginName *string) (Schema, error)
	// GetFullSchema retrieves the full schema of a plugin.
	// This makes the use of `/schemas` endpoint in Kong.
	GetFullSchema(ctx context.Context, pluginName *string) (Schema, error)
}

// PluginService handles Plugins in Kong.
type PluginService service

// GetFullSchema retrieves the full schema of a plugin.
func (s *PluginService) GetFullSchema(ctx context.Context,
	pluginName *string) (Schema, error) {
	if isEmptyString(pluginName) {
		return nil, fmt.Errorf("pluginName cannot be empty")
	}
	endpoint := fmt.Sprintf("/schemas/plugins/%v", *pluginName)
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

// GetSchema retrieves the config schema of a plugin
//
// Deprecated: Use GetPluginSchema instead
func (s *PluginService) GetSchema(ctx context.Context,
	pluginName *string) (Schema, error) {
	if isEmptyString(pluginName) {
		return nil, fmt.Errorf("pluginName cannot be empty")
	}
	endpoint := fmt.Sprintf("/plugins/schema/%v", *pluginName)
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

// Create creates a Plugin in Kong.
// If an ID is specified, it will be used to
// create a plugin in Kong, otherwise an ID
// is auto-generated.
func (s *PluginService) Create(ctx context.Context,
	plugin *Plugin) (*Plugin, error) {

	queryPath := "/plugins"
	method := "POST"
	if plugin.ID != nil {
		queryPath = queryPath + "/" + *plugin.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, queryPath, nil, plugin)
	if err != nil {
		return nil, err
	}

	var createdPlugin Plugin
	_, err = s.client.Do(ctx, req, &createdPlugin)
	if err != nil {
		return nil, err
	}
	return &createdPlugin, nil
}

// Get fetches a Plugin in Kong.
func (s *PluginService) Get(ctx context.Context,
	usernameOrID *string) (*Plugin, error) {

	if isEmptyString(usernameOrID) {
		return nil, fmt.Errorf("usernameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/plugins/%v", *usernameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var plugin Plugin
	_, err = s.client.Do(ctx, req, &plugin)
	if err != nil {
		return nil, err
	}
	return &plugin, nil
}

// Update updates a Plugin in Kong
func (s *PluginService) Update(ctx context.Context,
	plugin *Plugin) (*Plugin, error) {

	if isEmptyString(plugin.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/plugins/%v", *plugin.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, plugin)
	if err != nil {
		return nil, err
	}

	var updatedAPI Plugin
	_, err = s.client.Do(ctx, req, &updatedAPI)
	if err != nil {
		return nil, err
	}
	return &updatedAPI, nil
}

// Delete deletes a Plugin in Kong
func (s *PluginService) Delete(ctx context.Context,
	usernameOrID *string) error {

	if isEmptyString(usernameOrID) {
		return fmt.Errorf("usernameOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/plugins/%v", *usernameOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// Validate validates a Plugin against its schema
func (s *PluginService) Validate(ctx context.Context, plugin *Plugin) (bool, string, error) {
	endpoint := "/schemas/plugins/validate"
	req, err := s.client.NewRequest("POST", endpoint, nil, &plugin)
	if err != nil {
		return false, "", err
	}
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		// Arguably Kong should return a 422 Unprocessable Entity for a well-formed
		// HTTP request with a mangled plugin, but it doesn't, it returns a 400.
		// Hopefully (usually) we get a 400 because of a mangled plugin rather than
		// a mangled request, but we can't easily tell as messageFromBody masks errors
		if resp.StatusCode == http.StatusBadRequest {
			var apiError *APIError
			ok := errors.As(err, &apiError)
			if !ok {
				return false, "", err
			}
			return false, apiError.message, nil
		}
		return false, "", err
	}
	return resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK, "", nil
}

// listByPath fetches a list of Plugins in Kong
// on a specific path.
// This is a helper method for listing all plugins
// or plugins for specific entities.
func (s *PluginService) listByPath(ctx context.Context,
	path string, opt *ListOpt) ([]*Plugin, *ListOpt, error) {
	data, next, err := s.client.list(ctx, path, opt)
	if err != nil {
		return nil, nil, err
	}
	var plugins []*Plugin

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var plugin Plugin
		err = json.Unmarshal(b, &plugin)
		if err != nil {
			return nil, nil, err
		}
		plugins = append(plugins, &plugin)
	}

	return plugins, next, nil
}

// ListAll fetches all Plugins in Kong.
// This method can take a while if there
// a lot of Plugins present.
func (s *PluginService) listAllByPath(ctx context.Context,
	path string) ([]*Plugin, error) {
	var plugins, data []*Plugin
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.listByPath(ctx, path, opt)
		if err != nil {
			return nil, err
		}
		plugins = append(plugins, data...)
	}
	return plugins, nil
}

// List fetches a list of Plugins in Kong.
// opt can be used to control pagination.
func (s *PluginService) List(ctx context.Context,
	opt *ListOpt) ([]*Plugin, *ListOpt, error) {
	return s.listByPath(ctx, "/plugins", opt)
}

// ListAll fetches all Plugins in Kong.
// This method can take a while if there
// a lot of Plugins present.
func (s *PluginService) ListAll(ctx context.Context) ([]*Plugin, error) {
	return s.listAllByPath(ctx, "/plugins")
}

// ListAllForConsumer fetches all Plugins in Kong enabled for a consumer.
func (s *PluginService) ListAllForConsumer(ctx context.Context,
	consumerIDorName *string) ([]*Plugin, error) {
	if isEmptyString(consumerIDorName) {
		return nil, fmt.Errorf("consumerIDorName cannot be nil")
	}
	return s.listAllByPath(ctx, "/consumers/"+*consumerIDorName+"/plugins")
}

// ListAllForService fetches all Plugins in Kong enabled for a service.
func (s *PluginService) ListAllForService(ctx context.Context,
	serviceIDorName *string) ([]*Plugin, error) {
	if isEmptyString(serviceIDorName) {
		return nil, fmt.Errorf("serviceIDorName cannot be nil")
	}
	return s.listAllByPath(ctx, "/services/"+*serviceIDorName+"/plugins")
}

// ListAllForRoute fetches all Plugins in Kong enabled for a service.
func (s *PluginService) ListAllForRoute(ctx context.Context,
	routeID *string) ([]*Plugin, error) {
	if isEmptyString(routeID) {
		return nil, fmt.Errorf("routeID cannot be nil")
	}
	return s.listAllByPath(ctx, "/routes/"+*routeID+"/plugins")
}
