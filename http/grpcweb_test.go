package server_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"

	grpcserver "github.com/contiamo/goserver/grpc"
	"github.com/contiamo/goserver/grpc/test"
	. "github.com/contiamo/goserver/http"
)

var _ = Describe("Grpcweb", func() {
	It("should be possible to serve a grpc service over http", func() {
		grpcSrv, _ := grpcserver.New(&grpcserver.Config{
			Register: func(srv *grpc.Server) {
				test.RegisterPingPongServer(srv, test.NewPingPongServer())
			},
		})
		webServer, err := New(&Config{
			Options: []Option{
				WithGRPC(grpcSrv),
			},
		})
		Expect(err).NotTo(HaveOccurred())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go ListenAndServe(ctx, ":8000", webServer)
	})

})
