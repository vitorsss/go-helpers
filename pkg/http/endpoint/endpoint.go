package endpoint

import (
	"context"
	"net/http"

	"github.com/vitorsss/go-helpers/pkg/http/requester"
)

type Endpoint interface {
	Get(
		ctx context.Context,
		options ...EndpointOption,
	) (Response, error)
	Head(
		ctx context.Context,
		options ...EndpointOption,
	) (Response, error)
	Post(
		ctx context.Context,
		options ...EndpointOption,
	) (Response, error)
	Put(
		ctx context.Context,
		options ...EndpointOption,
	) (Response, error)
	Patch(
		ctx context.Context,
		options ...EndpointOption,
	) (Response, error)
	Delete(
		ctx context.Context,
		options ...EndpointOption,
	) (Response, error)
	Connect(
		ctx context.Context,
		options ...EndpointOption,
	) (Response, error)
	Options(
		ctx context.Context,
		options ...EndpointOption,
	) (Response, error)
	Trace(
		ctx context.Context,
		options ...EndpointOption,
	) (Response, error)
}

type endpoint struct {
	BaseOptions []EndpointOption
	Requester   requester.Requester
	URL         string
}

func NewEndpoint(
	baseURI string,
	pathURI string,
	requester requester.Requester,
	options ...EndpointOption,
) Endpoint {
	urlStr, err := joinURL(baseURI, pathURI)
	if err != nil {
		panic(err)
	}

	err = validateURLParams(urlStr)
	if err != nil {
		panic(err)
	}

	return &endpoint{
		BaseOptions: options,
		Requester:   requester,
		URL:         urlStr,
	}
}

func (e *endpoint) Get(
	ctx context.Context,
	options ...EndpointOption,
) (Response, error) {
	return e.do(ctx,
		http.MethodGet,
		options...,
	)
}

func (e *endpoint) Head(
	ctx context.Context,
	options ...EndpointOption,
) (Response, error) {
	return e.do(ctx,
		http.MethodHead,
		options...,
	)
}

func (e *endpoint) Post(
	ctx context.Context,
	options ...EndpointOption,
) (Response, error) {
	return e.do(ctx,
		http.MethodPost,
		options...,
	)
}

func (e *endpoint) Put(
	ctx context.Context,
	options ...EndpointOption,
) (Response, error) {
	return e.do(ctx,
		http.MethodPut,
		options...,
	)
}

func (e *endpoint) Patch(
	ctx context.Context,
	options ...EndpointOption,
) (Response, error) {
	return e.do(ctx,
		http.MethodPatch,
		options...,
	)
}

func (e *endpoint) Delete(
	ctx context.Context,
	options ...EndpointOption,
) (Response, error) {
	return e.do(ctx,
		http.MethodDelete,
		options...,
	)
}

func (e *endpoint) Connect(
	ctx context.Context,
	options ...EndpointOption,
) (Response, error) {
	return e.do(ctx,
		http.MethodConnect,
		options...,
	)
}

func (e *endpoint) Options(
	ctx context.Context,
	options ...EndpointOption,
) (Response, error) {
	return e.do(ctx,
		http.MethodOptions,
		options...,
	)
}

func (e *endpoint) Trace(
	ctx context.Context,
	options ...EndpointOption,
) (Response, error) {
	return e.do(ctx,
		http.MethodTrace,
		options...,
	)
}

func (e *endpoint) do(
	ctx context.Context,
	method string,
	options ...EndpointOption,
) (Response, error) {
	opts, err := newEndpointOptions(
		ctx,
		e.BaseOptions,
		options,
	)
	if err != nil {
		return nil, err
	}

	req, err := e.parseOptionsToRequest(
		ctx,
		method,
		opts,
	)
	if err != nil {
		return nil, err
	}

	res, err := e.Requester.Do(req)
	if err != nil {
		return nil, err
	}

	for _, errAnalyzer := range opts.errAnalyzers {
		err = errAnalyzer(res)
		if err != nil {
			return nil, err
		}
	}

	return &response{
		Response: res,
	}, nil
}
