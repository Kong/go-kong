package kong

// Developer represents a Developer in Kong.
// +k8s:deepcopy-gen=true
type Developer struct {
	CreatedAt *int      `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	Status    *int      `json:"status,omitempty" yaml:"status,omitempty"`
	Email     *string   `json:"email,omitempty" yaml:"email,omitempty"`
	CustomID  *string   `json:"custom_id,omitempty" yaml:"custom_id,omitempty"`
	UpdatedAt *int      `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Roles     []*string `json:"roles,omitempty" yaml:"roles,omitempty"`
	RbacUser  *RBACUser `json:"rbac_user,omitempty" yaml:"rbac_user,omitempty"`
	Meta      *string   `json:"meta,omitempty" yaml:"meta,omitempty"`
	Password  *string   `json:"password,omitempty" yaml:"password,omitempty"`
}

// DeveloperRole represents a Developer Role in Kong.
// +k8s:deepcopy-gen=true
type DeveloperRole struct {
	Comment   *string `json:"comment,omitempty" yaml:"comment,omitempty"`
	CreatedAt *int    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID        *string `json:"id,omitempty" yaml:"id,omitempty"`
	Name      *string `json:"name,omitempty" yaml:"name,omitempty"`
}
