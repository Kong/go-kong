package kong

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigurationDeepCopyInto(T *testing.T) {
	assert := assert.New(T)

	var c Configuration
	byt := []byte(`{"int":42,"float":4.2,"strings":["foo","bar"]}`)
	require.NoError(T, json.Unmarshal(byt, &c))

	c2 := c.DeepCopy()
	assert.Equal(c, c2)

	// Both are independent now
	c["int"] = 24
	assert.Equal(24, c["int"])
	assert.EqualValues(42, c2["int"])

	c["strings"] = []string{"fubar"}
	assert.Equal([]string{"fubar"}, c["strings"].([]string))
	assert.Equal([]interface{}{"foo", "bar"}, c2["strings"])
}
