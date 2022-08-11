package gopool

import "fmt"

var (
	defaultOptions = func() *Options {
		return &Options{
			panicHandler: func(err interface{}) {
				fmt.Printf("panic error:%+v", err)
			},
		}
	}
)

type Options struct {
	size         int
	queue        int
	spawn        int
	panicHandler func(err interface{})
}

type Option func(opts *Options)

func loadOptions(options ...Option) *Options {
	opts := defaultOptions()
	for _, option := range options {
		option(opts)
	}
	return opts
}

// WithPanicHandler set panicHandler
func WithPanicHandler(panicHandler func(err interface{})) Option {
	return func(opts *Options) {
		opts.panicHandler = panicHandler
	}
}

// WithQueueLength set queue length
func WithQueueLength(length int) Option {
	return func(opts *Options) {
		opts.queue = length
	}
}

// WithSpawnSize set spawn size
func WithSpawnSize(size int) Option {
	return func(opts *Options) {
		opts.spawn = size
	}
}
