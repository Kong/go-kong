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
	// CreateForService creates a Plugin in Kong.
	CreateForService(ctx context.Context, serviceIDorName *string, plugin *Plugin) (*Plugin, error)
	// CreateForRoute creates a Plugin in Kong.
	CreateForRoute(ctx context.Context, routeIDorName *string, plugin *Plugin) (*Plugin, error)
	// CreateForConsumerGroup creates a Plugin in Kong.
	CreateForConsumerGroup(ctx context.Context, cgIDorName *string, plugin *Plugin) (*Plugin, error)
	// Get fetches a Plugin in Kong.
	Get(ctx context.Context, usernameOrID *string) (*Plugin, error)
	// Update updates a Plugin in Kong
	Update(ctx context.Context, plugin *Plugin) (*Plugin, error)
	// UpdateForService updates a Plugin in Kong for a service
	UpdateForService(ctx context.Context, serviceIDorName *string, plugin *Plugin) (*Plugin, error)
	// UpdateForRoute updates a Plugin in Kong for a service
	UpdateForRoute(ctx context.Context, routeIDorName *string, plugin *Plugin) (*Plugin, error)
	// Delete deletes a Plugin in Kong
	Delete(ctx context.Context, usernameOrID *string) error
	// DeleteForService deletes a Plugin in Kong
	DeleteForService(ctx context.Context, serviceIDorName *string, pluginID *string) error
	// DeleteForRoute deletes a Plugin in Kong
	DeleteForRoute(ctx context.Context, routeIDorName *string, pluginID *string) error
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
	// ListAllForConsumerGroups fetches all Plugins in Kong enabled for a consumer group.
	ListAllForConsumerGroups(ctx context.Context, cgID *string) ([]*Plugin, error)
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
	pluginName *string,
) (Schema, error) {
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
// Deprecated: Use GetFullSchema instead
func (s *PluginService) GetSchema(ctx context.Context,
	pluginName *string,
) (Schema, error) {
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
	plugin *Plugin,
) (*Plugin, error) {
	queryPath := "/plugins"
	method := "POST"
	if plugin.ID != nil {
		queryPath = queryPath + "/" + *plugin.ID
		method = "PUT"
	}
	return s.sendRequest(ctx, plugin, queryPath, method)
}

// CreateForService creates a Plugin in Kong at Service level.
// If an ID is specified, it will be used to
// create a plugin in Kong, otherwise an ID
// is auto-generated.
func (s *PluginService) CreateForService(ctx context.Context,
	serviceIDorName *string, plugin *Plugin,
) (*Plugin, error) {
	queryPath := "/plugins"
	method := "POST"
	if plugin.ID != nil {
		queryPath = queryPath + "/" + *plugin.ID
		method = "PUT"
	}
	if isEmptyString(serviceIDorName) {
		return nil, fmt.Errorf("serviceIDorName cannot be nil")
	}

	return s.sendRequest(ctx, plugin, fmt.Sprintf("/services/%v"+queryPath, *serviceIDorName), method)
}

// CreateForRoute creates a Plugin in Kong at Route level.
// If an ID is specified, it will be used to
// create a plugin in Kong, otherwise an ID
// is auto-generated.
func (s *PluginService) CreateForRoute(ctx context.Context,
	routeIDorName *string, plugin *Plugin,
) (*Plugin, error) {
	queryPath := "/plugins"
	method := "POST"

	if plugin.ID != nil {
		queryPath = queryPath + "/" + *plugin.ID
		method = "PUT"
	}
	if isEmptyString(routeIDorName) {
		return nil, fmt.Errorf("routeIDorName cannot be nil")
	}

	return s.sendRequest(ctx, plugin, fmt.Sprintf("/routes/%v"+queryPath, *routeIDorName), method)
}

// CreateForConsumerGroup creates a Plugin in Kong at ConsumerGroup level.
// If an ID is specified, it will be used to
// create a plugin in Kong, otherwise an ID
// is auto-generated.
func (s *PluginService) CreateForConsumerGroup(ctx context.Context,
	cgIDorName *string, plugin *Plugin,
) (*Plugin, error) {
	queryPath := "/plugins"
	method := "POST"

	if plugin.ID != nil {
		queryPath = queryPath + "/" + *plugin.ID
		method = "PUT"
	}
	if isEmptyString(cgIDorName) {
		return nil, fmt.Errorf("cgIDorName cannot be nil")
	}

	return s.sendRequest(ctx, plugin, fmt.Sprintf("/consumer_groups/%v"+queryPath, *cgIDorName), method)
}

// Get fetches a Plugin in Kong.
func (s *PluginService) Get(ctx context.Context,
	usernameOrID *string,
) (*Plugin, error) {
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
	plugin *Plugin,
) (*Plugin, error) {
	if isEmptyString(plugin.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/plugins/%v", *plugin.ID)
	return s.sendRequest(ctx, plugin, endpoint, "PATCH")
}

// UpdateForService updates a Plugin in Kong at Service level.
func (s *PluginService) UpdateForService(ctx context.Context,
	serviceIDorName *string, plugin *Plugin,
) (*Plugin, error) {
	if isEmptyString(plugin.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}
	if isEmptyString(serviceIDorName) {
		return nil, fmt.Errorf("serviceIDorName cannot be nil")
	}

	endpoint := fmt.Sprintf("/services/%v/plugins/%v", *serviceIDorName, *plugin.ID)
	return s.sendRequest(ctx, plugin, endpoint, "PATCH")
}

// UpdateForRoute updates a Plugin in Kong at Route level.
func (s *PluginService) UpdateForRoute(ctx context.Context,
	routeIDorName *string, plugin *Plugin,
) (*Plugin, error) {
	if isEmptyString(plugin.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}
	if isEmptyString(routeIDorName) {
		return nil, fmt.Errorf("routeIDorName cannot be nil")
	}

	endpoint := fmt.Sprintf("/routes/%v/plugins/%v", *routeIDorName, *plugin.ID)
	return s.sendRequest(ctx, plugin, endpoint, "PATCH")
}

// Delete deletes a Plugin in Kong
func (s *PluginService) Delete(ctx context.Context,
	pluginID *string,
) error {
	if isEmptyString(pluginID) {
		return fmt.Errorf("pluginID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/plugins/%v", *pluginID)
	_, err := s.sendRequest(ctx, nil, endpoint, "DELETE")
	if err != nil {
		return err
	}
	return err
}

// DeleteForService deletes a Plugin in Kong at Service level.
func (s *PluginService) DeleteForService(ctx context.Context,
	serviceIDorName *string, pluginID *string,
) error {
	if isEmptyString(pluginID) {
		return fmt.Errorf("plugin ID cannot be nil for Delete operation")
	}
	if isEmptyString(serviceIDorName) {
		return fmt.Errorf("serviceIDorName cannot be nil")
	}

	endpoint := fmt.Sprintf("/services/%v/plugins/%v", *serviceIDorName, *pluginID)
	_, err := s.sendRequest(ctx, nil, endpoint, "DELETE")
	if err != nil {
		return err
	}
	return err
}

// DeleteForRoute deletes a Plugin in Kong at Route level.
func (s *PluginService) DeleteForRoute(ctx context.Context,
	routeIDorName *string, pluginID *string,
) error {
	if isEmptyString(pluginID) {
		return fmt.Errorf("plugin ID cannot be nil for Delete operation")
	}
	if isEmptyString(routeIDorName) {
		return fmt.Errorf("routeIDorName cannot be nil")
	}

	endpoint := fmt.Sprintf("/routes/%v/plugins/%v", *routeIDorName, *pluginID)
	_, err := s.sendRequest(ctx, nil, endpoint, "DELETE")
	if err != nil {
		return err
	}
	return nil
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
		if resp == nil {
			return false, "", err
		}
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
	path string, opt *ListOpt,
) ([]*Plugin, *ListOpt, error) {
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
	path string,
) ([]*Plugin, error) {
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
	opt *ListOpt,
) ([]*Plugin, *ListOpt, error) {
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
	consumerIDorName *string,
) ([]*Plugin, error) {
	if isEmptyString(consumerIDorName) {
		return nil, fmt.Errorf("consumerIDorName cannot be nil")
	}
	return s.listAllByPath(ctx, "/consumers/"+*consumerIDorName+"/plugins")
}

// ListAllForService fetches all Plugins in Kong enabled for a service.
func (s *PluginService) ListAllForService(ctx context.Context,
	serviceIDorName *string,
) ([]*Plugin, error) {
	if isEmptyString(serviceIDorName) {
		return nil, fmt.Errorf("serviceIDorName cannot be nil")
	}
	return s.listAllByPath(ctx, "/services/"+*serviceIDorName+"/plugins")
}

// ListAllForRoute fetches all Plugins in Kong enabled for a service.
func (s *PluginService) ListAllForRoute(ctx context.Context,
	routeID *string,
) ([]*Plugin, error) {
	if isEmptyString(routeID) {
		return nil, fmt.Errorf("routeID cannot be nil")
	}
	return s.listAllByPath(ctx, "/routes/"+*routeID+"/plugins")
}

// ListAllForConsumerGroups fetches all Plugins in Kong enabled for a consumer group.
func (s *PluginService) ListAllForConsumerGroups(ctx context.Context,
	cgID *string,
) ([]*Plugin, error) {
	if isEmptyString(cgID) {
		return nil, fmt.Errorf("cgID cannot be nil")
	}
	return s.listAllByPath(ctx, "/consumer_groups/"+*cgID+"/plugins")
}

func (s *PluginService) sendRequest(ctx context.Context, plugin *Plugin, endpoint, method string) (*Plugin, error) {
	var req *http.Request
	var err error
	if method == "DELETE" {
		req, err = s.client.NewRequest(method, endpoint, nil, nil)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = s.client.NewRequest(method, endpoint, nil, plugin)
		if err != nil {
			return nil, err
		}
	}
	var createdPlugin Plugin
	if method == "DELETE" {
		_, err = s.client.Do(ctx, req, nil)
		if err != nil {
			return nil, err
		}
	} else {
		_, err = s.client.Do(ctx, req, &createdPlugin)
		if err != nil {
			return nil, err
		}
	}
	return &createdPlugin, nil
}
