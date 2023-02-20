package wredis

import (
	"context"
	"github.com/guoyk93/winter"
	"github.com/redis/go-redis/v9"
)

// Get get previously injected [redis.Client]
func Get(ctx context.Context, opts ...Option) *redis.Client {
	opt := createOptions(opts...)
	return ctx.Value(opt.key).(*redis.Client)
}

// Install install component
func Install(a winter.App, opts ...Option) {
	opt := createOptions(opts...)

	var r *redis.Client

	a.Component("redis-" + string(opt.key)).
		Startup(func(ctx context.Context) (err error) {
			rOpts := &redis.Options{}
			if opt.opts != nil {
				rOpts = opt.opts
			} else if opt.url != "" {
				if rOpts, err = redis.ParseURL(opt.url); err != nil {
					return
				}
			}
			r = redis.NewClient(rOpts)
			return
		}).
		Check(func(ctx context.Context) error {
			return r.Ping(ctx).Err()
		}).
		Middleware(func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(c winter.Context) {
				c.Inject(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, opt.key, r)
				})
				h(c)
			}
		}).
		Shutdown(func(ctx context.Context) error {
			return r.Close()
		})

}
