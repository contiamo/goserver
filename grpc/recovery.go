package grpc

import (
	"runtime/debug"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// WithRecovery recovers the server from panics caused inside of the gRPC calls
func WithRecovery() Option {
	return &recoveryOption{}
}

type recoveryOption struct{}

func (opt *recoveryOption) GetOptions() (grpc.ServerOption, grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor, error) {
	recOpt := grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		switch err := p.(type) {
		case error:
			logrus.WithFields(logrus.Fields{
				"stacktrace": string(debug.Stack()),
			}).Error(err)
		}
		return status.Errorf(codes.Internal, "%s", p)
	})
	si := grpc_recovery.StreamServerInterceptor(recOpt)
	ui := grpc_recovery.UnaryServerInterceptor(recOpt)
	return nil, si, ui, nil
}

func (opt *recoveryOption) PostProcess(s *grpc.Server) error {
	return nil
}
