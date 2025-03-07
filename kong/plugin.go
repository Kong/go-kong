package kong

// Plugin represents a Plugin in Kong.
// Read https://docs.konghq.com/gateway/latest/admin-api/#plugin-object
// +k8s:deepcopy-gen=true
type Plugin struct {
	CreatedAt     *int            `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID            *string         `json:"id,omitempty" yaml:"id,omitempty"`
	Name          *string         `json:"name,omitempty" yaml:"name,omitempty"`
	InstanceName  *string         `json:"instance_name,omitempty" yaml:"instance_name,omitempty"`
	Route         *Route          `json:"route,omitempty" yaml:"route,omitempty"`
	Service       *Service        `json:"service,omitempty" yaml:"service,omitempty"`
	Consumer      *Consumer       `json:"consumer,omitempty" yaml:"consumer,omitempty"`
	ConsumerGroup *ConsumerGroup  `json:"consumer_group,omitempty" yaml:"consumer_group,omitempty"`
	Config        Configuration   `json:"config,omitempty" yaml:"config,omitempty"`
	Enabled       *bool           `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	RunOn         *string         `json:"run_on,omitempty" yaml:"run_on,omitempty"`
	Ordering      *PluginOrdering `json:"ordering,omitempty" yaml:"ordering,omitempty"`
	Protocols     []*string       `json:"protocols,omitempty" yaml:"protocols,omitempty"`
	Tags          []*string       `json:"tags,omitempty" yaml:"tags,omitempty"`
	Partials      []*PartialLink  `json:"partials,omitempty" yaml:"partials,omitempty"`
}

// +k8s:deepcopy-gen=true
type PartialLink struct {
	*Partial
	Path *string `json:"path,omitempty" yaml:"path,omitempty"`
}

// PluginOrdering contains before or after instructions for plugin execution order
// +k8s:deepcopy-gen=true
type PluginOrdering struct {
	Before PluginOrderingPhase `json:"before,omitempty"`
	After  PluginOrderingPhase `json:"after,omitempty"`
}

// TODO this explanation is bad, but the organization of the overall struct defies a good explanation at this level
// beyond "they're the things used in PluginOrdering. This is a map from a phase name (which can only be "access"
// in the initial 3.0 release) to a list of plugins that the plugin containing the PluginOrdering should run before
// or after

// PluginOrderingPhase indicates which plugins in a phase should affect the target plugin's order
// +k8s:deepcopy-gen=true
type PluginOrderingPhase map[string][]string

// FriendlyName returns the endpoint key name or ID.
func (p *Plugin) FriendlyName() string {
	if p.Name != nil {
		return *p.Name
	}
	if p.ID != nil {
		return *p.ID
	}
	return ""
}
