package wasynq

import "github.com/hibiken/asynq"

type options struct {
	redisURL string
	redisOpt asynq.RedisConnOpt
}

type Option = func(opts *options)

func WithRedisURL(s string) Option {
	return func(opts *options) {
		opts.redisURL = s
	}
}

func WithRedisOpt(o asynq.RedisConnOpt) Option {
	return func(opts *options) {
		opts.redisOpt = o
	}
}
