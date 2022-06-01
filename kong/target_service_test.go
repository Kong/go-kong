package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTargetsUpstream(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	target := &Target{
		Target: String("10.0.0.1"),
	}

	// upstream is required
	badTarget, err := client.Targets.Create(defaultCtx, nil, target)
	assert.NotNil(err)
	assert.Nil(badTarget)

	// create a upstream
	fixtureUpstream, err := client.Upstreams.Create(defaultCtx, &Upstream{
		Name: String("vhost.com"),
	})
	assert.NoError(err)
	assert.NotNil(fixtureUpstream)
	assert.NotNil(fixtureUpstream.ID)

	createdTarget, err := client.Targets.Create(defaultCtx,
		fixtureUpstream.ID, &Target{
			Target: String("10.0.0.1:80"),
		})
	assert.NoError(err)
	assert.NotNil(createdTarget)

	err = client.Targets.Delete(defaultCtx, fixtureUpstream.ID,
		createdTarget.ID)
	assert.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	target = &Target{
		ID:       String(id),
		Target:   String("10.0.0.3"),
		Upstream: fixtureUpstream,
	}

	createdTarget, err = client.Targets.Create(defaultCtx,
		fixtureUpstream.ID, target)
	assert.NoError(err)
	assert.NotNil(createdTarget)
	assert.Equal(id, *createdTarget.ID)

	err = client.Upstreams.Delete(defaultCtx, fixtureUpstream.ID)
	assert.NoError(err)
}

func TestTargetsUpdate(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	// create a upstream
	fixtureUpstream, err := client.Upstreams.Create(defaultCtx, &Upstream{
		Name: String("vhost.com"),
	})
	assert.NoError(err)
	assert.NotNil(fixtureUpstream)
	assert.NotNil(fixtureUpstream.ID)

	targetID := "0fa49cd2-ee93-492a-bedf-b80778d539ae"
	createdTarget, err := client.Targets.Create(defaultCtx,
		fixtureUpstream.ID, &Target{
			ID:     &targetID,
			Target: String("10.0.0.1:80"),
		})
	assert.NoError(err)
	assert.NotNil(createdTarget)
	assert.Equal(targetID, *createdTarget.ID)

	err = client.Targets.Delete(defaultCtx, fixtureUpstream.ID,
		createdTarget.ID)
	assert.NoError(err)

	createdTarget, err = client.Targets.Create(defaultCtx,
		fixtureUpstream.ID, &Target{
			ID:     &targetID,
			Target: String("10.0.0.2:80"),
		})
	assert.NoError(err)
	assert.NotNil(createdTarget)
	assert.Equal(targetID, *createdTarget.ID)
	assert.Equal("10.0.0.2:80", *createdTarget.Target)

	err = client.Upstreams.Delete(defaultCtx, fixtureUpstream.ID)
	assert.NoError(err)
}

func TestTargetWithTags(T *testing.T) {
	RunWhenKong(T, ">=1.1.0")
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	fixtureUpstream, err := client.Upstreams.Create(defaultCtx, &Upstream{
		Name: String("vhost.com"),
	})
	assert.NoError(err)

	createdTarget, err := client.Targets.Create(defaultCtx,
		fixtureUpstream.ID, &Target{
			Target: String("10.0.0.1:80"),
			Tags:   StringSlice("tag1", "tag2"),
		})
	assert.NoError(err)
	assert.NotNil(createdTarget)
	assert.Equal(StringSlice("tag1", "tag2"), createdTarget.Tags)

	err = client.Upstreams.Delete(defaultCtx, fixtureUpstream.ID)
	assert.NoError(err)
}

func TestTargetListEndpoint(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	upstream := &Upstream{
		Name: String("vhost2.com"),
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	assert.NoError(err)
	assert.NotNil(createdUpstream)

	// fixtures
	targets := []*Target{
		{
			Target:   String("10.42.1.2"),
			Upstream: createdUpstream,
		},
		{
			Target:   String("10.42.1.3"),
			Upstream: createdUpstream,
		},
		{
			Target:   String("10.42.1.4"),
			Upstream: createdUpstream,
		},
	}
	// create fixturs
	for i := 0; i < len(targets); i++ {
		target, err := client.Targets.Create(defaultCtx,
			createdUpstream.ID, targets[i])
		assert.NoError(err)
		assert.NotNil(target)
		targets[i] = target
	}

	targetsFromKong, next, err := client.Targets.List(defaultCtx,
		createdUpstream.ID, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(targetsFromKong)
	assert.Equal(3, len(targetsFromKong))

	// check if we see all targets
	assert.True(compareTargets(targets, targetsFromKong))

	// Test pagination
	targetsFromKong = []*Target{}

	// first page
	page1, next, err := client.Targets.List(defaultCtx,
		createdUpstream.ID, &ListOpt{Size: 1})
	assert.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	targetsFromKong = append(targetsFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.Targets.List(defaultCtx,
		createdUpstream.ID, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Equal(2, len(page2))
	targetsFromKong = append(targetsFromKong, page2...)

	assert.True(compareTargets(targets, targetsFromKong))

	targets, err = client.Targets.ListAll(defaultCtx, createdUpstream.ID)
	assert.NoError(err)
	assert.NotNil(targets)
	assert.Equal(3, len(targets))

	assert.NoError(client.Upstreams.Delete(defaultCtx, createdUpstream.ID))
}

func compareTargets(expected, actual []*Target) bool {
	var expectedUsernames, actualUsernames []string
	for _, target := range expected {
		expectedUsernames = append(expectedUsernames, *target.Target)
	}

	for _, target := range actual {
		actualUsernames = append(actualUsernames, *target.Target)
	}

	return (compareSlices(expectedUsernames, actualUsernames))
}

func TestTargetMarkHealthy(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	upstream := &Upstream{
		Name: String("vhost1.com"),
		Healthchecks: &Healthcheck{
			Passive: &PassiveHealthcheck{
				Unhealthy: &Unhealthy{
					HTTPFailures: Int(5),
				},
			},
		},
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	assert.NoError(err)
	assert.NotNil(createdUpstream)

	createdTarget, err := client.Targets.Create(defaultCtx,
		createdUpstream.ID, &Target{
			Target: String("10.0.0.1:80"),
		})
	assert.NoError(err)
	assert.NotNil(createdTarget)

	assert.NotNil(client.Targets.MarkHealthy(defaultCtx, createdTarget.Upstream.ID, nil))
	assert.NotNil(client.Targets.MarkHealthy(defaultCtx, nil, createdTarget))
	assert.NoError(client.Targets.MarkHealthy(defaultCtx, createdTarget.Upstream.ID, createdTarget))

	assert.NoError(client.Upstreams.Delete(defaultCtx, createdUpstream.ID))
}

func TestTargetMarkUnhealthy(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	upstream := &Upstream{
		Name: String("vhost1.com"),
		Healthchecks: &Healthcheck{
			Passive: &PassiveHealthcheck{
				Unhealthy: &Unhealthy{
					HTTPFailures: Int(5),
				},
			},
		},
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	assert.NoError(err)
	assert.NotNil(createdUpstream)

	createdTarget, err := client.Targets.Create(defaultCtx,
		createdUpstream.ID, &Target{
			Target: String("10.0.0.1:80"),
		})
	assert.NoError(err)
	assert.NotNil(createdTarget)

	assert.NotNil(client.Targets.MarkUnhealthy(defaultCtx, createdTarget.Upstream.ID, nil))
	assert.NotNil(client.Targets.MarkUnhealthy(defaultCtx, nil, createdTarget))
	assert.NoError(client.Targets.MarkUnhealthy(defaultCtx, createdTarget.Upstream.ID, createdTarget))

	assert.NoError(client.Upstreams.Delete(defaultCtx, createdUpstream.ID))
}
