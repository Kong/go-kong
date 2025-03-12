package kong

// Partial represents a Partial in Kong.
// +k8s:deepcopy-gen=true
type Partial struct {
	ID        *string       `json:"id,omitempty" yaml:"id,omitempty"`
	Name      *string       `json:"name,omitempty" yaml:"name,omitempty"`
	Type      *string       `json:"type,omitempty" yaml:"type,omitempty"`
	Config    Configuration `json:"config,omitempty" yaml:"config,omitempty"`
	CreatedAt *int          `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt *int          `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}

// FriendlyName returns the endpoint key name or ID.
func (p *Partial) FriendlyName() string {
	if p.Name != nil {
		return *p.Name
	}
	if p.ID != nil {
		return *p.ID
	}
	return ""
}
