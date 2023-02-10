package kong

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RBACEndpointPermission represents an RBAC Endpoint Permission in Kong Enterprise
// +k8s:deepcopy-gen=true
// Note: this type implements a custom JSON marshaler. Review the associated MarshalJSON()
// function if it does not marshal as expected.
type RBACEndpointPermission struct {
	CreatedAt *int      `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Workspace *string   `json:"workspace,omitempty" yaml:"workspace,omitempty"`
	Endpoint  *string   `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	Actions   []*string `json:"actions,omitempty" yaml:"actions,omitempty"`
	Negative  *bool     `json:"negative,omitempty" yaml:"negative,omitempty"`
	Role      *RBACRole `json:"role,omitempty" yaml:"role,omitempty"`
	Comment   *string   `json:"comment,omitempty" yaml:"comment,omitempty"`
}

// MarshalJSON marshals an endpoint permission into a suitable form for the Kong admin API
func (e *RBACEndpointPermission) MarshalJSON() ([]byte, error) {
	type ep struct {
		CreatedAt *int      `json:"created_at,omitempty" yaml:"created_at,omitempty"`
		Workspace *string   `json:"workspace,omitempty" yaml:"workspace,omitempty"`
		Endpoint  *string   `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
		Actions   *string   `json:"actions,omitempty" yaml:"actions,omitempty"`
		Negative  *bool     `json:"negative,omitempty" yaml:"negative,omitempty"`
		Role      *RBACRole `json:"role,omitempty" yaml:"role,omitempty"`
		Comment   *string   `json:"comment,omitempty" yaml:"comment,omitempty"`
	}
	var actions []string
	for _, action := range e.Actions {
		actions = append(actions, *action)
	}
	return json.Marshal(&ep{
		CreatedAt: e.CreatedAt,
		Workspace: e.Workspace,
		Endpoint:  e.Endpoint,
		Actions:   String(strings.Join(actions, ",")),
		Negative:  e.Negative,
		Comment:   e.Comment,
	})
}

// FriendlyName returns a composite Name base on Role , workspace, and endpoint
func (e *RBACEndpointPermission) FriendlyName() string {
	if e.Role != nil && e.Workspace != nil && e.Endpoint != nil {
		return fmt.Sprintf("%s-%s-%s", e.Role.FriendlyName(), *e.Workspace, *e.Endpoint)
	}
	return ""
}
