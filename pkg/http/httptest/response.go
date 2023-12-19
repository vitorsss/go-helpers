package httptest

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"

	"olympos.io/encoding/edn"
)

type response struct {
	fn       func(req *http.Request) (*http.Response, error)
	maxTimes int
	minTimes int
	times    int
}

func (r *response) exec(req *http.Request) (*http.Response, error) {
	r.times++
	res, err := r.fn(req)
	if res != nil && err != nil {
		return nil, errors.New("httptest.Response: no response provided")
	}
	return res, err
}

func (s *request) Return(
	status int,
	body []byte,
	header http.Header,
) ResponseTimesRecorder {
	return s.DoAndReturn(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: status,
			Body:       io.NopCloser(bytes.NewReader(body)),
			Header:     header,
		}, nil
	})
}

func (s *request) ReturnJSON(
	status int,
	body interface{},
	header http.Header,
) ResponseTimesRecorder {
	data, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	header.Set("Content-Type", "application/json")
	return s.Return(status, data, header)
}

func (s *request) ReturnEDN(
	status int,
	body interface{},
	header http.Header,
) ResponseTimesRecorder {
	data, err := edn.Marshal(body)
	if err != nil {
		panic(err)
	}
	header.Set("Content-Type", "application/edn")
	return s.Return(status, data, header)
}

func (s *request) DoAndReturn(
	fn func(req *http.Request) (*http.Response, error),
) ResponseTimesRecorder {
	s.currentResponse = &response{
		fn:       fn,
		maxTimes: 1,
		minTimes: 1,
		times:    0,
	}
	s.responses = append(s.responses, s.currentResponse)
	return s
}

func (s *request) Times(times int) ResponseRecorder {
	s.currentResponse.maxTimes = times
	s.currentResponse.minTimes = times
	return s
}

func (s *request) MaxTimes(times int) {
	s.currentResponse.maxTimes = times
	s.currentResponse.minTimes = 0
}

func (s *request) MinTimes(times int) {
	s.currentResponse.maxTimes = math.MaxInt
	s.currentResponse.minTimes = times
}

func (s *request) AnyTimes() {
	s.currentResponse.maxTimes = math.MaxInt
	s.currentResponse.minTimes = 0
}
