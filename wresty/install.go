package wresty

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/guoyk93/winter"
	"github.com/guoyk93/winter/wext"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

// R create a resty.Request from context
func R(ctx context.Context, altKeys ...string) *resty.Request {
	return Client(ctx, altKeys...).R().SetContext(ctx)
}

// Client get previously injected [resty.Client]
func Client(ctx context.Context, altKeys ...string) *resty.Client {
	return Ext.Instance(altKeys...).Get(ctx)
}

// Installer create [wext.Installer]
func Installer(opts ...Option) wext.Installer {
	o := Ext.Options(opts...)

	return wext.WrapInstaller(func(a winter.App, altKeys ...string) {
		ins := Ext.Instance(altKeys...)

		var rc *resty.Client

		a.Component(ins.Key()).
			Startup(func(ctx context.Context) (err error) {
				// create http.Client with otel
				hc := &http.Client{
					Transport: otelhttp.NewTransport(http.DefaultTransport),
				}
				if o.hcSetup != nil {
					hc = o.hcSetup(hc)
				}
				// create resty.Client
				rc = resty.NewWithClient(hc)
				if o.rSetup != nil {
					rc = o.rSetup(rc)
				}
				return
			}).
			Middleware(ins.Middleware(rc))
	})
}
