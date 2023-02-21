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

// Get get previously injected [resty.Client]
func Get(ctx context.Context, opts ...Option) *resty.Client {
	o := createOptions(opts...)
	return ctx.Value(o.key).(*resty.Client)
}

// Install install component
func Install(a winter.App, opts ...Option) {
	o := createOptions(opts...)

	var rc *resty.Client

	a.Component("resty-" + string(o.key)).
		Startup(func(ctx context.Context) (err error) {
			// using transport with otelhttp
			hc := &http.Client{
				Transport: otelhttp.NewTransport(http.DefaultTransport),
			}
			if o.hcSetup != nil {
				hc = o.hcSetup(hc)
			}
			rc = resty.NewWithClient(hc)
			if o.rSetup != nil {
				rc = o.rSetup(rc)
			}
			return
		}).
		Middleware(func(h winter.HandlerFunc) winter.HandlerFunc {
			return func(c winter.Context) {
				c.Inject(func(ctx context.Context) context.Context {
					return context.WithValue(ctx, o.key, rc)
				})
				h(c)
			}
		})
}
