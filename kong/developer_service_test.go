package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDevelopersService(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{Portal: true})
	// NOTE: Developer Portal is not available in Kong < 3.5.0. Requires special config/license to enable.
	RunWhenEnterprise(T, "<3.5.0", RequiredFeatures{Portal: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	testWs, err := NewTestWorkspace(client, "default")
	require.NoError(T, err)
	require.NoError(T, testWs.UpdateConfig(map[string]interface{}{
		"portal_auth":         "basic-auth",
		"portal_session_conf": map[string]interface{}{"secret": "garbage"},
		"portal":              true,
	}))

	developer := &Developer{
		Meta:     String("{\"full_name\": \"Foo BAR\"}"),
		Email:    String("foo.bar@example.com"),
		CustomID: String("custom_id_foo"),
		Password: String("F00P@ssword"),
	}

	createdDeveloper, err := client.Developers.Create(defaultCtx, developer)
	require.NoError(T, err)
	assert.NotNil(createdDeveloper)

	developer, err = client.Developers.Get(defaultCtx, createdDeveloper.ID)
	require.NoError(T, err)
	assert.NotNil(developer)

	developer, err = client.Developers.GetByCustomID(defaultCtx,
		String("does-not-exist"))
	require.Error(T, err)
	assert.Nil(developer)

	developer, err = client.Developers.GetByCustomID(defaultCtx,
		String("custom_id_foo"))
	require.NoError(T, err)
	assert.NotNil(developer)

	developer.Email = String("bar@example.com")
	developer, err = client.Developers.Update(defaultCtx, developer)
	require.NoError(T, err)
	assert.NotNil(developer)
	assert.Equal("bar@example.com", *developer.Email)

	err = client.Developers.Delete(defaultCtx, createdDeveloper.ID)
	require.NoError(T, err)

	// ID can be specified
	id := uuid.NewString()
	developer = &Developer{
		Meta:     String("{\"full_name\": \"Foo BAR\"}"),
		Email:    String("foo.bar@example.com"),
		Password: String("F00P@ssword"),
		ID:       String(id),
	}

	createdDeveloper, err = client.Developers.Create(defaultCtx, developer)
	require.NoError(T, err)
	assert.NotNil(createdDeveloper)
	assert.Equal(id, *createdDeveloper.ID)

	err = client.Developers.Delete(defaultCtx, createdDeveloper.ID)
	require.NoError(T, err)

	require.NoError(T, testWs.Reset())
}

func TestDeveloperListEndpoint(T *testing.T) {
	RunWhenEnterprise(T, ">=0.33.0", RequiredFeatures{Portal: true})
	// NOTE: Developer Portal is not available in Kong < 3.5.0. Requires special config/license to enable.
	RunWhenEnterprise(T, "<3.5.0", RequiredFeatures{Portal: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(T, err)
	assert.NotNil(client)

	testWs, err := NewTestWorkspace(client, "default")
	require.NoError(T, err)
	require.NoError(T, testWs.UpdateConfig(map[string]interface{}{
		"portal_auth":         "basic-auth",
		"portal_session_conf": map[string]interface{}{"secret": "garbage"},
		"portal":              true,
	}))

	// fixtures
	developers := []*Developer{
		{
			Email:    String("foo1.bar@example.com"),
			Meta:     String("{\"full_name\": \"Foo1 BAR\"}"),
			Password: String("F001P@ssword"),
		},
		{
			Email:    String("foo2@example.com"),
			Meta:     String("{\"full_name\": \"Foo2 BAR\"}"),
			Password: String("F002P@ssword"),
		},
		{
			Email:    String("foo3@example.com"),
			Meta:     String("{\"full_name\": \"Foo3 BAR\"}"),
			Password: String("F003P@ssword"),
		},
	}

	// create fixturs
	for i := 0; i < len(developers); i++ {
		developer, err := client.Developers.Create(defaultCtx, developers[i])
		require.NoError(T, err)
		assert.NotNil(developer)
		developers[i] = developer
	}

	developersFromKong, next, err := client.Developers.List(defaultCtx, nil)
	require.NoError(T, err)
	assert.Nil(next)
	assert.NotNil(developersFromKong)
	assert.Len(developersFromKong, 3)

	// check if we see all developers
	assert.True(compareDevelopers(developers, developersFromKong))

	// Test pagination
	developersFromKong = []*Developer{}

	// first page
	page1, next, err := client.Developers.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(T, err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)
	developersFromKong = append(developersFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.Developers.List(defaultCtx, next)
	require.NoError(T, err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 2)
	developersFromKong = append(developersFromKong, page2...)

	assert.True(compareDevelopers(developers, developersFromKong))

	developers, err = client.Developers.ListAll(defaultCtx)
	require.NoError(T, err)
	assert.NotNil(developers)
	assert.Len(developers, 3)

	for i := 0; i < len(developers); i++ {
		require.NoError(T, client.Developers.Delete(defaultCtx, developers[i].ID))
	}

	require.NoError(T, testWs.Reset())
}

func compareDevelopers(expected, actual []*Developer) bool {
	var expectedEmails, actualEmails []string
	for _, developer := range expected {
		expectedEmails = append(expectedEmails, *developer.Email)
	}

	for _, developer := range actual {
		actualEmails = append(actualEmails, *developer.Email)
	}

	return (compareSlices(expectedEmails, actualEmails))
}
