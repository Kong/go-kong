package kong

// Consumer represents a Consumer in Kong.
// Read https://docs.konghq.com/gateway/latest/admin-api/#consumer-object
// +k8s:deepcopy-gen=true
type Consumer struct {
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	CustomID  *string   `json:"custom_id,omitempty" yaml:"custom_id,omitempty"`
	Username  *string   `json:"username,omitempty" yaml:"username,omitempty"`
	CreatedAt *int64    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// FriendlyName returns the endpoint key name or ID.
func (c *Consumer) FriendlyName() string {
	if c.Username != nil {
		return *c.Username
	}
	if c.ID != nil {
		return *c.ID
	}
	return ""
}
