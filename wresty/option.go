package wresty

import (
	"github.com/go-resty/resty/v2"
	"net/http"
)

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

// WithRestySetup setup [resty.Client]
func WithRestySetup(fn func(r *resty.Client) *resty.Client) Option {
	return func(opts *options) {
		opts.rSetup = fn
	}
}

// WithHTTPClientSetup setup [http.HTTPClient]
func WithHTTPClientSetup(fn func(hc *http.Client) *http.Client) Option {
	return func(opts *options) {
		opts.hcSetup = fn
	}
}

type options struct {
	key     KeyType
	rSetup  func(r *resty.Client) *resty.Client
	hcSetup func(c *http.Client) *http.Client
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
