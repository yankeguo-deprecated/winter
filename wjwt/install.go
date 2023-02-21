package wjwt

import (
	"context"
	"encoding/json"
	"github.com/guoyk93/rg"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wjwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"time"
)

// Get get JWT Payload from Istio RequestAuthentication header
func Get(c winter.Context, opts ...Option) jwt.Token {
	o := c.Value(createOptions(opts...).key).(*options)
	t := jwt.New()
	var m map[string]any
	buf := []byte(c.Req().Header.Get(o.payloadHeader))
	rg.Must0(json.Unmarshal(buf, &m))
	for k, v := range m {
		t.Set(k, v)
	}
	return t
}

// Sign create a signed JWT
func Sign(ctx context.Context, fn func(b *jwt.Builder) *jwt.Builder, opts ...Option) string {
	o := ctx.Value(createOptions(opts...).key).(*options)
	k := wjwk.Get(ctx, wjwk.WithKey(string(o.jwkKey)))
	b := fn(jwt.NewBuilder().Issuer(o.issuer).IssuedAt(time.Now()))
	t := rg.Must(b.Build())
	signed := rg.Must(jwt.Sign(t, jwt.WithKey(k.Algorithm(), k)))
	return string(signed)
}

// Install install component
func Install(a winter.App, opts ...Option) {
	o := createOptions(opts...)

	a.Component("jwt-" + string(o.key)).
		Middleware(func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(c winter.Context) {
				c.Inject(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, o.key, o)
				})
				h(c)
			}
		})
}
