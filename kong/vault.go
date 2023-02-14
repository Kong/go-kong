package kong

// Vault represents a Vault in Kong.
// +k8s:deepcopy-gen=true
type Vault struct {
	ID          *string       `json:"id,omitempty" yaml:"id,omitempty"`
	Name        *string       `json:"name,omitempty" yaml:"name,omitempty"`
	Description *string       `json:"description,omitempty" yaml:"description,omitempty"`
	Prefix      *string       `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Config      Configuration `json:"config,omitempty" yaml:"config,omitempty"`
	CreatedAt   *int64        `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt   *int64        `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Tags        []*string     `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// FriendlyName returns the endpoint key prefix or ID.
func (s *Vault) FriendlyName() string {
	if s.Prefix != nil {
		return *s.Prefix
	}
	if s.ID != nil {
		return *s.ID
	}
	return ""
}
