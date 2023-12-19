package httptest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsss/go-helpers/pkg/assertutil"
	"github.com/vitorsss/go-helpers/pkg/test"
)

func Test_Server_ServeHTTP(t *testing.T) {
	type args struct {
		request *http.Request
	}

	type want struct {
		responseCode   int
		responseHeader http.Header
		responseBody   []byte
		err            error
		reporter       *regexp.Regexp
	}

	defaultHeader := http.Header{
		"some_header": []string{
			"header_value",
		},
	}

	emptyRegex := regexp.MustCompile("")

	tests := []struct {
		name string
		args func(baseURL string) *args
		want func() *want
		mock func(server Server, args *args, want *want)
	}{
		{
			name: "should return mocked GET response successfully using server URL",
			args: func(baseURL string) *args {
				url, err := url.Parse(fmt.Sprintf("%s/some/path?and=query", baseURL))
				if err != nil {
					panic(err)
				}
				return &args{
					request: &http.Request{
						Method: http.MethodGet,
						URL:    url,
						Header: defaultHeader,
					},
				}
			},
			want: func() *want {
				return &want{
					responseCode:   200,
					responseHeader: defaultHeader,
					responseBody:   []byte(`{"a": 1, "b": 2}`),
					err:            nil,
					reporter:       emptyRegex,
				}
			},
			mock: func(server Server, args *args, want *want) {
				server.Query(args.request.URL.Query()).
					Header(args.request.Header).
					Get(args.request.URL.Path).
					Return(
						want.responseCode,
						want.responseBody,
						want.responseHeader,
					).AnyTimes()
			},
		},
		{
			name: "should return mocked HEAD response successfully using server URL",
			args: func(baseURL string) *args {
				url, err := url.Parse(fmt.Sprintf("%s/some/path?and=query", baseURL))
				if err != nil {
					panic(err)
				}
				return &args{
					request: &http.Request{
						Method: http.MethodHead,
						URL:    url,
						Header: defaultHeader,
					},
				}
			},
			want: func() *want {
				return &want{
					responseCode:   200,
					responseHeader: defaultHeader,
					responseBody:   []byte(`{"a": 1, "b": 2}`),
					err:            nil,
					reporter:       emptyRegex,
				}
			},
			mock: func(server Server, args *args, want *want) {
				server.Query(args.request.URL.Query()).
					Header(args.request.Header).
					Head(args.request.URL.Path).
					Return(
						want.responseCode,
						want.responseBody,
						want.responseHeader,
					).MinTimes(1)
			},
		},
		{
			name: "should return mocked POST response successfully using server URL",
			args: func(baseURL string) *args {
				url, err := url.Parse(fmt.Sprintf("%s/some/path?and=query", baseURL))
				if err != nil {
					panic(err)
				}
				return &args{
					request: &http.Request{
						Method: http.MethodPost,
						URL:    url,
						Header: defaultHeader,
						Body:   io.NopCloser(bytes.NewReader([]byte(`{"a": 1, "b": 2, "c": 3}`))),
					},
				}
			},
			want: func() *want {
				return &want{
					responseCode:   200,
					responseHeader: defaultHeader,
					responseBody:   []byte(`{"a": 1, "b": 2}`),
					err:            nil,
					reporter:       emptyRegex,
				}
			},
			mock: func(server Server, args *args, want *want) {
				server.Query(args.request.URL.Query()).
					Header(args.request.Header).
					Body([]byte(`{"a": 1, "b": 2, "c": 3}`)).
					Post(args.request.URL.Path).
					Return(
						want.responseCode,
						want.responseBody,
						want.responseHeader,
					).MaxTimes(1)
			},
		},
		{
			name: "should return mocked PUT response successfully using server URL",
			args: func(baseURL string) *args {
				url, err := url.Parse(fmt.Sprintf("%s/some/path?and=query", baseURL))
				if err != nil {
					panic(err)
				}
				return &args{
					request: &http.Request{
						Method: http.MethodPut,
						URL:    url,
						Header: defaultHeader,
						Body:   io.NopCloser(bytes.NewReader([]byte(`{"a": 1, "b": 2, "c": 3}`))),
					},
				}
			},
			want: func() *want {
				return &want{
					responseCode:   200,
					responseHeader: defaultHeader,
					responseBody:   []byte(`{"a": 1, "b": 2}`),
					err:            nil,
					reporter:       emptyRegex,
				}
			},
			mock: func(server Server, args *args, want *want) {
				server.Query(args.request.URL.Query()).
					Header(args.request.Header).
					Body([]byte(`{"a": 1, "b": 2, "c": 3}`)).
					Put(args.request.URL.Path).
					Return(
						want.responseCode,
						want.responseBody,
						want.responseHeader,
					)
			},
		},
		{
			name: "should return mocked PATCH response successfully using server URL",
			args: func(baseURL string) *args {
				url, err := url.Parse(fmt.Sprintf("%s/some/path?and=query", baseURL))
				if err != nil {
					panic(err)
				}
				return &args{
					request: &http.Request{
						Method: http.MethodPatch,
						URL:    url,
						Header: defaultHeader,
						Body:   io.NopCloser(bytes.NewReader([]byte(`{"a": 1, "b": 2, "c": 3}`))),
					},
				}
			},
			want: func() *want {
				return &want{
					responseCode:   200,
					responseHeader: defaultHeader,
					responseBody:   []byte(`{"a": 1, "b": 2}`),
					err:            nil,
					reporter:       emptyRegex,
				}
			},
			mock: func(server Server, args *args, want *want) {
				server.Query(args.request.URL.Query()).
					Header(args.request.Header).
					Body([]byte(`{"a": 1, "b": 2, "c": 3}`)).
					Patch(args.request.URL.Path).
					Return(
						want.responseCode,
						want.responseBody,
						want.responseHeader,
					)
			},
		},
		{
			name: "should return mocked DELETE response successfully using server URL",
			args: func(baseURL string) *args {
				url, err := url.Parse(fmt.Sprintf("%s/some/path?and=query", baseURL))
				if err != nil {
					panic(err)
				}
				return &args{
					request: &http.Request{
						Method: http.MethodDelete,
						URL:    url,
						Header: defaultHeader,
					},
				}
			},
			want: func() *want {
				return &want{
					responseCode:   200,
					responseHeader: defaultHeader,
					responseBody:   []byte(`{"a": 1, "b": 2}`),
					err:            nil,
					reporter:       emptyRegex,
				}
			},
			mock: func(server Server, args *args, want *want) {
				server.Query(args.request.URL.Query()).
					Header(args.request.Header).
					Delete(args.request.URL.Path).
					Return(
						want.responseCode,
						want.responseBody,
						want.responseHeader,
					)
			},
		},
		{
			name: "should return mocked CONNECT response successfully using server URL",
			args: func(baseURL string) *args {
				url, err := url.Parse(fmt.Sprintf("%s/some/path?and=query", baseURL))
				if err != nil {
					panic(err)
				}
				return &args{
					request: &http.Request{
						Method: http.MethodConnect,
						URL:    url,
						Header: defaultHeader,
					},
				}
			},
			want: func() *want {
				return &want{
					responseCode:   200,
					responseHeader: defaultHeader,
					responseBody:   []byte(`{"a": 1, "b": 2}`),
					err:            nil,
					reporter:       emptyRegex,
				}
			},
			mock: func(server Server, args *args, want *want) {
				server.Query(args.request.URL.Query()).
					Header(args.request.Header).
					Connect(args.request.URL.Path).
					Return(
						want.responseCode,
						want.responseBody,
						want.responseHeader,
					)
			},
		},
		{
			name: "should return mocked OPTIONS response successfully using server URL",
			args: func(baseURL string) *args {
				url, err := url.Parse(fmt.Sprintf("%s/some/path?and=query", baseURL))
				if err != nil {
					panic(err)
				}
				return &args{
					request: &http.Request{
						Method: http.MethodOptions,
						URL:    url,
						Header: defaultHeader,
					},
				}
			},
			want: func() *want {
				return &want{
					responseCode:   200,
					responseHeader: defaultHeader,
					responseBody:   []byte(`{"a": 1, "b": 2}`),
					err:            nil,
					reporter:       emptyRegex,
				}
			},
			mock: func(server Server, args *args, want *want) {
				server.Query(args.request.URL.Query()).
					Header(args.request.Header).
					Options(args.request.URL.Path).
					Return(
						want.responseCode,
						want.responseBody,
						want.responseHeader,
					)
			},
		},
		{
			name: "should return mocked TRACE response successfully using server URL",
			args: func(baseURL string) *args {
				url, err := url.Parse(fmt.Sprintf("%s/some/path?and=query", baseURL))
				if err != nil {
					panic(err)
				}
				return &args{
					request: &http.Request{
						Method: http.MethodTrace,
						URL:    url,
						Header: defaultHeader,
					},
				}
			},
			want: func() *want {
				return &want{
					responseCode:   200,
					responseHeader: defaultHeader,
					responseBody:   []byte(`{"a": 1, "b": 2}`),
					err:            nil,
					reporter:       emptyRegex,
				}
			},
			mock: func(server Server, args *args, want *want) {
				server.Query(args.request.URL.Query()).
					Header(args.request.Header).
					Trace(args.request.URL.Path).
					Return(
						want.responseCode,
						want.responseBody,
						want.responseHeader,
					)
			},
		},
		{
			name: "should return mocked GET response successfully using generic URL",
			args: func(baseURL string) *args {
				url, err := url.Parse("http://test/some/path?and=query")
				if err != nil {
					panic(err)
				}
				return &args{
					request: &http.Request{
						Method: http.MethodGet,
						URL:    url,
						Header: defaultHeader,
					},
				}
			},
			want: func() *want {
				return &want{
					responseCode:   200,
					responseHeader: defaultHeader,
					responseBody:   []byte(`{"a": 1, "b": 2}`),
					err:            nil,
					reporter:       emptyRegex,
				}
			},
			mock: func(server Server, args *args, want *want) {
				server.Query(args.request.URL.Query()).
					Header(args.request.Header).
					Get(args.request.URL.Path).
					Return(
						want.responseCode,
						want.responseBody,
						want.responseHeader,
					)
			},
		},
		{
			name: "should return mocked GET response successfully using generic URL but reporting missing calls",
			args: func(baseURL string) *args {
				url, err := url.Parse("http://test/some/path?and=query")
				if err != nil {
					panic(err)
				}
				return &args{
					request: &http.Request{
						Method: http.MethodGet,
						URL:    url,
						Header: defaultHeader,
					},
				}
			},
			want: func() *want {
				return &want{
					responseCode:   200,
					responseHeader: defaultHeader,
					responseBody:   []byte(`{"a": 1, "b": 2}`),
					err:            nil,
					reporter:       regexp.MustCompile("\nMissing calls: 1 of 2 \n\tTrace: \t.+/go-helpers/pkg/http/httptest/server_test.go:\\d+\n\t       \t.+/go-helpers/pkg/http/httptest/server_test.go:\\d+\n\tMethod: GET - Path: /some/path Query: and=query - Header: map\\[some_header:\\[header_value\\]\\] - Body: \n"),
				}
			},
			mock: func(server Server, args *args, want *want) {
				server.Query(args.request.URL.Query()).
					Header(args.request.Header).
					Get(args.request.URL.Path).
					Return(
						want.responseCode,
						want.responseBody,
						want.responseHeader,
					).Times(2)
			},
		},
		{
			name: "should return mocked GET response successfully using generic URL but reporting too many calls",
			args: func(baseURL string) *args {
				url, err := url.Parse("http://test/some/path?and=query")
				if err != nil {
					panic(err)
				}
				return &args{
					request: &http.Request{
						Method: http.MethodGet,
						URL:    url,
						Header: defaultHeader,
					},
				}
			},
			want: func() *want {
				return &want{
					responseCode:   200,
					responseHeader: defaultHeader,
					responseBody:   []byte(`{"a": 1, "b": 2}`),
					err:            nil,
					reporter:       regexp.MustCompile("\nToo many calls: 1 of 0 \n\tTrace: \t.+/go-helpers/pkg/http/httptest/server_test.go:\\d+\n\t       \t.+/go-helpers/pkg/http/httptest/server_test.go:\\d+\n\tMethod: GET - Path: /some/path Query: and=query - Header: map\\[some_header:\\[header_value\\]\\] - Body: \n"),
				}
			},
			mock: func(server Server, args *args, want *want) {
				server.Query(args.request.URL.Query()).
					Header(args.request.Header).
					Get(args.request.URL.Path).
					Return(
						want.responseCode,
						want.responseBody,
						want.responseHeader,
					).Times(0)
			},
		},
		{
			name: "should return default error response for unmapped request",
			args: func(baseURL string) *args {
				url, err := url.Parse("http://test/some/path?and=query")
				if err != nil {
					panic(err)
				}
				return &args{
					request: &http.Request{
						Method: http.MethodGet,
						URL:    url,
						Header: defaultHeader,
					},
				}
			},
			want: func() *want {
				return &want{
					responseCode:   500,
					responseHeader: http.Header{},
					responseBody:   []byte(`httptest.Server: unmapped request`),
					err:            nil,
					reporter:       regexp.MustCompile("\nUnexpected calls: \n\tMethod: GET - Path: /some/path Query: and=query - Header: map\\[Some_header:\\[header_value\\]\\] - Body: \n"),
				}
			},
			mock: func(server Server, args *args, want *want) {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter := test.NewStringBuilderHelper()
			server := New(reporter)

			args, want := tt.args(server.BaseURL()), tt.want()

			tt.mock(server, args, want)

			response, err := server.Requester().Do(args.request)

			if assertutil.Error(t, want.err, err) {
				assert.Equal(t, want.responseCode, response.StatusCode)
				AssertHeader(t, want.responseHeader, response.Header)

				body, err := io.ReadAll(response.Body)
				if err != nil {
					panic(err)
				}

				assert.Equal(t, string(want.responseBody), string(body))

				assert.True(t, want.reporter.MatchString(reporter.Content()), "reporter content: ", reporter.Content())
			}
		})
	}
}
