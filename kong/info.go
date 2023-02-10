package kong

// Info represents the information concerning Kong.
type Info struct {
	Version       string                `json:"version,omitempty" yaml:"version,omitempty"`
	Configuration *RuntimeConfiguration `json:"configuration,omitempty" yaml:"configuration,omitempty"`
}

// RuntimeConfiguration represents the runtime configuration of Kong.
type RuntimeConfiguration struct {
	Database string `json:"database,omitempty" yaml:"database,omitempty"`
	Portal   bool   `json:"portal,omitempty" yaml:"portal,omitempty"`
	RBAC     string `json:"rbac,omitempty" yaml:"rbac,omitempty"`
}
