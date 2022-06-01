package kong

import (
	"errors"
	"fmt"
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
