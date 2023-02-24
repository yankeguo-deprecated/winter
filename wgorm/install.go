package wgorm

import (
	"context"
	"errors"
	"github.com/guoyk93/winter/wext"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	ext = wext.New[options, *injected]("gorm").
		Startup(func(ctx context.Context, opt *options) (inj *injected, err error) {
			inj = &injected{
				debug: opt.debug,
			}
			// create db
			if opt.mysqlConfig != nil {
				inj.db, err = gorm.Open(mysql.New(*opt.mysqlConfig), opt.gormOptions...)
			} else if opt.mysqlDSN != "" {
				inj.db, err = gorm.Open(mysql.Open(opt.mysqlDSN), opt.gormOptions...)
			} else {
				err = errors.New("failed to initialize gorm component")
				return
			}
			// instrument
			if err = inj.db.Use(otelgorm.NewPlugin(opt.tracingOpts...)); err != nil {
				return
			}
			return
		}).
		Check(func(ctx context.Context, inj *injected) error {
			return inj.db.WithContext(ctx).Select("SELECT 1").Error
		})
)

// DB get previously injected [gorm.DB]
func DB(ctx context.Context, altKeys ...string) *gorm.DB {
	inj := ext.Instance(altKeys...).Get(ctx)
	if inj.debug {
		return inj.db.WithContext(ctx).Debug()
	} else {
		return inj.db.WithContext(ctx)
	}
}

// Installer create [wext.Installer]
func Installer(opts ...Option) wext.Installer {
	return ext.Installer(opts...)
}
