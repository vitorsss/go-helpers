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
	paramOr    = "|"
	tokenizer  = regexp.MustCompile(
		fmt.Sprintf("%s([^%s]+)%s",
			paramStart,
			paramStart,
			paramEnd,
		),
	)
	paramLiteralRegex = regexp.MustCompile(`^"([^"]+)"$`)
	paramEndRegex     = regexp.MustCompile(paramEnd)
	defaultParamValue = "__empty__"
)

func (o *endpointOptions) replaceURLParams(urlString string) (string, error) {
	return replaceURLParams(urlString, o.params, false)
}

func replaceURLParams(urlString string, params map[string]string, defaultParam bool) (string, error) {
	missingKeys := []string{}
	parsedURL := tokenizer.ReplaceAllFunc([]byte(urlString), func(b []byte) []byte {
		param := string(b[1 : len(b)-1])
		keys := strings.Split(param, paramOr)
		for _, key := range keys {
			if value, ok := params[key]; ok {
				return []byte(value)
			}
			literal := paramLiteralRegex.FindSubmatch([]byte(key))
			if literal != nil {
				return literal[1]
			}
		}
		if defaultParam {
			return []byte(defaultParamValue)
		}
		missingKeys = append(missingKeys, param)
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
