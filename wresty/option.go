package wresty

import (
	"github.com/go-resty/resty/v2"
	"github.com/guoyk93/winter/wext"
	"net/http"
)

type options struct {
	rSetup  func(r *resty.Client) *resty.Client
	hcSetup func(c *http.Client) *http.Client
}

// Option option for installation
type Option = func(opts *options)

// Ext the [wext.Extension]
var Ext = wext.New[options, *resty.Client]("resty", func() *options {
	return &options{}
})

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
