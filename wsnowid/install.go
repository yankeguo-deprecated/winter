package wsnowid

import (
	"context"
	"github.com/guoyk93/rg"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wresty"
	"strconv"
)

// Next return a new id
func Next(ctx context.Context, opts ...Option) string {
	return NextN(ctx, 1, opts...)[0]
}

// NextN return n ids generated from snowid service
func NextN(ctx context.Context, size int, opts ...Option) []string {
	if size < 1 {
		winter.HaltString("wsnowid: invalid argument: size", winter.HaltWithBadRequest())
	}

	o := ctx.Value(createOptions(opts...).key).(*options)

	var ret []string
	res := rg.Must(
		wresty.R(ctx, wresty.WithKey(string(o.restyKey))).
			SetQueryParam("size", strconv.Itoa(size)).
			SetResult(&ret).
			Get(o.url),
	)
	if res.IsError() {
		winter.HaltString(res.String())
	}
	if len(ret) != size {
		winter.HaltString("wsnowid: invalid returns")
	}
	return ret
}

// Install install component
func Install(a winter.App, opts ...Option) {
	o := createOptions(opts...)

	a.Component("snowid-" + string(o.key)).
		Check(func(ctx context.Context) (err error) {
			defer rg.Guard(&err)
			_ = Next(ctx, opts...)
			return
		}).
		Middleware(func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(c winter.Context) {
				c.Inject(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, o.key, o)
				})
				h(c)
			}
		})
}
