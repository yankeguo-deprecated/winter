package wext

import (
	"context"
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
