package kong

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHasError(T *testing.T) {
	for _, tt := range []struct {
		name     string
		response http.Response
		want     error
	}{
		{
			name: "code 200",
			response: http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("")),
			},
		},
		{
			name: "code 404",
			response: http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(strings.NewReader(`{"message": "potayto pohtato", "some": "other field"}`)),
			},
			want: &APIError{
				httpCode: 404,
				message:  "potayto pohtato",
			},
		},
		{
			name: "code 404, message field missing",
			response: http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(strings.NewReader(`{"nothing": "nothing"}`)),
			},
			want: &APIError{
				httpCode: 404,
				message:  "",
			},
		},
		{
			name: "code 404, empty body",
			response: http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(strings.NewReader(``)),
			},
			want: &APIError{
				httpCode: 404,
				message:  "<failed to parse response body: unexpected end of JSON input>",
			},
		},
		{
			name: "code 404, unparseable json",
			response: http.Response{
				StatusCode: 404,
				Body:       io.NopCloser(strings.NewReader(`This is not json`)),
			},
			want: &APIError{
				httpCode: 404,
				message:  "<failed to parse response body: invalid character 'T' looking for beginning of value>",
			},
		},
		{
			name: "code 429 with retry-after header",
			response: http.Response{
				StatusCode: http.StatusTooManyRequests,
				Body:       io.NopCloser(strings.NewReader("")),
				Header: map[string][]string{
					"Retry-After": {"123"},
				},
			},
			want: &APIError{
				httpCode: http.StatusTooManyRequests,
				message:  "<failed to parse response body: unexpected end of JSON input>",
				details: ErrTooManyRequestsDetails{
					RetryAfter: time.Second * 123,
				},
			},
		},
		{
			name: "code 429 with no retry-after header",
			response: http.Response{
				StatusCode: http.StatusTooManyRequests,
				Body:       io.NopCloser(strings.NewReader("")),
			},
			want: &APIError{
				httpCode: http.StatusTooManyRequests,
				message:  "<failed to parse response body: unexpected end of JSON input>",
			},
		},
	} {
		T.Run(tt.name, func(T *testing.T) {
			tt := tt
			got := hasError(&tt.response)
			assert.Equal(T, tt.want, got)
		})
	}
}
