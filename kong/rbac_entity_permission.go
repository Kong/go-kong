package kong

import (
	"encoding/json"
	"strings"
)

// RBACEntityPermission represents an RBAC Entity Permission in Kong Enterprise
// +k8s:deepcopy-gen=true
// Note: this type implements a custom JSON marshaler. Review the associated MarshalJSON()
// function if it does not marshal as expected.
type RBACEntityPermission struct {
	CreatedAt  *int      `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	EntityID   *string   `json:"entity_id,omitempty" yaml:"entity_id,omitempty"`
	EntityType *string   `json:"entity_type,omitempty" yaml:"entity_type,omitempty"`
	Actions    []*string `json:"actions,omitempty" yaml:"actions,omitempty"`
	Negative   *bool     `json:"negative,omitempty" yaml:"negative,omitempty"`
	Role       *RBACRole `json:"role,omitempty" yaml:"role,omitempty"`
	Comment    *string   `json:"comment,omitempty" yaml:"comment,omitempty"`
}

// MarshalJSON marshals an endpoint permission into a suitable form for the Kong admin API
func (e *RBACEntityPermission) MarshalJSON() ([]byte, error) {
	type ep struct {
		CreatedAt  *int      `json:"created_at,omitempty" yaml:"created_at,omitempty"`
		EntityID   *string   `json:"entity_id,omitempty" yaml:"entity_id,omitempty"`
		EntityType *string   `json:"entity_type,omitempty" yaml:"entity_type,omitempty"`
		Actions    *string   `json:"actions,omitempty" yaml:"actions,omitempty"`
		Negative   *bool     `json:"negative,omitempty" yaml:"negative,omitempty"`
		Role       *RBACRole `json:"role,omitempty" yaml:"role,omitempty"`
		Comment    *string   `json:"comment,omitempty" yaml:"comment,omitempty"`
	}
	var actions []string
	for _, action := range e.Actions {
		actions = append(actions, *action)
	}
	return json.Marshal(&ep{
		CreatedAt:  e.CreatedAt,
		EntityID:   e.EntityID,
		EntityType: e.EntityType,
		Actions:    String(strings.Join(actions, ",")),
		Negative:   e.Negative,
		Comment:    e.Comment,
	})
}
