package wjwk

type options struct {
	raw []byte
}

// Option option for installation
type Option = func(opts *options)

// WithRaw set raw JWK
func WithRaw(buf []byte) Option {
	return func(opts *options) {
		opts.raw = buf
	}
}
