# goserver

[![Go Report Card](https://goreportcard.com/badge/github.com/contiamo/goserver)](https://goreportcard.com/report/github.com/contiamo/goserver) [![Documentation](https://godoc.org/github.com/contiamo/goserver?status.svg)](http://godoc.org/github.com/contiamo/goserver)

## Scope

This package provides helpers to setup HTTP and gRPC servers following best practices.
It includes helpers for

- gRPC and HTTP
  - logging
  - tracing
  - metrics collection
  - recovery
- only for gRPC
  - credential loading
  - reflection

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

  ctx, cancel := context.WithCancel(context.Background()
  c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

  // start server
  go grpcserver.ListenAndServe(ctx, ":3001", grpcServer)
  // start /metrics endpoint
  go goserver.ListenAndServeMonitoring(ctx, ":8080", nil)

  <-c
  cancel()
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
    Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      span, _ := opentracing.StartSpanFromContext(r.Context(), "logic")
      defer span.Finish()
      io.Copy(w, r.Body)
    }),
    Options: []httpserver.Option{
      httpserver.WithLogging("http-echo"),
      httpserver.WithTracing("localhost:6831", "example", nil, nil),
      httpserver.WithMetrics("http-echo", nil),
      httpserver.WithRecovery(os.Stderr, true),
    },
  })
  if err != nil {
    logrus.Fatal(err)
  }

  ctx, cancel := context.WithCancel(context.Background()
  c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

  // start server
  go httpserver.ListenAndServe(ctx, ":8000", httpServer)
  // start /metrics endpoint
  go goserver.ListenAndServeMonitoring(ctx, ":8080", nil)

  <-c
  cancel()
}
```

## Using goserver as Middleware
It's not necessary to use goserver's server component. It is also possible to use it just as middleware in other servers. For example,
to use goserver's recovery and logging middleware with [chi](https://github.com/go-chi/chi), you can do the following:
```go
package main

import (
  "io"
  "net/http"

  goserver "github.com/contiamo/goserver/http"

  "github.com/go-chi/chi"
)

func main() {
  r := chi.NewRouter()

  // initialize goserver middleware
  logging := goserver.WithLogging("application-name")
  recovery := goserver.WithRecovery(os.Stderr, true)

  // tell the chi router to use the goserver middleware
  r.Use(logging.WrapHandler)
  r.Use(recovery.WrapHandler)

  // ... setup your routes and handlers...

  // start the server
  http.ListenAndServe(":8080", r)
}
```

To run the server and metrics on a different port:

```
  go func() {
  	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
  }()

  err = goserver.ListenAndServeMetricsAndHealth(:8081, nil)
  if err != nil {
	log.Fatal(err)
  }
```
