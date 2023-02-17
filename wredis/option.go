package wredis

import "github.com/redis/go-redis/v9"

type KeyType string

const (
	Default KeyType = "default"
)

// Option function modifying options
type Option func(opts *options)

// WithKey change key for injection
func WithKey(k KeyType) Option {
	return func(opts *options) {
		opts.key = k
	}
}

// WithOptions with [redis.Options]
func WithOptions(rOpts *redis.Options) Option {
	return func(opts *options) {
		opts.opts = rOpts
	}
}

type options struct {
	key  KeyType
	opts *redis.Options
}

func createOptions(opts ...Option) *options {
	opt := &options{
		key: Default,
	}
	for _, item := range opts {
		item(opt)
	}
	return opt
}
