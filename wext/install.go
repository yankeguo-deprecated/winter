package wext

import "github.com/guoyk93/winter"

// Install install a [Installer]
func Install(a winter.App, installer Installer, altKeys ...string) {
	installer.Install(a, altKeys...)
}
