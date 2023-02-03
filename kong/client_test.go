package kong

import (
	"bytes"
	"context"
	"encoding/json"
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

func TestBaseRootURL(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		client, err := NewClient(nil, nil)
		require.NoError(t, err)
		require.NotNil(t, client)

		require.Equal(t, client.BaseRootURL(), "http://localhost:8001")
	})

	t.Run("set via env", func(t *testing.T) {
		t.Setenv("KONG_ADMIN_URL", "https://customkong.com")
		client, err := NewClient(nil, nil)
		require.NoError(t, err)
		require.NotNil(t, client)

		require.Equal(t, client.BaseRootURL(), "https://customkong.com")
	})

	t.Run("set via flag", func(t *testing.T) {
		client, err := NewClient(String("https://customkong2.com"), nil)
		require.NoError(t, err)
		require.NotNil(t, client)

		require.Equal(t, client.BaseRootURL(), "https://customkong2.com")
	})
}

func TestReloadDeclarativeRawConfig(t *testing.T) {
	RunWhenDBMode(t, "off")

	tests := []struct {
		name    string
		config  Configuration
		wantErr bool
	}{
		{
			name: "basic config works",
			config: Configuration{
				"_format_version": "1.1",
				"services": []Configuration{
					{
						"host":     "mockbin.com",
						"port":     443,
						"protocol": "https",
						"routes": []Configuration{
							{"paths": []string{"/"}},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing _format_version fails",
			config: Configuration{
				"services": []Configuration{
					{
						"host":     "mockbin.com",
						"port":     443,
						"protocol": "https",
						"routes": []Configuration{
							{"paths": []string{"/"}},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid config fails",
			config: Configuration{
				"dummy_key": []Configuration{
					{
						"host":     "mockbin.com",
						"port":     443,
						"protocol": "https",
					},
				},
			},
			wantErr: true,
		},
	}

	// POST /config actually interprets query string fields as additional fields within the config unless they're
	// explicitly stripped by the API code (no idea why this behavior exists). Without stripping "flatten_errors" in
	// https://github.com/Kong/kong/blob/master/kong/api/routes/config.lua#L71 Kong will actually reject configuration
	// because it thinks "flatten_errors" is an actual (and invalid) field inside the config.
	//
	// This is the one test where we want version-dependent behavior, but only for one value within the test. We test
	// config updates on all DB-less capable versions, but only set flattenErrors=true on 3.2+. To handle that, this
	// snippet is borrowed from RunWhenKong, to allow toggling that behavior only on or off depending on the version.
	var flattenErrors bool
	client, err := NewTestClient(nil, nil)
	if err != nil {
		t.Error(err)
	}
	info, err := client.Root(defaultCtx)
	if err != nil {
		t.Error(err)
	}
	version := VersionFromInfo(info)
	currentVersion, err := ParseSemanticVersion(version)
	if err != nil {
		t.Error(err)
	}
	r, err := NewRange(">=3.2.0")
	if err != nil {
		t.Error(err)
	}
	if r(currentVersion) {
		t.Log("Kong version >=3.2, enabling flattenErrors for ReloadDeclarativeRawConfig")
		flattenErrors = true
	} else {
		t.Log("Kong version <3.2, disabling flattenErrors for ReloadDeclarativeRawConfig")
	}

	for _, tt := range tests {
		client, err := NewTestClient(nil, nil)
		require.NoError(t, err)
		require.NotNil(t, client)

		tt := tt
		t.Run("with_schema/"+tt.name, func(t *testing.T) {
			ctx := context.Background()
			b, err := json.Marshal(tt.config)
			require.NoError(t, err)

			body, err := client.ReloadDeclarativeRawConfig(ctx, bytes.NewBuffer(b), true, flattenErrors)
			// We only get empty body when there's a transient network error or
			// we fail to read the response body which shouldn't happen in tests.
			assert.NotEmpty(t, string(body))

			if tt.wantErr {
				assert.Errorf(t, err, "Client.SendConfig() got unexpected error = %v", err)
			} else {
				assert.NoError(t, err)
			}

			// this is somewhat untrue: network or HTTP-level failures _can_ result in a nil response body. however,
			// none of our test cases should cause network or HTTP-level failures, so fail if they do occur. if this
			// _does_ encounter such a failure, we need to investigate and either update tests or fix some upstream bug
			// if it's not some transient issue with the testing environment
			require.NotNilf(t, body, "body was nil; should never be nil")
		})
	}
}
