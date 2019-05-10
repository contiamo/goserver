package goserver

import (
	"context"
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// InitTracer asserts that the global tracer is initialized.
//
// This will read the configuration from the "JAEGER_*"" environment variables.
// Overriding the empty values with the supplied server and app value.  If a
// sampler type is not configured via the environment variables, then InitTracer
// will be configured with the constant sampler.
func InitTracer(server, app string) error {
	global := opentracing.GlobalTracer()
	if _, ok := global.(opentracing.NoopTracer); ok {
		cfg, err := config.FromEnv()
		if err != nil {
			return err
		}
		if cfg.ServiceName == "" {
			cfg.ServiceName = app
		}

		if cfg.Sampler.Type == "" {
			cfg.Sampler.Type = "const"
			cfg.Sampler.Param = 1
		}

		if cfg.Reporter.BufferFlushInterval == 0 {
			cfg.Reporter.BufferFlushInterval = 1 * time.Second
		}
		if cfg.Reporter.LocalAgentHostPort == "" {
			cfg.Reporter.LocalAgentHostPort = server
		}

		_, err = cfg.InitGlobalTracer(app, config.Logger(jaeger.StdLogger))
		if err != nil {
			return err
		}
	}
	return nil
}

// ListenAndServeMetricsAndHealth starts up an HTTP server serving /metrics and /health
//
// Deprecated: use ListenAndServeMonitoring instead
func ListenAndServeMetricsAndHealth(addr string, healthHandler http.Handler) error {
	return monitoringServer(addr, healthHandler).ListenAndServe()
}

// ListenAndServeMonitoring starts up an HTTP server serving /metrics and /health.
//
// When the context is cancelled, the server will be gracefully shutdown.
func ListenAndServeMonitoring(ctx context.Context, addr string, healthHandler http.Handler) error {
	srv := monitoringServer(addr, healthHandler)

	go func() {
		<-ctx.Done()
		shutdownContext, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		srv.Shutdown(shutdownContext)
	}()

	return srv.ListenAndServe()
}

func monitoringServer(addr string, healthHandler http.Handler) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	if healthHandler != nil {
		mux.Handle("/health", healthHandler)
	} else {
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"msg":"ok"}`))
		})
	}
	return &http.Server{
		Addr:    addr,
		Handler: mux,
	}
}
