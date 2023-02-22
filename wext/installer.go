package wext

type Installer interface {
	Install(altKeys ...string)
}

type installer struct {
	fn func(altKeys ...string)
}

func (e installer) Install(altKeys ...string) {
	e.fn(altKeys...)
}

// WrapInstaller create an [Installer] from function
func WrapInstaller(fn func(altKeys ...string)) Installer {
	return installer{fn: fn}
}
