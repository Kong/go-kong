package kong

import (
	"testing"
	"time"

	"github.com/google/uuid"
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
	require.Equal(createResponse.ClientID, clientID)
	require.Equal(createResponse.CreatedAt, createdAt)
	require.Equal(createResponse.ConsumerGroups, consumerGroup)
	require.Equal(createResponse.Scopes, scopes)
	require.Equal(createResponse.ExhaustedScopes, exhaustedScopes)

	err = client.KonnectApplication.Delete(defaultCtx, createResponse.ID)
	require.NoError(err)
}

func TestKonnectApplicationService_ListAll(T *testing.T) {
	RunWhenDBMode(T, "postgres")
	RunWhenEnterprise(T, ">=3.6.0", RequiredFeatures{})

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
		
		t.Cleanup(func() {
			assert.NoError(client.KonnectApplication.Delete(context.Background(), createResponse.ID))
	})
	}

	listKonnectApplicationResponse, err := client.KonnectApplication.ListAll(defaultCtx)
	require.NoError(err)
	require.Len(listKonnectApplicationResponse, expectedKonnectApplications)
}
