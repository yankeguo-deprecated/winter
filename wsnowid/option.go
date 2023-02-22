package wsnowid

import "github.com/guoyk93/winter/wext"

type options struct {
	url       string
	restyKeys []string
}

// Option option for installation
type Option = func(opts *options)

var Ext = wext.New[options, *options]("snowid", func() *options {
	return &options{}
})

// WithRestyKey set [wresty.KeyType]
func WithRestyKey(k ...string) Option {
	return func(opts *options) {
		opts.restyKeys = k
	}
}

// WithURL set url of snowid service
func WithURL(u string) Option {
	return func(opts *options) {
		opts.url = u
	}
}
