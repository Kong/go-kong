package kong

import (
	"fmt"
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
	configuration, ok := info["configuration"].(map[string]interface{})
	if !ok {
		t.Errorf("failed to cast 'configuration' to map[string]interface{}")
	}

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
		t.Skip("Enterprise test Kong instance, skipping")
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

	dbMode, err := getKongConfigString(t, "database")
	if err != nil {
		t.Skip(err.Error())
	}
	if dbMode != dbmode {
		t.Skipf("detected Kong running in dbmode:%q but requested dbmode:%q", dbMode, dbmode)
	}
}

type RouterFlavor string

const (
	Traditional           RouterFlavor = "traditional"
	TraditionalCompatible RouterFlavor = "traditional_compatible"
	Expressions           RouterFlavor = "expressions"
)

func SkipWhenKongRouterFlavor(t *testing.T, flavor ...RouterFlavor) {
	t.Helper()

	routerFlavor, err := getKongConfigString(t, "router_flavor")
	if err != nil {
		t.Skip(err.Error())
	}
	for _, f := range flavor {
		if RouterFlavor(routerFlavor) == f {
			t.Skipf("router flavor:%q skipping", f)
		}
	}
}

func RunWhenKongRouterFlavor(t *testing.T, flavor RouterFlavor) {
	t.Helper()

	routerFlavor, err := getKongConfigString(t, "router_flavor")
	if err != nil {
		t.Skip(err.Error())
	}
	if RouterFlavor(routerFlavor) != flavor {
		t.Skipf("router flavor:%q, expecting %q, skipping", routerFlavor, flavor)
	}
}

func RunWhenKongConfigEnabled(t *testing.T, key string) {
	t.Helper()

	value, err := getKongConfigBoolean(t, key)
	if err != nil {
		t.Skip(err.Error())
	}

	if !value {
		t.Skipf("%q is not enabled", key)
	}
}

func getKongConfig(t *testing.T) map[string]any {
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
		t.Error("failed to find `configuration` key in kong admin response")
	}

	configuration, ok := config.(map[string]any)
	if !ok {
		t.Errorf("`configuration` is not a map but %T", config)
	}

	return configuration
}

func getKongConfigValue(t *testing.T, key string) (interface{}, error) {
	t.Helper()

	configuration := getKongConfig(t)
	value, ok := configuration[key]
	if !ok {
		return nil, fmt.Errorf("`configuration.%s` does not exist", key)
	}
	return value, nil
}

func getKongConfigString(t *testing.T, key string) (string, error) {
	t.Helper()

	value, err := getKongConfigValue(t, key)
	if err != nil {
		return "", err
	}

	str, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("`configuration.%s` is not a string but a %T", key, value)
	}
	return str, nil
}

func getKongConfigBoolean(t *testing.T, key string) (bool, error) {
	t.Helper()

	value, err := getKongConfigValue(t, key)
	if err != nil {
		return false, err
	}

	boolean, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("`configuration.%s` is not a boolean but a %T", key, value)
	}
	return boolean, nil
}
