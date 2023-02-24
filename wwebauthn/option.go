package wwebauthn

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

type options struct {
	cfg *webauthn.Config
}

// Option option for installation
type Option = func(opts *options)

// WithConfig set [webauthn.Config]
func WithConfig(cfg *webauthn.Config) Option {
	return func(opts *options) {
		opts.cfg = cfg
	}
}
