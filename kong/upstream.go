package kong

// Upstream represents an Upstream in Kong.
// +k8s:deepcopy-gen=true
type Upstream struct {
	ID                     *string      `json:"id,omitempty" yaml:"id,omitempty"`
	Name                   *string      `json:"name,omitempty" yaml:"name,omitempty"`
	HostHeader             *string      `json:"host_header,omitempty" yaml:"host_header,omitempty"`
	ClientCertificate      *Certificate `json:"client_certificate,omitempty" yaml:"client_certificate,omitempty"`
	Algorithm              *string      `json:"algorithm,omitempty" yaml:"algorithm,omitempty"`
	Slots                  *int         `json:"slots,omitempty" yaml:"slots,omitempty"`
	Healthchecks           *Healthcheck `json:"healthchecks,omitempty" yaml:"healthchecks,omitempty"`
	CreatedAt              *int64       `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	HashOn                 *string      `json:"hash_on,omitempty" yaml:"hash_on,omitempty"`
	HashFallback           *string      `json:"hash_fallback,omitempty" yaml:"hash_fallback,omitempty"`
	HashOnHeader           *string      `json:"hash_on_header,omitempty" yaml:"hash_on_header,omitempty"`
	HashFallbackHeader     *string      `json:"hash_fallback_header,omitempty" yaml:"hash_fallback_header,omitempty"`
	HashOnCookie           *string      `json:"hash_on_cookie,omitempty" yaml:"hash_on_cookie,omitempty"`
	HashOnCookiePath       *string      `json:"hash_on_cookie_path,omitempty" yaml:"hash_on_cookie_path,omitempty"`
	HashOnQueryArg         *string      `json:"hash_on_query_arg,omitempty" yaml:"hash_on_query_arg,omitempty"`
	HashFallbackQueryArg   *string      `json:"hash_fallback_query_arg,omitempty" yaml:"hash_fallback_query_arg,omitempty"` //nolint:lll
	HashOnURICapture       *string      `json:"hash_on_uri_capture,omitempty" yaml:"hash_on_uri_capture,omitempty"`
	HashFallbackURICapture *string      `json:"hash_fallback_uri_capture,omitempty" yaml:"hash_fallback_uri_capture,omitempty"` //nolint:lll
	Tags                   []*string    `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// Healthy configures thresholds and HTTP status codes
// to mark targets healthy for an upstream.
// +k8s:deepcopy-gen=true
type Healthy struct {
	HTTPStatuses []int `json:"http_statuses,omitempty" yaml:"http_statuses,omitempty"`
	// +kubebuilder:validation:Minimum=0
	Interval *int `json:"interval,omitempty" yaml:"interval,omitempty"`
	// +kubebuilder:validation:Minimum=0
	Successes *int `json:"successes,omitempty" yaml:"successes,omitempty"`
}

// Unhealthy configures thresholds and HTTP status codes
// to mark targets unhealthy.
// +k8s:deepcopy-gen=true
type Unhealthy struct {
	// +kubebuilder:validation:Minimum=0
	HTTPFailures *int  `json:"http_failures,omitempty" yaml:"http_failures,omitempty"`
	HTTPStatuses []int `json:"http_statuses,omitempty" yaml:"http_statuses,omitempty"`
	// +kubebuilder:validation:Minimum=0
	TCPFailures *int `json:"tcp_failures,omitempty" yaml:"tcp_failures,omitempty"`
	// +kubebuilder:validation:Minimum=0
	Timeouts *int `json:"timeouts,omitempty" yaml:"timeouts,omitempty"`
	// +kubebuilder:validation:Minimum=0
	Interval *int `json:"interval,omitempty" yaml:"interval,omitempty"`
}

// ActiveHealthcheck configures active health check probing.
// +k8s:deepcopy-gen=true
type ActiveHealthcheck struct {
	// +kubebuilder:validation:Minimum=1
	Concurrency *int     `json:"concurrency,omitempty" yaml:"concurrency,omitempty"`
	Healthy     *Healthy `json:"healthy,omitempty" yaml:"healthy,omitempty"`
	// +kubebuilder:validation:Pattern=^/.*$
	HTTPPath               *string `json:"http_path,omitempty" yaml:"http_path,omitempty"`
	HTTPSSni               *string `json:"https_sni,omitempty" yaml:"https_sni,omitempty"`
	HTTPSVerifyCertificate *bool   `json:"https_verify_certificate,omitempty" yaml:"https_verify_certificate,omitempty"`
	Type                   *string `json:"type,omitempty" yaml:"type,omitempty"`
	// +kubebuilder:validation:Minimum=0
	Timeout   *int       `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Unhealthy *Unhealthy `json:"unhealthy,omitempty" yaml:"unhealthy,omitempty"`
}

// PassiveHealthcheck configures passive checks around
// passive health checks.
// +k8s:deepcopy-gen=true
type PassiveHealthcheck struct {
	Healthy   *Healthy   `json:"healthy,omitempty" yaml:"healthy,omitempty"`
	Type      *string    `json:"type,omitempty" yaml:"type,omitempty"`
	Unhealthy *Unhealthy `json:"unhealthy,omitempty" yaml:"unhealthy,omitempty"`
}

// Healthcheck represents a health-check config of an upstream
// in Kong.
// +k8s:deepcopy-gen=true
type Healthcheck struct {
	Active    *ActiveHealthcheck  `json:"active,omitempty" yaml:"active,omitempty"`
	Passive   *PassiveHealthcheck `json:"passive,omitempty" yaml:"passive,omitempty"`
	Threshold *float64            `json:"threshold,omitempty" yaml:"threshold,omitempty"`
}

// HealthDataAddress represents the health data address of a target
// +k8s:deepcopy-gen=true
type HealthDataAddress struct {
	Port   *int    `json:"port,omitempty" yaml:"port,omitempty"`
	IP     *string `json:"ip,omitempty" yaml:"ip,omitempty"`
	Health *string `json:"health,omitempty" yaml:"data,omitempty"`
	Weight *int    `json:"weight,omitempty" yaml:"weight,omitempty"`
}

// HealthDataWeight represents the health data weight of a target
// +k8s:deepcopy-gen=true
type HealthDataWeight struct {
	Total       *int `json:"total,omitempty" yaml:"total,omitempty"`
	Available   *int `json:"available,omitempty" yaml:"available,omitempty"`
	Unavailable *int `json:"unavailable,omitempty" yaml:"unavailable,omitempty"`
}

// HealthData represents the health data of a target
// +k8s:deepcopy-gen=true
type HealthData struct {
	Host       *string              `json:"host,omitempty" yaml:"host,omitempty"`
	Port       *int                 `json:"port,omitempty" yaml:"port,omitempty"`
	NodeWeight *int                 `json:"nodeWeight,omitempty" yaml:"nodeWeight,omitempty"`
	Weight     *HealthDataWeight    `json:"weight,omitempty" yaml:"weight,omitempty"`
	Addresses  []*HealthDataAddress `json:"addresses,omitempty" yaml:"addresses,omitempty"`
	DNS        *string              `json:"dns,omitempty" yaml:"dns,omitempty"`
}

// UpstreamNodeHealth represents the node health of a upstream
// +k8s:deepcopy-gen=true
type UpstreamNodeHealth struct {
	ID        *string     `json:"id,omitempty" yaml:"id,omitempty"`
	CreatedAt *float64    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Data      *HealthData `json:"data,omitempty" yaml:"data,omitempty"`
	Health    *string     `json:"health,omitempty" yaml:"data,omitempty"`
	Target    *string     `json:"target,omitempty" yaml:"target,omitempty"`
	Upstream  *Upstream   `json:"upstream,omitempty" yaml:"upstream,omitempty"`
	Weight    *int        `json:"weight,omitempty" yaml:"weight,omitempty"`
	Tags      []*string   `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// FriendlyName returns the endpoint key name or ID.
func (u *Upstream) FriendlyName() string {
	if u.Name != nil {
		return *u.Name
	}
	if u.ID != nil {
		return *u.ID
	}
	return ""
}
