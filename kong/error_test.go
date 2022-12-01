package kong

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsNotFoundErr(T *testing.T) {
	assert := assert.New(T)
	var e error = NewAPIError(404, "")
	assert.True(IsNotFoundErr(e))
	assert.False(IsNotFoundErr(nil))

	err := fmt.Errorf("not a 404")
	assert.False(IsNotFoundErr(err))
}

func TestIsNotFoundErrE2E(T *testing.T) {
	assert := assert.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	consumer, err := client.Consumers.Get(defaultCtx, String("does-not-exists"))
	assert.Nil(consumer)
	assert.NotNil(err)
	assert.True(IsNotFoundErr(err))
}

func TestAPIError_Code(T *testing.T) {
	assert := assert.New(T)
	require := require.New(T)

	client, err := NewTestClient(nil, nil)
	assert.NoError(err)
	assert.NotNil(client)

	consumer, err := client.Consumers.Get(defaultCtx, String("does-not-exists"))
	assert.Nil(consumer)
	assert.NotNil(err)

	var kongErr *APIError
	ok := errors.As(err, &kongErr)
	require.True(ok)
	assert.True(kongErr.Code() == 404)
}

func TestIsForbiddenErrE2E(T *testing.T) {
	assert := assert.New(T)

	// Create the http client with a custom transport
	// always returning a 403
	forbiddenTransport := &forbiddenTransport{}
	httpClient := &http.Client{
		Transport: forbiddenTransport,
	}
	// if KONG_ADMIN_TOKEN is set, a different custom transport would be used
	T.Setenv("KONG_ADMIN_TOKEN", "")

	client, err := NewTestClient(nil, httpClient)
	require.NoError(T, err)
	require.NotNil(T, client)

	_, err = client.Consumers.ListAll(defaultCtx)
	require.Error(T, err)
	assert.True(IsForbiddenErr(err))
}

type forbiddenTransport struct{}

func (ft *forbiddenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Create a new response with 403 status code
	return nil, NewAPIError(
		http.StatusForbidden,
		"Enterprise license missing or expired")
}
