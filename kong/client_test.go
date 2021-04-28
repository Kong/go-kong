//nolint:unparam
package kong

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewTestClient(t *testing.T) {
	assert := assert.New(t)

	client, err := NewTestClient(String("foo/bar"), nil)
	assert.Nil(client)
	assert.NotNil(err)
}

func TestKongStatus(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	status, err := client.Status(defaultCtx)
	assert.Nil(err)
	assert.NotNil(status)
}

func TestRoot(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	root, err := client.Root(defaultCtx)
	assert.Nil(err)
	assert.NotNil(root)
	assert.NotNil(root["version"])
}

func TestDo(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	req, err := client.NewRequest("GET", "/does-not-exist", nil, nil)
	assert.Nil(err)
	assert.NotNil(req)
	resp, err := client.Do(context.Background(), req, nil)
	assert.True(IsNotFoundErr(err))
	assert.NotNil(resp)
	assert.Equal(404, resp.StatusCode)

	req, err = client.NewRequest("POST", "/", nil, nil)
	assert.Nil(err)
	assert.NotNil(req)
	resp, err = client.Do(context.Background(), req, nil)
	assert.NotNil(err)
	assert.NotNil(resp)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)
	assert.Empty(body)
	assert.Equal(405, resp.StatusCode)
}

func TestMain(m *testing.M) {
	// to test ListAll code for pagination
	pageSize = 1
	os.Exit(m.Run())
}

var currentVersion semver.Version

// runWhenKong skips the current test if the version of Kong doesn't
// fall in the semverRange.
// This helper function can be used in tests to write version specific
// tests for Kong.
func runWhenKong(t *testing.T, semverRange string) {
	if currentVersion.Major == 0 {
		client, err := NewTestClient(nil, nil)
		if err != nil {
			t.Error(err)
		}
		res, err := client.Kong(defaultCtx)
		if err != nil {
			t.Error(err)
		}
		currentVersion = res.Version
	}
	r, err := semver.ParseRange(semverRange)
	if err != nil {
		t.Error(err)
	}
	if !r(currentVersion) {
		t.Skip()
	}

}

type requiredFeatures struct {
	portal bool
	rbac   bool
}

// runWhenEnterprise skips a test if the version
// of Kong running is not enterprise edition. Skips
// the current test if the version of Kong doesn't
// fall within the semver range. If a test requires
// RBAC and RBAC is not enabled on Kong the test
// will be skipped
func runWhenEnterprise(t *testing.T, semverRange string, required requiredFeatures) {
	client, err := NewTestClient(nil, nil)
	if err != nil {
		t.Error(err)
	}
	res, err := client.Kong(defaultCtx)
	if err != nil {
		t.Error(err)
	}

	if !res.Enterprise {
		t.Log("non-Enterprise test Kong instance, skipping")
		t.Skip()
	}

	if required.rbac && !res.RBAC {
		t.Log("RBAC not enabled on test Kong instance, skipping")
		t.Skip()
	}

	if required.portal && !res.Portal {
		t.Log("Portal not enabled on test Kong instance, skipping")
		t.Skip()
	}

	runWhenKong(t, semverRange)

}

func TestRunWhenEnterprise(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", requiredFeatures{})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	root, err := client.Root(defaultCtx)
	assert.Nil(err)
	assert.NotNil(root)
	v := root["version"].(string)
	assert.Contains(v, "enterprise")
}

func NewTestClient(baseURL *string, client *http.Client) (*Client, error) {
	if value, exists := os.LookupEnv("KONG_ADMIN_TOKEN"); exists {
		c := &http.Client{}
		defaultTransport := http.DefaultTransport.(*http.Transport)
		c.Transport = defaultTransport
		c.Transport = &headerRoundTripper{
			headers: http.Header{
				"kong-admin-token": []string{value},
			},
			rt: defaultTransport,
		}
		return NewClient(baseURL, c)
	}
	return NewClient(baseURL, client)
}

type TestWorkspace struct {
	workspace      *Workspace
	client         *Client
	originalConfig map[string]interface{}
}

// NewTestWorkspace takes a client and workspace name and returns a TestWorkspace
// containing the workspace object for the workspace with that name, the workspace's
// original configuration, and the client. It returns an error if that workspace
// cannot be retrieved.
func NewTestWorkspace(client *Client, wsName string) (*TestWorkspace, error) {
	workspace, err := client.Workspaces.Get(defaultCtx, String(wsName))
	if err != nil {
		return nil, err
	}
	testWs := TestWorkspace{client: client, workspace: workspace}
	testWs.originalConfig = map[string]interface{}{}
	for k, v := range workspace.Config {
		testWs.originalConfig[k] = v
	}
	return &testWs, nil
}

// Reset restores a test workspace to its original configuration.
func (t *TestWorkspace) Reset() error {
	t.workspace.Config = t.originalConfig
	_, err := t.client.Workspaces.Update(defaultCtx, t.workspace)
	return err
}

// UpdateConfig patches a workspace's configuration with the provided values.
func (t *TestWorkspace) UpdateConfig(config map[string]interface{}) error {
	t.workspace.Config = config
	_, err := t.client.Workspaces.Update(defaultCtx, t.workspace)
	return err
}

func TestTestWorkspace(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", requiredFeatures{portal: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	wsName := "default"

	origWorkspace, err := client.Workspaces.Get(defaultCtx, String(wsName))
	assert.Nil(err)

	testWs, err := NewTestWorkspace(client, wsName)
	assert.Nil(err)
	assert.Equal(wsName, *testWs.workspace.Name)

	err = testWs.UpdateConfig(map[string]interface{}{"portal": true, "portal_auto_approve": true})
	assert.Nil(err)
	currWorkspace, err := client.Workspaces.Get(defaultCtx, String(wsName))
	assert.Nil(err)
	assert.Equal(currWorkspace.Config["portal"], true)
	assert.Equal(currWorkspace.Config["portal_auto_approve"], true)

	err = testWs.Reset()
	assert.Nil(err)
	currWorkspace, err = client.Workspaces.Get(defaultCtx, String(wsName))
	assert.Nil(err)
	assert.Equal(currWorkspace.Config, origWorkspace.Config)
}
