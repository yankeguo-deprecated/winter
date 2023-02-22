package wext

import (
	"strings"
)

type Extension[OPTS any, INJ any] interface {
	Kind() string

	Instance(altKeys ...string) Instance[OPTS, INJ]

	Options(opts ...func(opt *OPTS)) *OPTS
}

type extension[OPTS any, INJ any] struct {
	kind       string
	createOpts func() *OPTS
}

func (e *extension[OPTS, INJ]) Kind() string {
	return e.kind
}

func (e *extension[OPTS, INJ]) Instance(altKeys ...string) Instance[OPTS, INJ] {
	key := e.kind + "."
	if len(altKeys) == 0 {
		key = key + "default"
	} else {
		key = key + strings.Join(altKeys, ".")
	}
	return &instance[OPTS, INJ]{
		key: key,
	}
}

func (e *extension[OPTS, INJ]) Options(opts ...func(opt *OPTS)) *OPTS {
	o := e.createOpts()
	for _, opt := range opts {
		opt(o)
	}
	return o
}

func New[OPTS any, INJ any](kind string, createOpts func() *OPTS) Extension[OPTS, INJ] {
	return &extension[OPTS, INJ]{
		kind:       kind,
		createOpts: createOpts,
	}
}
