package wext

import (
	"context"
	"github.com/guoyk93/winter"
)

type ContextKey string

type Instance[OPTS any, INJ any] interface {
	Key() string

	ContextKey() ContextKey

	Set(ctx context.Context, inj INJ) context.Context

	Get(ctx context.Context) INJ

	Middleware(inj *INJ) func(h winter.HandlerFunc) winter.HandlerFunc
}

type instance[OPTS any, INJ any] struct {
	key string
}

func (e *instance[OPTS, INJ]) Key() string {
	return e.key
}

func (e *instance[OPTS, INJ]) ContextKey() ContextKey {
	return ContextKey(e.key)
}

func (e *instance[OPTS, INJ]) Set(ctx context.Context, inj INJ) context.Context {
	return context.WithValue(ctx, e.ContextKey(), inj)
}

func (e *instance[OPTS, INJ]) Get(ctx context.Context) INJ {
	return ctx.Value(e.ContextKey()).(INJ)
}

func (e *instance[OPTS, INJ]) Middleware(inj *INJ) func(h winter.HandlerFunc) winter.HandlerFunc {
	return func(h winter.HandlerFunc) winter.HandlerFunc {
		return func(c winter.Context) {
			c.Inject(func(ctx context.Context) context.Context {
				return e.Set(ctx, *inj)
			})
			h(c)
		}
	}
}
