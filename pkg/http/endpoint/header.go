package endpoint

import (
	"context"
	"fmt"
	"net/http"
)

type withHeaderParamEndpointOption struct {
	key   string
	value string
}

func (o *withHeaderParamEndpointOption) apply(
	ctx context.Context,
	opts *endpointOptions,
) error {
	if opts.headers == nil {
		opts.headers = http.Header{}
	}
	opts.headers.Add(o.key, o.value)
	return nil
}

func WithHeaderParam[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string | bool](key string, value T) EndpointOption {
	return &withHeaderParamEndpointOption{
		key:   key,
		value: fmt.Sprintf("%v", value),
	}
}

type withHeaderEndpointOption struct {
	headers http.Header
}

func (o *withHeaderEndpointOption) apply(
	ctx context.Context,
	opts *endpointOptions,
) error {
	if opts.headers == nil {
		opts.headers = o.headers
	} else {
		for key, values := range o.headers {
			for _, value := range values {
				opts.headers.Add(key, value)
			}
		}
	}
	return nil
}

func WithHeader(headers http.Header) EndpointOption {
	return &withHeaderEndpointOption{
		headers: headers,
	}
}
