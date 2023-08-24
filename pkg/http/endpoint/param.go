package endpoint

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var ErrInvalidURLParams = errors.New("endpoint: invalid url params")

type withParamEndpointOption struct {
	key   string
	value string
}

func (o *withParamEndpointOption) apply(
	ctx context.Context,
	opts *endpointOptions,
) error {
	opts.params[o.key] = o.value
	return nil
}

func WithParam[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string | bool](key string, value T) EndpointOption {
	return &withParamEndpointOption{
		key:   key,
		value: fmt.Sprintf("%v", value),
	}
}

var (
	paramStart = "{"
	paramEnd   = "}"
	tokenizer  = regexp.MustCompile(
		fmt.Sprintf("%s([^%s]+)%s",
			paramStart,
			paramStart,
			paramEnd,
		),
	)
	paramEndRegex = regexp.MustCompile(paramEnd)
)

func (o *endpointOptions) replaceURLParams(urlString string) (string, error) {
	missingKeys := []string{}
	parsedURL := tokenizer.ReplaceAllFunc([]byte(urlString), func(b []byte) []byte {
		key := string(b[1 : len(b)-1])
		if value, ok := o.params[key]; ok {
			return []byte(value)
		}
		missingKeys = append(missingKeys, key)
		return []byte{}
	})
	if len(missingKeys) > 0 {
		return "", fmt.Errorf("endpoint: missing params - %v", missingKeys)
	}
	return string(parsedURL), nil
}

func validateURLParams(urlStr string) error {
	paramsParts := strings.Split(urlStr, paramStart)
	for idx, part := range paramsParts {
		ends := paramEndRegex.FindAllIndex([]byte(part), 2)
		if idx == 0 {
			if len(ends) != 0 {
				return ErrInvalidURLParams
			}
		} else if len(ends) != 1 {
			return ErrInvalidURLParams
		}
	}
	return nil
}
