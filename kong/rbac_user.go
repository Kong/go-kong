package kong

// RBACUser represents an RBAC user in Kong Enterprise
// +k8s:deepcopy-gen=true
type RBACUser struct {
	CreatedAt      *int    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Comment        *string `json:"comment,omitempty" yaml:"comment,omitempty"`
	ID             *string `json:"id,omitempty" yaml:"id,omitempty"`
	Name           *string `json:"name,omitempty" yaml:"name,omitempty"`
	Enabled        *bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	UserToken      *string `json:"user_token,omitempty" yaml:"user_token,omitempty"`
	UserTokenIdent *string `json:"user_token_ident,omitempty" yaml:"user_token_ident,omitempty"`
}
