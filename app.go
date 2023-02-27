package winter

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"net/http"
	"net/http/pprof"
	"strings"
	"sync/atomic"
)

// HandlerFunc handler func with [Context] as argument
type HandlerFunc func(c Context)

// App the main interface of [summer]
type App interface {
	// Handler inherit [http.Handler]
	http.Handler

	// Registry inherit [Registry]
	Registry

	// HandleFunc register an action function with given path pattern
	//
	// This function is similar with [http.ServeMux.HandleFunc]
	HandleFunc(pattern string, fn HandlerFunc)
}

type app struct {
	Registry

	opts options

	hMain *http.ServeMux

	hProm http.Handler
	hProf http.Handler

	cc chan struct{}

	failed int64
}

func (a *app) HandleFunc(pattern string, fn HandlerFunc) {
	a.hMain.Handle(
		pattern,
		otelhttp.NewHandler(
			otelhttp.WithRouteTag(
				pattern,
				http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					c := newContext(rw, req)
					c.responseLogging = a.opts.responseLogging
					func() {
						defer c.Perform()
						a.Wrap(fn)(c)
					}()
				}),
			),
			pattern,
		),
	)
}

func (a *app) serveReadiness(rw http.ResponseWriter, req *http.Request) {
	c := newContext(rw, req)
	defer c.Perform()

	a.Wrap(func(c Context) {
		cr := &checkResult{}
		a.Check(c, cr.Collect)
		s, failed := cr.Result()

		status := http.StatusOK
		if failed {
			atomic.AddInt64(&a.failed, 1)
			status = http.StatusInternalServerError
		} else {
			atomic.StoreInt64(&a.failed, 0)
		}

		c.Code(status)
		c.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Text(s)
	})(c)
}

func (a *app) serveLiveness(rw http.ResponseWriter, req *http.Request) {
	c := newContext(rw, req)
	defer c.Perform()

	c.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if a.opts.readinessCascade > 0 &&
		atomic.LoadInt64(&a.failed) > a.opts.readinessCascade {
		c.Code(http.StatusInternalServerError)
		c.Text("CASCADED")
	} else {
		c.Code(http.StatusOK)
		c.Text("OK")
	}
}

func (a *app) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// alive, ready, metrics
	if req.URL.Path == a.opts.readinessPath {
		// support readinessPath == livenessPath
		a.serveReadiness(rw, req)
		return
	} else if req.URL.Path == a.opts.livenessPath {
		a.serveLiveness(rw, req)
		return
	} else if req.URL.Path == a.opts.metricsPath {
		a.hProm.ServeHTTP(rw, req)
		return
	}

	// pprof
	if strings.HasPrefix(req.URL.Path, "/debug/") {
		a.hProf.ServeHTTP(rw, req)
		return
	}

	// concurrency
	if a.cc != nil {
		<-a.cc
		defer func() {
			a.cc <- struct{}{}
		}()
	}

	a.hMain.ServeHTTP(rw, req)
}

// New create an [App] with [Option]
func New(opts ...Option) App {
	a := &app{
		opts: options{
			concurrency:      128,
			readinessCascade: 5,
			readinessPath:    DefaultReadinessPath,
			livenessPath:     DefaultLivenessPath,
			metricsPath:      DefaultMetricsPath,
		},
	}
	for _, opt := range opts {
		opt(&a.opts)
	}

	a.Registry = NewRegistry()

	{
		a.hMain = &http.ServeMux{}
		a.hProm = promhttp.Handler()
		m := &http.ServeMux{}
		m.HandleFunc("/debug/pprof/", pprof.Index)
		m.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		m.HandleFunc("/debug/pprof/profile", pprof.Profile)
		m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		m.HandleFunc("/debug/pprof/trace", pprof.Trace)
		a.hProf = m
	}

	// create concurrency controller
	if a.opts.concurrency > 0 {
		a.cc = make(chan struct{}, a.opts.concurrency)
		for i := 0; i < a.opts.concurrency; i++ {
			a.cc <- struct{}{}
		}
	}
	return a
}
