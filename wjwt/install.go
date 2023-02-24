package wjwt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/guoyk93/rg"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wext"
	"github.com/guoyk93/winter/wjwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"strings"
	"time"
)

var (
	ext = wext.Simple[options]("jwt").
		Default(func(opt *options) {
			opt.payloadHeader = "X-JWT-Payload"
		})
)

// Get get JWT Payload from Istio RequestAuthentication header
func Get(c winter.Context, altKeys ...string) jwt.Token {
	o := ext.Instance(altKeys...).Get(c)

	var pl string

	if o.debugPayload {
		splits := strings.Split(
			strings.TrimSpace(
				strings.TrimPrefix(
					strings.TrimSpace(c.Req().Header.Get("Authorization")),
					"Bearer",
				),
			),
			".",
		)

		if len(splits) == 3 {
			pl = splits[1]
		}
	} else {
		pl = c.Req().Header.Get(o.payloadHeader)
	}

	buf := rg.Must(base64.RawURLEncoding.DecodeString(pl))

	var m map[string]any
	rg.Must0(json.Unmarshal(buf, &m))

	t := jwt.New()
	for k, v := range m {
		t.Set(k, v)
	}
	return t
}

// Sign create a signed JWT
func Sign(ctx context.Context, fn func(b *jwt.Builder) *jwt.Builder, altKeys ...string) string {
	o := ext.Instance(altKeys...).Get(ctx)
	k := wjwk.Get(ctx, o.jwkKeys...)
	b := fn(jwt.NewBuilder().Issuer(o.issuer).IssuedAt(time.Now()))
	t := rg.Must(b.Build())
	signed := rg.Must(jwt.Sign(t, jwt.WithKey(k.Algorithm(), k)))
	return string(signed)
}

// Installer install component
func Installer(opts ...Option) wext.Installer {
	return ext.Installer(opts...)
}
