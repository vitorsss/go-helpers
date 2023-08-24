package endpoint

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
)

type withBasicAuthEndpointOption struct {
	authValue string
}

func (o *withBasicAuthEndpointOption) apply(
	ctx context.Context,
	opts *endpointOptions,
) error {
	if opts.headers == nil {
		opts.headers = http.Header{}
	}
	opts.headers.Add(
		"Authorization",
		o.authValue,
	)
	return nil
}

func WithBasicAuth(
	user string,
	pass string,
) EndpointOption {
	return &withBasicAuthEndpointOption{
		authValue: fmt.Sprintf(
			"Basic %s",
			base64.StdEncoding.EncodeToString([]byte(
				fmt.Sprintf("%s:%s",
					user,
					pass,
				),
			)),
		),
	}
}

func WithBearerTokenAuth(
	token string,
) EndpointOption {
	return &withBasicAuthEndpointOption{
		authValue: fmt.Sprintf(
			"Bearer %s",
			token,
		),
	}
}

type AuthFn func(ctx context.Context) (headerKey string, headerValue string, err error)

type withCustomAuthEndpointOption struct {
	authFn AuthFn
}

func (o *withCustomAuthEndpointOption) apply(
	ctx context.Context,
	opts *endpointOptions,
) error {
	key, value, err := o.authFn(ctx)
	if err != nil {
		return err
	}
	if opts.headers == nil {
		opts.headers = http.Header{}
	}
	opts.headers.Add(
		key,
		value,
	)
	return nil
}

func WithCustomAuth(authFn AuthFn) EndpointOption {
	return &withCustomAuthEndpointOption{
		authFn: authFn,
	}
}
