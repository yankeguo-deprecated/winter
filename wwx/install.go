package wwx

import (
	"context"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wext"
)

var (
	ext = wext.New[options, *app]("wwx").
		Startup(func(ctx context.Context, opt *options) (inj *app, err error) {
			inj = &app{
				opt: opt,
			}
			return
		})
)

func HandleCallback(c winter.Context, altKeys ...string) {
	ext.Instance(altKeys...).Get(c).HandleCallback(c)
}

func Installer(a winter.App, opts ...Option) wext.Installer {
	return ext.Installer(opts...)
}
