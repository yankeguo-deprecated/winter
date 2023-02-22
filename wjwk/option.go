package wjwk

import (
	"github.com/guoyk93/winter/wext"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type options struct {
	raw []byte
}

// Option option for installation
type Option = func(opts *options)

var Ext = wext.New[options, jwk.Key]("jwk", func() *options {
	return &options{}
})

// WithRaw set raw JWK
func WithRaw(buf []byte) Option {
	return func(opts *options) {
		opts.raw = buf
	}
}
