package endpoint

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitorsss/go-helpers/pkg/assertutil"
	"github.com/vitorsss/go-helpers/pkg/http/httptest"
)

type someType struct {
	ID int `json:"id" edn:"id"`
}

func Test_NewEndpoint(t *testing.T) {
	end := NewEndpoint(
		"http://example.com/base/path",
		"/api/v1/as/{paramA}/bs/{paramB}",
		nil,
	)
	assert.NotNil(t, end)
}

func Test_Endpoint_Get(t *testing.T) {
	server := httptest.New(t)
	end := NewEndpoint(
		server.BaseURL(),
		"/api/v1/as/{paramA}",
		server.Requester(),
		WithQueryParam("test", "aaa"),
	)

	expected := someType{
		ID: 42,
	}

	server.
		Query(url.Values{
			"test": []string{
				"aaa",
			},
		}).
		Get("/api/v1/as/some_value").
		ReturnJSON(
			200,
			expected,
			http.Header{},
		).
		Return(
			404,
			[]byte(`Not Found`),
			http.Header{},
		)

	response, err := end.Get(context.Background(),
		WithParam("paramA", "some_value"),
	)
	if !assertutil.Error(t, nil, err) {
		return
	}

	assert.Equal(t, 200, response.Status())

	result := someType{}

	err = response.Unmarshal(&result)
	if !assertutil.Error(t, nil, err) {
		return
	}

	assert.Equal(t, expected, result)

	response, err = end.Get(context.Background(),
		WithParam("paramA", "some_value"),
	)
	if !assertutil.Error(t, nil, err) {
		return
	}

	assert.Equal(t, 404, response.Status())
}
