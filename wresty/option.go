package wresty

import (
	"github.com/go-resty/resty/v2"
	"net/http"
)

type options struct {
	rSetup  func(r *resty.Client) *resty.Client
	hcSetup func(c *http.Client) *http.Client
}

// Option option for installation
type Option = func(opts *options)

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
