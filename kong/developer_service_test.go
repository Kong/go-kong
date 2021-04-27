package kong

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestDevelopersService(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", requiredFeatures{portal: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, "", nil)
	assert.Nil(err)
	assert.NotNil(client)

	testWs, err := NewTestWorkspace(client, "default")
	assert.Nil(err)
	assert.NoError(testWs.UpdateConfig(map[string]interface{}{"portal_auth": "basic-auth",
		"portal_session_conf": map[string]interface{}{"secret": "garbage"},
		"portal":              true}))

	developer := &Developer{
		Meta:     String("{\"full_name\": \"Foo BAR\"}"),
		Email:    String("foo.bar@example.com"),
		CustomID: String("custom_id_foo"),
		Password: String("F00P@ssword"),
	}

	createdDeveloper, err := client.Developers.Create(defaultCtx, developer)
	assert.Nil(err)
	assert.NotNil(createdDeveloper)

	developer, err = client.Developers.Get(defaultCtx, createdDeveloper.ID)
	assert.Nil(err)
	assert.NotNil(developer)

	developer, err = client.Developers.GetByCustomID(defaultCtx,
		String("does-not-exist"))
	assert.NotNil(err)
	assert.Nil(developer)

	developer, err = client.Developers.GetByCustomID(defaultCtx,
		String("custom_id_foo"))
	assert.Nil(err)
	assert.NotNil(developer)

	developer.Email = String("bar@example.com")
	developer, err = client.Developers.Update(defaultCtx, developer)
	assert.Nil(err)
	assert.NotNil(developer)
	assert.Equal("bar@example.com", *developer.Email)

	err = client.Developers.Delete(defaultCtx, createdDeveloper.ID)
	assert.Nil(err)

	// ID can be specified
	id := uuid.NewV4().String()
	developer = &Developer{
		Meta:     String("{\"full_name\": \"Foo BAR\"}"),
		Email:    String("foo.bar@example.com"),
		Password: String("F00P@ssword"),
		ID:       String(id),
	}

	createdDeveloper, err = client.Developers.Create(defaultCtx, developer)
	assert.Nil(err)
	assert.NotNil(createdDeveloper)
	assert.Equal(id, *createdDeveloper.ID)

	err = client.Developers.Delete(defaultCtx, createdDeveloper.ID)
	assert.Nil(err)

	assert.NoError(testWs.Reset())
}

func TestDeveloperListEndpoint(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", requiredFeatures{portal: true})
	assert := assert.New(T)

	client, err := NewTestClient(nil, "", nil)
	assert.Nil(err)
	assert.NotNil(client)

	testWs, err := NewTestWorkspace(client, "default")
	assert.Nil(err)
	assert.NoError(testWs.UpdateConfig(map[string]interface{}{"portal_auth": "basic-auth",
		"portal_session_conf": map[string]interface{}{"secret": "garbage"},
		"portal":              true}))

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
		assert.Nil(err)
		assert.NotNil(developer)
		developers[i] = developer
	}

	developersFromKong, next, err := client.Developers.List(defaultCtx, nil)
	assert.Nil(err)
	assert.Nil(next)
	assert.NotNil(developersFromKong)
	assert.Equal(3, len(developersFromKong))

	// check if we see all developers
	assert.True(compareDevelopers(developers, developersFromKong))

	// Test pagination
	developersFromKong = []*Developer{}

	// first page
	page1, next, err := client.Developers.List(defaultCtx, &ListOpt{Size: 1})
	assert.Nil(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	developersFromKong = append(developersFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.Developers.List(defaultCtx, next)
	assert.Nil(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Equal(2, len(page2))
	developersFromKong = append(developersFromKong, page2...)

	assert.True(compareDevelopers(developers, developersFromKong))

	developers, err = client.Developers.ListAll(defaultCtx)
	assert.Nil(err)
	assert.NotNil(developers)
	assert.Equal(3, len(developers))

	for i := 0; i < len(developers); i++ {
		assert.Nil(client.Developers.Delete(defaultCtx, developers[i].ID))
	}

	assert.NoError(testWs.Reset())
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
