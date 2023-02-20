package wsnowid

import (
	"context"
	"errors"
	"github.com/guoyk93/rg"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wresty"
)

// Next return a new id
func Next(ctx context.Context, opts ...Option) string {
	return NextN(ctx, 1, opts...)[0]
}

// NextN return n ids generated from snowid service
func NextN(ctx context.Context, size int, opts ...Option) []string {
	if size < 1 {
		winter.Halt(errors.New("wsnowid: invalid argument: size"))
	}
	opt := createOptions(opts...)
	opt = ctx.Value(opt.key).(*options)
	var ret []string
	res := rg.Must(wresty.R(ctx, wresty.WithKey(string(opt.rKey))).SetResult(&ret).Get(opt.url))
	if res.IsError() {
		winter.Halt(errors.New(res.String()))
	}
	if len(ret) != size {
		winter.Halt(errors.New("wsnowid: invalid returns"))
	}
	return ret
}

// Install install component
func Install(a winter.App, opts ...Option) {
	opt := createOptions(opts...)

	a.Component("snowid-" + string(opt.key)).
		Check(func(ctx context.Context) (err error) {
			defer rg.Guard(&err)
			_ = Next(ctx, opts...)
			return
		}).
		Middleware(func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(c winter.Context) {
				c.Inject(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, opt.key, opt)
				})
				h(c)
			}
		})
}
