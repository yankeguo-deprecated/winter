package wwebauthn

import (
	"context"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/guoyk93/winter"
)

func Get(ctx context.Context, opts ...Option) *webauthn.WebAuthn {
	opt := createOptions(opts...)
	return ctx.Value(opt.key).(*webauthn.WebAuthn)
}

func Install(a winter.App, opts ...Option) {
	opt := createOptions(opts...)

	var w *webauthn.WebAuthn

	a.Component("webauthn").
		Startup(func(ctx context.Context) (err error) {
			w, err = webauthn.New(opt.cfg)
			return
		}).
		Middleware(func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(c winter.Context) {
				c.Inject(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, opt.key, w)
				})
				h(c)
			}
		})
}
