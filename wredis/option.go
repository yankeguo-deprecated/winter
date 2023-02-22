package wredis

import (
	"github.com/guoyk93/winter/wext"
	"github.com/redis/go-redis/v9"
)

type options struct {
	opts *redis.Options
	url  string
}

var Ext = wext.New[options, *redis.Client]("redis", func() *options {
	return &options{}
})

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
