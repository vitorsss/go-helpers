package httptest

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"sort"
	"strings"

	"github.com/vitorsss/go-helpers/pkg/test"
)

type request struct {
	query           url.Values
	header          http.Header
	body            []byte
	responses       []*response
	currentResponse *response
	callerInfo      []string
}

func (r *request) traceString(prefix string) string {
	sb := &strings.Builder{}
	sb.WriteString(fmt.Sprintf("%sTrace: ", prefix))
	tab := "\t"
	for idx, callerInfo := range r.callerInfo {
		sb.WriteString(fmt.Sprintf("%s%s\n", tab, callerInfo))
		if idx == 0 {
			tab = fmt.Sprintf("%s       \t", prefix)
		}
	}
	return sb.String()
}

func (r *request) line(key requestKey) string {
	return fmt.Sprintf("Method: %s - Path: %s %s",
		key.method,
		key.path,
		r.String(),
	)
}

func (r *request) String() string {
	return fmt.Sprintf("Query: %s - Header: %v - Body: %s",
		r.query.Encode(),
		r.header,
		string(r.body),
	)
}

func matchURLValues(
	expected url.Values,
	current url.Values,
) bool {
	if len(expected) != len(current) {
		return false
	}
	for key, expectedValues := range expected {
		currentValues, ok := current[key]
		if !ok {
			return false
		}
		sort.Strings(expectedValues)
		sort.Strings(currentValues)
		if !slices.Equal(expectedValues, currentValues) {
			return false
		}
	}
	return true
}

func matchHeader(
	expected http.Header,
	current http.Header,
) bool {
	for key, expectedValues := range expected {
		currentValues := current.Values(key)
		sort.Strings(expectedValues)
		sort.Strings(currentValues)
		if !slices.Equal(expectedValues, currentValues) {
			return false
		}
	}
	return true
}

func (r *request) match(
	query url.Values,
	header http.Header,
	body []byte,
) bool {
	return matchURLValues(r.query, query) &&
		matchHeader(r.header, header) &&
		bytes.Equal(r.body, body)
}

func (r *request) countTimes() (minTimes, maxTimes, times int) {
	for _, response := range r.responses {
		minTimes += response.minTimes
		maxTimes += response.maxTimes
		times += response.times
	}
	return
}

type requestKey struct {
	method string
	path   string
}

func makeRequestKey(
	method string,
	path string,
) requestKey {
	return requestKey{
		method: method,
		path:   path,
	}
}

func (s *server) Header(header http.Header) RequestRecorder {
	if s.currentRequest == nil {
		s.currentRequest = &request{}
	}
	s.currentRequest.header = header
	return s
}

func (s *server) Query(query url.Values) RequestRecorder {
	if s.currentRequest == nil {
		s.currentRequest = &request{}
	}
	s.currentRequest.query = query
	return s
}

func (s *server) Body(body []byte) RequestRecorder {
	if s.currentRequest == nil {
		s.currentRequest = &request{}
	}
	s.currentRequest.body = body
	return s
}

func (s *server) record(
	method string,
	path string,
) *request {
	if s.currentRequest == nil {
		s.currentRequest = &request{}
	}
	s.currentRequest.callerInfo = test.CallerInfo("httptest")
	key := makeRequestKey(method, path)
	request := s.currentRequest
	s.currentRequest = nil
	s.requests[key] = append(s.requests[key], request)
	return request
}

func (s *server) Get(
	path string,
) ResponseRecorder {
	return s.record(http.MethodGet, path)
}

func (s *server) Head(
	path string,
) ResponseRecorder {
	return s.record(http.MethodHead, path)
}

func (s *server) Post(
	path string,
) ResponseRecorder {
	return s.record(http.MethodPost, path)
}

func (s *server) Put(
	path string,
) ResponseRecorder {
	return s.record(http.MethodPut, path)
}

func (s *server) Patch(
	path string,
) ResponseRecorder {
	return s.record(http.MethodPatch, path)
}

func (s *server) Delete(
	path string,
) ResponseRecorder {
	return s.record(http.MethodDelete, path)
}

func (s *server) Connect(
	path string,
) ResponseRecorder {
	return s.record(http.MethodConnect, path)
}

func (s *server) Options(
	path string,
) ResponseRecorder {
	return s.record(http.MethodOptions, path)
}

func (s *server) Trace(
	path string,
) ResponseRecorder {
	return s.record(http.MethodTrace, path)
}
