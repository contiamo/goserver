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
			grpc.MaxSendMsgSize(1 << 12),
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
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span, _ := opentracing.StartSpanFromContext(r.Context(), "logic")
			defer span.Finish()
			io.Copy(w, r.Body)
		}),
		Options: []httpserver.Option{
			httpserver.WithLogging("http-echo"),
			httpserver.WithTracing("localhost:6831", "example", map[string]string{"example.namespace": "test"}, nil),
			httpserver.WithMetrics("http-echo"),
			httpserver.WithRecovery(os.Stderr, true),
		},
	})
	if err != nil {
		logrus.Fatal(err)
	}

	// start servers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go httpserver.ListenAndServe(ctx, ":8000", httpServer)
	go grpcserver.ListenAndServe(ctx, ":3001", grpcServer)
	go goserver.ListenAndServeMetricsAndHealth(":8080", nil)
	select {}
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
