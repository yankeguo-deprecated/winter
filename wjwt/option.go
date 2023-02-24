package wjwt

type options struct {
	jwkKeys       []string
	payloadHeader string
	issuer        string
	debugPayload  bool
}

// Option option for installation
type Option = func(opts *options)

// WithJWKKey set key for JWK
func WithJWKKey(s ...string) Option {
	return func(opts *options) {
		opts.jwkKeys = s
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

// WithDebugPayload set debugPayload mode, when on, will extract Payload from 'Authorization' header without validating
func WithDebugPayload(d bool) Option {
	return func(opts *options) {
		opts.debugPayload = d
	}
}
