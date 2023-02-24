package wgorm

import (
	"github.com/guoyk93/winter/wext"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type options struct {
	mysqlDSN    string
	mysqlConfig *mysql.Config
	gormOptions []gorm.Option
	otelOptions []otelgorm.Option
	debug       bool
}

type injected struct {
	db    *gorm.DB
	debug bool
}

// Option option for installation
type Option = func(opts *options)

// Ext the [wext.Extension]
var Ext = wext.New[options, *injected]("gorm", func() *options {
	return &options{}
})

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

// WithGORMOptions add [gorm.Option]
func WithGORMOptions(os ...gorm.Option) Option {
	return func(opts *options) {
		opts.gormOptions = append(opts.gormOptions, os...)
	}
}

// WithOTELOptions add [otelgorm.Option]
func WithOTELOptions(os ...otelgorm.Option) Option {
	return func(opts *options) {
		opts.otelOptions = append(opts.otelOptions, os...)
	}
}

// WithDebug set debug
func WithDebug(d bool) Option {
	return func(opts *options) {
		opts.debug = d
	}
}
