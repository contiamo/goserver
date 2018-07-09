package http

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// WithMetrics configures metrics collection
func WithMetrics(app string) Option {
	return &metricsOption{app}
}

type metricsOption struct{ app string }

func (opt *metricsOption) WrapHandler(handler http.Handler) (http.Handler, error) {
	latency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "latencies",
		Help:        "How long it took to process the request, partitioned by status code, method and HTTP path.",
		ConstLabels: prometheus.Labels{"service": opt.app},
		Buckets:     []float64{300, 1200, 5000},
	},
		[]string{"code", "method"},
	)
	requestsCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:        "requests",
			Help:        "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
			ConstLabels: prometheus.Labels{"service": opt.app},
		},
		[]string{"code", "method"},
	)
	prometheus.Unregister(latency)
	prometheus.Unregister(requestsCounter)
	if err := prometheus.Register(latency); err != nil {
		return nil, err
	}
	if err := prometheus.Register(requestsCounter); err != nil {
		return nil, err
	}
	handler = promhttp.InstrumentHandlerDuration(latency, handler)
	handler = promhttp.InstrumentHandlerCounter(requestsCounter, handler)
	return handler, nil
}
