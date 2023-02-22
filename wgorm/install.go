package wgorm

import (
	"context"
	"errors"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wext"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB get previously injected [gorm.DB]
func DB(ctx context.Context, altKeys ...string) *gorm.DB {
	inj := Ext.Instance(altKeys...).Get(ctx)
	if inj.debug {
		return inj.db.WithContext(ctx).Debug()
	} else {
		return inj.db.WithContext(ctx)
	}
}

// Installer create [wext.Installer]
func Installer(opts ...Option) wext.Installer {
	o := Ext.Options(opts...)

	return wext.WrapInstaller(func(a winter.App, altKeys ...string) {
		ins := Ext.Instance(altKeys...)

		inj := &injected{
			debug: o.debug,
		}

		a.Component(ins.Key()).
			Startup(func(ctx context.Context) (err error) {
				if o.mysqlConfig != nil {
					inj.db, err = gorm.Open(mysql.New(*o.mysqlConfig), o.gormOptions...)
					return
				} else if o.mysqlDSN != "" {
					inj.db, err = gorm.Open(mysql.Open(o.mysqlDSN), o.gormOptions...)
					return
				} else {
					err = errors.New("failed to initialize gorm component")
				}
				return
			}).
			Check(func(ctx context.Context) error {
				return inj.db.WithContext(ctx).Select("SELECT 1").Error
			}).
			Middleware(ins.Middleware(&inj))
	})
}
