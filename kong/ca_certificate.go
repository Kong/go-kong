package kong

// CACertificate represents a CACertificate in Kong.
// +k8s:deepcopy-gen=true
type CACertificate struct {
	ID         *string   `json:"id,omitempty" yaml:"id,omitempty"`
	Cert       *string   `json:"cert,omitempty" yaml:"cert,omitempty"`
	CertDigest *string   `json:"cert_digest,omitempty" yaml:"cert_digest,omitempty"`
	CreatedAt  *int64    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Tags       []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}
