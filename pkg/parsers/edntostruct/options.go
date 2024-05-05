package edntostruct

type options struct {
	tagTypes map[string]TypeFn
}

func defaultOptions() *options {
	return &options{
		tagTypes: map[string]TypeFn{
			"inst": TimeTypeFn,
			"uuid": UUIDTypeFn,
		},
	}
}

type Option func(opt *options)

func WithTagTypeFn(tag string, fn TypeFn) Option {
	return func(opt *options) {
		opt.tagTypes[tag] = fn
	}
}
