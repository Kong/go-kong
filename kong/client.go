package kong

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/kong/go-kong/kong/custom"
)

const (
	defaultBaseURL = "http://localhost:8001"
	// DefaultTimeout is the timeout used for network connections and requests
	// including TCP, TLS and HTTP layers.
	DefaultTimeout = 60 * time.Second
)

var pageSize = 1000

type service struct {
	client *Client
}

var defaultCtx = context.Background()

// Client talks to the Admin API or control plane of a
// Kong cluster
type Client struct {
	client                  *http.Client
	defaultRootURL          string
	workspace               string       // Do not access directly. Use Workspace()/SetWorkspace().
	workspaceLock           sync.RWMutex // Synchronizes access to workspace.
	common                  service
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

	Schemas AbstractSchemaService

	logger         io.Writer
	debug          bool
	CustomEntities AbstractCustomEntityService

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

// NewClient returns a Client which talks to Admin API of Kong
func NewClient(baseURL *string, client *http.Client) (*Client, error) {
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
	if baseURL != nil {
		rootURL = *baseURL
	} else if urlFromEnv := os.Getenv("KONG_ADMIN_URL"); urlFromEnv != "" {
		rootURL = urlFromEnv
	} else {
		rootURL = defaultBaseURL
	}
	url, err := url.ParseRequestURI(rootURL)
	if err != nil {
		return nil, fmt.Errorf("parsing URL: %w", err)
	}
	kong.defaultRootURL = url.String()

	kong.common.client = kong
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

	kong.credentials = (*credentialService)(&kong.common)
	kong.KeyAuths = (*KeyAuthService)(&kong.common)
	kong.BasicAuths = (*BasicAuthService)(&kong.common)
	kong.HMACAuths = (*HMACAuthService)(&kong.common)
	kong.JWTAuths = (*JWTAuthService)(&kong.common)
	kong.MTLSAuths = (*MTLSAuthService)(&kong.common)
	kong.ACLs = (*ACLService)(&kong.common)

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
		return c.defaultRootURL + "/" + workspace
	}
	return c.defaultRootURL
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

// Do executes an HTTP request and returns a kong.Response
func (c *Client) Do(ctx context.Context, req *http.Request,
	v interface{}) (*Response, error) {
	resp, err := c.DoRAW(ctx, req)
	if err != nil {
		return nil, err
	}

	// log the response
	err = c.logResponse(resp)
	if err != nil {
		return nil, err
	}

	response := newResponse(resp)

	///check for API errors
	if err = hasError(resp); err != nil {
		return response, err
	}

	// response
	if v != nil {
		if writer, ok := v.(io.Writer); ok {
			_, err = io.Copy(writer, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}
	return response, err
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
