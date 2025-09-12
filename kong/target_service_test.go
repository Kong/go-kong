package kong

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTargetsUpstream(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	target := &Target{
		Target: String("10.0.0.1"),
	}

	// upstream is required
	badTarget, err := client.Targets.Create(defaultCtx, nil, target)
	require.Error(err)
	assert.Nil(badTarget)

	// create a upstream
	fixtureUpstream, err := client.Upstreams.Create(defaultCtx, &Upstream{
		Name: String("vhost.com"),
	})
	require.NoError(err)
	require.NotNil(fixtureUpstream)
	assert.NotNil(fixtureUpstream.ID)

	createdTarget, err := client.Targets.Create(defaultCtx,
		fixtureUpstream.ID, &Target{
			Target: String("10.0.0.1:80"),
		})
	require.NoError(err)
	assert.NotNil(createdTarget)

	err = client.Targets.Delete(defaultCtx, fixtureUpstream.ID,
		createdTarget.ID)
	require.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	target = &Target{
		ID:       String(id),
		Target:   String("10.0.0.3"),
		Upstream: fixtureUpstream,
	}

	createdTarget, err = client.Targets.Create(defaultCtx,
		fixtureUpstream.ID, target)
	require.NoError(err)
	assert.NotNil(createdTarget)
	assert.Equal(id, *createdTarget.ID)

	err = client.Upstreams.Delete(defaultCtx, fixtureUpstream.ID)
	require.NoError(err)
}

func TestTargetsUpdate(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	// create a upstream
	fixtureUpstream, err := client.Upstreams.Create(defaultCtx, &Upstream{
		Name: String("vhost.com"),
	})
	T.Cleanup(func() {
		err = client.Upstreams.Delete(defaultCtx, fixtureUpstream.ID)
		require.NoError(err)
	})
	require.NoError(err)
	require.NotNil(fixtureUpstream)
	require.NotNil(fixtureUpstream.ID)

	targetID := "0fa49cd2-ee93-492a-bedf-b80778d539ae"
	createdTarget, err := client.Targets.Create(defaultCtx,
		fixtureUpstream.ID, &Target{
			ID:     &targetID,
			Target: String("10.0.0.1:80"),
		})
	require.NoError(err)
	require.NotNil(createdTarget)
	require.NotNil(createdTarget.ID)
	assert.Equal(targetID, *createdTarget.ID)

	err = client.Targets.Delete(defaultCtx, fixtureUpstream.ID,
		createdTarget.ID)
	require.NoError(err)

	createdTarget, err = client.Targets.Create(defaultCtx,
		fixtureUpstream.ID, &Target{
			ID:     &targetID,
			Weight: Int(1),
			Target: String("10.0.0.2:80"),
		})
	T.Cleanup(func() {
		err = client.Targets.Delete(defaultCtx, fixtureUpstream.ID, createdTarget.ID)
		require.NoError(err)
	})
	require.NoError(err)
	require.NotNil(createdTarget)
	require.NotNil(createdTarget.ID)
	require.NotNil(createdTarget.Target)
	require.NotNil(createdTarget.Weight)
	assert.Equal(targetID, *createdTarget.ID)
	assert.Equal(1, *createdTarget.Weight)
	assert.Equal("10.0.0.2:80", *createdTarget.Target)

	updatedTarget, err := client.Targets.Create(defaultCtx,
		fixtureUpstream.ID, &Target{
			ID:     &targetID,
			Weight: Int(2),
			Target: String("10.0.0.3:80"),
		})
	require.NoError(err)
	require.NotNil(updatedTarget)
	require.NotNil(updatedTarget.ID)
	require.NotNil(updatedTarget.Target)
	require.NotNil(updatedTarget.Weight)
	assert.Equal(targetID, *updatedTarget.ID)
	assert.Equal(2, *updatedTarget.Weight)
	assert.Equal("10.0.0.3:80", *updatedTarget.Target)
}

func TestTargetWithTags(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenKong(T, ">=1.1.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	fixtureUpstream, err := client.Upstreams.Create(defaultCtx, &Upstream{
		Name: String("vhost.com"),
	})
	require.NoError(err)

	createdTarget, err := client.Targets.Create(defaultCtx,
		fixtureUpstream.ID, &Target{
			Target: String("10.0.0.1:80"),
			Tags:   StringSlice("tag1", "tag2"),
		})
	require.NoError(err)
	require.NotNil(createdTarget)
	assert.Equal(StringSlice("tag1", "tag2"), createdTarget.Tags)

	err = client.Upstreams.Delete(defaultCtx, fixtureUpstream.ID)
	require.NoError(err)
}

func TestTargetWithFailover(T *testing.T) {
	RunWhenKong(T, ">=3.12.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	fixtureUpstream, err := client.Upstreams.Create(defaultCtx, &Upstream{
		Name: String("vhost.com"),
	})
	require.NoError(err)

	createdTarget, err := client.Targets.Create(defaultCtx,
		fixtureUpstream.ID, &Target{
			Target:   String("10.0.0.1:80"),
			Failover: Bool(true),
		})
	require.NoError(err)
	require.NotNil(createdTarget)
	assert.True(*createdTarget.Failover)

	err = client.Upstreams.Delete(defaultCtx, fixtureUpstream.ID)
	require.NoError(err)
}

func TestTargetWithFailoverDefault(T *testing.T) {
	RunWhenKong(T, ">=3.12.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	fixtureUpstream, err := client.Upstreams.Create(defaultCtx, &Upstream{
		Name: String("vhost.com"),
	})
	require.NoError(err)

	createdTarget, err := client.Targets.Create(defaultCtx,
		fixtureUpstream.ID, &Target{
			Target: String("10.0.0.1:80"),
		})
	require.NoError(err)
	require.NotNil(createdTarget)
	assert.False(*createdTarget.Failover)

	err = client.Upstreams.Delete(defaultCtx, fixtureUpstream.ID)
	require.NoError(err)
}

func TestTargetListEndpoint(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	upstream := &Upstream{
		Name: String("vhost2.com"),
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	require.NoError(err)
	require.NotNil(createdUpstream)

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
		require.NoError(err)
		assert.NotNil(target)
		targets[i] = target
	}

	targetsFromKong, next, err := client.Targets.List(defaultCtx,
		createdUpstream.ID, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(targetsFromKong)
	assert.Len(targetsFromKong, 3)

	// check if we see all targets
	assert.True(compareTargets(targets, targetsFromKong))

	// Test pagination
	targetsFromKong = []*Target{}

	// first page
	page1, next, err := client.Targets.List(defaultCtx,
		createdUpstream.ID, &ListOpt{Size: 1})
	require.NoError(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)
	targetsFromKong = append(targetsFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.Targets.List(defaultCtx,
		createdUpstream.ID, next)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 2)
	targetsFromKong = append(targetsFromKong, page2...)

	assert.True(compareTargets(targets, targetsFromKong))

	targets, err = client.Targets.ListAll(defaultCtx, createdUpstream.ID)
	require.NoError(err)
	assert.NotNil(targets)
	assert.Len(targets, 3)

	require.NoError(client.Upstreams.Delete(defaultCtx, createdUpstream.ID))
}

func TestTargetsUpdatePatch(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)
	client.SetDebugMode(true)

	// create a upstream
	fixtureUpstream, err := client.Upstreams.Create(defaultCtx, &Upstream{
		Name: String("vhost.com"),
	})
	require.NoError(err)
	require.NotNil(fixtureUpstream)
	require.NotNil(fixtureUpstream.ID)

	targetID := "0fa49cd2-ee93-492a-bedf-b80778d539ae"
	createdTarget, err := client.Targets.Create(defaultCtx,
		fixtureUpstream.ID, &Target{
			ID:     &targetID,
			Weight: Int(100),
			Target: String("10.0.0.1:80"),
		})
	T.Cleanup(func() {
		err = client.Upstreams.Delete(defaultCtx, fixtureUpstream.ID)
		require.NoError(err)
	})
	require.NoError(err)
	require.NotNil(createdTarget)
	require.NotNil(createdTarget.ID)
	require.NotNil(createdTarget.Weight)
	assert.Equal(targetID, *createdTarget.ID)
	assert.Equal(100, *createdTarget.Weight)

	updatedTarget, err := client.Targets.Update(defaultCtx,
		fixtureUpstream.ID, createdTarget.ID, &Target{
			ID:     createdTarget.ID,
			Weight: Int(10000), // Update weight
			Target: createdTarget.Target,
		})
	require.NoError(err)
	require.NotNil(updatedTarget)
	require.NotNil(updatedTarget.ID)
	require.NotNil(updatedTarget.Target)
	require.NotNil(updatedTarget.Weight)
	assert.Equal(targetID, *updatedTarget.ID)
	assert.Equal("10.0.0.1:80", *updatedTarget.Target)
	assert.Equal(10000, *updatedTarget.Weight)
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
	RunWhenDBMode(T, "postgres")

	// TODO https://github.com/Kong/go-kong/issues/213 this does not yet work on 3.x
	RunWhenKong(T, "<3.0.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
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
	require.NoError(err)
	require.NotNil(createdUpstream)

	createdTarget, err := client.Targets.Create(defaultCtx,
		createdUpstream.ID, &Target{
			Target: String("10.0.0.1:80"),
		})
	require.NoError(err)
	assert.NotNil(createdTarget)

	require.Error(client.Targets.MarkHealthy(defaultCtx, createdTarget.Upstream.ID, nil))
	require.Error(client.Targets.MarkHealthy(defaultCtx, nil, createdTarget))
	assert.Eventually(func() bool {
		err := client.Targets.MarkHealthy(defaultCtx, createdTarget.Upstream.ID, createdTarget)
		if err != nil {
			T.Logf("failed marking target %s healthy", *createdTarget.ID)
			return false
		}
		return true
	}, 5*time.Second, 100*time.Millisecond)

	require.NoError(client.Upstreams.Delete(defaultCtx, createdUpstream.ID))
}

func TestTargetMarkUnhealthy(T *testing.T) {
	RunWhenDBMode(T, "postgres")

	// TODO https://github.com/Kong/go-kong/issues/213 this does not yet work on 3.x
	RunWhenKong(T, "<3.0.0")
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

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
	require.NoError(err)
	require.NotNil(createdUpstream)

	createdTarget, err := client.Targets.Create(defaultCtx,
		createdUpstream.ID, &Target{
			Target: String("10.0.0.1:80"),
		})
	require.NoError(err)
	require.NotNil(createdTarget)

	require.Error(client.Targets.MarkUnhealthy(defaultCtx, createdTarget.Upstream.ID, nil))
	require.Error(client.Targets.MarkUnhealthy(defaultCtx, nil, createdTarget))

	assert.Eventually(func() bool {
		err := client.Targets.MarkUnhealthy(defaultCtx, createdTarget.Upstream.ID, createdTarget)
		if err != nil {
			T.Logf("failed marking target %s healthy", *createdTarget.ID)
			return false
		}
		return true
	}, 5*time.Second, 100*time.Millisecond)

	require.NoError(client.Upstreams.Delete(defaultCtx, createdUpstream.ID))
}
