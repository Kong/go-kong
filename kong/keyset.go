package kong

// KeySet represents a set of keys in Kong.
// +k8s:deepcopy-gen=true
type KeySet struct {
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	CreatedAt *int64    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt *int64    `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Name      *string   `json:"name,omitempty" yaml:"name,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}
