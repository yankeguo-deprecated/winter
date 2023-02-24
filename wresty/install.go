package wresty

import (
	"context"
	"github.com/go-resty/resty/v2"
	"github.com/guoyk93/winter/wext"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
)

var (
	ext = wext.New[options, *resty.Client]("resty").
		Startup(func(ctx context.Context, opt *options) (inj *resty.Client, err error) {
			// create http.Client with otel
			hc := &http.Client{
				Transport: otelhttp.NewTransport(http.DefaultTransport),
			}
			if opt.hcSetup != nil {
				hc = opt.hcSetup(hc)
			}
			// create resty.Client
			inj = resty.NewWithClient(hc)
			if opt.rSetup != nil {
				inj = opt.rSetup(inj)
			}
			return
		})
)

// R create a resty.Request from context
func R(ctx context.Context, altKeys ...string) *resty.Request {
	return Client(ctx, altKeys...).R().SetContext(ctx)
}

// Client get previously injected [resty.Client]
func Client(ctx context.Context, altKeys ...string) *resty.Client {
	return ext.Instance(altKeys...).Get(ctx)
}

// Installer create [wext.Installer]
func Installer(opts ...Option) wext.Installer {
	return ext.Installer(opts...)
}
