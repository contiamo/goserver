package server

import (
	"github.com/contiamo/goserver"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

// WithTracing configures the server to support tracing in the opentracing format.
// It attaches a root span to the context of all incoming requests and configures the tracer to send the traces to the configured opentracing server.
func WithTracing(server, serviceName string) Option {
	return &tracingOption{server, serviceName}
}

type tracingOption struct {
	server, serviceName string
}

func (opt *tracingOption) GetOptions() (grpc.ServerOption, grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor, error) {
	if err := goserver.InitTracer(opt.server, opt.serviceName); err != nil {
		return nil, nil, nil, err
	}
	ui := otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer(), otgrpc.SpanDecorator(func(
		span opentracing.Span,
		method string,
		req, resp interface{},
		grpcError error) {
	}))
	return nil, nil, ui, nil
}

func (opt *tracingOption) PostProcess(s *grpc.Server) error {
	return nil
}
