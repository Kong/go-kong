package kong

import "encoding/json"

// FilterChain represents a FilterChain in Kong.
// Read https://docs.konghq.com/gateway/latest/admin-api/#filter-chain-object
// +k8s:deepcopy-gen=true
type FilterChain struct {
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	Name      *string   `json:"name,omitempty" yaml:"name,omitempty"`
	Enabled   *bool     `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Route     *Route    `json:"route,omitempty" yaml:"route,omitempty"`
	Service   *Service  `json:"service,omitempty" yaml:"service,omitempty"`
	Filters   []*Filter `json:"filters,omitempty" yaml:"filters,omitempty"`
	CreatedAt *int      `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt *int      `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// Filter contains information about each filter in the chain
// +k8s:deepcopy-gen=true
type Filter struct {
	Name    *string          `json:"name,omitempty" yaml:"name,omitempty"`
	Config  *json.RawMessage `json:"config,omitempty" yaml:"config,omitempty"`
	Enabled *bool            `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

// FriendlyName returns the endpoint key name or ID.
func (f *FilterChain) FriendlyName() string {
	if f.Name != nil {
		return *f.Name
	}
	if f.ID != nil {
		return *f.ID
	}
	return ""
}
