package http

import (
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

// WithGRPC configures the webserver to serve grpc-web requests as specified in https://github.com/improbable-eng/grpc-web
func WithGRPC(srv *grpc.Server) Option {
	return &grpcwebOption{srv}
}

type grpcwebOption struct{ srv *grpc.Server }

func (opt *grpcwebOption) WrapHandler(handler http.Handler) (http.Handler, error) {
	wrappedGrpc := grpcweb.WrapServer(opt.srv)
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		if wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
			return
		}
		handler.ServeHTTP(resp, req)
	}), nil
}
