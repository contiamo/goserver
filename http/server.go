package server

import (
	"net/http"
	"time"
)

// Config contains the http servers config options
type Config struct {
	Addr    string
	Handler http.Handler
	Options []Option
}

// New creates a new http server
//
// Example:
//
// srv, _ := server.New(&server.Config{
//  Addr: ":8080",
//  Handler: http.DefaultServeMux,
//  Options: []server.Option{
//    server.WithLogging("my-server"),
//    server.WithMetrics("my-server"),
//    server.WithRecovery(),
//    server.WithTracing("opentracing-server:6831", "my-server"),
//  }
// })
//
// srv.ListenAndServe()
func New(cfg *Config) (*http.Server, error) {
	var (
		h   = cfg.Handler
		err error
	)
	for _, opt := range cfg.Options {
		h, err = opt.WrapHandler(h)
		if err != nil {
			return nil, err
		}
	}
	return &http.Server{
		Handler:        h,
		Addr:           cfg.Addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}, nil
}

type Option interface {
	WrapHandler(handler http.Handler) (http.Handler, error)
}
