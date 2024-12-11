package kong

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKonnectApplicationService_Create(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenEnterprise(T, ">=3.6.0", RequiredFeatures{})

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	var (
		clientID        = uuid.NewString()
		consumerGroup   = []string{uuid.NewString()}
		scopes          = []string{"/auth"}
		authStrategyID  = uuid.NewString()
		exhaustedScopes = []string{"/eauth"}
		orgID           = uuid.NewString()
		developerID     = uuid.NewString()
		createdAt       = time.Now().Unix()
	)

	kaa := &KonnectApplication{
		ID:              &clientID,
		ClientID:        clientID,
		ConsumerGroups:  consumerGroup,
		Scopes:          scopes,
		AuthStrategyID:  &authStrategyID,
		ExhaustedScopes: exhaustedScopes,
		ApplicationContext: &ApplicationContext{
			OrganizationID: &orgID,
			DeveloperID:    &developerID,
		},
		CreatedAt: createdAt,
	}
	createResponse, err := client.KonnectApplication.Create(defaultCtx, kaa)
	require.NoError(err)
	require.NotNil(createResponse)

	T.Cleanup(func() {
		assert.NoError(client.KonnectApplication.Delete(context.Background(), createResponse.ID))
	})

	require.Equal(createResponse.ClientID, clientID)
	require.Equal(createResponse.CreatedAt, createdAt)
	require.Equal(createResponse.ConsumerGroups, consumerGroup)
	require.Equal(createResponse.Scopes, scopes)
	require.Equal(createResponse.ExhaustedScopes, exhaustedScopes)
}

func TestKonnectApplicationService_ListAll(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenEnterprise(T, ">=3.6.0", RequiredFeatures{})

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	var (
		expectedKonnectApplications = 10
		consumerGroup               = []string{uuid.NewString()}
		scopes                      = []string{"/auth"}
		authStrategyID              = uuid.NewString()
		exhaustedScopes             = []string{"/eauth"}
		orgID                       = uuid.NewString()
		developerID                 = uuid.NewString()
		createdAt                   = time.Now().Unix()
	)

	kaa := &KonnectApplication{
		ConsumerGroups:  consumerGroup,
		Scopes:          scopes,
		AuthStrategyID:  &authStrategyID,
		ExhaustedScopes: exhaustedScopes,
		ApplicationContext: &ApplicationContext{
			OrganizationID: &orgID,
			DeveloperID:    &developerID,
		},
		CreatedAt: createdAt,
	}

	for i := 0; i < expectedKonnectApplications; i++ {
		clientID := uuid.NewString()
		kaa.ID = &clientID
		kaa.ClientID = clientID
		createResponse, err := client.KonnectApplication.Create(defaultCtx, kaa)
		require.NoError(err)
		require.NotNil(createResponse)

		T.Cleanup(func() {
			assert.NoError(client.KonnectApplication.Delete(context.Background(), createResponse.ID))
		})
	}

	listKonnectApplicationResponse, err := client.KonnectApplication.ListAll(defaultCtx)
	require.NoError(err)
	require.Len(listKonnectApplicationResponse, expectedKonnectApplications)
}

func TestKonnectApplicationService_List(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenEnterprise(T, ">=3.6.0", RequiredFeatures{})

	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	var (
		expectedKonnectApplications = 10
		consumerGroup               = []string{uuid.NewString()}
		scopes                      = []string{"/auth"}
		authStrategyID              = uuid.NewString()
		exhaustedScopes             = []string{"/eauth"}
		orgID                       = uuid.NewString()
		developerID                 = uuid.NewString()
		createdAt                   = time.Now().Unix()
		tagA                        = "list1"
		tagB                        = "list2"
	)

	kaa := &KonnectApplication{
		ConsumerGroups:  consumerGroup,
		Scopes:          scopes,
		AuthStrategyID:  &authStrategyID,
		ExhaustedScopes: exhaustedScopes,
		ApplicationContext: &ApplicationContext{
			OrganizationID: &orgID,
			DeveloperID:    &developerID,
		},
		CreatedAt: createdAt,
	}

	for i := 0; i < expectedKonnectApplications/2; i++ {
		clientID := uuid.NewString()
		kaa.ID = &clientID
		kaa.ClientID = clientID
		kaa.Tags = lo.ToPtr([]string{tagA})
		createResponse, err := client.KonnectApplication.Create(defaultCtx, kaa)
		require.NoError(err)
		require.NotNil(createResponse)

		T.Cleanup(func() {
			assert.NoError(client.KonnectApplication.Delete(context.Background(), createResponse.ID))
		})
	}

	for i := 0; i < expectedKonnectApplications/2; i++ {
		clientID := uuid.NewString()
		kaa.ID = &clientID
		kaa.ClientID = clientID
		kaa.Tags = lo.ToPtr([]string{tagB})
		createResponse, err := client.KonnectApplication.Create(defaultCtx, kaa)
		require.NoError(err)
		require.NotNil(createResponse)

		T.Cleanup(func() {
			assert.NoError(client.KonnectApplication.Delete(context.Background(), createResponse.ID))
		})
	}

	// Filter by tag listA
	listKonnectApplicationResponseByTagA, _, err := client.KonnectApplication.List(defaultCtx, &ListOpt{
		Size: 10,
		Tags: []*string{lo.ToPtr(tagA)},
	})
	require.NoError(err)
	require.Len(listKonnectApplicationResponseByTagA, 5)

	// Filter by tag listB
	listKonnectApplicationResponseByTagB, _, err := client.KonnectApplication.List(defaultCtx, &ListOpt{
		Size: 10,
		Tags: []*string{lo.ToPtr(tagB)},
	})
	require.NoError(err)
	require.Len(listKonnectApplicationResponseByTagB, 5)

	size := 2
	// Filter by size
	listKonnectApplicationResponseBySize, _, err := client.KonnectApplication.List(defaultCtx, &ListOpt{
		Size: size,
	})
	require.NoError(err)
	require.Len(listKonnectApplicationResponseBySize, size)
}
