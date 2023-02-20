package wresty

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/guoyk93/winter"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

// R create a resty.Request from context
func R(ctx context.Context, opts ...Option) *resty.Request {
	return Get(ctx, opts...).R().SetContext(ctx)
}

// Get get [resty.Client]
func Get(ctx context.Context, opts ...Option) *resty.Client {
	opt := createOptions(opts...)
	return ctx.Value(opt.key).(*resty.Client)
}

// Install install component
func Install(a winter.App, opts ...Option) {
	opt := createOptions(opts...)

	var rc *resty.Client

	a.Component("resty-" + string(opt.key)).
		Startup(func(ctx context.Context) (err error) {
			// using transport with otelhttp
			hc := &http.Client{
				Transport: otelhttp.NewTransport(http.DefaultTransport),
			}
			if opt.hcSetup != nil {
				hc = opt.hcSetup(hc)
			}
			rc = resty.NewWithClient(hc)
			if opt.rSetup != nil {
				rc = opt.rSetup(rc)
			}
			return
		}).
		Middleware(func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(c winter.Context) {
				c.Inject(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, opt.key, rc)
				})
				h(c)
			}
		})
}
