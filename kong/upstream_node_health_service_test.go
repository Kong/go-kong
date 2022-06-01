package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpstreamNodeHealthService(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	// create a upstream
	fixtureUpstream, err := client.Upstreams.Create(
		defaultCtx,
		&Upstream{
			Name: String("vhost.com"),
		},
	)
	assert.NoError(err)
	require.NotNil(fixtureUpstream)
	assert.NotNil(fixtureUpstream.ID)

	// create a target
	createdTarget, err := client.Targets.Create(
		defaultCtx,
		fixtureUpstream.ID,
		&Target{
			Target: String("10.0.0.1:80"),
		},
	)
	assert.NoError(err)
	assert.NotNil(createdTarget)

	// upstream node health
	nodeHealths, err := client.UpstreamNodeHealth.ListAll(
		defaultCtx, fixtureUpstream.ID,
	)
	assert.NoError(err)
	assert.NotNil(nodeHealths)

	// cleanup targets
	err = client.Targets.Delete(
		defaultCtx, fixtureUpstream.ID, createdTarget.ID,
	)
	assert.NoError(err)

	// cleanup upstream
	err = client.Upstreams.Delete(defaultCtx, fixtureUpstream.ID)
	assert.NoError(err)
}
