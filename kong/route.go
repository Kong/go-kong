package kong

// Route represents a Route in Kong.
// Read https://getkong.org/docs/0.13.x/admin-api/#Route-object
// +k8s:deepcopy-gen=true
type Route struct {
	CreatedAt     *int                `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Expression    *string             `json:"expression,omitempty" yaml:"expression,omitempty"`
	Hosts         []*string           `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Headers       map[string][]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	ID            *string             `json:"id,omitempty" yaml:"id,omitempty"`
	Name          *string             `json:"name,omitempty" yaml:"name,omitempty"`
	Methods       []*string           `json:"methods,omitempty" yaml:"methods,omitempty"`
	Paths         []*string           `json:"paths,omitempty" yaml:"paths,omitempty"`
	PathHandling  *string             `json:"path_handling,omitempty" yaml:"path_handling,omitempty"`
	PreserveHost  *bool               `json:"preserve_host,omitempty" yaml:"preserve_host,omitempty"`
	Priority      *int                `json:"priority,omitempty" yaml:"priority,omitempty"`
	Protocols     []*string           `json:"protocols,omitempty" yaml:"protocols,omitempty"`
	RegexPriority *int                `json:"regex_priority,omitempty" yaml:"regex_priority,omitempty"`
	Service       *Service            `json:"service,omitempty" yaml:"service,omitempty"`
	StripPath     *bool               `json:"strip_path,omitempty" yaml:"strip_path,omitempty"`
	UpdatedAt     *int                `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	SNIs          []*string           `json:"snis,omitempty" yaml:"snis,omitempty"`
	Sources       []*CIDRPort         `json:"sources,omitempty" yaml:"sources,omitempty"`
	Destinations  []*CIDRPort         `json:"destinations,omitempty" yaml:"destinations,omitempty"`
	Tags          []*string           `json:"tags,omitempty" yaml:"tags,omitempty"`

	HTTPSRedirectStatusCode *int `json:"https_redirect_status_code,omitempty" yaml:"https_redirect_status_code,omitempty"`

	// Kong buffers requests and responses by default. Buffering is not always
	// desired, for instance if large payloads are being proxied using HTTP 1.1
	// chunked encoding.
	//
	// The request and response route buffering options are enabled by default
	// and allow the user to disable buffering if desired for their use case.
	//
	// SEE ALSO:
	// - https://github.com/Kong/kong/pull/6057
	// - https://docs.konghq.com/2.2.x/admin-api/#route-object
	//
	RequestBuffering  *bool `json:"request_buffering,omitempty" yaml:"request_buffering,omitempty"`
	ResponseBuffering *bool `json:"response_buffering,omitempty" yaml:"response_buffering,omitempty"`
}

// CIDRPort represents a set of CIDR and a port.
// +k8s:deepcopy-gen=true
type CIDRPort struct {
	IP   *string `json:"ip,omitempty" yaml:"ip,omitempty"`
	Port *int    `json:"port,omitempty" yaml:"port,omitempty"`
}

// FriendlyName returns the endpoint key name or ID.
func (r *Route) FriendlyName() string {
	if r.Name != nil {
		return *r.Name
	}
	if r.ID != nil {
		return *r.ID
	}
	return ""
}
