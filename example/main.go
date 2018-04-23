package main

import (
	"context"
	"io"
	"net/http"
	"os"

	"github.com/contiamo/goserver"
	grpcserver "github.com/contiamo/goserver/grpc"
	httpserver "github.com/contiamo/goserver/http"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	// setup grpc server with options
	grpcServer, err := grpcserver.New(&grpcserver.Config{
		Options: []grpcserver.Option{
			grpcserver.WithTracing("localhost:6831", "example"),
			grpcserver.WithLogging("grpc-echo"),
			grpcserver.WithMetrics(),
			grpcserver.WithRecovery(),
			grpcserver.WithReflection(),
		},
		Extras: []grpc.ServerOption{
			grpc.MaxSendMsgSize(4 << 12),
		},
		Register: func(srv *grpc.Server) {
			RegisterEchoServer(srv, &echoServer{})
		},
	})
	if err != nil {
		logrus.Fatal(err)
	}

	// setup http server with options
	httpServer, err := httpserver.New(&httpserver.Config{
		Addr: ":8000",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span, _ := opentracing.StartSpanFromContext(r.Context(), "logic")
			defer span.Finish()
			io.Copy(w, r.Body)
		}),
		Options: []httpserver.Option{
			httpserver.WithLogging("http-echo"),
			httpserver.WithTracing("localhost:6831", "example"),
			httpserver.WithMetrics("http-echo"),
			httpserver.WithRecovery(os.Stderr, true),
		},
	})
	if err != nil {
		logrus.Fatal(err)
	}

	// start servers
	go httpServer.ListenAndServe()
	go grpcserver.ListenAndServe(":3001", grpcServer)

	// start /metrics endpoint
	goserver.ListenAndServeMetricsAndHealth(":8080")
}

// example grpc server
type echoServer struct{}

func (srv *echoServer) Echo(ctx context.Context, req *EchoRequest) (*EchoResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "logic")
	defer span.Finish()
	return &EchoResponse{
		Data: req.Data,
	}, nil
}
