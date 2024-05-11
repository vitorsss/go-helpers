package edntostruct

import "go/types"

type options struct {
	tagTypes map[string]TypeFn
}

func defaultOptions() *options {
	return &options{
		tagTypes: map[string]TypeFn{
			"inst":    TimeTypeFn,
			"uuid":    UUIDTypeFn,
			"string":  genericTypeFn(types.Typ[types.String]),
			"bool":    genericTypeFn(types.Typ[types.Bool]),
			"float64": genericTypeFn(types.Typ[types.Float64]),
			"int64":   genericTypeFn(types.Typ[types.Int64]),
		},
	}
}

type Option func(opt *options)

func WithTagTypeFn(tag string, fn TypeFn) Option {
	return func(opt *options) {
		opt.tagTypes[tag] = fn
	}
}
