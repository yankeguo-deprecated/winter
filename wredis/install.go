package wredis

import (
	"context"
	"github.com/guoyk93/winter/wext"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

var (
	ext = wext.New[options, *redis.Client]("redis").
		Startup(func(ctx context.Context, opt *options) (inj *redis.Client, err error) {
			// create options
			rOpts := &redis.Options{}
			if opt.opts != nil {
				rOpts = opt.opts
			} else if opt.url != "" {
				if rOpts, err = redis.ParseURL(opt.url); err != nil {
					return
				}
			}
			// create client
			inj = redis.NewClient(rOpts)
			// instrument
			if err = redisotel.InstrumentTracing(inj, opt.tracingOpts...); err != nil {
				return
			}
			return
		}).
		Check(func(ctx context.Context, inj *redis.Client) error {
			return inj.Ping(ctx).Err()
		}).
		Shutdown(func(ctx context.Context, inj *redis.Client) error {
			return inj.Close()
		})
)

// Client get previously injected [redis.Client]
func Client(ctx context.Context, altKeys ...string) *redis.Client {
	return ext.Instance(altKeys...).Get(ctx)
}

// Installer install component
func Installer(opts ...Option) wext.Installer {
	return ext.Installer(opts...)
}
