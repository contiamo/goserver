package grpc

import (
	"github.com/contiamo/go-base/v2/pkg/tracing"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	opentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

// WithTracing configures the server to support tracing in the opentracing format.
// It attaches a root span to the context of all incoming requests.
//
// The tracer is configured via the JAEGER_* env variables.
//
// The server variable is deprecated and will be ignored.
func WithTracing(serverDeprecated, serviceName string) Option {
	if err := tracing.InitJaeger(serviceName); err != nil {
		panic(err)
	}
	return &tracingOption{
		server:      serverDeprecated,
		serviceName: serviceName,
	}
}

type tracingOption struct {
	server, serviceName string
}

func (opt *tracingOption) GetOptions() (grpc.ServerOption, grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor, error) {
	tracer := opentracing.GlobalTracer()

	ui := otgrpc.OpenTracingServerInterceptor(tracer)
	si := otgrpc.OpenTracingStreamServerInterceptor(tracer)
	return nil, si, ui, nil
}

func (opt *tracingOption) PostProcess(s *grpc.Server) error {
	return nil
}
