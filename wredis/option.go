package wredis

import "github.com/redis/go-redis/v9"

type KeyType string

const (
	Default KeyType = "default"
)

type options struct {
	key  KeyType
	opts *redis.Options
	url  string
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

// Option option for installation
type Option func(opts *options)

// WithKey set key for injection
func WithKey(k string) Option {
	return func(opts *options) {
		opts.key = KeyType(k)
	}
}

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
