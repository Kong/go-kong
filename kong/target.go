package kong

// Target represents a Target in Kong.
// +k8s:deepcopy-gen=true
type Target struct {
	CreatedAt *float64  `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	Target    *string   `json:"target,omitempty" yaml:"target,omitempty"`
	Upstream  *Upstream `json:"upstream,omitempty" yaml:"upstream,omitempty"`
	Weight    *int      `json:"weight,omitempty" yaml:"weight,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// FriendlyName returns the endpoint key name or ID.
func (t *Target) FriendlyName() string {
	if t.Target != nil {
		return *t.Target
	}
	if t.ID != nil {
		return *t.ID
	}
	return ""
}
