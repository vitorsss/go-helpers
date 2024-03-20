package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"olympos.io/encoding/edn"
)

type withRawBodyEndpointOption struct {
	body        io.ReadSeeker
	contentType string
	err         error
}

func (o *withRawBodyEndpointOption) apply(
	ctx context.Context,
	opts *endpointOptions,
) error {
	if o.err != nil {
		return o.err
	}
	_, err := o.body.Seek(0, io.SeekStart)
	if err != nil {
		o.err = errors.Wrap(err, "failed to reset body reader")
		return o.err
	}
	opts.body = o.body
	if o.contentType != "" {
		if opts.headers == nil {
			opts.headers = http.Header{}
		}
		opts.headers.Set("Content-Type", o.contentType)
	}
	return nil
}

func WithRawBody(contentType string, body io.ReadSeeker) EndpointOption {
	return &withRawBodyEndpointOption{
		contentType: contentType,
		body:        body,
	}
}

func WithJSONBody(content any) EndpointOption {
	data, err := json.Marshal(content)
	return &withRawBodyEndpointOption{
		contentType: "application/json",
		body:        bytes.NewReader(data),
		err:         errors.Wrap(err, "failed to marshal json body"),
	}
}

func WithEDNBody(content any) EndpointOption {
	data, err := edn.Marshal(content)
	return &withRawBodyEndpointOption{
		contentType: "application/edn",
		body:        bytes.NewReader(data),
		err:         errors.Wrap(err, "failed to marshal edn body"),
	}
}

func WithFormURLEncodedBody(content url.Values) EndpointOption {
	return &withRawBodyEndpointOption{
		contentType: "application/x-www-form-urlencoded",
		body:        strings.NewReader(content.Encode()),
	}
}
