package endpoint

import (
	"context"
	"io"
	"net/http"

	"github.com/pkg/errors"
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
	if res.StatusCode >= http.StatusBadRequest && res.StatusCode != http.StatusNotFound {
		content, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.Errorf("endpoint: response error - %s", string(content))
	}
	return nil
}

func WithDefaultErrAnalyzer() EndpointOption {
	return &withCustomErrAnalyzerEndpointOption{
		errAnalyzer: defaultErrAnalyzer,
	}
}

func notFoundErrAnalyzer(res *http.Response) error {
	if res.StatusCode == http.StatusNotFound {
		content, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.Errorf("endpoint: response error - %s", string(content))
	}
	return nil
}

func WithNotFoundErrAnalyzer() EndpointOption {
	return &withCustomErrAnalyzerEndpointOption{
		errAnalyzer: notFoundErrAnalyzer,
	}
}
