package wext

import "github.com/guoyk93/winter"

type Installer interface {
	Install(a winter.App, altKeys ...string)
}

type installer struct {
	fn func(a winter.App, altKeys ...string)
}

func (e installer) Install(a winter.App, altKeys ...string) {
	e.fn(a, altKeys...)
}

// WrapInstaller create an [Installer] from function
func WrapInstaller(fn func(a winter.App, altKeys ...string)) Installer {
	return installer{fn: fn}
}
