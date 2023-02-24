package wwebauthn

import (
	"context"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/guoyk93/winter/wext"
)

var (
	ext = wext.New[options, *webauthn.WebAuthn]("webauthn").
		Startup(func(ctx context.Context, opt *options) (*webauthn.WebAuthn, error) {
			return webauthn.New(opt.cfg)
		})
)

// Get get previously injected [webauthn.WebAuthn]
func Get(ctx context.Context, altKeys ...string) *webauthn.WebAuthn {
	return ext.Instance(altKeys...).Get(ctx)
}

// Installer install component
func Installer(opts ...Option) wext.Installer {
	return ext.Installer(opts...)
}
