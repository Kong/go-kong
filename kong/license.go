package kong

// License represents a License in Kong.
// Read https://docs.konghq.com/gateway/latest/admin-api/#consumer-object
// +k8s:deepcopy-gen=true
type License struct {
	ID        *string `json:"id,omitempty" yaml:"id,omitempty"`
	Payload   *string `json:"payload,omitempty" yaml:"payload,omitempty"`
	CreatedAt *int64  `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt *int64  `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
}

// FriendlyName returns the endpoint key name or ID.
func (c *License) FriendlyName() string {
	if c.ID != nil {
		return *c.ID
	}
	return ""
}
