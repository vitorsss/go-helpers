package endpoint

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type ErrAnalyzerFn func(res *http.Response) error

type withCustomErrAnalyzerEndpointOption struct {
	errAnalyzer ErrAnalyzerFn
}

func (o *withCustomErrAnalyzerEndpointOption) apply(
	ctx context.Context,
	opts *endpointOptions,
) error {
	opts.errAnalyzers = append(opts.errAnalyzers, o.errAnalyzer)
	return nil
}

func WithCustomErrAnalyzer(errAnalyzer ErrAnalyzerFn) EndpointOption {
	return &withCustomErrAnalyzerEndpointOption{
		errAnalyzer: errAnalyzer,
	}
}

func defaultErrAnalyzer(res *http.Response) error {
	if res.StatusCode >= http.StatusBadRequest {
		content, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("endpoint: response error - %s", string(content))
	}
	return nil
}

func WithDefaultErrAnalyzer() EndpointOption {
	return &withCustomErrAnalyzerEndpointOption{
		errAnalyzer: defaultErrAnalyzer,
	}
}
