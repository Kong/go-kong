package kong

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestDevelopersService(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", false)
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	developer := &Developer{
		Email:    String("foo@example.com"),
		CustomID: String("custom_id_foo"),
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
		Email: String("foo@example.com"),
		ID:    String(id),
	}

	createdDeveloper, err = client.Developers.Create(defaultCtx, developer)
	assert.Nil(err)
	assert.NotNil(createdDeveloper)
	assert.Equal(id, *createdDeveloper.ID)

	err = client.Developers.Delete(defaultCtx, createdDeveloper.ID)
	assert.Nil(err)
}

func TestDeveloperListEndpoint(T *testing.T) {
	runWhenEnterprise(T, ">=0.33.0", false)
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	// fixtures
	developers := []*Developer{
		{
			Email: String("foo1@example.com"),
		},
		{
			Email: String("foo2@example.com"),
		},
		{
			Email: String("foo3@example.com"),
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
