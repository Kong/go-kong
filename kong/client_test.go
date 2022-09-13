package kong

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.NoError(err)
	assert.NotNil(client)

	status, err := client.Status(defaultCtx)
	assert.NoError(err)
	assert.NotNil(status)
}

func TestRoot(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	root, err := client.Root(defaultCtx)
	assert.NoError(err)
	assert.NotNil(root)
	assert.NotNil(root["version"])
}

func TestRootJSON(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	root, err := client.RootJSON(defaultCtx)
	assert.NoError(err)
	assert.NotEmpty(root)
	assert.Contains(string(root), `"version"`)
}

func TestDo(T *testing.T) {
	testcases := []struct {
		name           string
		httpClientFunc func() *http.Client
	}{
		{
			name:           "nil http.Client",
			httpClientFunc: func() *http.Client { return nil },
		},
		{
			name:           "default/uninitialized http.Client",
			httpClientFunc: func() *http.Client { return &http.Client{} },
		},
		{
			name:           "default/uninitialized http.Client with HTTPClientWithHeaders",
			httpClientFunc: func() *http.Client { return HTTPClientWithHeaders(&http.Client{}, nil) },
		},
	}

	for _, tc := range testcases {
		tc := tc

		T.Run(tc.name, func(T *testing.T) {
			assert := assert.New(T)
			require := require.New(T)

			client, err := NewTestClient(nil, tc.httpClientFunc())
			require.NoError(err)
			require.NotNil(client)

			req, err := client.NewRequest("GET", "/does-not-exist", nil, nil)
			assert.NoError(err)
			require.NotNil(req)
			resp, err := client.Do(context.Background(), req, nil)
			assert.True(IsNotFoundErr(err), "got %v", err)
			require.NotNil(resp)
			assert.Equal(404, resp.StatusCode)

			req, err = client.NewRequest("POST", "/", nil, nil)
			assert.NoError(err)
			require.NotNil(req)
			resp, err = client.Do(context.Background(), req, nil)
			require.NotNil(err)
			require.NotNil(resp)
			body, err := io.ReadAll(resp.Body)
			assert.NoError(err)
			assert.Empty(body)
			assert.Equal(405, resp.StatusCode)
		})
	}
}

func TestMain(m *testing.M) {
	// to test ListAll code for pagination
	pageSize = 1
	os.Exit(m.Run())
}

func TestRunWhenEnterprise(T *testing.T) {
	// TODO refactor this to test that a version is Enterprise without relying on the IsKongGatewayEnterprise function
	// that this calls https://github.com/Kong/go-kong/issues/212
	RunWhenEnterprise(T, ">=0.33.0 <3.0.0", RequiredFeatures{})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	root, err := client.Root(defaultCtx)
	assert.NoError(err)
	assert.NotNil(root)
	v := root["version"].(string)
	assert.Contains(v, "enterprise")
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
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{Portal: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	wsName := "default"

	origWorkspace, err := client.Workspaces.Get(defaultCtx, String(wsName))
	assert.NoError(err)

	testWs, err := NewTestWorkspace(client, wsName)
	assert.NoError(err)
	assert.Equal(wsName, *testWs.workspace.Name)

	err = testWs.UpdateConfig(map[string]interface{}{"portal": true, "portal_auto_approve": true})
	assert.NoError(err)
	currWorkspace, err := client.Workspaces.Get(defaultCtx, String(wsName))
	assert.NoError(err)
	assert.Equal(currWorkspace.Config["portal"], true)
	assert.Equal(currWorkspace.Config["portal_auto_approve"], true)

	err = testWs.Reset()
	assert.NoError(err)
	currWorkspace, err = client.Workspaces.Get(defaultCtx, String(wsName))
	assert.NoError(err)
	assert.Equal(currWorkspace.Config, origWorkspace.Config)
}
