package httptest

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertHeader(t *testing.T, expected http.Header, current http.Header) bool {
	for key, expectedValues := range expected {
		values := current.Values(key)
		if !assert.Equal(t, expectedValues, values, fmt.Sprintf("header: %s", key)) {
			return false
		}
	}
	return true
}
