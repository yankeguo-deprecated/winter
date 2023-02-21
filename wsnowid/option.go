package wsnowid

import "github.com/guoyk93/winter/wresty"

type KeyType string

const (
	Default KeyType = "default"
)

type options struct {
	key      KeyType
	url      string
	restyKey wresty.KeyType
}

func createOptions(opts ...Option) *options {
	opt := &options{
		key:      Default,
		restyKey: wresty.Default,
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

// WithRestyKey set [wresty.KeyType]
func WithRestyKey(k string) Option {
	return func(opts *options) {
		opts.restyKey = wresty.KeyType(k)
	}
}

// WithURL set url of snowid service
func WithURL(u string) Option {
	return func(opts *options) {
		opts.url = u
	}
}
