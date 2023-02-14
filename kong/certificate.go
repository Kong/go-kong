package kong

// Certificate represents a Certificate in Kong.
// Read https://getkong.org/docs/0.14.x/admin-api/#certificate-object
// +k8s:deepcopy-gen=true
type Certificate struct {
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	Cert      *string   `json:"cert,omitempty" yaml:"cert,omitempty"`
	CertAlt   *string   `json:"cert_alt,omitempty" yaml:"cert_alt,omitempty"`
	Key       *string   `json:"key,omitempty" yaml:"key,omitempty"`
	KeyAlt    *string   `json:"key_alt,omitempty" yaml:"key_alt,omitempty"`
	CreatedAt *int64    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	SNIs      []*string `json:"snis,omitempty" yaml:"snis,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// FriendlyName returns the endpoint key name or ID.
func (c *Certificate) FriendlyName() string {
	if c.ID != nil {
		return *c.ID
	}
	if c.Cert != nil {
		return *c.Cert
	}
	return ""
}
