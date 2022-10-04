package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVaultsService(t *testing.T) {
	RunWhenEnterprise(t, ">=3.0.0", RequiredFeatures{})
	assert := assert.New(t)
	require := require.New(t)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	vault := &Vault{
		Name:        String("env"),
		Prefix:      String("my-env-vault"),
		Description: String("ENV vault for secrets"),
		Config: Configuration{
			"prefix": "SECURE_",
		},
	}

	createdVault, err := client.Vaults.Create(defaultCtx, vault)
	assert.NoError(err)
	require.NotNil(createdVault)

	vault, err = client.Vaults.Get(defaultCtx, createdVault.ID)
	assert.NoError(err)
	assert.NotNil(vault)

	vault.Prefix = String("my-new-env-vault")
	vault.Description = String("new ENV vault for secrets")

	vault, err = client.Vaults.Update(defaultCtx, vault)
	assert.NoError(err)
	assert.NotNil(vault)
	assert.Equal("env", *vault.Name)
	assert.Equal("new ENV vault for secrets", *vault.Description)
	assert.Equal("my-new-env-vault", *vault.Prefix)

	err = client.Vaults.Delete(defaultCtx, vault.ID)
	assert.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	vault = &Vault{
		Name:        String("aws"),
		Prefix:      String("my-aws-vault"),
		ID:          String(id),
		Description: String("aws vault for secrets"),
		Config: Configuration{
			"region": "us-east-2",
		},
	}

	createdVault, err = client.Vaults.Create(defaultCtx, vault)
	assert.NoError(err)
	assert.NotNil(createdVault)
	assert.Equal(id, *createdVault.ID)
	assert.Equal("aws", *createdVault.Name)
	assert.Equal("aws vault for secrets", *createdVault.Description)
	assert.Equal(Configuration{"region": "us-east-2"}, createdVault.Config)

	err = client.Vaults.Delete(defaultCtx, createdVault.ID)
	assert.NoError(err)

	_, err = client.Vaults.Create(defaultCtx, nil)
	assert.NotNil(err)

	_, err = client.Vaults.Update(defaultCtx, nil)
	assert.NotNil(err)
}

func TestVaultWithTags(t *testing.T) {
	RunWhenEnterprise(t, ">=3.0.0", RequiredFeatures{})
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	vault := &Vault{
		Name:        String("env"),
		Prefix:      String("my-env-vault"),
		Description: String("ENV vault for secrets"),
		Config: Configuration{
			"prefix": "SECURE_",
		},
		Tags: StringSlice("tag1", "tag2"),
	}

	createdVault, err := client.Vaults.Create(defaultCtx, vault)
	assert.NoError(err)
	assert.NotNil(createdVault)
	assert.Equal(StringSlice("tag1", "tag2"), createdVault.Tags)

	err = client.Vaults.Delete(defaultCtx, createdVault.ID)
	assert.NoError(err)
}

func TestVaultListEndpoint(t *testing.T) {
	RunWhenEnterprise(t, ">=3.0.0", RequiredFeatures{})
	assert := assert.New(t)
	require := require.New(t)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	// fixtures
	vaults := []*Vault{
		{
			Name:        String("env"),
			Prefix:      String("my-env-vault"),
			Description: String("ENV vault for secrets"),
			Config: Configuration{
				"prefix": "SECURE_",
			},
		},
		{
			Name:        String("aws"),
			Prefix:      String("my-aws-vault"),
			Description: String("aws vault for secrets"),
			Config: Configuration{
				"region": "us-east-2",
			},
		},
		{
			Name:        String("hcv"),
			Prefix:      String("my-hcv-vault"),
			Description: String("hcv vault for secrets"),
			Config: Configuration{
				"token": "super-secret-token",
			},
		},
	}

	// create fixturs
	for i := 0; i < len(vaults); i++ {
		vault, err := client.Vaults.Create(defaultCtx, vaults[i])
		assert.NoError(err)
		assert.NotNil(vault)
		vaults[i] = vault
	}

	vaultsFromKong, next, err := client.Vaults.List(defaultCtx, nil)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(vaultsFromKong)
	assert.Equal(3, len(vaultsFromKong))

	// check if we see all vaults
	assert.True(compareVaults(t, vaults, vaultsFromKong))

	// Test pagination
	vaultsFromKong = []*Vault{}

	// first page
	page1, next, err := client.Vaults.List(defaultCtx, &ListOpt{Size: 1})
	assert.NoError(err)
	require.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	vaultsFromKong = append(vaultsFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.Vaults.List(defaultCtx, next)
	assert.NoError(err)
	assert.Nil(next)
	assert.NotNil(page2)
	assert.Equal(2, len(page2))
	vaultsFromKong = append(vaultsFromKong, page2...)

	assert.True(compareVaults(t, vaults, vaultsFromKong))

	vaults, err = client.Vaults.ListAll(defaultCtx)
	assert.NoError(err)
	assert.NotNil(vaults)
	assert.Equal(3, len(vaults))

	for i := 0; i < len(vaults); i++ {
		assert.NoError(client.Vaults.Delete(defaultCtx, vaults[i].ID))
	}
}

func compareVaults(t *testing.T, expected, actual []*Vault) bool {
	var expectedPrefixes, actualPrefixes []string
	for _, vault := range expected {
		if !assert.NotNil(t, vault) {
			continue
		}
		expectedPrefixes = append(expectedPrefixes, *vault.Prefix)
	}

	for _, vault := range actual {
		actualPrefixes = append(actualPrefixes, *vault.Prefix)
	}

	return (compareSlices(expectedPrefixes, actualPrefixes))
}
