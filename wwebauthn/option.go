package wwebauthn

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/guoyk93/winter/wext"
)

type options struct {
	cfg *webauthn.Config
}

// Option option for installation
type Option = func(opts *options)

var Ext = wext.New[options, *webauthn.WebAuthn]("webauthn", func() *options {
	return &options{}
})

// WithConfig set [webauthn.Config]
func WithConfig(cfg *webauthn.Config) Option {
	return func(opts *options) {
		opts.cfg = cfg
	}
}
