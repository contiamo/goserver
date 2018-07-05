package grpc

import (
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// WithLogging configures a logrus logger to log all gRPC requests with duration and return status
func WithLogging(app string) Option {
	return &loggingOption{app}
}

type loggingOption struct{ app string }

func (opt *loggingOption) GetOptions() (grpc.ServerOption, grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor, error) {
	logrusEntry := logrus.WithFields(logrus.Fields{"component": opt.app})
	logOpts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(func(duration time.Duration) (key string, value interface{}) {
			return "grpc.time_ns", duration.Nanoseconds()
		}),
	}
	grpc_logrus.ReplaceGrpcLogger(logrusEntry)
	si := grpc_logrus.StreamServerInterceptor(logrusEntry, logOpts...)
	ui := grpc_logrus.UnaryServerInterceptor(logrusEntry, logOpts...)
	return nil, si, ui, nil
}

func (opt *loggingOption) PostProcess(s *grpc.Server) error {
	return nil
}
