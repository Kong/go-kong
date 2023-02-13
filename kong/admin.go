package kong

// Admin represents an Admin in Kong.
// +k8s:deepcopy-gen=true
type Admin struct {
	CreatedAt        *int    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID               *string `json:"id,omitempty" yaml:"id,omitempty"`
	Email            *string `json:"email,omitempty" yaml:"email,omitempty"`
	Username         *string `json:"username,omitempty" yaml:"username,omitempty"`
	Password         *string `json:"password,omitempty" yaml:"password,omitempty"`
	CustomID         *string `json:"custom_id,omitempty" yaml:"custom_id,omitempty"`
	RBACTokenEnabled *bool   `json:"rbac_token_enabled,omitempty" yaml:"rbac_token_enabled,omitempty"`
	Status           *int    `json:"status,omitempty" yaml:"status,omitempty"`
	Token            *string `json:"token,omitempty" yaml:"token,omitempty"`
}
