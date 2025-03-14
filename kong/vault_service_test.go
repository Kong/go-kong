package kong

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestVaultsService(t *testing.T) {
	RunWhenEnterprise(t, ">=3.0.0", RequiredFeatures{})
	require := require.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	vault := &Vault{
		Name:        String("env"),
		Prefix:      String("my-env-vault"),
		Description: String("ENV vault for secrets"),
		Config: Configuration{
			"prefix": "SECURE_",
		},
	}

	createdVault, err := client.Vaults.Create(defaultCtx, vault)
	require.NoError(err)
	require.NotNil(createdVault)
	t.Cleanup(func() {
		// Note the assert here as we might want more logic to be run as part of the cleanup,
		// regardless of the return value of removing this particular object.
		require.NoError(client.Vaults.Delete(defaultCtx, createdVault.ID))
	})

	vault, err = client.Vaults.Get(defaultCtx, createdVault.ID)
	require.NoError(err)
	require.NotNil(vault)

	vault.Prefix = String("my-new-env-vault")
	vault.Description = String("new ENV vault for secrets")

	vault, err = client.Vaults.Update(defaultCtx, vault)
	require.NoError(err)
	require.NotNil(vault)
	require.Equal("env", *vault.Name)
	require.Equal("new ENV vault for secrets", *vault.Description)
	require.Equal("my-new-env-vault", *vault.Prefix)

	err = client.Vaults.Delete(defaultCtx, vault.ID)
	require.NoError(err)

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
	require.NoError(err)
	require.NotNil(createdVault)
	t.Cleanup(func() {
		// Note the assert here as we might want more logic to be run as part of the cleanup,
		// regardless of the return value of removing this particular object.
		require.NoError(client.Vaults.Delete(defaultCtx, createdVault.ID))
	})
	require.Equal(id, *createdVault.ID)
	require.Equal("aws", *createdVault.Name)
	require.Equal("aws vault for secrets", *createdVault.Description)
	region, ok := createdVault.Config["region"]
	require.True(ok)
	require.Equal("us-east-2", region)

	err = client.Vaults.Delete(defaultCtx, createdVault.ID)
	require.NoError(err)

	_, err = client.Vaults.Create(defaultCtx, nil)
	require.EqualError(err, "cannot create a nil vault")

	_, err = client.Vaults.Update(defaultCtx, nil)
	require.EqualError(err, "cannot update a nil vault")
}

func TestVaultWithTags(t *testing.T) {
	RunWhenEnterprise(t, ">=3.0.0", RequiredFeatures{})
	require := require.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

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
	require.NoError(err)
	require.NotNil(createdVault)
	t.Cleanup(func() {
		// Note the assert here as we might want more logic to be run as part of the cleanup,
		// regardless of the return value of removing this particular object.
		require.NoError(client.Vaults.Delete(defaultCtx, createdVault.ID))
	})
	require.Equal(StringSlice("tag1", "tag2"), createdVault.Tags)

	err = client.Vaults.Delete(defaultCtx, createdVault.ID)
	require.NoError(err)
}

func TestVaultListEndpoint(t *testing.T) {
	RunWhenEnterprise(t, ">=3.0.0", RequiredFeatures{})
	require := require.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

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
		require.NoError(err)
		require.NotNil(vault)
		vaults[i] = vault

		t.Cleanup(func() {
			// Note the assert here as we might want more logic to be run as part of the cleanup,
			// regardless of the return value of removing this particular object.
			require.NoError(client.Vaults.Delete(defaultCtx, vault.ID))
		})
	}

	vaultsFromKong, next, err := client.Vaults.List(defaultCtx, nil)
	require.NoError(err)
	require.Nil(next)
	require.NotNil(vaultsFromKong)
	require.Len(vaultsFromKong, 3)

	// check if we see all vaults
	require.True(compareVaults(t, vaults, vaultsFromKong))

	// Test pagination
	vaultsFromKong = []*Vault{}

	// first page
	page1, next, err := client.Vaults.List(defaultCtx, &ListOpt{Size: 1})
	require.NoError(err)
	require.NotNil(next)
	require.NotNil(page1)
	require.Len(page1, 1)
	vaultsFromKong = append(vaultsFromKong, page1...)

	// last page
	next.Size = 2
	page2, next, err := client.Vaults.List(defaultCtx, next)
	require.NoError(err)
	require.Nil(next)
	require.NotNil(page2)
	require.Len(page2, 2)
	vaultsFromKong = append(vaultsFromKong, page2...)

	require.True(compareVaults(t, vaults, vaultsFromKong))

	vaults, err = client.Vaults.ListAll(defaultCtx)
	require.NoError(err)
	require.NotNil(vaults)
	require.Len(vaults, 3)
}

func compareVaults(t *testing.T, expected, actual []*Vault) bool {
	var expectedPrefixes, actualPrefixes []string
	for _, vault := range expected {
		require.NotNil(t, vault)
		expectedPrefixes = append(expectedPrefixes, *vault.Prefix)
	}

	for _, vault := range actual {
		actualPrefixes = append(actualPrefixes, *vault.Prefix)
	}

	return compareSlices(expectedPrefixes, actualPrefixes)
}
