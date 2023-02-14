package kong

// Workspace represents a Workspace in Kong.
type Workspace struct {
	CreatedAt *int                   `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID        *string                `json:"id,omitempty" yaml:"id,omitempty"`
	Name      *string                `json:"name,omitempty" yaml:"name,omitempty"`
	Comment   *string                `json:"comment,omitempty" yaml:"comment,omitempty"`
	Config    map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty"`
	Meta      map[string]interface{} `json:"meta,omitempty" yaml:"meta,omitempty"`
}

// WorkspaceEntity represents a WorkspaceEntity in Kong
// +k8s:deepcopy-gen=true
type WorkspaceEntity struct {
	EntityID         *string `json:"entity_id,omitempty" yaml:"entity_id,omitempty"`
	EntityType       *string `json:"entity_type,omitempty" yaml:"entity_type,omitempty"`
	UniqueFieldName  *string `json:"unique_field_name,omitempty" yaml:"unique_field_name,omitempty"`
	UniqueFieldValue *string `json:"unique_field_value,omitempty" yaml:"unique_field_value,omitempty"`
	WorkspaceID      *string `json:"workspace_id,omitempty" yaml:"workspace_id,omitempty"`
	WorkspaceName    *string `json:"workspace_name,omitempty" yaml:"workspace_name,omitempty"`
}
