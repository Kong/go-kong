package kong

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpstreamsService(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	upstream := &Upstream{
		Name: String("virtual-host1"),
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	assert.Nil(err)
	assert.NotNil(createdUpstream)

	upstream, err = client.Upstreams.Get(defaultCtx, createdUpstream.ID)
	assert.Nil(err)
	assert.NotNil(upstream)

	upstream.Name = String("virtual-host2")
	upstream, err = client.Upstreams.Update(defaultCtx, upstream)
	assert.Nil(err)
	assert.NotNil(upstream)
	assert.Equal("virtual-host2", *upstream.Name)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.Nil(err)

	// ID can be specified
	id := uuid.NewString()
	upstream = &Upstream{
		Name: String("key-auth"),
		ID:   String(id),
	}

	createdUpstream, err = client.Upstreams.Create(defaultCtx, upstream)
	assert.Nil(err)
	assert.NotNil(createdUpstream)
	assert.Equal(id, *createdUpstream.ID)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.Nil(err)
}

func TestUpstreamWithTags(T *testing.T) {
	runWhenKong(T, ">=1.1.0")
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	upstream := &Upstream{
		Name: String("key-auth"),
		Tags: StringSlice("tag1", "tag2"),
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	assert.Nil(err)
	assert.NotNil(createdUpstream)
	assert.Equal(StringSlice("tag1", "tag2"), createdUpstream.Tags)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.Nil(err)
}

// regression test for #6
func TestUpstreamWithActiveUnHealthyInterval(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	upstream := &Upstream{
		Name: String("upstream-foo"),
		Healthchecks: &Healthcheck{
			Active: &ActiveHealthcheck{
				Unhealthy: &Unhealthy{
					Interval: Int(5),
				},
			},
		},
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	assert.Nil(err)
	assert.NotNil(createdUpstream)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.Nil(err)
}

// regression test for #6
func TestUpstreamWithPassiveUnHealthyInterval(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	upstream := &Upstream{
		Name: String("upstream-foo"),
		Healthchecks: &Healthcheck{
			Passive: &PassiveHealthcheck{
				Unhealthy: &Unhealthy{
					Interval: Int(5),
				},
			},
		},
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	assert.NotNil(err)
	assert.Nil(createdUpstream)
}

func TestUpstreamWithPassiveHealthy(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	upstream := &Upstream{
		Name: String("upstream-foo"),
		Healthchecks: &Healthcheck{
			Passive: &PassiveHealthcheck{
				Type: String("http"),
				Healthy: &Healthy{
					HTTPStatuses: []int{200, 201},
					Successes:    Int(3),
				},
			},
		},
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	assert.Nil(err)
	assert.NotNil(createdUpstream)
	assert.Equal("http", *createdUpstream.Healthchecks.Passive.Type)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.Nil(err)
}

func TestUpstreamWithAlgorithm(T *testing.T) {
	runWhenKong(T, ">=1.3.0")
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	upstream := &Upstream{
		Name:      String("upstream1"),
		Algorithm: String("least-connections"),
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	assert.Nil(err)
	assert.NotNil(createdUpstream)
	assert.Equal("least-connections", *createdUpstream.Algorithm)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.Nil(err)
}

func TestUpstreamListEndpoint(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	// fixtures
	upstreams := []*Upstream{
		{
			Name: String("vhost1.com"),
		},
		{
			Name: String("vhost2.com"),
		},
		{
			Name: String("vhost3.com"),
		},
	}

	// create fixturs
	for i := 0; i < len(upstreams); i++ {
		upstream, err := client.Upstreams.Create(defaultCtx, upstreams[i])
		assert.Nil(err)
		assert.NotNil(upstream)
		upstreams[i] = upstream
	}

	upstreamsFromKong, next, err := client.Upstreams.List(defaultCtx, nil)
	assert.Nil(err)
	assert.Nil(next)
	assert.NotNil(upstreamsFromKong)
	assert.Equal(3, len(upstreamsFromKong))

	// check if we see all upstreams
	assert.True(compareUpstreams(upstreams, upstreamsFromKong))

	// Test pagination
	upstreamsFromKong = []*Upstream{}

	// first page
	page1, next, err := client.Upstreams.List(defaultCtx, &ListOpt{Size: 1})
	assert.Nil(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	upstreamsFromKong = append(upstreamsFromKong, page1...)

	// second page
	page2, next, err := client.Upstreams.List(defaultCtx, next)
	assert.Nil(err)
	assert.NotNil(next)
	assert.NotNil(page2)
	assert.Equal(1, len(page2))
	upstreamsFromKong = append(upstreamsFromKong, page2...)

	// last page
	page3, next, err := client.Upstreams.List(defaultCtx, next)
	assert.Nil(err)
	assert.Nil(next)
	assert.NotNil(page3)
	assert.Equal(1, len(page3))
	upstreamsFromKong = append(upstreamsFromKong, page3...)

	assert.True(compareUpstreams(upstreams, upstreamsFromKong))

	upstreams, err = client.Upstreams.ListAll(defaultCtx)
	assert.Nil(err)
	assert.NotNil(upstreams)
	assert.Equal(3, len(upstreams))

	for i := 0; i < len(upstreams); i++ {
		assert.Nil(client.Upstreams.Delete(defaultCtx, upstreams[i].ID))
	}
}

func compareUpstreams(expected, actual []*Upstream) bool {
	var expectedNames, actualNames []string
	for _, upstream := range expected {
		expectedNames = append(expectedNames, *upstream.Name)
	}

	for _, upstream := range actual {
		actualNames = append(actualNames, *upstream.Name)
	}

	return (compareSlices(expectedNames, actualNames))
}

func TestUpstreamsWithHostHeader(T *testing.T) {
	runWhenKong(T, ">=1.4.0")
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	upstream := &Upstream{
		Name:       String("upstream-with-host-header"),
		HostHeader: String("example.com"),
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	assert.Nil(err)
	assert.NotNil(createdUpstream)
	assert.Equal("example.com", *createdUpstream.HostHeader)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.Nil(err)
}

func TestUpstreamGetFullSchema(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	schema, err := client.Upstreams.GetFullSchema(defaultCtx)
	_, ok := schema["fields"]
	assert.True(ok)
	assert.Nil(err)
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
			name: "name only",
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
			name: "name algorithm hashon",
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
			fullSchema, err := client.Upstreams.GetFullSchema(defaultCtx)
			assert.Nil(err)
			assert.NotNil(fullSchema)
			if err = FillUpstreamsDefaults(u, fullSchema); err != nil {
				t.Errorf(err.Error())
			}
			// Ignore fields to make tests pass despite small differences
			// across EE releases.
			opts := cmpopts.IgnoreFields(Healthcheck{}, "Threshold")
			if diff := cmp.Diff(u, tc.expected, opts); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}
