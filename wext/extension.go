package wext

import (
	"context"
	"github.com/guoyk93/winter"
	"strings"
)

type contextKey string

// Instance an instance of an [Extension]
type Instance[OPTS any, INJ any] interface {
	// Set inject a value into [context.Context]
	Set(ctx context.Context, inj INJ) context.Context

	// Get get the value from [context.Context]
	Get(ctx context.Context) INJ
}

type instance[OPTS any, INJ any] struct {
	key contextKey
}

func (e *instance[OPTS, INJ]) Set(ctx context.Context, inj INJ) context.Context {
	return context.WithValue(ctx, e.key, inj)
}

func (e *instance[OPTS, INJ]) Get(ctx context.Context) INJ {
	return ctx.Value(e.key).(INJ)
}

// Installer simple interface wrapping an install method
type Installer interface {
	Install(a winter.App, altKeys ...string)
}

type installerFunc func(a winter.App, altKeys ...string)

func (f installerFunc) Install(a winter.App, altKeys ...string) {
	f(a, altKeys...)
}

var (
	_ Installer = installerFunc(func(_ winter.App, _ ...string) {})
)

// Extension a abstraction of extension for wboot application
type Extension[OPTS any, INJ any] interface {
	// Instance create a instance of extension, one [winter.App] can have multiple instances
	Instance(altKeys ...string) Instance[OPTS, INJ]

	// Installer create [Installer]
	Installer(optFns ...func(opt *OPTS)) Installer

	// Default setup options by default
	Default(fn func(opt *OPTS)) Extension[OPTS, INJ]

	// Startup set startup function
	Startup(fn func(ctx context.Context, opt *OPTS) (inj INJ, err error)) Extension[OPTS, INJ]

	// Check set check func
	Check(fn func(ctx context.Context, inj INJ) (err error)) Extension[OPTS, INJ]

	// Middleware set middleware generation function
	Middleware(fn func(ins Instance[OPTS, INJ], opt *OPTS, inj *INJ) winter.MiddlewareFunc) Extension[OPTS, INJ]

	// Shutdown set the shutdown func
	Shutdown(fn func(ctx context.Context, inj INJ) (err error)) Extension[OPTS, INJ]
}

type extension[OPTS any, INJ any] struct {
	kind       string
	defaultFn  func(opt *OPTS)
	startupFn  func(ctx context.Context, opt *OPTS) (inj INJ, err error)
	checkFn    func(ctx context.Context, inj INJ) (err error)
	injectFn   func(ins Instance[OPTS, INJ], opt *OPTS, inj *INJ) winter.MiddlewareFunc
	shutdownFn func(ctx context.Context, inj INJ) (err error)
}

func (e *extension[OPTS, INJ]) createOptions(fns ...func(opt *OPTS)) *OPTS {
	var opt OPTS

	if e.defaultFn != nil {
		e.defaultFn(&opt)
	}

	for _, fn := range fns {
		fn(&opt)
	}
	return &opt
}

func (e *extension[OPTS, INJ]) createKey(altKeys ...string) string {
	key := e.kind + "."
	if len(altKeys) == 0 {
		key = key + "default"
	} else {
		key = key + strings.Join(altKeys, ".")
	}
	return key
}

func (e *extension[OPTS, INJ]) Instance(altKeys ...string) Instance[OPTS, INJ] {
	return e.instance(altKeys...)
}

func (e *extension[OPTS, INJ]) instance(altKeys ...string) *instance[OPTS, INJ] {
	key := e.createKey(altKeys...)
	return &instance[OPTS, INJ]{
		key: contextKey(key),
	}
}

func (e *extension[OPTS, INJ]) Installer(optFns ...func(opt *OPTS)) Installer {
	opt := e.createOptions(optFns...)

	return installerFunc(func(a winter.App, altKeys ...string) {
		ins := e.instance(altKeys...)

		var inj INJ

		c := a.Component(string(ins.key))

		if e.startupFn != nil {
			c = c.Startup(func(ctx context.Context) (err error) {
				inj, err = e.startupFn(ctx, opt)
				return
			})
		}
		if e.checkFn != nil {
			c = c.Check(func(ctx context.Context) error {
				return e.checkFn(ctx, inj)
			})
		}
		if e.shutdownFn != nil {
			c = c.Shutdown(func(ctx context.Context) (err error) {
				return e.shutdownFn(ctx, inj)
			})
		}

		if e.injectFn != nil {
			c = c.Middleware(e.injectFn(ins, opt, &inj))
		}

	})
}

func (e *extension[OPTS, INJ]) Default(fn func(opt *OPTS)) Extension[OPTS, INJ] {
	e.defaultFn = fn
	return e
}

func (e *extension[OPTS, INJ]) Startup(fn func(ctx context.Context, opt *OPTS) (inj INJ, err error)) Extension[OPTS, INJ] {
	e.startupFn = fn
	return e
}

func (e *extension[OPTS, INJ]) Middleware(fn func(ins Instance[OPTS, INJ], opt *OPTS, inj *INJ) winter.MiddlewareFunc) Extension[OPTS, INJ] {
	e.injectFn = fn
	return e
}

func (e *extension[OPTS, INJ]) Check(fn func(ctx context.Context, inj INJ) (err error)) Extension[OPTS, INJ] {
	e.checkFn = fn
	return e
}

func (e *extension[OPTS, INJ]) Shutdown(fn func(ctx context.Context, inj INJ) (err error)) Extension[OPTS, INJ] {
	e.shutdownFn = fn
	return e
}

// New create a new [Extension] with a kind and options generator
func New[OPTS any, INJ any](kind string) Extension[OPTS, INJ] {
	e := &extension[OPTS, INJ]{
		kind: kind,
	}
	return e.Middleware(func(ins Instance[OPTS, INJ], opt *OPTS, inj *INJ) winter.MiddlewareFunc {
		return func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(c winter.Context) {
				c.Inject(func(ctx context.Context) context.Context {
					return ins.Set(ctx, *inj)
				})
				h(c)
			}
		}
	})
}

// Simple create a new [Extension] with options injected
func Simple[OPTS any](kind string) Extension[OPTS, *OPTS] {
	return New[OPTS, *OPTS](kind).
		Startup(func(ctx context.Context, opt *OPTS) (inj *OPTS, err error) {
			inj = opt
			return
		})
}
