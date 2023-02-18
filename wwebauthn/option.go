package wwebauthn

import (
	"github.com/go-webauthn/webauthn/webauthn"
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

// WithOptions update options
func WithOptions(cfg *webauthn.Config) Option {
	return func(opts *options) {
		opts.cfg = cfg
	}
}

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
