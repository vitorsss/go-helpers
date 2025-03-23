package edntostruct

import "go/types"

type options struct {
	tagTypes   map[string]TypeFn
	namedTypes map[string]NamedTypeFn
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
			"any":     genericTypeFn(types.NewInterfaceType(nil, nil)),
			"keyword": KeywordTypeFn,
		},
		namedTypes: map[string]NamedTypeFn{},
	}
}

type Option func(opt *options)

func WithTagTypeFn(tag string, fn TypeFn) Option {
	return func(opt *options) {
		opt.tagTypes[tag] = fn
	}
}

func WithNamedTypeFn(name string, fn NamedTypeFn) Option {
	return func(opt *options) {
		opt.namedTypes[name] = fn
	}
}
