package kong

// RBACRole represents an RBAC Role in Kong.
// +k8s:deepcopy-gen=true
type RBACRole struct {
	CreatedAt *int    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID        *string `json:"id,omitempty" yaml:"id,omitempty"`
	Name      *string `json:"name,omitempty" yaml:"name,omitempty"`
	Comment   *string `json:"comment,omitempty" yaml:"comment,omitempty"`
	IsDefault *bool   `json:"is_default,omitempty" yaml:"is_default,omitempty"`
}

// RBACPermissionsList is a list of permissions, both endpoint and entity, associated with a Role.
type RBACPermissionsList struct {
	Endpoints map[string]interface{} `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	Entities  map[string]interface{} `json:"entities,omitempty" yaml:"entities,omitempty"`
}

// FriendlyName returns the endpoint key name or ID.
func (r *RBACRole) FriendlyName() string {
	if r.Name != nil {
		return *r.Name
	}
	if r.ID != nil {
		return *r.ID
	}
	return ""
}
