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

type errorTransport struct {
	code    int
	message string
}

func (et *errorTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	return nil, NewAPIError(et.code, et.message)
}

func TestIsForbiddenErrE2E(T *testing.T) {
	assert := assert.New(T)

	// Create the http client with a custom transport
	// always returning a 403
	forbiddenTransport := &errorTransport{
		code:    http.StatusForbidden,
		message: "Enterprise license missing or expired",
	}
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

func TestIsMethodNotAllowedErrE2E(T *testing.T) {
	assert := assert.New(T)

	// Create the http client with a custom transport
	// always returning a 405
	methodNotAllowedTransport := &errorTransport{
		code:    http.StatusMethodNotAllowed,
		message: "Method Not Allowed",
	}
	httpClient := &http.Client{
		Transport: methodNotAllowedTransport,
	}
	// if KONG_ADMIN_TOKEN is set, a different custom transport would be used
	T.Setenv("KONG_ADMIN_TOKEN", "")

	client, err := NewTestClient(nil, httpClient)
	require.NoError(T, err)
	require.NotNil(T, client)

	_, err = client.Consumers.ListAll(defaultCtx)
	require.Error(T, err)
	assert.True(IsMethodNotAllowedErr(err))
}

func TestIsBadRequestErrE2E(T *testing.T) {
	assert := assert.New(T)

	// Create the http client with a custom transport
	// always returning a 405
	badRequestTransport := &errorTransport{
		code:    http.StatusBadRequest,
		message: "Bad Request",
	}
	httpClient := &http.Client{
		Transport: badRequestTransport,
	}
	// if KONG_ADMIN_TOKEN is set, a different custom transport would be used
	T.Setenv("KONG_ADMIN_TOKEN", "")

	client, err := NewTestClient(nil, httpClient)
	require.NoError(T, err)
	require.NotNil(T, client)

	_, err = client.Consumers.ListAll(defaultCtx)
	require.Error(T, err)
	assert.True(IsBadRequestErr(err))
}
