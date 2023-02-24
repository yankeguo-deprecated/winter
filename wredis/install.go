package wredis

import (
	"context"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wext"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

// Client get previously injected [redis.Client]
func Client(ctx context.Context, altKeys ...string) *redis.Client {
	return Ext.Instance(altKeys...).Get(ctx)
}

// Installer install component
func Installer(opts ...Option) wext.Installer {
	o := Ext.Options(opts...)

	return wext.WrapInstaller(func(a winter.App, altKeys ...string) {
		ins := Ext.Instance(altKeys...)

		var r *redis.Client

		a.Component(ins.Key()).
			Startup(func(ctx context.Context) (err error) {
				// create options
				rOpts := &redis.Options{}
				if o.opts != nil {
					rOpts = o.opts
				} else if o.url != "" {
					if rOpts, err = redis.ParseURL(o.url); err != nil {
						return
					}
				}
				// create client
				r = redis.NewClient(rOpts)
				// instrument
				if err = redisotel.InstrumentTracing(r, o.tracingOpts...); err != nil {
					return
				}
				return
			}).
			Check(func(ctx context.Context) error {
				return r.Ping(ctx).Err()
			}).
			Middleware(ins.Middleware(&r)).
			Shutdown(func(ctx context.Context) error {
				return r.Close()
			})

	})
}
