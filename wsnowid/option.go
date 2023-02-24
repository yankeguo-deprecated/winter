package wsnowid

type options struct {
	url       string
	restyKeys []string
}

// Option option for installation
type Option = func(opts *options)

// WithRestyKey set [wresty.KeyType]
func WithRestyKey(k ...string) Option {
	return func(opts *options) {
		opts.restyKeys = k
	}
}

// WithURL set url of snowid service
func WithURL(u string) Option {
	return func(opts *options) {
		opts.url = u
	}
}
