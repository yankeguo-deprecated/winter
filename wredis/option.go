package wredis

import (
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type options struct {
	opts        *redis.Options
	url         string
	tracingOpts []redisotel.TracingOption
}

// Option option for installation
type Option = func(opts *options)

// WithURL set redis url
func WithURL(k string) Option {
	return func(opts *options) {
		opts.url = k
	}
}

// WithOptions set [redis.Options]
func WithOptions(rOpts *redis.Options) Option {
	return func(opts *options) {
		opts.opts = rOpts
	}
}

// WithTracingOptions set [redisotel.TracingOption]
func WithTracingOptions(os ...redisotel.TracingOption) Option {
	return func(opts *options) {
		opts.tracingOpts = append(opts.tracingOpts, os...)
	}
}
