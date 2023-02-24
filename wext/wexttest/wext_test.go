package wexttest

import (
	"context"
	"errors"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wext"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestExtension(t *testing.T) {
	type options struct {
		forDefault  bool
		forOption   bool
		forStartup  bool
		forCheck    bool
		forShutdown bool
	}

	var handled bool

	type injected struct {
		opts *options
		v    int
	}

	var (
		_opt *options
		_inj *injected
	)

	ext := wext.New[options, *injected]("test-ext").
		Default(func(opt *options) {
			_opt = opt
			opt.forDefault = true
		}).
		Startup(func(ctx context.Context, opt *options) (*injected, error) {
			require.True(t, opt.forOption)
			opt.forStartup = true
			_inj = &injected{
				opts: opt,
				v:    11,
			}
			return _inj, nil
		}).
		Check(func(ctx context.Context, inj *injected) (err error) {
			inj.opts.forCheck = true
			return errors.New("CHECK")
		}).
		Shutdown(func(ctx context.Context, inj *injected) (err error) {
			inj.opts.forShutdown = true
			return errors.New("SHUTDOWN")
		})

	var invoked int

	Run(t, ext, Options[options, *injected]{
		AltKeys: []string{"test-altkey"},
		Options: []func(opt *options){
			func(opt *options) {
				require.True(t, opt.forDefault)
				opt.forOption = true

				invoked++
			},
			func(opt *options) {
				invoked++
			},
		},
		BeforeInstaller: func() {
			invoked++
			require.Nil(t, _opt)
		},
		AfterInstaller: func() {
			invoked++
			require.NotNil(t, _opt)
		},
		BeforeInstall: func() {
			invoked++
			require.True(t, _opt.forDefault)
		},
		AfterInstall: func() {
			invoked++
			require.Nil(t, _inj)
		},
		BeforeStartup: func() {
			invoked++
			require.Nil(t, _inj)
			require.False(t, _opt.forStartup)
		},
		AfterStartup: func() {
			invoked++
			require.True(t, _inj.opts.forStartup)
		},
		BeforeCheck: func() {
			invoked++
			require.False(t, _inj.opts.forCheck)
		},
		OnCheck: func(code int, body []byte) {
			invoked++
			require.Equal(t, http.StatusInternalServerError, code)
			require.Equal(t, "test-ext.test-altkey: CHECK", string(body))
		},
		AfterCheck: func() {
			invoked++
			require.True(t, _inj.opts.forCheck)
		},
		BeforeShutdown: func() {
			invoked++
			require.False(t, _inj.opts.forShutdown)
		},
		AfterShutdown: func() {
			invoked++
			require.True(t, _inj.opts.forShutdown)
		},
		BeforeHandle: func() {
			invoked++
			require.False(t, handled)
		},
		Handle: func(c winter.Context, ins wext.Instance[options, *injected]) {
			invoked++
			inj := ins.Get(c)
			require.Equal(t, 11, inj.v)
			c.Code(401)
			c.Text("BBB")
		},
		OnHandle: func(code int, body []byte) {
			invoked++
			handled = true
			require.Equal(t, 401, code)
			require.Equal(t, "BBB", string(body))
		},
		AfterHandle: func() {
			invoked++
			require.True(t, handled)
		},
	})

	require.Equal(t, invoked, 17)
}
