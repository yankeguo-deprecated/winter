package wboot

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/guoyk93/rg"
	"github.com/guoyk93/winter"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func setupOTEL() (err error) {
	defer rg.Guard(&err)

	// clear error handler
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {}))
	// using zipkin
	otel.SetTracerProvider(trace.NewTracerProvider(
		trace.WithBatcher(rg.Must(zipkin.New("", zipkin.WithLogr(logr.Discard())))),
	))
	// using b3
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
			b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader|b3.B3SingleHeader)),
		),
	)

	// re-initialize http client
	otelhttp.DefaultClient = &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	http.DefaultClient = otelhttp.DefaultClient

	return
}

// Main best practice of running a [winter.App]
func Main(fn func() (a winter.App, err error)) {
	var err error
	defer func() {
		if err == nil {
			return
		}
		log.Println("exited with error:", err.Error())
		os.Exit(1)
	}()
	defer rg.Guard(&err)

	log.SetOutput(os.Stdout)

	rg.Must0(setupOTEL())

	ctx := context.Background()

	a := rg.Must(fn())
	rg.Must0(a.Startup(ctx))
	defer a.Shutdown(ctx)

	s := &http.Server{
		Addr:    EnvStr("BIND") + ":" + EnvStrOr("PORT", "8080"),
		Handler: a,
	}

	chErr := make(chan error, 1)
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		chErr <- s.ListenAndServe()
	}()

	select {
	case err = <-chErr:
		return
	case sig := <-chSig:
		log.Println("signal caught:", sig.String())
		time.Sleep(time.Second * 3)
	}

	err = s.Shutdown(ctx)
}
