package kong

// CustomPluginDefintion represents a custom-plugin definition in Kong.
// +k8s:deepcopy-gen=true
type CustomPluginDefintion struct {
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	Name      *string   `json:"name,omitempty" yaml:"name,omitempty"`
	Handler   *string   `json:"handler,omitempty" yaml:"handler,omitempty"`
	Schema    *string   `json:"schema,omitempty" yaml:"schema,omitempty"`
	CreatedAt *int64    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt *int64    `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// FriendlyName returns the entity name or ID.
func (c *CustomPluginDefintion) FriendlyName() string {
	if c.Name != nil {
		return *c.Name
	}
	if c.ID != nil {
		return *c.ID
	}
	return ""
}
