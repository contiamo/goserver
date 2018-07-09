package grpc

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

// WithMetrics configures the server to collect usage metrics in prometheus format
func WithMetrics() Option {
	return &metricsOption{}
}

type metricsOption struct{}

func (opt *metricsOption) GetOptions() (grpc.ServerOption, grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor, error) {
	si := grpc_prometheus.StreamServerInterceptor
	ui := grpc_prometheus.UnaryServerInterceptor
	return nil, si, ui, nil
}

func (opt *metricsOption) PostProcess(s *grpc.Server) error {
	grpc_prometheus.Register(s)
	return nil
}
