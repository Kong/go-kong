package kong

// AIModel represents an AI Model in Kong.
// +k8s:deepcopy-gen=true
type AIModel struct {
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	Name      *string   `json:"name,omitempty" yaml:"name,omitempty"`
	Alias     *string   `json:"alias,omitempty" yaml:"alias,omitempty"`
	CreatedAt *int64    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt *int64    `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// FriendlyName returns the endpoint key name or ID.
func (a *AIModel) FriendlyName() string {
	if a.Name != nil {
		return *a.Name
	}
	if a.ID != nil {
		return *a.ID
	}
	return ""
}
