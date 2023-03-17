package kong

// DegraphqlRoute represents a route from an exposed URI to
// a specific GraphQL query on the given Service.
// +k8s:deepcopy-gen=true
type DegraphqlRoute struct {
	ID      *string   `json:"id,omitempty" yaml:"id,omitempty"`
	Service *Service  `json:"service,omitempty" yaml:"service,omitempty"`
	Methods []*string `json:"methods,omitempty" yaml:"methods,omitempty"`
	URI     *string   `json:"uri,omitempty" yaml:"uri,omitempty"`
	Query   *string   `json:"query,omitempty" yaml:"query,omitempty"`
}
