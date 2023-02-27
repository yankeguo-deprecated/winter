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

func Handler(altKeys ...string) winter.HandlerFunc {
	return func(c winter.Context) {
		ext.Instance(altKeys...).Get(c).Handle(c)
	}
}

func Installer(opts ...Option) wext.Installer {
	return ext.Installer(opts...)
}
