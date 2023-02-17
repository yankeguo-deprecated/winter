package wredis

type KeyType string

const (
	Default KeyType = "default"
)

// Option function modifying options
type Option func(opts *options)

// WithKey change key for injection
func WithKey(k KeyType) Option {
	return func(opts *options) {
		opts.key = k
	}
}

type options struct {
	key KeyType
}

func buildOptions(opts ...Option) *options {
	opt := &options{
		key: Default,
	}
	for _, item := range opts {
		item(opt)
	}
	return opt
}
