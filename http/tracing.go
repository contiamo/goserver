package server

import (
	"net/http"

	"github.com/contiamo/goserver"
	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
)

func WithTracing(server, app string) Option {
	return &tracingOption{server, app}
}

type tracingOption struct{ server, app string }

func (opt *tracingOption) WrapHandler(handler http.Handler) (http.Handler, error) {
	if err := goserver.InitTracer(opt.server, opt.app); err != nil {
		return nil, err
	}
	mw := nethttp.Middleware(
		opentracing.GlobalTracer(),
		handler,
		nethttp.OperationNameFunc(func(r *http.Request) string {
			return "HTTP " + r.Method + " " + r.URL.Path
		}),
		nethttp.MWSpanObserver(func(sp opentracing.Span, r *http.Request) {
			sp.SetTag("http.uri", r.URL.EscapedPath())
		}),
	)
	return mw, nil
}
