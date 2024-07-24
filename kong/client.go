package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/kong/go-kong/kong/custom"
)

const (
	// defaultBaseURL is the endpoint for admin API.
	// ref: https://docs.konghq.com/gateway/latest/production/networking/default-ports/	
	defaultBaseURL = "http://localhost:8001"
	// defaultStatusURL is the endpoint for status API
	// By default, the Status API listens on localhost.
	// If you need to request it from elsewhere,
	// please modify the `KONG_STATUS_LISTEN` environment variable of Gateway.
	defaultStatusURL = "http://localhost:8007"
	// DefaultTimeout is the timeout used for network connections and requests
	// including TCP, TLS and HTTP layers.
	DefaultTimeout = 60 * time.Second
)

var pageSize = 1000

type service struct {
	client *Client
}

// Doer is the function signature for a Client request dispatcher.
type Doer func(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error)

var defaultCtx = context.Background()

// Client talks to the Admin API or control plane of a
// Kong cluster
type Client struct {
	client                  *http.Client
	baseRootURL             string
	statusURL               string
	workspace               string       // Do not access directly. Use Workspace()/SetWorkspace().
	UserAgent               string       // User-Agent for the client.
	workspaceLock           sync.RWMutex // Synchronizes access to workspace.
	common                  service
	ConsumerGroupConsumers  AbstractConsumerGroupConsumerService
	ConsumerGroups          AbstractConsumerGroupService
	Consumers               AbstractConsumerService
	Developers              AbstractDeveloperService
	DeveloperRoles          AbstractDeveloperRoleService
	Services                AbstractSvcService
	Routes                  AbstractRouteService
	CACertificates          AbstractCACertificateService
	Certificates            AbstractCertificateService
	Plugins                 AbstractPluginService
	SNIs                    AbstractSNIService
	Upstreams               AbstractUpstreamService
	UpstreamNodeHealth      AbstractUpstreamNodeHealthService
	Targets                 AbstractTargetService
	Workspaces              AbstractWorkspaceService
	Admins                  AbstractAdminService
	RBACUsers               AbstractRBACUserService
	RBACRoles               AbstractRBACRoleService
	RBACEndpointPermissions AbstractRBACEndpointPermissionService
	RBACEntityPermissions   AbstractRBACEntityPermissionService
	Vaults                  AbstractVaultService
	Keys                    AbstractKeyService
	KeySets                 AbstractKeySetService
	Licenses                AbstractLicenseService
	FilterChains            AbstractFilterChainService

	credentials       abstractCredentialService
	KeyAuths          AbstractKeyAuthService
	BasicAuths        AbstractBasicAuthService
	HMACAuths         AbstractHMACAuthService
	JWTAuths          AbstractJWTAuthService
	MTLSAuths         AbstractMTLSAuthService
	ACLs              AbstractACLService
	Oauth2Credentials AbstractOauth2Service
	Tags              AbstractTagService
	Info              AbstractInfoService

	GraphqlRateLimitingCostDecorations AbstractGraphqlRateLimitingCostDecorationService
	DegraphqlRoutes                    AbstractDegraphqlRouteService

	Schemas AbstractSchemaService

	logger         io.Writer
	debug          bool
	CustomEntities AbstractCustomEntityService
	doer           Doer

	custom.Registry
}

// Status respresents current status of a Kong node.
type Status struct {
	Database struct {
		Reachable bool `json:"reachable"`
	} `json:"database"`
	Server struct {
		ConnectionsAccepted int `json:"connections_accepted"`
		ConnectionsActive   int `json:"connections_active"`
		ConnectionsHandled  int `json:"connections_handled"`
		ConnectionsReading  int `json:"connections_reading"`
		ConnectionsWaiting  int `json:"connections_waiting"`
		ConnectionsWriting  int `json:"connections_writing"`
		TotalRequests       int `json:"total_requests"`
	} `json:"server"`
	ConfigurationHash string `json:"configuration_hash,omitempty" yaml:"configuration_hash,omitempty"`
}

type StatusMessage struct {
	Message string `json:"message"`
}

type RequestOptions struct {
	BaseURL   *string
	StatusURL *string
}

func parseStatusListen(listen string) string {
	re := regexp.MustCompile(`^([\w\.:]+)\s*(.*)?`)
	matches := re.FindStringSubmatch(listen)

	if len(matches) == 0 {
		return ""
	}

	address := matches[1]
	extraParams := matches[2]

	// use http protocol by default
	protocol := "http://"

	// if the listen address contains ssl, use https protocol
	if strings.Contains(extraParams, "ssl") {
		protocol = "https://"
	}

	return fmt.Sprintf("%s%s", protocol, address)
}

// NewClientWithOpts returns a Client which talks to Kong's Admin API and Status API.
func NewClientWithOpts(requestOpts RequestOptions, client *http.Client) (*Client, error) {
	if client == nil {
		transport := &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: DefaultTimeout,
			}).DialContext,
			TLSHandshakeTimeout: DefaultTimeout,
		}
		client = &http.Client{
			Timeout:   DefaultTimeout,
			Transport: transport,
		}
	}
	kong := new(Client)
	kong.client = client
	var rootURL string
	if requestOpts.BaseURL != nil {
		rootURL = *requestOpts.BaseURL
	} else if urlFromEnv := os.Getenv("KONG_ADMIN_URL"); urlFromEnv != "" {
		rootURL = urlFromEnv
	} else {
		rootURL = defaultBaseURL
	}
	parsedRootURL, err := url.ParseRequestURI(rootURL)
	if err != nil {
		return nil, fmt.Errorf("parsing URL: %w", err)
	}
	kong.baseRootURL = parsedRootURL.String()

	var statusURL string
	if requestOpts.StatusURL != nil {
		statusURL = *requestOpts.StatusURL
	} else if listenFromEnv := os.Getenv("KONG_STATUS_LISTEN"); listenFromEnv != "" {
		// KONG_STATUS_LISTEN supports the configuration formats of Kong/Nginx.
		// Only the most commonly used format is handled here.
		// TODO: Support more formats.
		// https://github.com/Kong/kong/blob/2384d2e129d223010fb8a4bb686afb028dca972f/kong.conf.default#L643-L663
		statusURL = parseStatusListen(listenFromEnv)
	} else {
		statusURL = defaultStatusURL
	}
	parsedStatusURL, err := url.ParseRequestURI(statusURL)
	if err != nil {
		return nil, fmt.Errorf("parsing statusURL: %w", err)
	}
	kong.statusURL = parsedStatusURL.String()

	kong.common.client = kong
	kong.ConsumerGroupConsumers = (*ConsumerGroupConsumerService)(&kong.common)
	kong.ConsumerGroups = (*ConsumerGroupService)(&kong.common)
	kong.Consumers = (*ConsumerService)(&kong.common)
	kong.Developers = (*DeveloperService)(&kong.common)
	kong.DeveloperRoles = (*DeveloperRoleService)(&kong.common)
	kong.Services = (*Svcservice)(&kong.common)
	kong.Routes = (*RouteService)(&kong.common)
	kong.Plugins = (*PluginService)(&kong.common)
	kong.Certificates = (*CertificateService)(&kong.common)
	kong.CACertificates = (*CACertificateService)(&kong.common)
	kong.SNIs = (*SNIService)(&kong.common)
	kong.Upstreams = (*UpstreamService)(&kong.common)
	kong.UpstreamNodeHealth = (*UpstreamNodeHealthService)(&kong.common)
	kong.Targets = (*TargetService)(&kong.common)
	kong.Workspaces = (*WorkspaceService)(&kong.common)
	kong.Admins = (*AdminService)(&kong.common)
	kong.RBACUsers = (*RBACUserService)(&kong.common)
	kong.RBACRoles = (*RBACRoleService)(&kong.common)
	kong.RBACEndpointPermissions = (*RBACEndpointPermissionService)(&kong.common)
	kong.RBACEntityPermissions = (*RBACEntityPermissionService)(&kong.common)
	kong.Vaults = (*VaultService)(&kong.common)
	kong.Keys = (*KeyService)(&kong.common)
	kong.KeySets = (*KeySetService)(&kong.common)
	kong.Licenses = (*LicenseService)(&kong.common)
	kong.FilterChains = (*FilterChainService)(&kong.common)

	kong.credentials = (*credentialService)(&kong.common)
	kong.KeyAuths = (*KeyAuthService)(&kong.common)
	kong.BasicAuths = (*BasicAuthService)(&kong.common)
	kong.HMACAuths = (*HMACAuthService)(&kong.common)
	kong.JWTAuths = (*JWTAuthService)(&kong.common)
	kong.MTLSAuths = (*MTLSAuthService)(&kong.common)
	kong.ACLs = (*ACLService)(&kong.common)

	kong.GraphqlRateLimitingCostDecorations = (*GraphqlRateLimitingCostDecorationService)(&kong.common)
	kong.DegraphqlRoutes = (*DegraphqlRouteService)(&kong.common)

	kong.Schemas = (*SchemaService)(&kong.common)

	kong.Oauth2Credentials = (*Oauth2Service)(&kong.common)
	kong.Tags = (*TagService)(&kong.common)
	kong.Info = (*InfoService)(&kong.common)

	kong.CustomEntities = (*CustomEntityService)(&kong.common)
	kong.Registry = custom.NewDefaultRegistry()

	for i := 0; i < len(defaultCustomEntities); i++ {
		err := kong.Register(defaultCustomEntities[i].Type(),
			&defaultCustomEntities[i])
		if err != nil {
			return nil, err
		}
	}
	kong.logger = os.Stderr
	return kong, nil
}

// NewClient returns a Client which talks to Admin API of Kong
func NewClient(baseURL *string, client *http.Client) (*Client, error) {
	return NewClientWithOpts(RequestOptions{BaseURL: baseURL}, client)
}

// SetDoer sets a Doer implementation to be used for custom request dispatching.
func (c *Client) SetDoer(doer Doer) *Client {
	c.doer = doer
	return c
}

// Doer returns the Doer used by this client.
func (c *Client) Doer() Doer {
	return c.doer
}

// SetWorkspace sets the Kong Enteprise workspace in the client.
// Calling this function with an empty string resets the workspace to default workspace.
func (c *Client) SetWorkspace(workspace string) {
	c.workspaceLock.Lock()
	defer c.workspaceLock.Unlock()
	c.workspace = workspace
}

// Workspace return the workspace
func (c *Client) Workspace() string {
	c.workspaceLock.RLock()
	defer c.workspaceLock.RUnlock()
	return c.workspace
}

// baseURL build the base URL from the rootURL and the workspace
func (c *Client) workspacedBaseURL(workspace string) string {
	if len(workspace) > 0 {
		return c.baseRootURL + "/" + workspace
	}
	return c.baseRootURL
}

// DoRAW executes an HTTP request and returns an http.Response
// the caller is responsible for closing the response body.
func (c *Client) DoRAW(ctx context.Context, req *http.Request) (*http.Response, error) {
	var err error
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	// log the request
	err = c.logRequest(req)
	if err != nil {
		return nil, err
	}

	// Make the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making HTTP request: %w", err)
	}

	return resp, err
}

// Do executes an HTTP request and returns a Response.
//
// The caller can optionally provide v parameter, which when provided will contain
// the response body. Do supports wither an io.Writer (which will contain the
// response body verbatim) or anything else which the body should be unmarshalled
// into.
//
// By default, Do() calls DoRaw() to send the request and return the response before unmarshalling, logging,
// and error handling. The Client's WithDoer() method allows overriding this to inject custom behavior.
func (c *Client) Do(
	ctx context.Context,
	req *http.Request,
	v interface{},
) (*Response, error) {
	if c.UserAgent != "" && req != nil {
		req.Header.Add("User-Agent", c.UserAgent)
	}

	var resp *http.Response
	var err error

	if c.doer != nil {
		resp, err = c.doer(ctx, c.client, req) //nolint:bodyclose
	} else {
		resp, err = c.DoRAW(ctx, req) //nolint:bodyclose
	}

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = c.logResponse(resp); err != nil {
		return nil, err
	}

	response := newResponse(resp)

	// Check for API errors.
	// If an error status code was returned, then parse the body and create
	// an API Error out of it.
	if err = hasError(resp); err != nil {
		return response, err
	}

	if v != nil {
		switch v := v.(type) {
		case io.Writer:
			_, err = io.Copy(v, resp.Body)
			if err != nil {
				return nil, fmt.Errorf("failed copying response body: %w", err)
			}
			return response, nil
		default:
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, fmt.Errorf("failed decoding response body: %w", err)
			}
			return response, nil
		}
	}

	return response, nil
}

// ErrorOrResponseError helps to handle the case where
// there might not be a "hard" (connection) error but the
// response itself represents an error.
func ErrorOrResponseError(res *Response, err error) error {
	if err != nil {
		return err
	}
	if res.StatusCode >= http.StatusBadRequest { // errors start at 400
		return fmt.Errorf("unexpected response: %q", res.Status)
	}
	return nil
}

// SetDebugMode enables or disables logging of
// the request to the logger set by SetLogger().
// By default, debug logging is disabled.
func (c *Client) SetDebugMode(enableDebug bool) {
	c.debug = enableDebug
}

func (c *Client) logRequest(r *http.Request) error {
	if !c.debug {
		return nil
	}
	dump, err := httputil.DumpRequestOut(r, true)
	if err != nil {
		return err
	}
	_, err = c.logger.Write(append(dump, '\n'))
	return err
}

func (c *Client) logResponse(r *http.Response) error {
	if !c.debug {
		return nil
	}
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		return err
	}
	_, err = c.logger.Write(append(dump, '\n'))
	return err
}

// SetLogger sets the debug logger, defaults to os.StdErr
func (c *Client) SetLogger(w io.Writer) {
	if w == nil {
		return
	}
	c.logger = w
}

// Status returns the status of a Kong node
func (c *Client) Status(ctx context.Context) (*Status, error) {
	req, err := c.NewRequest("GET", "/status", nil, nil)
	if err != nil {
		return nil, err
	}

	var s Status
	_, err = c.Do(ctx, req, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// Ready returns 200 only after the Kong node has configured itself and is ready to start proxying traffic.
func (c *Client) Ready(ctx context.Context) (*StatusMessage, error) {
	req, err := http.NewRequest("GET", c.statusURL+"/status/ready", nil)
	if err != nil {
		return nil, err
	}

	var sm StatusMessage
	_, err = c.Do(ctx, req, &sm)
	if err != nil {
		return nil, err
	}
	return &sm, nil
}

// Config gets the specified config from the configured Admin API endpoint
// and should contain the JSON serialized body that adheres to the configuration
// format specified at:
// https://docs.konghq.com/gateway/latest/production/deployment-topologies/db-less-and-declarative-config/#declarative-configuration-format
// It returns the response body and an error, if it encounters any.
func (c *Client) Config(ctx context.Context) ([]byte, error) {
	req, err := c.NewRequest("GET", "/config", nil, nil)
	if err != nil {
		return nil, err
	}

	var configWrapper map[string]string
	_, err = c.Do(ctx, req, &configWrapper)
	if err != nil {
		return nil, err
	}
	config, ok := configWrapper["config"]
	if !ok {
		return nil, errors.New("config field not found in GET /config response body")
	}

	return []byte(config), nil
}

// Root returns the response of GET request on root of Admin API (GET / or /kong with a workspace).
func (c *Client) Root(ctx context.Context) (map[string]interface{}, error) {
	endpoint := "/"
	ws := c.Workspace()
	if len(ws) > 0 {
		endpoint = "/kong"
	}
	req, err := c.NewRequestRaw("GET", c.workspacedBaseURL(ws), endpoint, nil, nil)
	if err != nil {
		return nil, err
	}
	var info map[string]interface{}
	_, err = c.Do(ctx, req, &info)
	if err != nil {
		return nil, err
	}
	return info, nil
}

// RootJSON returns the response of GET request on the root of the Admin API
// (GET / or /kong with a workspace) returning the raw JSON response data.
func (c *Client) RootJSON(ctx context.Context) ([]byte, error) {
	endpoint := "/"
	ws := c.Workspace()
	if len(ws) > 0 {
		endpoint = "/kong"
	}

	req, err := c.NewRequestRaw("GET", c.workspacedBaseURL(ws), endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.DoRAW(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) BaseRootURL() string {
	return c.baseRootURL
}

// ReloadDeclarativeRawConfig sends out the specified config to configured Admin
// API endpoint using the provided reader which should contain the JSON
// serialized body that adheres to the configuration format specified at:
// https://docs.konghq.com/gateway/latest/production/deployment-topologies/db-less-and-declarative-config/#declarative-configuration-format
// It returns APIError with a response body in case it receives a valid HTTP response with <200 or >=400 status codes.
func (c *Client) ReloadDeclarativeRawConfig(
	ctx context.Context,
	config io.Reader,
	checkHash bool,
	flattenErrors bool,
) error {
	type sendConfigParams struct {
		CheckHash     int `url:"check_hash,omitempty"`
		FlattenErrors int `url:"flatten_errors,omitempty"`
	}
	var checkHashI int
	if checkHash {
		checkHashI = 1
	}
	var flattenErrorsI int
	if flattenErrors {
		flattenErrorsI = 1
	}
	req, err := c.NewRequest(
		"POST",
		"/config",
		sendConfigParams{CheckHash: checkHashI, FlattenErrors: flattenErrorsI},
		config,
	)
	if err != nil {
		return fmt.Errorf("creating new HTTP request for /config: %w", err)
	}

	resp, err := c.DoRAW(ctx, req)
	if err != nil {
		return fmt.Errorf("failed posting new config to /config: %w", err)
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read /config %d status response body: %w", resp.StatusCode, err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return NewAPIErrorWithRaw(resp.StatusCode, "failed posting new config to /config", b)
	}

	return nil
}
