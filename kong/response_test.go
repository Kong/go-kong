package kong

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageFromBody(T *testing.T) {
	for _, tt := range []struct {
		name     string
		response http.Response
		want     error
	}{
		{
			name: "code 200",
			response: http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			},
		},
		{
			name: "code 404",
			response: http.Response{
				StatusCode: 404,
				Body:       ioutil.NopCloser(strings.NewReader(`{"message": "potayto pohtato", "some": "other field"}`)),
			},
			want: &kongAPIError{
				httpCode: 404,
				message:  "potayto pohtato",
			},
		},
	} {
		T.Run(tt.name, func(T *testing.T) {
			got := hasError(&tt.response)
			assert.Equal(T, tt.want, got)
		})
	}
}
