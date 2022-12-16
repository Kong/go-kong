package kong

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRequestBody(t *testing.T) {
	t.Run("body can be string", func(t *testing.T) {
		cl, err := NewClient(nil, nil)
		require.NoError(t, err)

		body := `{"_format_version":"1.1","services":[{"host":"example.com","name":"foo"}]}`

		req, err := cl.NewRequest("POST", "/", nil, body)
		require.NoError(t, err)

		b, err := io.ReadAll(req.Body)
		require.NoError(t, err)

		assert.Equal(t,
			`{"_format_version":"1.1","services":[{"host":"example.com","name":"foo"}]}`,
			string(b),
		)
	})

	t.Run("body can be []byte", func(t *testing.T) {
		cl, err := NewClient(nil, nil)
		require.NoError(t, err)

		body := []byte(`{"_format_version":"1.1","services":[{"host":"example.com","name":"foo"}]}`)

		req, err := cl.NewRequest("POST", "/", nil, body)
		require.NoError(t, err)

		b, err := io.ReadAll(req.Body)
		require.NoError(t, err)

		assert.Equal(t,
			`{"_format_version":"1.1","services":[{"host":"example.com","name":"foo"}]}`,
			string(b),
		)
	})

	t.Run("body can be a bytes.Buffer", func(t *testing.T) {
		cl, err := NewClient(nil, nil)
		require.NoError(t, err)

		body := bytes.NewBufferString(`{"_format_version":"1.1","services":[{"host":"example.com","name":"foo"}]}`)

		req, err := cl.NewRequest("POST", "/", nil, body)
		require.NoError(t, err)

		b, err := io.ReadAll(req.Body)
		require.NoError(t, err)

		assert.Equal(t,
			`{"_format_version":"1.1","services":[{"host":"example.com","name":"foo"}]}`,
			string(b),
		)
	})

	t.Run("body can be a map", func(t *testing.T) {
		cl, err := NewClient(nil, nil)
		require.NoError(t, err)

		body := map[string]any{
			"_format_version": "1.1",
			"services": []map[string]any{
				{
					"host": "example.com",
					"name": "foo",
				},
			},
		}

		req, err := cl.NewRequest("POST", "/", nil, body)
		require.NoError(t, err)

		b, err := io.ReadAll(req.Body)
		require.NoError(t, err)

		assert.Equal(t,
			`{"_format_version":"1.1","services":[{"host":"example.com","name":"foo"}]}`,
			string(b),
		)
	})
}
