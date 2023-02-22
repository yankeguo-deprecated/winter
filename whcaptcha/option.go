package whcaptcha

import (
	"github.com/guoyk93/winter/wext"
)

type options struct {
	restyKeys []string
	siteKey   string
	secret    string
}

// Option option for installation
type Option = func(opts *options)

// Ext the [wext.Extension]
var Ext = wext.New[options, *options]("hcaptcha", func() *options {
	return &options{}
})

// WithRestyKey set key for [wresty] extraction
func WithRestyKey(k ...string) Option {
	return func(opts *options) {
		opts.restyKeys = k
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
