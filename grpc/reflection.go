package server

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// WithReflection configures the server to provide an API for external tools to get informations about the used services and types.
// It enables a service which returns the proto definitions for the methods, so that you can use tools like [grpcurl](https://github.com/fullstorydev/grpcurl) without the need to supply the protofile to interact with the service.
func WithReflection() Option {
	return &reflectionOption{}
}

type reflectionOption struct{}

func (opt *reflectionOption) GetOptions() (grpc.ServerOption, grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor, error) {
	return nil, nil, nil, nil
}

func (opt *reflectionOption) PostProcess(s *grpc.Server) error {
	reflection.Register(s)
	return nil
}
