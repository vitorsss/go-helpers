package httptest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/vitorsss/go-helpers/pkg/http/requester"
	"github.com/vitorsss/go-helpers/pkg/test"
)

type server struct {
	internal         *httptest.Server
	lock             *sync.Mutex
	testHelper       test.TestHelper
	requests         map[requestKey][]*request
	unmappedRequests map[requestKey][]request
	currentRequest   *request
}

func New(reporter test.TestReporter) Server {
	server := &server{
		lock:             &sync.Mutex{},
		requests:         map[requestKey][]*request{},
		unmappedRequests: map[requestKey][]request{},
		testHelper:       test.AsHelper(reporter),
	}
	server.testHelper.Helper()

	server.internal = httptest.NewServer(server)

	if cleanuper, ok := test.IsCleanuper(reporter); ok {
		cleanuper.Cleanup(server.Cleanup)
	}

	return server
}

func (s *server) Cleanup() {
	s.testHelper.Helper()
	for key, requests := range s.requests {
		for _, request := range requests {
			min, max, times := request.countTimes()
			if times < min {
				s.testHelper.Errorf("\nMissing calls: %d of %d \n%s\t%s\n", times, min, request.traceString("\t"), request.line(key))
			} else if times > max {
				s.testHelper.Errorf("\nToo many calls: %d of %d \n%s\t%s\n", times, max, request.traceString("\t"), request.line(key))
			}
		}
	}
	if len(s.unmappedRequests) > 0 {
		unmappedRequestsSB := &strings.Builder{}
		unmappedRequestsSB.WriteString("\nUnexpected calls: \n")
		for key, requests := range s.unmappedRequests {
			for _, request := range requests {
				unmappedRequestsSB.WriteString(
					fmt.Sprintf("\t%s\n", request.line(key)),
				)
			}
		}

		s.testHelper.Errorf(unmappedRequestsSB.String())
	}

	s.internal.Close()
}

func (s *server) Requester() requester.Requester {
	return s
}

func (s *server) RoundTripper() http.RoundTripper {
	return s
}

func (s *server) BaseURL() string {
	return s.internal.URL
}

func (s *server) Do(req *http.Request) (*http.Response, error) {
	return s.RoundTrip(req)
}

func (s *server) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL == nil {
		return nil, errors.New("http: nil Request.URL")
	}

	req.Header = CanonicalizeHeader(req.Header)

	recorder := httptest.NewRecorder()
	s.ServeHTTP(recorder, req)
	return recorder.Result(), nil
}

func (s *server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	response, err := s.searchResponse(req)
	if err != nil {
		response = &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader(err.Error())),
		}
	}
	for key, values := range response.Header {
		for _, value := range values {
			rw.Header().Add(key, value)
		}
	}
	rw.WriteHeader(response.StatusCode)
	if response.Body != nil {
		_, err := io.Copy(rw, response.Body)
		if err != nil {
			panic(err)
		}
	}
}

func (s *server) searchResponse(req *http.Request) (*http.Response, error) {
	key := makeRequestKey(req.Method, req.URL.Path)

	var body []byte
	var err error

	if req.Body != nil {
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
	}
	request := &request{
		query:  req.URL.Query(),
		header: req.Header,
		body:   body,
	}
	requests := s.requests[key]

	for _, existingRequest := range requests {
		if existingRequest.match(request.query, request.header, request.body) {
			request = existingRequest
			break
		}
	}

	if len(request.responses) == 0 {
		s.unmappedRequests[key] = append(s.unmappedRequests[key], *request)
		return nil, errors.New("httptest.Server: unmapped request")
	}

	var response *response
	for _, response = range request.responses {
		if response.times < response.maxTimes {
			break
		}
	}

	req.Body = io.NopCloser(bytes.NewReader(body))
	return response.exec(req)
}
