package kong

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestStringArrayToString(t *testing.T) {
	assert := assert.New(t)

	arr := StringSlice("foo", "bar")
	s := stringArrayToString(arr)
	assert.Equal("[ foo, bar ]", s)

	arr = StringSlice("foo")
	s = stringArrayToString(arr)
	assert.Equal("[ foo ]", s)

	assert.Equal(stringArrayToString(nil), "nil")
}

func TestString(t *testing.T) {
	assert := assert.New(t)

	s := String("foo")
	assert.Equal("foo", *s)
}

func TestBool(t *testing.T) {
	assert := assert.New(t)

	b := Bool(true)
	assert.Equal(true, *b)
}

func TestInt(t *testing.T) {
	assert := assert.New(t)

	i := Int(42)
	assert.Equal(42, *i)
}

func TestStringSlice(t *testing.T) {
	assert := assert.New(t)

	arrp := StringSlice()
	assert.Empty(arrp)

	arrp = StringSlice("foo", "bar")
	assert.Equal(2, len(arrp))
	assert.Equal("foo", *arrp[0])
	assert.Equal("bar", *arrp[1])
}

func TestFixVersion(t *testing.T) {
	validVersions := map[string]string{
		"0.14.1":                          "0.14.1",
		"0.14.2rc":                        "0.14.2-rc",
		"0.14.2rc1":                       "0.14.2-rc1",
		"0.14.2preview":                   "0.14.2-preview",
		"0.14.2preview1":                  "0.14.2-preview1",
		"0.33-enterprise-edition":         "0.33.0-enterprise",
		"0.33-1-enterprise-edition":       "0.33.1-enterprise",
		"1.3.0.0-enterprise-edition-lite": "1.3.0-0-enterprise-lite",
		"1.3.0.0-enterprise-lite":         "1.3.0-0-enterprise-lite",
	}
	for inputVersion, expectedVersion := range validVersions {
		v, err := ParseSemanticVersion(inputVersion)
		if err != nil {
			t.Errorf("error converting %s: %v", inputVersion, err)
		} else if v.String() != expectedVersion {
			t.Errorf("converting %s, expecting %s, getting %s", inputVersion, expectedVersion, v.String())
		}
	}

	invalidVersions := []string{
		"",
		"0-1-1",
	}
	for _, inputVersion := range invalidVersions {
		_, err := ParseSemanticVersion(inputVersion)
		if err == nil {
			t.Errorf("expecting error converting %s, getting no errors", inputVersion)
		}
	}
}

func Test_requestWithHeaders(t *testing.T) {
	type args struct {
		req     *http.Request
		headers http.Header
	}
	tests := []struct {
		name string
		args args
		want *http.Request
	}{
		{
			name: "returns request as is if no headers are set",
			args: args{
				req: &http.Request{
					Method: "GET",
					Header: http.Header{
						"foo": []string{"bar", "baz"},
					},
				},
				headers: http.Header{},
			},
			want: &http.Request{
				Method: "GET",
				Header: http.Header{
					"foo": []string{"bar", "baz"},
				},
			},
		},
		{
			name: "returns request with headers added",
			args: args{
				req: &http.Request{
					Method: "GET",
					Header: http.Header{
						"foo": []string{"bar", "baz"},
					},
				},
				headers: http.Header{
					"password": []string{"my-secret-key"},
					"key-with": []string{"multiple", "values"},
				},
			},
			want: &http.Request{
				Method: "GET",
				Header: http.Header{
					"foo":      []string{"bar", "baz"},
					"Password": []string{"my-secret-key"},
					"Key-With": []string{"multiple", "values"},
				},
			},
		},
		{
			name: "returns nil when input request is nil",
			args: args{
				req: nil,
				headers: http.Header{
					"password": []string{"my-secret-key"},
					"key-with": []string{"multiple", "values"},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := requestWithHeaders(tt.args.req, tt.args.headers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("requestWithHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFillRoutesDefaults(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	tests := []struct {
		name     string
		route    *Route
		expected *Route
	}{
		{
			name: "fills defaults for all fields except paths, leaves name unchanged",
			route: &Route{
				Name:  String("r1"),
				Paths: []*string{String("/r1")},
			},
			expected: &Route{
				Name:                    String("r1"),
				Paths:                   []*string{String("/r1")},
				PreserveHost:            Bool(false),
				Protocols:               []*string{String("http"), String("https")},
				RegexPriority:           Int(0),
				StripPath:               Bool(true),
				HTTPSRedirectStatusCode: Int(426),
			},
		},
		{
			name: "fills defaults for all fields except paths and protocols, leaves name unchanged",
			route: &Route{
				Name:      String("r1"),
				Paths:     []*string{String("/r1")},
				Protocols: []*string{String("grpc")},
			},
			expected: &Route{
				Name:                    String("r1"),
				Paths:                   []*string{String("/r1")},
				PreserveHost:            Bool(false),
				Protocols:               []*string{String("grpc")},
				RegexPriority:           Int(0),
				StripPath:               Bool(true),
				HTTPSRedirectStatusCode: Int(426),
			},
		},
	}

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			r := tc.route
			fullSchema, err := client.Schemas.Get(defaultCtx, "routes")
			assert.Nil(err)
			assert.NotNil(fullSchema)
			if err = FillEntityDefaults(r, fullSchema); err != nil {
				t.Errorf(err.Error())
			}
			// Ignore fields to make tests pass despite small differences across releases.
			opts := cmpopts.IgnoreFields(
				Route{},
				"RequestBuffering", "ResponseBuffering", "PathHandling",
			)
			if diff := cmp.Diff(r, tc.expected, opts); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestFillServiceDefaults(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	tests := []struct {
		name     string
		service  *Service
		expected *Service
	}{
		{
			name: "fills defaults for all fields, leaves name and host unchanged",
			service: &Service{
				Name: String("svc1"),
				Host: String("mockbin.org"),
			},
			expected: &Service{
				Name:           String("svc1"),
				Host:           String("mockbin.org"),
				Port:           Int(80),
				Protocol:       String("http"),
				ConnectTimeout: Int(60000),
				ReadTimeout:    Int(60000),
				Retries:        Int(5),
				WriteTimeout:   Int(60000),
			},
		},
		{
			name: "fills defaults for all fields except port, leaves name and host unchanged",
			service: &Service{
				Name: String("svc1"),
				Host: String("mockbin.org"),
				Port: Int(8080),
			},
			expected: &Service{
				Name:           String("svc1"),
				Host:           String("mockbin.org"),
				Port:           Int(8080),
				Protocol:       String("http"),
				ConnectTimeout: Int(60000),
				ReadTimeout:    Int(60000),
				Retries:        Int(5),
				WriteTimeout:   Int(60000),
			},
		},
		{
			name: "fills defaults for all fields except port, leaves name, tags and host unchanged",
			service: &Service{
				Name: String("svc1"),
				Host: String("mockbin.org"),
				Port: Int(8080),
				Tags: []*string{String("tag1"), String("tag2")},
			},
			expected: &Service{
				Name:           String("svc1"),
				Host:           String("mockbin.org"),
				Port:           Int(8080),
				Protocol:       String("http"),
				ConnectTimeout: Int(60000),
				ReadTimeout:    Int(60000),
				Retries:        Int(5),
				WriteTimeout:   Int(60000),
				Tags:           []*string{String("tag1"), String("tag2")},
			},
		},
	}

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			s := tc.service
			fullSchema, err := client.Schemas.Get(defaultCtx, "services")
			assert.Nil(err)
			assert.NotNil(fullSchema)
			if err := FillEntityDefaults(s, fullSchema); err != nil {
				t.Errorf(err.Error())
			}
			opt := []cmp.Option{
				cmpopts.IgnoreFields(Service{}, "Enabled"),
			}
			if diff := cmp.Diff(s, tc.expected, opt...); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestFillTargetDefaults(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	tests := []struct {
		name     string
		target   *Target
		expected *Target
	}{
		{
			name:   "fills default for weight",
			target: &Target{},
			expected: &Target{
				Weight: Int(100),
			},
		},
	}

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			target := tc.target
			fullSchema, err := client.Schemas.Get(defaultCtx, "targets")
			assert.Nil(err)
			assert.NotNil(fullSchema)
			if err := FillEntityDefaults(target, fullSchema); err != nil {
				t.Errorf(err.Error())
			}
			if diff := cmp.Diff(target, tc.expected); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func TestFillUpstreamsDefaults(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	tests := []struct {
		name     string
		upstream *Upstream
		expected *Upstream
	}{
		{
			name: "fills defaults for all fields, leaves name unchanged",
			upstream: &Upstream{
				Name: String("upstream1"),
			},
			expected: &Upstream{
				Name:      String("upstream1"),
				Algorithm: String("round-robin"),
				Slots:     Int(10000),
				Healthchecks: &Healthcheck{
					Active: &ActiveHealthcheck{
						Concurrency: Int(10),
						Healthy: &Healthy{
							HTTPStatuses: []int{200, 302},
							Interval:     Int(0),
							Successes:    Int(0),
						},
						HTTPPath:               String("/"),
						HTTPSVerifyCertificate: Bool(true),
						Type:                   String("http"),
						Timeout:                Int(1),
						Unhealthy: &Unhealthy{
							HTTPFailures: Int(0),
							HTTPStatuses: []int{
								429, 404,
								500, 501, 502, 503, 504, 505,
							},
							TCPFailures: Int(0),
							Timeouts:    Int(0),
							Interval:    Int(0),
						},
					},
					Passive: &PassiveHealthcheck{
						Healthy: &Healthy{
							HTTPStatuses: []int{
								200, 201, 202, 203, 204, 205, 206, 207, 208, 226,
								300, 301, 302, 303, 304, 305, 306, 307, 308,
							},
							Successes: Int(0),
						},
						Type: String("http"),
						Unhealthy: &Unhealthy{
							HTTPFailures: Int(0),
							HTTPStatuses: []int{429, 500, 503},
							TCPFailures:  Int(0),
							Timeouts:     Int(0),
						},
					},
				},
				HashOn:           String("none"),
				HashFallback:     String("none"),
				HashOnCookiePath: String("/"),
			},
		},
		{
			name: "fills defaults for all fields except algorithm and hash_on, leaves name unchanged",
			upstream: &Upstream{
				Name:      String("upstream1"),
				Algorithm: String("consistent-hashing"),
				HashOn:    String("ip"),
			},
			expected: &Upstream{
				Name:      String("upstream1"),
				Algorithm: String("consistent-hashing"),
				Slots:     Int(10000),
				Healthchecks: &Healthcheck{
					Active: &ActiveHealthcheck{
						Concurrency: Int(10),
						Healthy: &Healthy{
							HTTPStatuses: []int{200, 302},
							Interval:     Int(0),
							Successes:    Int(0),
						},
						HTTPPath:               String("/"),
						HTTPSVerifyCertificate: Bool(true),
						Type:                   String("http"),
						Timeout:                Int(1),
						Unhealthy: &Unhealthy{
							HTTPFailures: Int(0),
							HTTPStatuses: []int{
								429, 404,
								500, 501, 502, 503, 504, 505,
							},
							TCPFailures: Int(0),
							Timeouts:    Int(0),
							Interval:    Int(0),
						},
					},
					Passive: &PassiveHealthcheck{
						Healthy: &Healthy{
							HTTPStatuses: []int{
								200, 201, 202, 203, 204, 205, 206, 207, 208, 226,
								300, 301, 302, 303, 304, 305, 306, 307, 308,
							},
							Successes: Int(0),
						},
						Type: String("http"),
						Unhealthy: &Unhealthy{
							HTTPFailures: Int(0),
							HTTPStatuses: []int{429, 500, 503},
							TCPFailures:  Int(0),
							Timeouts:     Int(0),
						},
					},
				},
				HashOn:           String("ip"),
				HashFallback:     String("none"),
				HashOnCookiePath: String("/"),
			},
		},
	}

	for _, tc := range tests {
		T.Run(tc.name, func(t *testing.T) {
			u := tc.upstream
			fullSchema, err := client.Schemas.Get(defaultCtx, "upstreams")
			assert.Nil(err)
			assert.NotNil(fullSchema)
			if err = FillEntityDefaults(u, fullSchema); err != nil {
				t.Errorf(err.Error())
			}
			// Ignore fields to make tests pass despite small differences across releases.
			opts := cmpopts.IgnoreFields(Healthcheck{}, "Threshold")
			if diff := cmp.Diff(u, tc.expected, opts); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
