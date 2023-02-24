package wjwk

import (
	"context"
	"errors"
	"github.com/guoyk93/winter/wext"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

var (
	ext = wext.New[options, jwk.Key]("jwk").
		Startup(func(ctx context.Context, opt *options) (inj jwk.Key, err error) {
			if len(opt.raw) != 0 {
				inj, err = jwk.ParseKey(opt.raw)
			} else {
				err = errors.New("wjwk: missing key source")
			}
			return
		})
)

// Get get previously injected [jwk.Key]
func Get(ctx context.Context, altKeys ...string) jwk.Key {
	return ext.Instance(altKeys...).Get(ctx)
}

// Installer install component
func Installer(opts ...Option) wext.Installer {
	return ext.Installer(opts...)
}
