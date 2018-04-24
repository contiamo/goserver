package server

import (
	"context"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Config contains the configuration for your server. See New() for example usage
type Config struct {
	Options  []Option
	Extras   []grpc.ServerOption
	Register func(*grpc.Server)
}

// New creates a new gRPC server with the given options
//
// Example:
//
// srv, _ := server.New(&server.Config{
//   Options: []server.Option{
//    server.WithCredentials("cert.pem","key.pem","ca.pem"),
//    server.WithTracing("opentracing:6818", "my-service")
//    server.WithLogging("my-service"),
//    server.WithMetrics("my-service"),
//    server.WithRecovery(),
//    server.WithReflection()),
//   },
//   Extras: []grpc.ServerOption{
//    grpc.MaxReceiveMsgSize(4<<12),
//    grpc.MaxSendMsgSize(4<<12),
//   },
//   Register: func(srv, *grpc.Server){
//    myservice.RegisterMyFooServiceServer(srv, myFooServiceImpl)
//    myservice.RegisterMyBarServiceServer(srv, myBarServiceImpl)
//   },
// }
//
// server.ListenAndServe(":3001", srv)
func New(cfg *Config) (*grpc.Server, error) {
	streamServerInterceptors := []grpc.StreamServerInterceptor{}
	unaryServerInterceptors := []grpc.UnaryServerInterceptor{}
	serverOptions := []grpc.ServerOption{}
	for _, opt := range cfg.Options {
		o, si, ui, err := opt.GetOptions()
		if err != nil {
			return nil, err
		}
		if si != nil {
			streamServerInterceptors = append(streamServerInterceptors, si)
		}
		if ui != nil {
			unaryServerInterceptors = append(unaryServerInterceptors, ui)
		}
		if o != nil {
			serverOptions = append(serverOptions, o)
		}
	}
	if len(streamServerInterceptors) > 0 {
		serverOptions = append(serverOptions, grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamServerInterceptors...)))
	}
	if len(unaryServerInterceptors) > 0 {
		serverOptions = append(serverOptions, grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryServerInterceptors...)))
	}
	if len(cfg.Extras) > 0 {
		serverOptions = append(serverOptions, cfg.Extras...)
	}
	srv := grpc.NewServer(serverOptions...)
	cfg.Register(srv)
	for _, opt := range cfg.Options {
		if err := opt.PostProcess(srv); err != nil {
			return nil, err
		}
	}
	return srv, nil
}

// ListenAndServe serves an gRPC server over TCP
func ListenAndServe(ctx context.Context, addr string, srv *grpc.Server) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()
	logrus.Info("start listening for gRPC requests on " + addr)
	return srv.Serve(lis)
}

// Option is the interface for the package supplied configuration helpers
type Option interface {
	GetOptions() (grpc.ServerOption, grpc.StreamServerInterceptor, grpc.UnaryServerInterceptor, error)
	PostProcess(s *grpc.Server) error
}
