package server_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	. "github.com/contiamo/goserver/grpc"
	"github.com/contiamo/goserver/grpc/test"
)

var _ = Describe("Recovery", func() {
	It("should recover from panic when recovery option is set", func() {
		srv, err := createServerWithOptions([]Option{WithRecovery()})
		Expect(err).NotTo(HaveOccurred())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go ListenAndServe(ctx, ":3004", srv)
		cli, err := createPlaintextTestClient(ctx, "localhost:3004")
		Expect(err).NotTo(HaveOccurred())
		_, err = cli.Panic(ctx, &test.PingReq{})
		Expect(err).To(HaveOccurred())
		Expect(grpc.Code(err)).To(Equal(codes.Internal))
	})
})
