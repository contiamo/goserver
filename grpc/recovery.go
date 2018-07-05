package grpc

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

// WithRecovery recovers the server from panics caused inside of the gRPC calls
func WithRecovery() Option {
	return &recoveryOption{}
}

type recoveryOption struct{}

func (opt *recoveryOption) GetOptions() (grpc.ServerOption, grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor, error) {
	si := grpc_recovery.StreamServerInterceptor()
	ui := grpc_recovery.UnaryServerInterceptor()
	return nil, si, ui, nil
}

func (opt *recoveryOption) PostProcess(s *grpc.Server) error {
	return nil
}
