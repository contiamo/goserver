package grpc

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/contiamo/goserver/grpc/test"
)

var _ = Describe("Credentials", func() {
	It("should be possible to load credentials via the WithCredentials() option", func() {
		crt := "test/pki/test-server.crt"
		key := "test/pki/test-server.key"
		ca := "test/pki/ca.crt"
		srv, err := createServerWithOptions([]Option{WithCredentials(crt, key, ca)})
		Expect(err).NotTo(HaveOccurred())
		Expect(srv).NotTo(BeNil())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go ListenAndServe(ctx, ":3001", srv)
		cli, err := createTestClient(ctx, "localhost:3001")
		Expect(err).NotTo(HaveOccurred())
		resp, err := cli.Ping(ctx, &test.PingReq{Msg: "test"})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Msg).To(Equal("test"))
	})

	It("should be possible to load credentials via the WithCredentials() option without a CA certificate", func() {
		crt := "test/pki/test-server.crt"
		key := "test/pki/test-server.key"
		ca := ""
		srv, err := createServerWithOptions([]Option{WithCredentials(crt, key, ca)})
		Expect(err).NotTo(HaveOccurred())
		Expect(srv).NotTo(BeNil())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go ListenAndServe(ctx, ":3002", srv)
		cli, err := createTestClient(ctx, "localhost:3002")
		Expect(err).NotTo(HaveOccurred())
		resp, err := cli.Ping(ctx, &test.PingReq{Msg: "test"})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Msg).To(Equal("test"))
	})

})
