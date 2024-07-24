package kong

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
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

func TestKongReady(t *testing.T) {
	kongImageTag := os.Getenv("KONG_IMAGE_TAG")
	if kongImageTag == "" {
		t.Skip("KONG_IMAGE_TAG environment variable is not set")
	}

	currentVersion, err := semver.NewVersion(kongImageTag)
	if err != nil {
		// We have set the KONG_IMAGE_TAG env var to master when running the test for nightly builds.
		if kongImageTag != "master" {
			t.Fatalf("Failed to parse KONG_IMAGE_TAG: %v", err)
		}
	} else {

		// This API was only made available since Kong v3.3.
		// ref: https://docs.konghq.com/gateway/api/status/latest/#/default/get_status_ready
		minVersion, err := semver.NewVersion("3.3")
		if err != nil {
			t.Fatalf("Failed to parse the minimum required version: %v", err)
		}

		// Compare the versions and skip the test if the current version is less than the minimum version
		if currentVersion.LessThan(minVersion) {
			t.Skipf("Skipping test because KONG_IMAGE_TAG %s is less than 3.3", kongImageTag)
		}
	}

	assert := assert.New(t)
	require := require.New(t)

	client, err := NewTestClientWithOpts(
		RequestOptions{
			BaseURL:   String("http://localhost:8001"),
			StatusURL: String("http://localhost:8100"),
		}, nil,
	)
	require.NoError(err)
	require.NotNil(client)

	sm, err := client.Ready(defaultCtx)
	if err != nil {
		// for dbless mode, the ready endpoint returns 503
		assert.Equal("HTTP status 503 (message: \"no configuration available (empty configuration present)\")", err.Error())
		assert.Nil(sm)
	} else {
		// for db-mode, the ready endpoint returns 200
		assert.Equal("ready", sm.Message)
	}
}

func TestKongConfig(t *testing.T) {
	RunWhenDBMode(t, "off")
	client, err := NewTestClient(nil, nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	config, err := client.Config(defaultCtx)
	require.NoError(t, err)
	require.NotNil(t, config)

	var configStruct map[string]any
	require.NoError(t, yaml.Unmarshal(config, &configStruct))
	require.Contains(t, configStruct, "_format_version")
	require.Contains(t, configStruct, "_transform")
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
		name              string
		kongClientFunc    func() (*Client, error)
		requestAssertions func(t *testing.T, req *http.Request)
	}{
		{
			name: "nil http.Client",
			kongClientFunc: func() (*Client, error) {
				return NewTestClient(nil, nil)
			},
		},
		{
			name: "default/uninitialized http.Client",
			kongClientFunc: func() (*Client, error) {
				return NewTestClient(nil, &http.Client{})
			},
		},
		{
			name: "default/uninitialized http.Client with HTTPClientWithHeaders",
			kongClientFunc: func() (*Client, error) {
				return NewTestClient(nil, HTTPClientWithHeaders(&http.Client{}, nil))
			},
		},
		{
			name: "default http.Client and custom doer",
			kongClientFunc: func() (*Client, error) {
				httpClient := &http.Client{}

				doer := func(_ context.Context, httpClient *http.Client, req *http.Request) (*http.Response, error) {
					req.Header.Add("Doer", "doer")

					return httpClient.Do(req)
				}

				client, err := NewTestClient(nil, httpClient)
				if err != nil {
					return client, err
				}

				return client.SetDoer(doer), nil
			},
			requestAssertions: func(t *testing.T, req *http.Request) {
				headers := http.Header{
					"Doer": []string{"doer"},
				}

				assert.NotNil(t, req)
				assertHeadersExist(t, req, headers)
			},
		},
	}

	for _, tc := range testcases {
		tc := tc

		T.Run(tc.name, func(T *testing.T) {
			assert := assert.New(T)
			require := require.New(T)

			client, err := tc.kongClientFunc()
			require.NoError(err)
			require.NotNil(client)

			req, err := client.NewRequest("GET", "/does-not-exist", nil, nil)
			require.NoError(err)
			require.NotNil(req)
			resp, err := client.Do(context.Background(), req, nil)
			assert.True(IsNotFoundErr(err), "got %v", err)
			require.NotNil(resp)
			assert.Equal(404, resp.StatusCode)

			if tc.requestAssertions != nil {
				tc.requestAssertions(T, req)
			}

			req, err = client.NewRequest("POST", "/", nil, nil)
			require.NoError(err)
			require.NotNil(req)
			resp, err = client.Do(context.Background(), req, nil)
			require.NotNil(err)
			require.NotNil(resp)
			assert.Equal(405, resp.StatusCode)

			if tc.requestAssertions != nil {
				tc.requestAssertions(T, req)
			}
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
	SkipWhenKongRouterFlavor(t, Expressions)
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
	require.NoError(t, err)
	info, err := client.Root(defaultCtx)
	require.NoError(t, err)
	version := VersionFromInfo(info)
	currentVersion, err := ParseSemanticVersion(version)
	require.NoError(t, err)
	r, err := NewRange(">=3.2.0")
	require.NoError(t, err)

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

			err = client.ReloadDeclarativeRawConfig(ctx, bytes.NewBuffer(b), true, flattenErrors)

			if tt.wantErr {
				assert.Error(t, err)
				apiErr := &APIError{}
				assert.True(t, errors.As(err, &apiErr))
				assert.NotEmpty(t, apiErr.Raw(), "expected non-empty response body in APIError")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestReloadDeclarativeRawConfig_NetworkErrorDoesntReturnAPIError(t *testing.T) {
	RunWhenDBMode(t, "off")
	SkipWhenKongRouterFlavor(t, Expressions)

	// Point to a non-reachable URL so that we get a network error.
	client, err := NewClient(String("http://non-reachable.url"), nil)
	require.NoError(t, err)
	require.NotNil(t, client)

	err = client.ReloadDeclarativeRawConfig(context.Background(), bytes.NewReader([]byte("dummy-config")), true, true)
	require.Error(t, err)
	require.False(t, errors.Is(err, &APIError{}), "expected error to not be an APIError")
}

func assertHeadersExist(t *testing.T, request *http.Request, headers http.Header) {
	for k, v := range headers {
		assert.Contains(t, request.Header, k)
		for _, val := range v {
			assert.Contains(t, request.Header.Values(k), val)
		}
	}
}
