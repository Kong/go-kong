package kong

import (
	"fmt"
)

type KongAPIError struct {
	httpCode int
	message  string
}

func (e *KongAPIError) Error() string {
	return fmt.Sprintf("HTTP status %d (message: %q)", e.httpCode, e.message)
}

func (e *KongAPIError) Code() int {
	return e.httpCode
}

// IsNotFoundErr returns true if the error or it's cause is
// a 404 response from Kong.
func IsNotFoundErr(e error) bool {
	switch e := e.(type) {
	case *KongAPIError:
		return e.httpCode == 404
	default:
		return false
	}
}
