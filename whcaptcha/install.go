package whcaptcha

import (
	"errors"
	"github.com/guoyk93/rg"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/pkg/clientip"
	"github.com/guoyk93/winter/wext"
	"github.com/guoyk93/winter/wresty"
)

var (
	ErrInvalidToken = errors.New("hcaptcha: invalid token")
)

// Validate validate hcaptcha
func Validate(c winter.Context, token string, altKeys ...string) {
	o := Ext.Instance(altKeys...).Get(c)

	var ret struct {
		Success bool `json:"success"`
	}
	res := rg.Must(wresty.R(c, o.restyKeys...).SetFormData(map[string]string{
		"sitekey":  o.siteKey,
		"secret":   o.secret,
		"response": token,
		"remoteip": clientip.Get(c),
	}).SetResult(&ret).Post("https://hcaptcha.com/siteverify"))

	if res.IsError() {
		winter.HaltString(res.String())
	}
	if !ret.Success {
		winter.Halt(ErrInvalidToken)
	}
}

// Installer install component
func Installer(opts ...Option) wext.Installer {
	o := Ext.Options(opts...)

	return wext.WrapInstaller(func(a winter.App, altKeys ...string) {
		ins := Ext.Instance(altKeys...)

		a.Component(ins.Key()).Middleware(ins.Middleware(&o))
	})
}
