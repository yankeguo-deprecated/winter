package wresty

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/guoyk93/winter"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

// Get get [resty.Client]
func Get(c winter.Context, opts ...Option) *resty.Client {
	opt := createOptions(opts...)
	return c.Value(opt.key).(*resty.Client)
}

// R create a resty.Request from context
func R(c winter.Context, opts ...Option) *resty.Request {
	return Get(c, opts...).R().SetContext(c)
}

// Install install component
func Install(a winter.App, opts ...Option) {
	opt := createOptions(opts...)

	var c *resty.Client

	a.Component("resty").
		Startup(func(ctx context.Context) (err error) {
			// using transport with otelhttp
			hc := &http.Client{
				Transport: otelhttp.NewTransport(http.DefaultTransport),
			}
			if opt.hcSetup != nil {
				hc = opt.hcSetup(hc)
			}
			c = resty.NewWithClient(hc)
			if opt.rSetup != nil {
				c = opt.rSetup(c)
			}
			return
		}).
		Middleware(func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(c winter.Context) {
				c.Inject(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, opt.key, c)
				})
				h(c)
			}
		})
}
