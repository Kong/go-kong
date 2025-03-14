package custom

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultRegistry(t *testing.T) {
	assert := assert.New(t)
	r := NewDefaultRegistry()

	require.NotNil(t, r)
	var typ Type = "foo"
	entitiy := EntityCRUDDefinition{
		Name: typ,
	}
	err := r.Register(typ, &entitiy)
	require.NoError(t, err)
	err = r.Register(typ, &entitiy)
	require.Error(t, err)

	e := r.Lookup(typ)
	assert.NotNil(e)
	assert.Equal(e, &entitiy)
	e = r.Lookup("NotExists")
	assert.Nil(e)

	err = r.Unregister("NotExists)")
	require.Error(t, err)

	err = r.Unregister(typ)
	require.NoError(t, err)
}
