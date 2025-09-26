package kong

// Service represents a Service in Kong.
// Read https://docs.konghq.com/gateway/latest/admin-api/#service-object
// +k8s:deepcopy-gen=true
type Service struct {
	ClientCertificate *Certificate `json:"client_certificate,omitempty" yaml:"client_certificate,omitempty"`
	ConnectTimeout    *int         `json:"connect_timeout,omitempty" yaml:"connect_timeout,omitempty"`
	CreatedAt         *int         `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Enabled           *bool        `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Host              *string      `json:"host,omitempty" yaml:"host,omitempty"`
	ID                *string      `json:"id,omitempty" yaml:"id,omitempty"`
	Name              *string      `json:"name,omitempty" yaml:"name,omitempty"`
	Path              *string      `json:"path,omitempty" yaml:"path,omitempty"`
	Port              *int         `json:"port,omitempty" yaml:"port,omitempty"`
	Protocol          *string      `json:"protocol,omitempty" yaml:"protocol,omitempty"`
	ReadTimeout       *int         `json:"read_timeout,omitempty" yaml:"read_timeout,omitempty"`
	Retries           *int         `json:"retries,omitempty" yaml:"retries,omitempty"`
	UpdatedAt         *int         `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	URL               *string      `json:"url,omitempty" yaml:"url,omitempty"`
	WriteTimeout      *int         `json:"write_timeout,omitempty" yaml:"write_timeout,omitempty"`
	Tags              []*string    `json:"tags,omitempty" yaml:"tags,omitempty"`
	TLSVerify         *bool        `json:"tls_verify,omitempty" yaml:"tls_verify,omitempty"`
	TLSVerifyDepth    *int         `json:"tls_verify_depth,omitempty" yaml:"tls_verify_depth,omitempty"`
	CACertificates    []*string    `json:"ca_certificates,omitempty" yaml:"ca_certificates,omitempty"`
	TLSSANs           *TLSSANs     `json:"tls_sans,omitempty" yaml:"tls_sans,omitempty"`
}

// TLSSANs represents additional Subject Alternative Names that can be matched on Upstream server's TLS certificate.
// +k8s:deepcopy-gen=true
type TLSSANs struct {
	DNSNames []*string `json:"dnsnames,omitempty" yaml:"dnsnames,omitempty"`
	URIs     []*string `json:"uris,omitempty" yaml:"uris,omitempty"`
}

// FriendlyName returns the endpoint key name or ID.
func (s *Service) FriendlyName() string {
	if s.Name != nil {
		return *s.Name
	}
	if s.ID != nil {
		return *s.ID
	}
	return ""
}
