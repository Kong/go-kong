package kong

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Service represents a Service in Kong.
// Read https://getkong.org/docs/0.13.x/admin-api/#Service-object
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
}

// CIDRPort represents a set of CIDR and a port.
// +k8s:deepcopy-gen=true
type CIDRPort struct {
	IP   *string `json:"ip,omitempty" yaml:"ip,omitempty"`
	Port *int    `json:"port,omitempty" yaml:"port,omitempty"`
}

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

// Consumer represents a Consumer in Kong.
// Read https://getkong.org/docs/0.13.x/admin-api/#consumer-object
// +k8s:deepcopy-gen=true
type Consumer struct {
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	CustomID  *string   `json:"custom_id,omitempty" yaml:"custom_id,omitempty"`
	Username  *string   `json:"username,omitempty" yaml:"username,omitempty"`
	CreatedAt *int64    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// ConsumerGroupObject represents a ConsumerGroup in Kong.
// +k8s:deepcopy-gen=true
type ConsumerGroupObject struct {
	ConsumerGroup *ConsumerGroup         `json:"consumer_group,omitempty" yaml:"consumer_group,omitempty"`
	Consumers     []*Consumer            `json:"consumers,omitempty" yaml:"consumers,omitempty"`
	Plugins       []*ConsumerGroupPlugin `json:"plugins,omitempty" yaml:"plugins,omitempty"`
}

// ConsumerGroup represents a ConsumerGroup in Kong.
// +k8s:deepcopy-gen=true
type ConsumerGroup struct {
	ID        *string `json:"id,omitempty" yaml:"id,omitempty"`
	Name      *string `json:"name,omitempty" yaml:"name,omitempty"`
	CreatedAt *int64  `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

// ConsumerGroupConsumer represents a ConsumerGroupConsumer in Kong.
// +k8s:deepcopy-gen=true
type ConsumerGroupConsumer struct {
	Consumer      *Consumer      `json:"consumer,omitempty" yaml:"consumer,omitempty"`
	ConsumerGroup *ConsumerGroup `json:"consumer_group,omitempty" yaml:"consumer_group,omitempty"`
	CreatedAt     *int64         `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

// ConsumerGroupRLA represents a ConsumerGroupRLA in Kong.
// +k8s:deepcopy-gen=true
type ConsumerGroupRLA struct {
	ConsumerGroup *string       `json:"consumer_group,omitempty" yaml:"consumer_group,omitempty"`
	Config        Configuration `json:"config,omitempty" yaml:"config,omitempty"`
	Plugin        *string       `json:"plugin,omitempty" yaml:"plugin,omitempty"`
}

// ConsumerGroupPlugin represents a ConsumerGroupPlugin in Kong.
// +k8s:deepcopy-gen=true
type ConsumerGroupPlugin struct {
	ID            *string        `json:"id,omitempty" yaml:"id,omitempty"`
	Name          *string        `json:"name,omitempty" yaml:"name,omitempty"`
	CreatedAt     *int64         `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Config        Configuration  `json:"config,omitempty" yaml:"config,omitempty"`
	ConsumerGroup *ConsumerGroup `json:"consumer_group,omitempty" yaml:"consumer_group,omitempty"`
}

// FriendlyName returns the endpoint key name or ID.
func (s *ConsumerGroup) FriendlyName() string {
	if s.Name != nil {
		return *s.Name
	}
	if s.ID != nil {
		return *s.ID
	}
	return ""
}

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

// SNI represents a SNI in Kong.
// Read https://getkong.org/docs/0.14.x/admin-api/#sni-object
// +k8s:deepcopy-gen=true
type SNI struct {
	ID          *string      `json:"id,omitempty" yaml:"id,omitempty"`
	Name        *string      `json:"name,omitempty" yaml:"name,omitempty"`
	CreatedAt   *int64       `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Certificate *Certificate `json:"certificate,omitempty" yaml:"certificate,omitempty"`
	Tags        []*string    `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// Healthy configures thresholds and HTTP status codes
// to mark targets healthy for an upstream.
// +k8s:deepcopy-gen=true
type Healthy struct {
	HTTPStatuses []int `json:"http_statuses,omitempty" yaml:"http_statuses,omitempty"`
	//+kubebuilder:validation:Minimum=0
	Interval *int `json:"interval,omitempty" yaml:"interval,omitempty"`
	//+kubebuilder:validation:Minimum=0
	Successes *int `json:"successes,omitempty" yaml:"successes,omitempty"`
}

// Unhealthy configures thresholds and HTTP status codes
// to mark targets unhealthy.
// +k8s:deepcopy-gen=true
type Unhealthy struct {
	//+kubebuilder:validation:Minimum=0
	HTTPFailures *int  `json:"http_failures,omitempty" yaml:"http_failures,omitempty"`
	HTTPStatuses []int `json:"http_statuses,omitempty" yaml:"http_statuses,omitempty"`
	//+kubebuilder:validation:Minimum=0
	TCPFailures *int `json:"tcp_failures,omitempty" yaml:"tcp_failures,omitempty"`
	//+kubebuilder:validation:Minimum=0
	Timeouts *int `json:"timeouts,omitempty" yaml:"timeouts,omitempty"`
	//+kubebuilder:validation:Minimum=0
	Interval *int `json:"interval,omitempty" yaml:"interval,omitempty"`
}

// ActiveHealthcheck configures active health check probing.
// +k8s:deepcopy-gen=true
type ActiveHealthcheck struct {
	//+kubebuilder:validation:Minimum=1
	Concurrency *int     `json:"concurrency,omitempty" yaml:"concurrency,omitempty"`
	Healthy     *Healthy `json:"healthy,omitempty" yaml:"healthy,omitempty"`
	//+kubebuilder:validation:Pattern=^/.*$
	HTTPPath               *string `json:"http_path,omitempty" yaml:"http_path,omitempty"`
	HTTPSSni               *string `json:"https_sni,omitempty" yaml:"https_sni,omitempty"`
	HTTPSVerifyCertificate *bool   `json:"https_verify_certificate,omitempty" yaml:"https_verify_certificate,omitempty"`
	Type                   *string `json:"type,omitempty" yaml:"type,omitempty"`
	//+kubebuilder:validation:Minimum=0
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

// Target represents a Target in Kong.
// +k8s:deepcopy-gen=true
type Target struct {
	CreatedAt *float64  `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	Target    *string   `json:"target,omitempty" yaml:"target,omitempty"`
	Upstream  *Upstream `json:"upstream,omitempty" yaml:"upstream,omitempty"`
	Weight    *int      `json:"weight,omitempty" yaml:"weight,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// Configuration represents a config of a plugin in Kong.
type Configuration map[string]interface{}

// KongPluginOrdering contains before or after instructions for plugin execution order
// +k8s:deepcopy-gen=true
type PluginOrdering struct {
	Before PluginOrderingPhase `json:"before,omitempty"`
	After  PluginOrderingPhase `json:"after,omitempty"`
}

// TODO this explanation is bad, but the organization of the overall struct defies a good explanation at this level
// beyond "they're the things used in PluginOrdering. This is a map from a phase name (which can only be "access"
// in the initial 3.0 release) to a list of plugins that the plugin containing the PluginOrdering should run before
// or after

// PluginOrderingPhase indicates which plugins in a phase should affect the target plugin's order
// +k8s:deepcopy-gen=true
type PluginOrderingPhase map[string][]string

// DeepCopyInto copies the receiver, writing into out. in must be non-nil.
func (in Configuration) DeepCopyInto(out *Configuration) {
	// Resorting to JSON since interface{} cannot be DeepCopied easily.
	// This could be replaced using reflection-fu.
	// XXX Ignoring errors
	b, _ := json.Marshal(&in)
	_ = json.Unmarshal(b, out)
}

// DeepCopy copies the receiver, creating a new Configuration.
func (in Configuration) DeepCopy() Configuration {
	if in == nil {
		return nil
	}
	out := new(Configuration)
	in.DeepCopyInto(out)
	return *out
}

// Plugin represents a Plugin in Kong.
// Read https://getkong.org/docs/0.13.x/admin-api/#Plugin-object
// +k8s:deepcopy-gen=true
type Plugin struct {
	CreatedAt *int            `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID        *string         `json:"id,omitempty" yaml:"id,omitempty"`
	Name      *string         `json:"name,omitempty" yaml:"name,omitempty"`
	Route     *Route          `json:"route,omitempty" yaml:"route,omitempty"`
	Service   *Service        `json:"service,omitempty" yaml:"service,omitempty"`
	Consumer  *Consumer       `json:"consumer,omitempty" yaml:"consumer,omitempty"`
	Config    Configuration   `json:"config,omitempty" yaml:"config,omitempty"`
	Enabled   *bool           `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	RunOn     *string         `json:"run_on,omitempty" yaml:"run_on,omitempty"`
	Ordering  *PluginOrdering `json:"ordering,omitempty" yaml:"ordering,omitempty"`
	Protocols []*string       `json:"protocols,omitempty" yaml:"protocols,omitempty"`
	Tags      []*string       `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// Enterprise Entities

// Workspace represents a Workspace in Kong.
type Workspace struct {
	CreatedAt *int                   `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID        *string                `json:"id,omitempty" yaml:"id,omitempty"`
	Name      *string                `json:"name,omitempty" yaml:"name,omitempty"`
	Comment   *string                `json:"comment,omitempty" yaml:"comment,omitempty"`
	Config    map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty"`
	Meta      map[string]interface{} `json:"meta,omitempty" yaml:"meta,omitempty"`
}

// Admin represents an Admin in Kong.
// +k8s:deepcopy-gen=true
type Admin struct {
	CreatedAt        *int    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID               *string `json:"id,omitempty" yaml:"id,omitempty"`
	Email            *string `json:"email,omitempty" yaml:"email,omitempty"`
	Username         *string `json:"username,omitempty" yaml:"username,omitempty"`
	Password         *string `json:"password,omitempty" yaml:"password,omitempty"`
	CustomID         *string `json:"custom_id,omitempty" yaml:"custom_id,omitempty"`
	RBACTokenEnabled *bool   `json:"rbac_token_enabled,omitempty" yaml:"rbac_token_enabled,omitempty"`
	Status           *int    `json:"status,omitempty" yaml:"status,omitempty"`
	Token            *string `json:"token,omitempty" yaml:"token,omitempty"`
}

// RBACUser represents an RBAC user in Kong Enterprise
// +k8s:deepcopy-gen=true
type RBACUser struct {
	CreatedAt      *int    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Comment        *string `json:"comment,omitempty" yaml:"comment,omitempty"`
	ID             *string `json:"id,omitempty" yaml:"id,omitempty"`
	Name           *string `json:"name,omitempty" yaml:"name,omitempty"`
	Enabled        *bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	UserToken      *string `json:"user_token,omitempty" yaml:"user_token,omitempty"`
	UserTokenIdent *string `json:"user_token_ident,omitempty" yaml:"user_token_ident,omitempty"`
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

// RBACRole represents an RBAC Role in Kong.
// +k8s:deepcopy-gen=true
type RBACRole struct {
	CreatedAt *int    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID        *string `json:"id,omitempty" yaml:"id,omitempty"`
	Name      *string `json:"name,omitempty" yaml:"name,omitempty"`
	Comment   *string `json:"comment,omitempty" yaml:"comment,omitempty"`
	IsDefault *bool   `json:"is_default,omitempty" yaml:"is_default,omitempty"`
}

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

// RBACEntityPermission represents an RBAC Entity Permission in Kong Enterprise
// +k8s:deepcopy-gen=true
// Note: this type implements a custom JSON marshaler. Review the associated MarshalJSON()
// function if it does not marshal as expected.
type RBACEntityPermission struct {
	CreatedAt  *int      `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	EntityID   *string   `json:"entity_id,omitempty" yaml:"entity_id,omitempty"`
	EntityType *string   `json:"entity_type,omitempty" yaml:"entity_type,omitempty"`
	Actions    []*string `json:"actions,omitempty" yaml:"actions,omitempty"`
	Negative   *bool     `json:"negative,omitempty" yaml:"negative,omitempty"`
	Role       *RBACRole `json:"role,omitempty" yaml:"role,omitempty"`
	Comment    *string   `json:"comment,omitempty" yaml:"comment,omitempty"`
}

// MarshalJSON marshals an endpoint permission into a suitable form for the Kong admin API
func (e *RBACEntityPermission) MarshalJSON() ([]byte, error) {
	type ep struct {
		CreatedAt  *int      `json:"created_at,omitempty" yaml:"created_at,omitempty"`
		EntityID   *string   `json:"entity_id,omitempty" yaml:"entity_id,omitempty"`
		EntityType *string   `json:"entity_type,omitempty" yaml:"entity_type,omitempty"`
		Actions    *string   `json:"actions,omitempty" yaml:"actions,omitempty"`
		Negative   *bool     `json:"negative,omitempty" yaml:"negative,omitempty"`
		Role       *RBACRole `json:"role,omitempty" yaml:"role,omitempty"`
		Comment    *string   `json:"comment,omitempty" yaml:"comment,omitempty"`
	}
	var actions []string
	for _, action := range e.Actions {
		actions = append(actions, *action)
	}
	return json.Marshal(&ep{
		CreatedAt:  e.CreatedAt,
		EntityID:   e.EntityID,
		EntityType: e.EntityType,
		Actions:    String(strings.Join(actions, ",")),
		Negative:   e.Negative,
		Comment:    e.Comment,
	})
}

// RBACPermissionsList is a list of permissions, both endpoint and entity, associated with a Role.
type RBACPermissionsList struct {
	Endpoints map[string]interface{} `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	Entities  map[string]interface{} `json:"entities,omitempty" yaml:"entities,omitempty"`
}

// Developer represents a Developer in Kong.
// +k8s:deepcopy-gen=true
type Developer struct {
	CreatedAt *int      `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	Status    *int      `json:"status,omitempty" yaml:"status,omitempty"`
	Email     *string   `json:"email,omitempty" yaml:"email,omitempty"`
	CustomID  *string   `json:"custom_id,omitempty" yaml:"custom_id,omitempty"`
	UpdatedAt *int      `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Roles     []*string `json:"roles,omitempty" yaml:"roles,omitempty"`
	RbacUser  *RBACUser `json:"rbac_user,omitempty" yaml:"rbac_user,omitempty"`
	Meta      *string   `json:"meta,omitempty" yaml:"meta,omitempty"`
	Password  *string   `json:"password,omitempty" yaml:"password,omitempty"`
}

// DeveloperRole represents a Developer Role in Kong.
// +k8s:deepcopy-gen=true
type DeveloperRole struct {
	Comment   *string `json:"comment,omitempty" yaml:"comment,omitempty"`
	CreatedAt *int    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ID        *string `json:"id,omitempty" yaml:"id,omitempty"`
	Name      *string `json:"name,omitempty" yaml:"name,omitempty"`
}

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

// Vault represents a Vault in Kong.
// +k8s:deepcopy-gen=true
type Vault struct {
	ID          *string       `json:"id,omitempty" yaml:"id,omitempty"`
	Name        *string       `json:"name,omitempty" yaml:"name,omitempty"`
	Description *string       `json:"description,omitempty" yaml:"description,omitempty"`
	Prefix      *string       `json:"prefix,omitempty" yaml:"prefix,omitempty"`
	Config      Configuration `json:"config,omitempty" yaml:"config,omitempty"`
	CreatedAt   *int64        `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt   *int64        `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Tags        []*string     `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// Key represents a JWK or PEM key in Kong.
// +k8s:deepcopy-gen=true
type Key struct {
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	CreatedAt *int64    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt *int64    `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Set       *KeySet   `json:"set,omitempty" yaml:"set,omitempty"`
	Name      *string   `json:"name,omitempty" yaml:"name,omitempty"`
	KID       *string   `json:"kid,omitempty" yaml:"kid,omitempty"`
	JWK       *string   `json:"jwk,omitempty" yaml:"jwk,omitempty"`
	PEM       *PEM      `json:"pem,omitempty" yaml:"pem,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// PEM represents a PEM formatted key in Kong.
// +k8s:deepcopy-gen=true
type PEM struct {
	PublicKey  *string `json:"public_key,omitempty" yaml:"public_key,omitempty"`
	PrivateKey *string `json:"private_key,omitempty" yaml:"private_key,omitempty"`
}

// KeySet represents a set of keys in Kong.
// +k8s:deepcopy-gen=true
type KeySet struct {
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	CreatedAt *int64    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt *int64    `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Name      *string   `json:"name,omitempty" yaml:"name,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// FriendlyName returns the endpoint key prefix or ID.
func (s *Vault) FriendlyName() string {
	if s.Prefix != nil {
		return *s.Prefix
	}
	if s.ID != nil {
		return *s.ID
	}
	return ""
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

// FriendlyName returns the endpoint key name or ID.
func (t *Target) FriendlyName() string {
	if t.Target != nil {
		return *t.Target
	}
	if t.ID != nil {
		return *t.ID
	}
	return ""
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

// FriendlyName returns the endpoint key name or ID.
func (s *SNI) FriendlyName() string {
	if s.Name != nil {
		return *s.Name
	}
	if s.ID != nil {
		return *s.ID
	}
	return ""
}

// FriendlyName returns the endpoint key name or ID.
func (p *Plugin) FriendlyName() string {
	if p.Name != nil {
		return *p.Name
	}
	if p.ID != nil {
		return *p.ID
	}
	return ""
}

// FriendlyName returns the endpoint key name or ID.
func (c *Consumer) FriendlyName() string {
	if c.Username != nil {
		return *c.Username
	}
	if c.ID != nil {
		return *c.ID
	}
	return ""
}

// FriendlyName returns the endpoint key name or ID.
func (c *CACertificate) FriendlyName() string {
	if c.ID != nil {
		return *c.ID
	}
	if c.Cert != nil {
		return *c.Cert
	}
	return ""
}

// FriendlyName returns the endpoint key name or ID.
func (r *RBACRole) FriendlyName() string {
	if r.Name != nil {
		return *r.Name
	}
	if r.ID != nil {
		return *r.ID
	}
	return ""
}

// FriendlyName returns a composite Name base on Role , workspace, and endpoint
func (e *RBACEndpointPermission) FriendlyName() string {
	if e.Role != nil && e.Workspace != nil && e.Endpoint != nil {
		return fmt.Sprintf("%s-%s-%s", e.Role.FriendlyName(), *e.Workspace, *e.Endpoint)
	}
	return ""
}
