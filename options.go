package micro

import (
	"context"
	"sitesoft/coder.mvd/micro/meta"
)

type Option func(*Options)

type Options struct {
	Name string

	// Before and After funcs
	BeforeStart []func(s *Service) error
	BeforeStop  []func(s *Service) error
	AfterStart  []func(s *Service) error
	AfterStop   []func(s *Service) error

	// Other options for implementations of the interface
	// can be stored in a context
	Context  context.Context
	Metadata map[string]string
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Context:  context.Background(),
		Metadata: meta.Metadata{},
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func Name(c string) Option {
	return func(o *Options) {
		o.Name = c
	}
}

func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Metadata = md
	}
}

// Before and Afters

func BeforeStart(fn func(s *Service) error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

func BeforeStop(fn func(s *Service) error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

func AfterStart(fn func(s *Service) error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

func AfterStop(fn func(s *Service) error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}
