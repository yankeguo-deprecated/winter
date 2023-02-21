package wgorm

import "gorm.io/gorm"

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

// WithMySQLDSN set env key for redis options loading
func WithMySQLDSN(k string) Option {
	return func(opts *options) {
		opts.mysqlDSN = k
	}
}

// WithGORMOption add [gorm.Option]
func WithGORMOption(o gorm.Option) Option {
	return func(opts *options) {
		opts.gormOptions = append(opts.gormOptions, o)
	}
}

type options struct {
	key         KeyType
	mysqlDSN    string
	gormOptions []gorm.Option
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
