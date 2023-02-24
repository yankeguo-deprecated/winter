package wexttest

import (
	"context"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wext"
	"github.com/stretchr/testify/require"
	"io"
	"net/http/httptest"
	"testing"
)

type Options[OPTS any, INJ any] struct {
	AltKeys []string
	Options []func(opt *OPTS)

	BeforeInstaller func()
	AfterInstaller  func()

	BeforeInstall func()
	AfterInstall  func()

	BeforeStartup func()
	AfterStartup  func()

	BeforeCheck func()
	OnCheck     func(code int, body []byte)
	AfterCheck  func()

	BeforeShutdown func()
	AfterShutdown  func()

	BeforeHandle func()
	Handle       func(c winter.Context, ins wext.Instance[OPTS, INJ])
	OnHandle     func(code int, body []byte)
	AfterHandle  func()
}

// Run run tests with extension and hooks
func Run[OPTS any, INJ any](t *testing.T, ext wext.Extension[OPTS, INJ], opts Options[OPTS, INJ]) {
	a := winter.New()

	if opts.BeforeInstaller != nil {
		opts.BeforeInstaller()
	}

	insr := ext.Installer(opts.Options...)

	if opts.AfterInstaller != nil {
		opts.AfterInstaller()
	}

	if opts.BeforeInstall != nil {
		opts.BeforeInstall()
	}

	insr.Install(a, opts.AltKeys...)

	if opts.AfterInstall != nil {
		opts.AfterInstall()
	}

	if opts.Handle != nil {
		a.HandleFunc("/test", func(c winter.Context) {
			opts.Handle(c, ext.Instance(opts.AltKeys...))
		})
	}

	if opts.BeforeStartup != nil {
		opts.BeforeStartup()
	}

	ctx := context.Background()

	err := a.Startup(ctx)
	require.NoError(t, err)

	if opts.AfterStartup != nil {
		opts.AfterStartup()
	}

	defer func() {
		if opts.BeforeShutdown != nil {
			opts.BeforeShutdown()
		}

		a.Shutdown(ctx)

		if opts.AfterShutdown != nil {
			opts.AfterShutdown()
		}
	}()

	s := httptest.NewServer(a)
	defer s.Close()

	func() {
		if opts.BeforeCheck != nil {
			opts.BeforeCheck()
		}

		res, err := s.Client().Get(s.URL + winter.DefaultReadinessPath)
		require.NoError(t, err)
		defer res.Body.Close()

		buf, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		if opts.OnCheck != nil {
			opts.OnCheck(res.StatusCode, buf)
		}

		if opts.AfterCheck != nil {
			opts.AfterCheck()
		}
	}()

	func() {
		if opts.BeforeHandle != nil {
			opts.BeforeHandle()
		}

		res, err := s.Client().Get(s.URL + "/test")
		require.NoError(t, err)
		defer res.Body.Close()

		buf, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		if opts.OnHandle != nil {
			opts.OnHandle(res.StatusCode, buf)
		}

		if opts.AfterHandle != nil {
			opts.AfterHandle()
		}
	}()
}
