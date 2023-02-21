package wgorm

import (
	"context"
	"errors"
	"github.com/guoyk93/winter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Get get previously injected [gorm.DB]
func Get(ctx context.Context, opts ...Option) *gorm.DB {
	o := createOptions(opts...)
	return ctx.Value(o.key).(*gorm.DB).WithContext(ctx)
}

// Install install component
func Install(a winter.App, opts ...Option) {
	o := createOptions(opts...)

	var (
		db *gorm.DB
	)

	a.Component("gorm-" + string(o.key)).
		Startup(func(ctx context.Context) (err error) {
			if o.mysqlConfig != nil {
				db, err = gorm.Open(mysql.New(*o.mysqlConfig), o.gormOptions...)
				return
			} else if o.mysqlDSN != "" {
				db, err = gorm.Open(mysql.Open(o.mysqlDSN), o.gormOptions...)
				return
			} else {
				err = errors.New("failed to initialize gorm component")
			}
			return
		}).
		Check(func(ctx context.Context) error {
			return db.WithContext(ctx).Select("SELECT 1").Error
		}).
		Middleware(func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(c winter.Context) {
				c.Inject(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, o.key, db)
				})
				h(c)
			}
		})

}
