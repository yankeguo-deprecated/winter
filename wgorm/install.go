package wgorm

import (
	"context"
	"github.com/guoyk93/winter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Get get previously injected [gorm.DB]
func Get(ctx context.Context, opts ...Option) *gorm.DB {
	opt := createOptions(opts...)
	return ctx.Value(opt.key).(*gorm.DB).WithContext(ctx)
}

// Install install component
func Install(a winter.App, opts ...Option) {
	opt := createOptions(opts...)

	var db *gorm.DB

	a.Component("gorm-" + string(opt.key)).
		Startup(func(ctx context.Context) (err error) {
			db, err = gorm.Open(mysql.Open(opt.mysqlDSN), opt.gormOptions...)
			return
		}).
		Check(func(ctx context.Context) error {
			return db.WithContext(ctx).Select("SELECT 1").Error
		}).
		Middleware(func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(c winter.Context) {
				c.Inject(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, opt.key, db)
				})
				h(c)
			}
		})

}
