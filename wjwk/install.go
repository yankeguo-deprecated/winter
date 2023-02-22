package wjwk

import (
	"context"
	"errors"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wext"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

// Get get previously injected [jwk.Key]
func Get(ctx context.Context, altKeys ...string) jwk.Key {
	return Ext.Instance(altKeys...).Get(ctx)
}

// Installer install component
func Installer(opts ...Option) wext.Installer {
	o := Ext.Options(opts...)

	return wext.WrapInstaller(func(a winter.App, altKeys ...string) {
		ins := Ext.Instance(altKeys...)

		var k jwk.Key

		a.Component(ins.Key()).
			Startup(func(ctx context.Context) (err error) {
				if len(o.raw) != 0 {
					k, err = jwk.ParseKey(o.raw)
				} else {
					err = errors.New("wjwk: missing key source")
				}
				return
			}).
			Middleware(ins.Middleware(&k))
	})

}
