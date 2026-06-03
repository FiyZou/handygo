package safego

import (
	"context"
	"fmt"
	"runtime/debug"
)

type Logger interface {
	Errorw(msg string, keysAndValues ...any)
}

type Option func(*options)

type options struct {
	logger  Logger
	onPanic func(recovered any, stack []byte)
}

func WithLogger(logger Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

func WithPanicHandler(handler func(recovered any, stack []byte)) Option {
	return func(opts *options) {
		opts.onPanic = handler
	}
}

func Go(fn func(), opts ...Option) {
	option := apply(opts...)
	go run(option, fn)
}

func GoContext(ctx context.Context, fn func(context.Context), opts ...Option) {
	option := apply(opts...)
	go run(option, func() {
		fn(ctx)
	})
}

func Run(fn func(), opts ...Option) {
	run(apply(opts...), fn)
}

func run(opts options, fn func()) {
	defer func() {
		if recovered := recover(); recovered != nil {
			stack := debug.Stack()
			if opts.logger != nil {
				opts.logger.Errorw("goroutine panic recovered", "panic", fmt.Sprint(recovered), "stack", string(stack))
			}
			if opts.onPanic != nil {
				opts.onPanic(recovered, stack)
			}
		}
	}()
	fn()
}

func apply(opts ...Option) options {
	option := options{}
	for _, opt := range opts {
		opt(&option)
	}
	return option
}
