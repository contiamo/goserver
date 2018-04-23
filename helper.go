package goserver

import (
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// InitTracer asserts that the global tracer is initialized
func InitTracer(server, app string) error {
	global := opentracing.GlobalTracer()
	if _, ok := global.(opentracing.NoopTracer); ok {
		cfg := config.Configuration{
			Sampler: &config.SamplerConfig{
				Type:  "const",
				Param: 1,
			},
			Reporter: &config.ReporterConfig{
				LogSpans:            false,
				BufferFlushInterval: 1 * time.Second,
				LocalAgentHostPort:  server,
			},
		}

		tracer, _, err := cfg.New(app, config.Logger(jaeger.StdLogger))
		if err != nil {
			return err
		}
		opentracing.SetGlobalTracer(tracer)
	}
	return nil
}

// ListenAndServeMetricsAndHealth starts up an HTTP server serving /metrics and /health
func ListenAndServeMetricsAndHealth(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"msg":"ok"}`))
	})
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return srv.ListenAndServe()
}
