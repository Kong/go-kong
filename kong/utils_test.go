package kong

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
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
		v, err := ParseSemanticVersion(inputVersion)
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
		_, err := ParseSemanticVersion(inputVersion)
		if err == nil {
			t.Errorf("expecting error converting %s, getting no errors", inputVersion)
		}
	}
}

func Test_requestWithHeaders(t *testing.T) {
	type args struct {
		req     *http.Request
		headers http.Header
	}
	tests := []struct {
		name string
		args args
		want *http.Request
	}{
		{
			name: "returns request as is if no headers are set",
			args: args{
				req: &http.Request{
					Method: "GET",
					Header: http.Header{
						"foo": []string{"bar", "baz"},
					},
				},
				headers: http.Header{},
			},
			want: &http.Request{
				Method: "GET",
				Header: http.Header{
					"foo": []string{"bar", "baz"},
				},
			},
		},
		{
			name: "returns request with headers added",
			args: args{
				req: &http.Request{
					Method: "GET",
					Header: http.Header{
						"foo": []string{"bar", "baz"},
					},
				},
				headers: http.Header{
					"password": []string{"my-secret-key"},
					"key-with": []string{"multiple", "values"},
				},
			},
			want: &http.Request{
				Method: "GET",
				Header: http.Header{
					"foo":      []string{"bar", "baz"},
					"Password": []string{"my-secret-key"},
					"Key-With": []string{"multiple", "values"},
				},
			},
		},
		{
			name: "returns nil when input request is nil",
			args: args{
				req: nil,
				headers: http.Header{
					"password": []string{"my-secret-key"},
					"key-with": []string{"multiple", "values"},
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := requestWithHeaders(tt.args.req, tt.args.headers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("requestWithHeaders() = %v, want %v", got, tt.want)
			}
		})
	}

}
