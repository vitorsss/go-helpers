package httptest

import "net/http"

func CanonicalizeHeader(header http.Header) http.Header {
	result := http.Header{}
	for key, values := range header {
		for _, value := range values {
			result.Add(key, value)
		}
	}
	return result
}
