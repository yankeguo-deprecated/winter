package whcaptcha

import "github.com/guoyk93/winter/wresty"

type KeyType string

const (
	Default KeyType = "default"
)

type options struct {
	key      KeyType
	restyKey wresty.KeyType
	siteKey  string
	secret   string
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

// WithRestyKey set key for [wresty] extraction
func WithRestyKey(k string) Option {
	return func(opts *options) {
		opts.restyKey = wresty.KeyType(k)
	}
}

// WithSiteKey set siteKey of hcpatcha service
func WithSiteKey(u string) Option {
	return func(opts *options) {
		opts.siteKey = u
	}
}

// WithSecret set secret of hcaptcha service
func WithSecret(u string) Option {
	return func(opts *options) {
		opts.secret = u
	}
}
