package kong

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/blang/semver/v4"
)

var (
	versionParts   = 4
	kong140Version = semver.MustParse("1.4.0")
)

// String returns pointer to s.
func String(s string) *string {
	return &s
}

// Bool returns a pointer to b.
func Bool(b bool) *bool {
	return &b
}

// Int returns a pointer to i.
func Int(i int) *int {
	return &i
}

func isEmptyString(s *string) bool {
	return s == nil || strings.TrimSpace(*s) == ""
}

// StringSlice converts a slice of string to a
// slice of *string
func StringSlice(elements ...string) []*string {
	var res []*string
	for _, element := range elements {
		e := element
		res = append(res, &e)
	}
	return res
}
func stringArrayToString(arr []*string) string {
	if arr == nil {
		return "nil"
	}

	var buf bytes.Buffer
	buf.WriteString("[ ")
	l := len(arr)
	for i, el := range arr {
		buf.WriteString(*el)
		if i != l-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(" ]")
	return buf.String()
}

// headerRoundTripper injects Headers into requests
// made via RT.
type headerRoundTripper struct {
	headers http.Header
	rt      http.RoundTripper
}

// RoundTrip satisfies the RoundTripper interface.
func (t headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	newRequest := new(http.Request)
	*newRequest = *req
	newRequest.Header = make(http.Header, len(req.Header))
	for k, s := range req.Header {
		newRequest.Header[k] = append([]string(nil), s...)
	}
	for k, v := range t.headers {
		newRequest.Header[k] = v
	}
	return t.rt.RoundTrip(newRequest)
}

// RoundTripperWithHTTPHeaders returns a client which injects headers
// before sending any request.
func HTTPClientWithHeaders(client *http.Client,
	headers http.Header) http.Client {
	var res http.Client
	if client == nil {
		defaultTransport := http.DefaultTransport.(*http.Transport)
		res.Transport = defaultTransport
	} else {
		res = *client
	}
	res.Transport = headerRoundTripper{
		headers: headers,
		rt:      client.Transport,
	}
	return res
}

func cleanSemVer(v string) (semver.Version, error) {
	// fix enterprise edition semver adding patch number
	// fix enterprise edition version with dash
	// fix bad version formats like 0.13.0preview1
	re := regexp.MustCompile(`(\d+\.\d+)(?:[\.-](\d+))?(?:\-?(.+)$|$)`)
	m := re.FindStringSubmatch(v)
	if len(m) != versionParts {
		return semver.Version{}, fmt.Errorf("Unknown Kong version : '%v'", v)
	}
	if m[2] == "" {
		m[2] = "0"
	}
	if m[3] != "" {
		m[3] = "-" + strings.Replace(m[3], "enterprise-edition", "enterprise", 1)
		m[3] = strings.Replace(m[3], ".", "", -1)
	}
	v = fmt.Sprintf("%s.%s%s", m[1], m[2], m[3])
	return semver.Make(v)
}

func getKong(root map[string]interface{}) (*Kong, error) {
	version := root["version"].(string)
	semVer, err := cleanSemVer(version)
	if err != nil {
		return nil, err
	}
	kong := new(Kong)
	kong.Version = semVer.String()
	kong.Enterprise = strings.Contains(version, "enterprise")
	kong.Database = root["configuration"].(map[string]interface{})["database"].(string)
	kong.Credentials.minVersion = kong140Version.String()
	kong.Credentials.hasTagSupport = semVer.GTE(kong140Version)
	return kong, nil
}
