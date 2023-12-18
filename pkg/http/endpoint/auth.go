package endpoint

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/icholy/digest"
	"github.com/vitorsss/go-helpers/pkg/http/requester"
)

type AuthFn func(
	ctx context.Context,
	requester requester.Requester,
	request *http.Request,
) (*http.Response, error)

type withCustomAuthEndpointOption struct {
	authFn AuthFn
}

func (o *withCustomAuthEndpointOption) apply(
	ctx context.Context,
	opts *endpointOptions,
) error {
	opts.authFn = o.authFn
	return nil
}

func WithCustomAuth(authFn AuthFn) EndpointOption {
	return &withCustomAuthEndpointOption{
		authFn: authFn,
	}
}

type withBasicAuthEndpointOption struct {
	authValue string
}

func (o *withBasicAuthEndpointOption) apply(
	ctx context.Context,
	opts *endpointOptions,
) error {
	opts.authFn = func(
		ctx context.Context,
		requester requester.Requester,
		request *http.Request,
	) (*http.Response, error) {
		request.Header.Add("Authorization", o.authValue)
		return nil, nil
	}
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

type AuthHeaderFn func(ctx context.Context) (headerKey string, headerValue string, err error)

type withCustomAuthHeaderEndpointOption struct {
	authFn AuthHeaderFn
}

func (o *withCustomAuthHeaderEndpointOption) apply(
	ctx context.Context,
	opts *endpointOptions,
) error {
	opts.authFn = func(
		ctx context.Context,
		requester requester.Requester,
		request *http.Request,
	) (*http.Response, error) {
		key, value, err := o.authFn(ctx)
		if err != nil {
			return nil, err
		}
		request.Header.Add(key, value)
		return nil, nil
	}
	return nil
}

func WithCustomAuthHeader(authFn AuthHeaderFn) EndpointOption {
	return &withCustomAuthHeaderEndpointOption{
		authFn: authFn,
	}
}

type withDigestAuthEndpointOption struct {
	user string
	pass string
}

func (o *withDigestAuthEndpointOption) apply(
	ctx context.Context,
	opts *endpointOptions,
) error {
	opts.authFn = func(
		ctx context.Context,
		requester requester.Requester,
		request *http.Request,
	) (*http.Response, error) {
		res, err := requester.Do(request)
		if err != nil {
			return res, err
		}
		if res == nil {
			return nil, errors.New("WithDigestAuth.authFn: invalid response state")
		}
		if res.StatusCode != http.StatusUnauthorized {
			return res, nil
		}

		header := res.Header.Get("WWW-Authenticate")

		chal, err := digest.ParseChallenge(header)
		if err != nil {
			return nil, err
		}

		cred, err := digest.Digest(chal, digest.Options{
			Username: o.user,
			Password: o.pass,
			Method:   request.Method,
			URI:      request.URL.RequestURI(),
			GetBody:  request.GetBody,
			Count:    1,
		})
		if err != nil {
			return nil, err
		}
		request.Header.Set("Authorization", cred.String())

		return nil, nil
	}
	return nil
}

func WithDigestAuth(
	user string,
	pass string,
) EndpointOption {
	return &withDigestAuthEndpointOption{
		user: user,
		pass: pass,
	}
}
