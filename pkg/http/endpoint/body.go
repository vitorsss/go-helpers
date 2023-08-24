package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/url"
	"strings"

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
		o.err = err
		return o.err
	}
	opts.body = o.body
	if o.contentType != "" {
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
		err:         err,
	}
}

func WithEDNBody(content any) EndpointOption {
	data, err := edn.Marshal(content)
	return &withRawBodyEndpointOption{
		contentType: "application/json",
		body:        bytes.NewReader(data),
		err:         err,
	}
}

func WithFormURLEncodedBody(content url.Values) EndpointOption {
	return &withRawBodyEndpointOption{
		contentType: "application/x-www-form-urlencoded",
		body:        strings.NewReader(content.Encode()),
	}
}
