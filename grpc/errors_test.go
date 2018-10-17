package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc/codes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/contiamo/goserver/grpc/test"
)

var _ = Describe("ErrorScrubbing", func() {
	It("should hide error messages", func() {

		srv, err := createServerWithOptions([]Option{WithErrorScrubbing(nil)})
		Expect(err).NotTo(HaveOccurred())
		Expect(srv).NotTo(BeNil())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go ListenAndServe(ctx, ":3003", srv)
		// it takes some time to run the server, can't be accessed immediately
		time.Sleep(100 * time.Millisecond)

		cli, err := createPlaintextTestClient(ctx, "localhost:3003")
		Expect(err).NotTo(HaveOccurred())
		_, err = cli.Err(ctx, &test.ErrorReq{Code: uint32(codes.Internal), Msg: "test"})
		Expect(err).To(HaveOccurred())

		st, _ := status.FromError(err)
		Expect(st.Message()).To(Equal("Internal Server Error"))

		_, err = cli.Err(ctx, &test.ErrorReq{Code: uint32(codes.Unknown), Msg: "test"})
		Expect(err).To(HaveOccurred())
		st, _ = status.FromError(err)
		Expect(st.Message()).To(Equal("Internal Server Error"))
	})

	It("should show error messages with the NoOpFunc", func() {

		srv, err := createServerWithOptions([]Option{WithErrorScrubbing(NoopErrorScrubber)})
		Expect(err).NotTo(HaveOccurred())
		Expect(srv).NotTo(BeNil())
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go ListenAndServe(ctx, ":3003", srv)
		// it takes some time to run the server, can't be accessed immediately
		time.Sleep(100 * time.Millisecond)

		cli, err := createPlaintextTestClient(ctx, "localhost:3003")
		Expect(err).NotTo(HaveOccurred())
		_, err = cli.Err(ctx, &test.ErrorReq{Code: uint32(codes.Internal), Msg: "test 1"})
		Expect(err).To(HaveOccurred())

		st, _ := status.FromError(err)
		Expect(st.Message()).To(Equal("test 1"))

		_, err = cli.Err(ctx, &test.ErrorReq{Code: uint32(codes.Unknown), Msg: "test 2"})
		Expect(err).To(HaveOccurred())

		st, _ = status.FromError(err)
		Expect(st.Message()).To(Equal("test 2"))
	})
})
