package wwebauthn

import (
	"context"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wext"
)

// Get get previously injected [webauthn.WebAuthn]
func Get(ctx context.Context, altKeys ...string) *webauthn.WebAuthn {
	return Ext.Instance(altKeys...).Get(ctx)
}

// Installer install component
func Installer(opts ...Option) wext.Installer {
	o := Ext.Options(opts...)

	return wext.WrapInstaller(func(a winter.App, altKeys ...string) {
		ins := Ext.Instance(altKeys...)

		var w *webauthn.WebAuthn

		a.Component(ins.Key()).
			Startup(func(ctx context.Context) (err error) {
				w, err = webauthn.New(o.cfg)
				return
			}).
			Middleware(ins.Middleware(&w))
	})
}
