package whcaptcha

import "github.com/guoyk93/winter/wresty"

type KeyType string

const (
	Default KeyType = "default"
)

// Option function modifying options
type Option func(opts *options)

// WithKey change key for injection
func WithKey(k string) Option {
	return func(opts *options) {
		opts.key = KeyType(k)
	}
}

// WithRestyKey change resty client key
func WithRestyKey(k string) Option {
	return func(opts *options) {
		opts.rKey = wresty.KeyType(k)
	}
}

// WithSiteKey set siteKey of snowid service
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

type options struct {
	key     KeyType
	rKey    wresty.KeyType
	siteKey string
	secret  string
}

func createOptions(opts ...Option) *options {
	opt := &options{
		key:  Default,
		rKey: wresty.Default,
	}
	for _, item := range opts {
		item(opt)
	}
	return opt
}
