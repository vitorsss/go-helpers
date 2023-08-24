package endpoint

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)

type endpointOptions struct {
	params       map[string]string
	body         io.Reader
	query        url.Values
	headers      http.Header
	errAnalyzers []ErrAnalyzerFn
}

type EndpointOption interface {
	apply(
		ctx context.Context,
		opts *endpointOptions,
	) error
}

func newEndpointOptions(
	ctx context.Context,
	baseOptions []EndpointOption,
	reqOptions []EndpointOption,
) (*endpointOptions, error) {
	opts := &endpointOptions{
		params: map[string]string{},
	}

	for _, opt := range baseOptions {
		err := opt.apply(ctx, opts)
		if err != nil {
			return nil, err
		}
	}

	for _, opt := range reqOptions {
		err := opt.apply(ctx, opts)
		if err != nil {
			return nil, err
		}
	}

	return opts, nil
}

func (e *endpoint) parseOptionsToRequest(
	ctx context.Context,
	method string,
	opts *endpointOptions,
) (*http.Request, error) {
	parsedURL, err := opts.replaceURLParams(e.URL)
	if err != nil {
		return nil, err
	}
	body := opts.body
	if body == nil {
		body = bytes.NewReader([]byte{})
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		parsedURL,
		body,
	)
	if err != nil {
		return nil, err
	}

	if opts.headers != nil {
		req.Header = opts.headers
	}
	if opts.query != nil {
		req.URL.RawQuery = opts.query.Encode()
	}

	return req, nil
}
