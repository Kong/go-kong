package kong

// GraphqlRateLimitingCostDecoration represents a decoration that
// Kong uses to estimate the cost of a portion of a Graphql request.
// +k8s:deepcopy-gen=true
type GraphqlRateLimitingCostDecoration struct {
	ID           *string   `json:"id,omitempty" yaml:"id,omitempty"`
	TypePath     *string   `json:"type_path,omitempty" yaml:"type_path,omitempty"`
	AddConstant  *float64  `json:"add_constant,omitempty" yaml:"add_constant,omitempty"`
	AddArguments []*string `json:"add_arguments,omitempty" yaml:"add_arguments,omitempty"`
	MulConstant  *float64  `json:"mul_constant,omitempty" yaml:"mul_constant,omitempty"`
	MulArguments []*string `json:"mul_arguments,omitempty" yaml:"mul_arguments,omitempty"`
}
