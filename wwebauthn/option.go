package wwebauthn

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

type KeyType string

const (
	Default KeyType = "default"
)

type options struct {
	key KeyType
	cfg *webauthn.Config
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

// Option option for installation
type Option func(opts *options)

// WithKey set key for injection
func WithKey(k string) Option {
	return func(opts *options) {
		opts.key = KeyType(k)
	}
}

// WithConfig set [webauthn.Config]
func WithConfig(cfg *webauthn.Config) Option {
	return func(opts *options) {
		opts.cfg = cfg
	}
}
