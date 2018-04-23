goserver
========

## Scope

This package provides helpers to setup HTTP and gRPC servers following best practices.
It includes helpers for

* gRPC and HTTP
  * logging
  * tracing
  * metrics collection
  * recovery
* only for gRPC
  * credential loading
  * reflection

## Example

### gRPC
```go
package main

import (
	"context"
	"github.com/contiamo/goserver"
	grpcserver "github.com/contiamo/goserver/grpc"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	// setup grpc server with options
	grpcServer, err := grpcserver.New(&grpcserver.Config{
		Options: []grpcserver.Option{
          grpcserver.WithCredentials("cert.pem","key.pem","ca.pem"),
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

	// start server
	go grpcserver.ListenAndServe(":3001", grpcServer)
	// start /metrics endpoint
	goserver.ListenAndServeMetricsAndHealth(":8080", nil)
}
```

### HTTP
```go
package main

import (
	"io"
	"net/http"
	"github.com/contiamo/goserver"
	httpserver "github.com/contiamo/goserver/http"
	"github.com/sirupsen/logrus"
)

func main() {
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

	// start server
	go httpServer.ListenAndServe()

	// start /metrics endpoint
	goserver.ListenAndServeMetricsAndHealth(":8080", nil)
}
```
