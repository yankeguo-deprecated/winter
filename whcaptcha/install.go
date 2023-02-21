package whcaptcha

import (
	"context"
	"errors"
	"github.com/guoyk93/rg"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wclientip"
	"github.com/guoyk93/winter/wresty"
)

var (
	ErrInvalidToken = errors.New("hcaptcha: invalid token")
)

// Validate validate hcaptcha
func Validate(c winter.Context, token string, opts ...Option) {
	o := c.Value(createOptions(opts...).key).(*options)
	var ret struct {
		Success bool `json:"success"`
	}
	res := rg.Must(wresty.R(c, wresty.WithKey(string(o.restyKey))).SetFormData(map[string]string{
		"sitekey":  o.siteKey,
		"secret":   o.secret,
		"response": token,
		"remoteip": wclientip.Get(c),
	}).SetResult(&ret).Post("https://hcaptcha.com/siteverify"))

	if res.IsError() {
		winter.HaltString(res.String())
	}
	if !ret.Success {
		winter.Halt(ErrInvalidToken)
	}
}

// Install install component
func Install(a winter.App, opts ...Option) {
	opt := createOptions(opts...)

	a.Component("hcaptcha-" + string(opt.key)).
		Middleware(func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(c winter.Context) {
				c.Inject(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, opt.key, opt)
				})
				h(c)
			}
		})
}
