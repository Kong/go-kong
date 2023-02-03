package kong

import (
	"net/http"
	"os"
	"testing"
)

type RequiredFeatures struct {
	Portal bool
	RBAC   bool
}

// RunWhenKong skips the current test if the version of Kong doesn't
// fall in the versionRange.
// This helper function can be used in tests to write version specific
// tests for Kong.
func RunWhenKong(t *testing.T, versionRange string) {
	t.Helper()

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
	r, err := NewRange(versionRange)
	if err != nil {
		t.Error(err)
	}
	if !r(currentVersion) {
		t.Skipf("kong version %s not in range %s", version, versionRange)
	}
}

// RunWhenEnterprise skips a test if the version
// of Kong running is not enterprise edition. Skips
// the current test if the version of Kong doesn't
// fall within the version range. If a test requires
// RBAC and RBAC is not enabled on Kong the test
// will be skipped
func RunWhenEnterprise(t *testing.T, versionRange string, required RequiredFeatures) {
	t.Helper()

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

	if !currentVersion.IsKongGatewayEnterprise() {
		t.Skip("non-Enterprise test Kong instance, skipping")
	}
	configuration := info["configuration"].(map[string]interface{})

	if required.RBAC && configuration["rbac"].(string) != "on" {
		t.Skip("RBAC not enabled on test Kong instance, skipping")
	}

	if required.Portal && !configuration["portal"].(bool) {
		t.Skip("Portal not enabled on test Kong instance, skipping")
	}

	r, err := NewRange(versionRange)
	if err != nil {
		t.Error(err)
	}
	if !r(currentVersion) {
		t.Skipf("kong version %s not in range %s", version, versionRange)
	}
}

// SkipWhenEnterprise skips a test if the Kong version is an Enterprise version
func SkipWhenEnterprise(t *testing.T) {
	t.Helper()

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

	if currentVersion.IsKongGatewayEnterprise() {
		t.Skip("non-Enterprise test Kong instance, skipping")
	}
}

func NewTestClient(baseURL *string, client *http.Client) (*Client, error) {
	if value, exists := os.LookupEnv("KONG_ADMIN_TOKEN"); exists && value != "" {
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

func RunWhenDBMode(t *testing.T, dbmode string) {
	t.Helper()

	client, err := NewTestClient(nil, nil)
	if err != nil {
		t.Error(err)
	}
	info, err := client.Root(defaultCtx)
	if err != nil {
		t.Error(err)
	}

	config, ok := info["configuration"]
	if !ok {
		t.Skip("failed to find 'configuration' config key in kong configuration")
	}

	configuration, ok := config.(map[string]any)
	if !ok {
		t.Skipf("'configuration' key is not a map but %T", config)
	}

	dbConfig, ok := configuration["database"]
	if !ok {
		t.Skip("failed to find 'database' config key in kong confiration")
	}

	dbMode, ok := dbConfig.(string)
	if !ok {
		t.Skipf("'database' config key is not a string but %T", dbConfig)
	}

	if dbMode != dbmode {
		t.Skipf("detected Kong running in dbmode:%q but requested dbmode:%q", dbMode, dbmode)
	}
}
