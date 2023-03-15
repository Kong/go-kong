package kong

// SNI represents a SNI in Kong.
// Read https://docs.konghq.com/gateway/latest/admin-api/#sni-object
// +k8s:deepcopy-gen=true
type SNI struct {
	ID          *string      `json:"id,omitempty" yaml:"id,omitempty"`
	Name        *string      `json:"name,omitempty" yaml:"name,omitempty"`
	CreatedAt   *int64       `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Certificate *Certificate `json:"certificate,omitempty" yaml:"certificate,omitempty"`
	Tags        []*string    `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// FriendlyName returns the endpoint key name or ID.
func (s *SNI) FriendlyName() string {
	if s.Name != nil {
		return *s.Name
	}
	if s.ID != nil {
		return *s.ID
	}
	return ""
}
