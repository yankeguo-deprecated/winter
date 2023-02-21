package wjwk

type KeyType string

const (
	Default KeyType = "default"
)

type options struct {
	key KeyType

	raw []byte
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

// WithRaw set raw JWK
func WithRaw(buf []byte) Option {
	return func(opts *options) {
		opts.raw = buf
	}
}
