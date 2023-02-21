package wwebauthn

import (
	"context"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/guoyk93/winter"
)

// Get get previously injected [webauthn.WebAuthn]
func Get(ctx context.Context, opts ...Option) *webauthn.WebAuthn {
	o := createOptions(opts...)
	return ctx.Value(o.key).(*webauthn.WebAuthn)
}

// Install install component
func Install(a winter.App, opts ...Option) {
	o := createOptions(opts...)

	var w *webauthn.WebAuthn

	a.Component("webauthn-" + string(o.key)).
		Startup(func(ctx context.Context) (err error) {
			w, err = webauthn.New(o.cfg)
			return
		}).
		Middleware(func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(c winter.Context) {
				c.Inject(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, o.key, w)
				})
				h(c)
			}
		})
}
