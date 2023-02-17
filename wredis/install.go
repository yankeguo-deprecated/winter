package wredis

import (
	"context"
	"github.com/guoyk93/winter"
	"github.com/redis/go-redis/v9"
	"strings"
)

// Get get component
func Get(ctx winter.Context, opts ...Option) *redis.Client {
	opt := createOptions(opts...)
	return ctx.Value(opt.key).(*redis.Client)
}

// Install install component
func Install(a winter.App, opts ...Option) {
	opt := createOptions(opts...)

	var r *redis.Client

	a.Component("redis").
		Startup(func(ctx context.Context) (err error) {
			var rOpts *redis.Options
			if opt.opts != nil {
				rOpts = opt.opts
			} else if envRedisURL := strings.TrimSpace("REDIS_URL"); envRedisURL == "" {
				rOpts = &redis.Options{}
			} else if rOpts, err = redis.ParseURL(envRedisURL); err != nil {
				return
			}
			r = redis.NewClient(rOpts)
			return
		}).
		Check(func(ctx context.Context) error {
			return r.Ping(ctx).Err()
		}).
		Middleware(func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(ctx winter.Context) {
				ctx.Inject(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, opt.key, r)
				})
				h(ctx)
			}
		}).
		Shutdown(func(ctx context.Context) error {
			return r.Close()
		})

}
