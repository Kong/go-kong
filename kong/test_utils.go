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
		t.Skip()
	}
}

// RunWhenEnterprise skips a test if the version
// of Kong running is not enterprise edition. Skips
// the current test if the version of Kong doesn't
// fall within the version range. If a test requires
// RBAC and RBAC is not enabled on Kong the test
// will be skipped
func RunWhenEnterprise(t *testing.T, versionRange string, required RequiredFeatures) {
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
		t.Log("non-Enterprise test Kong instance, skipping")
		t.Skip()
	}
	configuration := info["configuration"].(map[string]interface{})

	if required.RBAC && configuration["rbac"].(string) != "on" {
		t.Log("RBAC not enabled on test Kong instance, skipping")
		t.Skip()
	}

	if required.Portal && !configuration["portal"].(bool) {
		t.Log("Portal not enabled on test Kong instance, skipping")
		t.Skip()
	}

	r, err := NewRange(versionRange)
	if err != nil {
		t.Error(err)
	}
	if !r(currentVersion) {
		t.Skip()
	}
}

// SkipWhenEnterprise skips a test if the Kong version is an Enterprise version
func SkipWhenEnterprise(t *testing.T) {
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
		t.Log("non-Enterprise test Kong instance, skipping")
		t.Skip()
	}
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
