package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAIModelsCRUD(T *testing.T) {
	RunWhenAIGateway(T, ">=2.0.0")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	aiModel := &AIModel{
		Name:  String("gpt-5-2"),
		Alias: String("@openai/gpt-5.2"),
	}

	createdAIModel, err := client.AIModels.Create(defaultCtx, aiModel)
	require.NoError(err)
	assert.NotNil(createdAIModel)
	assert.NotNil(createdAIModel.ID)
	assert.Equal("gpt-5-2", *createdAIModel.Name)
	assert.Equal("@openai/gpt-5.2", *createdAIModel.Alias)

	aiModel, err = client.AIModels.Get(defaultCtx, createdAIModel.ID)
	require.NoError(err)
	assert.NotNil(aiModel)
	assert.Equal(*createdAIModel.ID, *aiModel.ID)

	aiModel.Name = String("gpt-5-3")
	aiModel.Alias = String("@openai/gpt-5.3")
	aiModel, err = client.AIModels.Update(defaultCtx, aiModel)
	require.NoError(err)
	assert.NotNil(aiModel)
	assert.Equal("gpt-5-3", *aiModel.Name)
	assert.Equal("@openai/gpt-5.3", *aiModel.Alias)

	err = client.AIModels.Delete(defaultCtx, createdAIModel.ID)
	require.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	aiModel = &AIModel{
		Name:  String("gpt-4-turbo"),
		Alias: String("@openai/gpt-4-turbo"),
		ID:    String(id),
	}

	createdAIModel, err = client.AIModels.Create(defaultCtx, aiModel)
	require.NoError(err)
	assert.NotNil(createdAIModel)
	assert.Equal(id, *createdAIModel.ID)

	err = client.AIModels.Delete(defaultCtx, createdAIModel.ID)
	require.NoError(err)
}

func TestAIModelWithTags(T *testing.T) {
	RunWhenAIGateway(T, ">=2.0.0")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	createdAIModel, err := client.AIModels.Create(defaultCtx, &AIModel{
		Name:  String("gpt-5"),
		Alias: String("@openai/gpt-5"),
		Tags:  StringSlice("tag1", "tag2"),
	})
	require.NoError(err)
	require.NotNil(createdAIModel)
	assert.Equal(StringSlice("tag1", "tag2"), createdAIModel.Tags)

	err = client.AIModels.Delete(defaultCtx, createdAIModel.ID)
	require.NoError(err)
}

func TestAIModelListEndpoint(T *testing.T) {
	RunWhenAIGateway(T, ">=2.0.0")

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	assert.NotNil(client)

	// fixtures
	aiModels := []*AIModel{
		{
			Name:  String("model1"),
			Alias: String("@provider/model1"),
		},
		{
			Name:  String("model2"),
			Alias: String("@provider/model2"),
		},
		{
			Name:  String("model3"),
			Alias: String("@provider/model3"),
		},
	}

	// create fixtures
	for i := 0; i < len(aiModels); i++ {
		aiModel, err := client.AIModels.Create(defaultCtx, aiModels[i])
		require.NoError(err)
		assert.NotNil(aiModel)
		aiModels[i] = aiModel
	}

	aiModelsFromKong, next, err := client.AIModels.List(defaultCtx, nil)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(aiModelsFromKong)
	assert.Len(aiModelsFromKong, 3)

	// check if we see all ai models
	assert.True(compareAIModels(T, aiModels, aiModelsFromKong))

	// Test pagination
	aiModelsFromKong = []*AIModel{}

	// first page
	page1, next, err := client.AIModels.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(err)
	require.NotNil(next)
	assert.NotNil(page1)
	assert.Len(page1, 1)
	aiModelsFromKong = append(aiModelsFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.AIModels.List(defaultCtx, next)
	require.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Len(page2, 2)
	aiModelsFromKong = append(aiModelsFromKong, page2...)

	assert.True(compareAIModels(T, aiModels, aiModelsFromKong))

	aiModels, err = client.AIModels.ListAll(defaultCtx)
	require.NoError(err)
	assert.NotNil(aiModels)
	assert.Len(aiModels, 3)

	for i := 0; i < len(aiModels); i++ {
		require.NoError(client.AIModels.Delete(defaultCtx, aiModels[i].ID))
	}
}

func compareAIModels(T *testing.T, expected, actual []*AIModel) bool {
	var expectedNames, actualNames []string
	for _, aiModel := range expected {
		if !assert.NotNil(T, aiModel) {
			continue
		}
		expectedNames = append(expectedNames, *aiModel.Name)
	}

	for _, aiModel := range actual {
		actualNames = append(actualNames, *aiModel.Name)
	}

	return compareSlices(expectedNames, actualNames)
}
