package wgorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type KeyType string

const (
	Default KeyType = "default"
)

type options struct {
	key         KeyType
	mysqlDSN    string
	mysqlConfig *mysql.Config
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

// Option option for installation
type Option func(opts *options)

// WithKey set key for injection
func WithKey(k string) Option {
	return func(opts *options) {
		opts.key = KeyType(k)
	}
}

// WithMySQLDSN set MySQL DSN
func WithMySQLDSN(k string) Option {
	return func(opts *options) {
		opts.mysqlDSN = k
	}
}

// WithMySQLConfig set MySQL config
func WithMySQLConfig(cfg *mysql.Config) Option {
	return func(opts *options) {
		opts.mysqlConfig = cfg
	}
}

// WithGORMOption add [gorm.Option]
func WithGORMOption(o gorm.Option) Option {
	return func(opts *options) {
		opts.gormOptions = append(opts.gormOptions, o)
	}
}
