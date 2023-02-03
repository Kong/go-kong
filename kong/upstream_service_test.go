package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpstreamsService(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	upstream := &Upstream{
		Name: String("virtual-host1"),
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	assert.NoError(err)
	require.NotNil(createdUpstream)

	upstream, err = client.Upstreams.Get(defaultCtx, createdUpstream.ID)
	assert.NoError(err)
	assert.NotNil(upstream)

	upstream.Name = String("virtual-host2")
	upstream, err = client.Upstreams.Update(defaultCtx, upstream)
	assert.NoError(err)
	assert.NotNil(upstream)
	assert.Equal("virtual-host2", *upstream.Name)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	upstream = &Upstream{
		Name: String("key-auth"),
		ID:   String(id),
	}

	createdUpstream, err = client.Upstreams.Create(defaultCtx, upstream)
	assert.NoError(err)
	assert.NotNil(createdUpstream)
	assert.Equal(id, *createdUpstream.ID)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.NoError(err)
}

func TestUpstreamWithTags(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=1.1.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	upstream := &Upstream{
		Name: String("key-auth"),
		Tags: StringSlice("tag1", "tag2"),
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	require.NoError(err)
	require.NotNil(createdUpstream)
	assert.Equal(StringSlice("tag1", "tag2"), createdUpstream.Tags)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.NoError(err)
}

// regression test for #6
func TestUpstreamWithActiveUnHealthyInterval(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(createdUpstream)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.NoError(err)
}

// regression test for #6
func TestUpstreamWithPassiveUnHealthyInterval(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
	assert.NoError(err)
	require.NotNil(createdUpstream)
	assert.Equal("http", *createdUpstream.Healthchecks.Passive.Type)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.NoError(err)
}

func TestUpstreamWithAlgorithm(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=1.3.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	upstream := &Upstream{
		Name:      String("upstream1"),
		Algorithm: String("least-connections"),
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	require.NoError(err)
	require.NotNil(createdUpstream)
	assert.Equal("least-connections", *createdUpstream.Algorithm)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.NoError(err)
}

func TestUpstreamListEndpoint(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
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
		assert.NoError(err)
		assert.NotNil(upstream)
		upstreams[i] = upstream
	}

	upstreamsFromKong, next, err := client.Upstreams.List(defaultCtx, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(upstreamsFromKong)
	assert.Equal(3, len(upstreamsFromKong))

	// check if we see all upstreams
	assert.True(compareUpstreams(T, upstreams, upstreamsFromKong))

	// Test pagination
	upstreamsFromKong = []*Upstream{}

	// first page
	page1, next, err := client.Upstreams.List(defaultCtx, &ListOpt{Size: 1})
	assert.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	upstreamsFromKong = append(upstreamsFromKong, page1...)

	// second page
	page2, next, err := client.Upstreams.List(defaultCtx, next)
	assert.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page2)
	assert.Equal(1, len(page2))
	upstreamsFromKong = append(upstreamsFromKong, page2...)

	// last page
	page3, next, err := client.Upstreams.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(page3)
	assert.Equal(1, len(page3))
	upstreamsFromKong = append(upstreamsFromKong, page3...)

	assert.True(compareUpstreams(T, upstreams, upstreamsFromKong))

	upstreams, err = client.Upstreams.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(upstreams)
	assert.Equal(3, len(upstreams))

	for i := 0; i < len(upstreams); i++ {
		assert.NoError(client.Upstreams.Delete(defaultCtx, upstreams[i].ID))
	}
}

func compareUpstreams(T *testing.T, expected, actual []*Upstream) bool {
	var expectedNames, actualNames []string
	for _, upstream := range expected {
		if !assert.NotNil(T, upstream) {
			continue
		}
		expectedNames = append(expectedNames, *upstream.Name)
	}

	for _, upstream := range actual {
		actualNames = append(actualNames, *upstream.Name)
	}

	return (compareSlices(expectedNames, actualNames))
}

func TestUpstreamsWithHostHeader(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=1.4.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	upstream := &Upstream{
		Name:       String("upstream-with-host-header"),
		HostHeader: String("example.com"),
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	require.NoError(err)
	require.NotNil(createdUpstream)
	assert.Equal("example.com", *createdUpstream.HostHeader)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.NoError(err)
}
