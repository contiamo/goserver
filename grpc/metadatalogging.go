package grpc

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// WithMDLogging configures a logrus logger to log all gRPC requests with duration and return status
func WithMDLogging() Option {
	return &mdLoggingOption{}
}

type mdLoggingOption struct{}

func (opt *mdLoggingOption) GetOptions() (grpc.ServerOption, grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor, error) {
	si := StreamServerInterceptor()
	ui := UnaryServerInterceptor()
	return nil, si, ui, nil
}

func (opt *mdLoggingOption) PostProcess(s *grpc.Server) error {
	return nil
}

// UnaryServerInterceptor returns a new unary server interceptors that adds logrus.Entry to the context.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		fields := extractLoggingFields(ctx)
		logrus.WithFields(fields).Debug("processed request")
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that adds logrus.Entry to the context.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		fields := extractLoggingFields(stream.Context())
		logrus.WithFields(fields).Debug("processed stream")
		return handler(srv, stream)
	}
}

func extractLoggingFields(ctx context.Context) logrus.Fields {
	fields := logrus.Fields{}

	if ctxMd, ok := metadata.FromIncomingContext(ctx); ok {
		for field, values := range ctxMd {
			if isSecure(field) {
				fields[field] = "****"
				continue
			}
			fields[field] = strings.Join(values, ",")
		}
	}
	return fields
}


func isSecure(name string) bool {
	name = strings.ToLower(name)
	return strings.Contains(name, "auth") ||
		strings.Contains(name, "token") ||
		strings.Contains(name, "password") ||
		strings.Contains(name, "secret")
}
