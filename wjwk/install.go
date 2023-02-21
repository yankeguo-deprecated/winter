package wjwk

import (
	"context"
	"errors"
	"github.com/guoyk93/winter"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

// Get get previously injected [jwk.Key]
func Get(ctx context.Context, opts ...Option) jwk.Key {
	o := createOptions(opts...)
	return ctx.Value(o.key).(jwk.Key)
}

// Install install component
func Install(a winter.App, opts ...Option) {
	o := createOptions(opts...)

	var k jwk.Key

	a.Component("jwk-" + string(o.key)).
		Startup(func(ctx context.Context) (err error) {
			if len(o.raw) != 0 {
				k, err = jwk.ParseKey(o.raw)
			} else {
				err = errors.New("failed loading jwk")
			}
			return
		}).
		Middleware(func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(c winter.Context) {
				c.Inject(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, o.key, k)
				})
				h(c)
			}
		})
}
