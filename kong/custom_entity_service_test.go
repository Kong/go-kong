package kong

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kong/go-kong/kong/custom"
)

func TestCustomEntityService(T *testing.T) {
	T.Skip()
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)
	// fixture consumer
	consumer, err := client.Consumers.Create(defaultCtx,
		&Consumer{Username: String("foo")})
	assert.NoError(err)
	require.NotNil(consumer)

	// create a key associated with the consumer
	k1 := custom.NewEntityObject("key-auth")
	k1.AddRelation("consumer_id", *consumer.ID)
	e1, err := client.CustomEntities.Create(defaultCtx, k1)
	assert.NotNil(e1)
	assert.NoError(err)

	// look up the key
	se := custom.NewEntityObject("key-auth")
	se.AddRelation("consumer_id", *consumer.ID)
	se.SetObject(map[string]interface{}{"id": e1.Object()["id"]})
	gotE, err := client.CustomEntities.Get(defaultCtx, se)
	assert.NotNil(gotE)
	assert.Equal(e1.Object()["key"], gotE.Object()["key"])
	assert.NoError(err)

	gotE.Object()["key"] = "my-secret"
	e1, err = client.CustomEntities.Update(defaultCtx, gotE)
	assert.NotNil(e1)
	assert.NoError(err)
	assert.Equal("my-secret", e1.Object()["key"])

	// PUT request
	k2 := custom.NewEntityObject("key-auth")
	id := "fc3898d9-4b4d-4491-a834-8358646e2d20"
	k2.SetObject(map[string]interface{}{
		"id":  id,
		"key": "super-secret",
	})
	k2.AddRelation("consumer_id", *consumer.ID)
	e2, err := client.CustomEntities.Create(defaultCtx, k2)
	assert.NotNil(e2)
	assert.NoError(err)
	assert.Equal("super-secret", e2.Object()["key"])
	assert.Equal(id, e2.Object()["id"])

	se = custom.NewEntityObject("key-auth")
	se.AddRelation("consumer_id", *consumer.ID)
	keyAuths, _, err := client.CustomEntities.List(defaultCtx, nil, se)

	assert.NoError(err)
	assert.Equal(2, len(keyAuths))

	// list endpoint
	keyAuths, err = client.CustomEntities.ListAll(defaultCtx, se)
	assert.NoError(err)
	assert.Equal(2, len(keyAuths))

	expectedKeys := []string{
		e1.Object()["key"].(string),
		e2.Object()["key"].(string),
	}
	actualKeys := []string{
		keyAuths[0].Object()["key"].(string),
		keyAuths[1].Object()["key"].(string),
	}
	sort.Strings(expectedKeys)
	sort.Strings(actualKeys)
	assert.Equal(expectedKeys, actualKeys)
	assert.NoError(client.CustomEntities.Delete(defaultCtx, e1))
	assert.NoError(client.CustomEntities.Delete(defaultCtx, e2))

	// delete fixture consumer
	assert.NoError(client.Consumers.Delete(defaultCtx, consumer.ID))
}
