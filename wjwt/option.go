package wjwt

import "github.com/guoyk93/winter/wjwk"

type KeyType string

const (
	Default KeyType = "default"
)

type options struct {
	key KeyType

	jwkKey        wjwk.KeyType
	payloadHeader string
	issuer        string
}

func createOptions(opts ...Option) *options {
	opt := &options{
		key:           Default,
		jwkKey:        wjwk.Default,
		payloadHeader: "X-JWT-Payload",
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

// WithJWKKey set key for JWK
func WithJWKKey(s string) Option {
	return func(opts *options) {
		opts.jwkKey = wjwk.KeyType(s)
	}
}

// WithIssuer set issuer
func WithIssuer(s string) Option {
	return func(opts *options) {
		opts.issuer = s
	}
}

// WithPayloadHeader set payload header set by Istio RequestAuthentication
func WithPayloadHeader(s string) Option {
	return func(opts *options) {
		opts.payloadHeader = s
	}
}
