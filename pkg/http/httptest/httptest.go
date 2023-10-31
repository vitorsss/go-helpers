package httptest

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/vitorsss/go-helpers/pkg/http/requester"
)

type Server interface {
	Requester() requester.Requester
	BaseURL(testing.T) string

	RequestRecorder
}

type RequestRecorder interface {
	Header(header http.Header) RequestRecorder
	Query(query url.Values) RequestRecorder
	Body(body []byte) RequestRecorder

	Get(
		path string,
	) ResponseRecorder
	Head(
		path string,
	) ResponseRecorder
	Post(
		path string,
	) ResponseRecorder
	Put(
		path string,
	) ResponseRecorder
	Patch(
		path string,
	) ResponseRecorder
	Delete(
		path string,
	) ResponseRecorder
	Connect(
		path string,
	) ResponseRecorder
	Options(
		path string,
	) ResponseRecorder
	Trace(
		path string,
	) ResponseRecorder
}

type ResponseRecorder interface {
	Return(
		status int,
		body []byte,
		header http.Header,
	) ResponseRecorder
	DoAndReturn(
		func(req *http.Request) (*http.Response, error),
	) ResponseRecorder
}
