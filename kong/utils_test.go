package kong

import (
	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestStringArrayToString(t *testing.T) {
	assert := assert.New(t)

	arr := StringSlice("foo", "bar")
	s := stringArrayToString(arr)
	assert.Equal("[ foo, bar ]", s)

	arr = StringSlice("foo")
	s = stringArrayToString(arr)
	assert.Equal("[ foo ]", s)

	assert.Equal(stringArrayToString(nil), "nil")
}

func TestString(t *testing.T) {
	assert := assert.New(t)

	s := String("foo")
	assert.Equal("foo", *s)
}

func TestBool(t *testing.T) {
	assert := assert.New(t)

	b := Bool(true)
	assert.Equal(true, *b)
}

func TestInt(t *testing.T) {
	assert := assert.New(t)

	i := Int(42)
	assert.Equal(42, *i)
}

func TestStringSlice(t *testing.T) {
	assert := assert.New(t)

	arrp := StringSlice()
	assert.Empty(arrp)

	arrp = StringSlice("foo", "bar")
	assert.Equal(2, len(arrp))
	assert.Equal("foo", *arrp[0])
	assert.Equal("bar", *arrp[1])
}

func TestFixVersion(t *testing.T) {
	validVersions := map[string]string{
		"0.14.1":                          "0.14.1",
		"0.14.2rc":                        "0.14.2-rc",
		"0.14.2rc1":                       "0.14.2-rc1",
		"0.14.2preview":                   "0.14.2-preview",
		"0.14.2preview1":                  "0.14.2-preview1",
		"0.33-enterprise-edition":         "0.33.0-enterprise",
		"0.33-1-enterprise-edition":       "0.33.1-enterprise",
		"1.3.0.0-enterprise-edition-lite": "1.3.0-0-enterprise-lite",
		"1.3.0.0-enterprise-lite":         "1.3.0-0-enterprise-lite",
	}
	for inputVersion, expectedVersion := range validVersions {
		v, err := cleanSemVer(inputVersion)
		if err != nil {
			t.Errorf("error converting %s: %v", inputVersion, err)
		} else if v.String() != expectedVersion {
			t.Errorf("converting %s, expecting %s, getting %s", inputVersion, expectedVersion, v.String())
		}
	}

	invalidVersions := []string{
		"",
		"0-1-1",
	}
	for _, inputVersion := range invalidVersions {
		_, err := cleanSemVer(inputVersion)
		if err == nil {
			t.Errorf("expecting error converting %s, getting no errors", inputVersion)
		}
	}
}

func Test_getKong(t *testing.T) {

	kongWithoutCredentialsSupport := new(Kong)
	kongWithoutCredentialsSupport.Credentials.minVersion = "1.4.0"
	kongWithoutCredentialsSupport.Credentials.hasTagSupport = false

	kongWithCredentialsSupport := new(Kong)
	kongWithCredentialsSupport.Credentials.minVersion = "1.4.0"
	kongWithCredentialsSupport.Credentials.hasTagSupport = true

	tests := []struct {
		name     string
		root     map[string]interface{}
		expected *Kong
	}{
		{
			root: map[string]interface{}{
				"version": "0.33-1-enterprise-edition",
				"configuration": map[string]interface{}{
					"database": "off",
				},
			},
			expected: &Kong{
				Version:     semver.MustParse("0.33.1-enterprise"),
				Enterprise:  true,
				Database:    "off",
				Credentials: kongWithoutCredentialsSupport.Credentials,
			},
		},
		{
			root: map[string]interface{}{
				"version": "2.3.2.0",
				"configuration": map[string]interface{}{
					"database": "cassandra",
				},
			},
			expected: &Kong{
				Version:     semver.MustParse("2.3.2-0"),
				Enterprise:  false,
				Database:    "cassandra",
				Credentials: kongWithCredentialsSupport.Credentials,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := getKong(tt.root)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, but is %v", tt.expected, result)
			}
		})
	}

}
