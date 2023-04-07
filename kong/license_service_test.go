package kong

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLicensesService(t *testing.T) {
	RunWhenEnterprise(t, ">=2.3.0", RequiredFeatures{})
	licenseString := os.Getenv("KONG_LICENSE_DATA")
	if licenseString == "" {
		t.Skip("KONG_LICENSE_DATA not set, skipping")
	}
	require := require.New(t)
	assert := assert.New(t)

	client, err := NewTestClient(nil, nil)
	require.NoError(err)
	require.NotNil(client)

	license := &License{
		Payload: String(licenseString),
	}

	createdLicense, err := client.Licenses.Create(defaultCtx, license)
	require.NoError(err)
	require.NotNil(createdLicense)
	t.Cleanup(func() {
		// Note the assert here as we might want more logic to be run as part of the cleanup,
		// regardless of the return value of removing this particular object.
		assert.NoError(client.Licenses.Delete(defaultCtx, createdLicense.ID))
	})

	license, err = client.Licenses.Get(defaultCtx, createdLicense.ID)
	require.NoError(err)
	require.NotNil(license)
	origTime := license.UpdatedAt

	license.Payload = String(licenseString)

	license, err = client.Licenses.Update(defaultCtx, license)
	require.NoError(err)
	require.NotNil(license)
	require.NotEqual(origTime, *license.UpdatedAt)

	err = client.Licenses.Delete(defaultCtx, license.ID)
	require.NoError(err)

	// ID can be specified
	id := uuid.NewString()
	license = &License{
		ID:      String(id),
		Payload: String(licenseString),
	}

	createdLicense, err = client.Licenses.Create(defaultCtx, license)
	require.NoError(err)
	require.NotNil(createdLicense)
	t.Cleanup(func() {
		// Note the assert here as we might want more logic to be run as part of the cleanup,
		// regardless of the return value of removing this particular object.
		assert.NoError(client.Licenses.Delete(defaultCtx, createdLicense.ID))
	})
	require.Equal(id, *createdLicense.ID)

	err = client.Licenses.Delete(defaultCtx, createdLicense.ID)
	require.NoError(err)

	_, err = client.Licenses.Create(defaultCtx, nil)
	require.EqualError(err, "cannot create a nil license")

	_, err = client.Licenses.Update(defaultCtx, nil)
	require.EqualError(err, "cannot update a nil license")
}

// Note: no list test as we only have the one test license allowed in public repos (TODO confirm)
// Not that much of a concern as we don't expect to use multiple in practice.
