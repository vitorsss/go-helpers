package endpoint

import (
	"context"
	"fmt"
	"net/url"
)

type withQueryParamEndpointOption struct {
	key   string
	value string
}

func (o *withQueryParamEndpointOption) apply(
	ctx context.Context,
	opts *endpointOptions,
) error {
	if opts.query == nil {
		opts.query = url.Values{}
	}
	opts.query.Add(o.key, o.value)
	return nil
}

func WithQueryParam[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string | bool](key string, value T) EndpointOption {
	return &withQueryParamEndpointOption{
		key:   key,
		value: fmt.Sprintf("%v", value),
	}
}

type withQueryEndpointOption struct {
	query url.Values
}

func (o *withQueryEndpointOption) apply(
	ctx context.Context,
	opts *endpointOptions,
) error {
	if opts.query == nil {
		opts.query = o.query
	} else {
		for key, values := range o.query {
			for _, value := range values {
				opts.query.Add(key, value)
			}
		}
	}
	return nil
}

func WithQuery(query url.Values) EndpointOption {
	return &withQueryEndpointOption{
		query: query,
	}
}
