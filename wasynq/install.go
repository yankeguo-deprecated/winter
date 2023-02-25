package wasynq

import (
	"context"
	"errors"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wext"
	"github.com/hibiken/asynq"
)

var (
	ext = wext.New[options, *asynq.Client]("asynq").
		Startup(func(ctx context.Context, opt *options) (inj *asynq.Client, err error) {
			if opt.redisOpt != nil {
				inj = asynq.NewClient(*opt.redisOpt)
				return
			} else if opt.redisURL != "" {
				var ao asynq.RedisConnOpt
				if ao, err = asynq.ParseRedisURI(opt.redisURL); err != nil {
					return
				}
				inj = asynq.NewClient(ao)
				return
			} else {
				err = errors.New("missing asynq options")
				return
			}
		}).
		Shutdown(func(ctx context.Context, inj *asynq.Client) error {
			return inj.Close()
		})
)

type Interface interface {
	Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

type clientWithContext struct {
	c   *asynq.Client
	ctx context.Context
}

func (c *clientWithContext) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return c.c.EnqueueContext(c.ctx, task, opts...)
}

func Client(ctx context.Context, altKeys ...string) Interface {
	c := ext.Instance(altKeys...).Get(ctx)
	return &clientWithContext{
		c:   c,
		ctx: ctx,
	}
}

func Installer(a winter.App, opts ...Option) wext.Installer {
	return ext.Installer(opts...)
}
