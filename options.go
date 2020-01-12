package template

import (
	"net/http"
	"text/template"
)

type Options struct {
	Name     string
	Path     string
	GoSource bool
	FuncMap template.FuncMap
	Filesystem http.FileSystem
}

func Name(name string) Option {
	return func(opts *Options) {
		opts.Name = name
	}
}

// File sets the path to the source template file.
// The path is either relative to the user's CWD or, if an external http.Filesystem is provided, relative to that filesystem.
func Path(templatePath string) Option {
	return func(opts *Options) {
		opts.Path = templatePath
	}
}

func GoSource(isGoSource bool) Option {
	return func(opts *Options) {
		opts.GoSource = isGoSource
	}
}

func FuncMap(funcMap template.FuncMap) Option {
    return func(opts *Options) {
        opts.FuncMap = funcMap
    }
}

func UseFilesystem(fs http.FileSystem) Option {
    return func(opts *Options) {
    	opts.Filesystem = fs
    }
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Name:     "default_template",
		FuncMap: template.FuncMap{},
		GoSource: false,
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

type Option func(options *Options)
