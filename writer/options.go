package writer

type Options struct {
	Overwrite bool
	Append    bool
}

func Overwrite(overwrite bool) Option {
    return func(opts *Options) {
        opts.Overwrite = overwrite
    }
}

func Append(appendWrite bool) Option {
    return func(opts *Options) {
        opts.Append = appendWrite
    }
}

func newOptions(options ...Option) Options {
	opt := Options{
		Overwrite: false,
		Append:    false,
	}

	for _, o := range options {
		o(&opt)
	}

	return opt
}
